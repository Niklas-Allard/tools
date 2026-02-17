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
	rg.GET("/status/:jobId", h.GetStatus)
	rg.GET("/downloads", h.ListDownloads)
	rg.GET("/download-file/:jobId", h.DownloadFile)
	rg.GET("/download-folder/:jobId", h.DownloadFolder)
	rg.DELETE("/download/:jobId", h.DeleteDownload)
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
