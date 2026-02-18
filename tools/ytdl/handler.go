package ytdl

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := NewHandler(db)
	rg.POST("/download", h.StartDownload)
	rg.POST("/batch-download", h.StartBatchDownload)
	rg.GET("/status/:jobId", h.GetStatus)
	rg.GET("/batch-status/:batchId", h.GetBatchStatus)
	rg.GET("/downloads", h.ListDownloads)
	rg.GET("/download-file/:jobId", h.DownloadFile)
	rg.GET("/download-folder/:jobId", h.DownloadFolder)
	rg.GET("/download-batch/:batchId", h.DownloadBatch)
	rg.DELETE("/download/:jobId", h.DeleteDownload)
	rg.DELETE("/batch/:batchId", h.DeleteBatch)
}

func (h *Handler) StartDownload(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jobID := fmt.Sprintf("%d", time.Now().UnixNano())
	download := Download{
		JobID:      jobID,
		URL:        req.URL,
		Status:     "pending",
		FilePath:   filepath.Join(DownloadDir, jobID),
		IsPlaylist: req.Playlist,
		CreatedAt:  time.Now(),
		Progress:   0,
	}
	h.db.Create(&download)

	go PerformDownload(h.db, download.ID, req)

	c.JSON(200, gin.H{
		"job_id":  jobID,
		"message": "Download gestartet",
	})
}

func (h *Handler) GetStatus(c *gin.Context) {
	jobID := c.Param("jobId")
	var download Download
	if err := h.db.Where("job_id = ?", jobID).First(&download).Error; err != nil {
		c.JSON(404, gin.H{"error": "Job nicht gefunden"})
		return
	}
	c.JSON(200, download)
}

func (h *Handler) ListDownloads(c *gin.Context) {
	var downloads []Download
	h.db.Order("created_at DESC").Find(&downloads)
	c.JSON(200, downloads)
}

func (h *Handler) DownloadFile(c *gin.Context) {
	jobID := c.Param("jobId")
	var download Download
	if err := h.db.Where("job_id = ?", jobID).First(&download).Error; err != nil {
		c.JSON(404, gin.H{"error": "Download nicht gefunden"})
		return
	}

	if download.Status != "completed" {
		c.JSON(400, gin.H{"error": "Download nicht abgeschlossen"})
		return
	}

	files, _ := filepath.Glob(filepath.Join(download.FilePath, "*"))
	var mainFile string
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f))
		if ext == ".mp4" || ext == ".mkv" || ext == ".webm" ||
			ext == ".mp3" || ext == ".m4a" || ext == ".opus" {
			mainFile = f
			break
		}
	}

	if mainFile == "" && len(files) > 0 {
		mainFile = files[0]
	}
	if mainFile == "" {
		c.JSON(404, gin.H{"error": "Keine Dateien gefunden"})
		return
	}

	c.FileAttachment(mainFile, filepath.Base(mainFile))
}

func (h *Handler) DownloadFolder(c *gin.Context) {
	jobID := c.Param("jobId")
	var download Download
	if err := h.db.Where("job_id = ?", jobID).First(&download).Error; err != nil {
		c.JSON(404, gin.H{"error": "Download nicht gefunden"})
		return
	}

	if download.Status != "completed" {
		c.JSON(400, gin.H{"error": "Download nicht abgeschlossen"})
		return
	}

	zipPath := filepath.Join(DownloadDir, jobID+".zip")

	func() {
		zipFile, err := os.Create(zipPath)
		if err != nil {
			return
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		filepath.Walk(download.FilePath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			zipEntry, err := zipWriter.Create(filepath.Base(path))
			if err != nil {
				return nil
			}
			io.Copy(zipEntry, file)
			return nil
		})
	}()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		c.JSON(500, gin.H{"error": "ZIP konnte nicht erstellt werden"})
		return
	}

	c.FileAttachment(zipPath, jobID+".zip")

	go func() {
		time.Sleep(10 * time.Second)
		os.Remove(zipPath)
	}()
}

func (h *Handler) DeleteDownload(c *gin.Context) {
	jobID := c.Param("jobId")
	var download Download
	if err := h.db.Where("job_id = ?", jobID).First(&download).Error; err != nil {
		c.JSON(404, gin.H{"error": "Download nicht gefunden"})
		return
	}

	os.RemoveAll(download.FilePath)
	h.db.Delete(&download)
	c.JSON(200, gin.H{"message": "Download gelöscht"})
}

func (h *Handler) StartBatchDownload(c *gin.Context) {
	var req BatchDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(req.URLs) == 0 {
		c.JSON(400, gin.H{"error": "Keine URLs angegeben"})
		return
	}

	batchID := fmt.Sprintf("batch_%d", time.Now().UnixNano())
	sharedFolder := filepath.Join(DownloadDir, batchID)
	os.MkdirAll(sharedFolder, 0755)

	jobIDs := []string{}
	for i, url := range req.URLs {
		if strings.TrimSpace(url) == "" {
			continue
		}
		jobID := fmt.Sprintf("%d_%d", time.Now().UnixNano(), i)
		download := Download{
			JobID:     jobID,
			URL:       strings.TrimSpace(url),
			Status:    "pending",
			FilePath:  sharedFolder,
			BatchID:   batchID,
			CreatedAt: time.Now(),
			Progress:  0,
		}
		h.db.Create(&download)

		singleReq := DownloadRequest{
			URL:                strings.TrimSpace(url),
			Format:             req.Format,
			Quality:            req.Quality,
			AudioOnly:          req.AudioOnly,
			AudioFormat:        req.AudioFormat,
			RemuxVideo:         req.RemuxVideo,
			FPS:                req.FPS,
			Playlist:           req.Playlist,
			MaxDownloads:       req.MaxDownloads,
			EmbedChapters:      req.EmbedChapters,
			EmbedMetadata:      req.EmbedMetadata,
			EmbedThumbnail:     req.EmbedThumbnail,
			EmbedSubtitles:     req.EmbedSubtitles,
			SubtitlesLang:      req.SubtitlesLang,
			WriteInfoJSON:      req.WriteInfoJSON,
			WriteDescription:   req.WriteDescription,
			WriteThumbnail:     req.WriteThumbnail,
			SplitChapters:      req.SplitChapters,
			MaxFilesize:        req.MaxFilesize,
			MinFilesize:        req.MinFilesize,
			RateLimit:          req.RateLimit,
			CookiesFromBrowser: req.CookiesFromBrowser,
			Sponsorblock:       req.Sponsorblock,
			AudioMultistreams:  req.AudioMultistreams,
			OutputTemplate:     req.OutputTemplate,
		}
		go PerformDownload(h.db, download.ID, singleReq)
		jobIDs = append(jobIDs, jobID)
	}

	c.JSON(200, gin.H{
		"batch_id": batchID,
		"job_ids":  jobIDs,
		"message":  fmt.Sprintf("Batch Download gestartet (%d URLs)", len(jobIDs)),
	})
}

func (h *Handler) GetBatchStatus(c *gin.Context) {
	batchID := c.Param("batchId")
	var downloads []Download
	if err := h.db.Where("batch_id = ?", batchID).Find(&downloads).Error; err != nil {
		c.JSON(404, gin.H{"error": "Batch nicht gefunden"})
		return
	}

	completed, failed, total := 0, 0, len(downloads)
	for _, dl := range downloads {
		if dl.Status == "completed" {
			completed++
		} else if dl.Status == "failed" {
			failed++
		}
	}

	overallStatus := "downloading"
	if completed+failed == total && total > 0 {
		if failed == total {
			overallStatus = "failed"
		} else {
			overallStatus = "completed"
		}
	}

	c.JSON(200, gin.H{
		"batch_id":  batchID,
		"downloads": downloads,
		"total":     total,
		"completed": completed,
		"failed":    failed,
		"status":    overallStatus,
	})
}

func (h *Handler) DownloadBatch(c *gin.Context) {
	batchID := c.Param("batchId")
	var download Download
	if err := h.db.Where("batch_id = ?", batchID).First(&download).Error; err != nil {
		c.JSON(404, gin.H{"error": "Batch nicht gefunden"})
		return
	}

	zipPath := filepath.Join(DownloadDir, batchID+".zip")

	func() {
		zipFile, err := os.Create(zipPath)
		if err != nil {
			return
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		filepath.Walk(download.FilePath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			zipEntry, err := zipWriter.Create(filepath.Base(path))
			if err != nil {
				return nil
			}
			io.Copy(zipEntry, file)
			return nil
		})
	}()

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		c.JSON(500, gin.H{"error": "ZIP konnte nicht erstellt werden"})
		return
	}

	c.FileAttachment(zipPath, batchID+".zip")

	go func() {
		time.Sleep(10 * time.Second)
		os.Remove(zipPath)
	}()
}

func (h *Handler) DeleteBatch(c *gin.Context) {
	batchID := c.Param("batchId")
	var downloads []Download
	if err := h.db.Where("batch_id = ?", batchID).Find(&downloads).Error; err != nil || len(downloads) == 0 {
		c.JSON(404, gin.H{"error": "Batch nicht gefunden"})
		return
	}

	// All batch downloads share the same folder
	os.RemoveAll(downloads[0].FilePath)
	h.db.Where("batch_id = ?", batchID).Delete(&Download{})
	c.JSON(200, gin.H{"message": "Batch gelöscht"})
}
