package ytdl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func PerformDownload(db *gorm.DB, downloadID uint, req DownloadRequest) {
	var download Download
	db.First(&download, downloadID)

	download.Status = "downloading"
	download.Progress = 10
	db.Save(&download)

	outputPath := download.FilePath
	os.MkdirAll(outputPath, 0755)

	args := []string{}

	// Output Template
	if req.OutputTemplate != "" {
		args = append(args, "-o", filepath.Join(outputPath, req.OutputTemplate))
	} else {
		args = append(args, "-o", filepath.Join(outputPath, "%(title)s.%(ext)s"))
	}

	// Format & Quality
	if req.AudioOnly {
		// Select only audio streams – prevents any video file from being downloaded.
		// Without "-f bestaudio", yt-dlp defaults to bestvideo+bestaudio and then
		// tries to extract audio from the merged video, which often fails on Windows.
		args = append(args, "-f", "bestaudio/best")
		args = append(args, "-x")
		if req.AudioFormat != "" {
			args = append(args, "--audio-format", req.AudioFormat)
		}
		args = append(args, "--audio-quality", "0")
	} else {
		format := req.Format
		if format == "" {
			if req.Quality != "" {
				format = fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best", req.Quality)
			} else {
				format = "bestvideo+bestaudio/best"
			}
		}
		args = append(args, "-f", format)

		if req.RemuxVideo != "" {
			args = append(args, "--remux-video", req.RemuxVideo)
		} else {
			args = append(args, "--merge-output-format", "mp4")
		}
	}

	// FPS Filter
	if req.FPS != "" {
		args = append(args, "-S", fmt.Sprintf("fps:%s", req.FPS))
	}

	// Playlist
	if req.Playlist {
		args = append(args, "--yes-playlist")
		if req.MaxDownloads > 0 {
			args = append(args, "--max-downloads", strconv.Itoa(req.MaxDownloads))
		}
	} else {
		args = append(args, "--no-playlist")
	}

	// Embedding
	if req.EmbedChapters {
		args = append(args, "--embed-chapters")
	}
	if req.EmbedMetadata {
		args = append(args, "--embed-metadata")
	}
	if req.EmbedThumbnail {
		args = append(args, "--embed-thumbnail")
	}
	if req.EmbedSubtitles {
		args = append(args, "--embed-subs")
	}

	// Subtitles
	if req.SubtitlesLang != "" {
		args = append(args, "--write-subs", "--sub-lang", req.SubtitlesLang)
		if !req.EmbedSubtitles {
			args = append(args, "--write-auto-subs")
		}
	}

	// Additional Files
	if req.WriteInfoJSON {
		args = append(args, "--write-info-json")
	}
	if req.WriteDescription {
		args = append(args, "--write-description")
	}
	if req.WriteThumbnail {
		args = append(args, "--write-thumbnail")
	}
	if req.SplitChapters {
		args = append(args, "--split-chapters")
	}

	// Size Limits
	if req.MaxFilesize != "" {
		args = append(args, "--max-filesize", req.MaxFilesize)
	}
	if req.MinFilesize != "" {
		args = append(args, "--min-filesize", req.MinFilesize)
	}

	// Rate Limit
	if req.RateLimit != "" {
		args = append(args, "--limit-rate", req.RateLimit)
	}

	// Cookies
	if req.CookiesFromBrowser != "" {
		args = append(args, "--cookies-from-browser", req.CookiesFromBrowser)
	}

	// SponsorBlock
	if req.Sponsorblock {
		args = append(args, "--sponsorblock-remove", "all")
	}

	// Audio Multistreams
	if req.AudioMultistreams {
		args = append(args, "--audio-multistreams")
	}

	args = append(args, "--newline", "--no-warnings")
	args = append(args, req.URL)

	log.Printf("[ytdl] AudioOnly=%v Command: ./yt-dlp.exe %s", req.AudioOnly, strings.Join(args, " "))

	cmd := exec.Command("./yt-dlp.exe", args...)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()

	log.Printf("[ytdl] Output:\n%s", string(output))

	if err != nil {
		log.Printf("[ytdl] Error: %v", err)
		download.Status = "failed"
		download.Title = "Fehler: " + err.Error() + " | " + string(output)
		download.Progress = 0
		db.Save(&download)
		return
	}

	// After audio-only downloads, remove leftover intermediate video files.
	// yt-dlp sometimes fails to clean up video files on Windows when using
	// SponsorBlock or thumbnail embedding together with -x (audio extraction).
	if req.AudioOnly {
		videoExts := map[string]bool{
			".mp4": true, ".mkv": true, ".webm": true,
			".avi": true, ".mov": true, ".flv": true,
		}
		if files, globErr := filepath.Glob(filepath.Join(outputPath, "*")); globErr == nil {
			for _, f := range files {
				if videoExts[strings.ToLower(filepath.Ext(f))] {
					os.Remove(f)
				}
			}
		}
	}

	// Parse title from output
	lines := strings.Split(string(output), "\n")
	title := "Download erfolgreich"
	for _, line := range lines {
		if strings.Contains(line, "[download] Destination:") {
			parts := strings.Split(line, "Destination:")
			if len(parts) > 1 {
				title = strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(line, "[Merger]") {
			parts := strings.Split(line, "into")
			if len(parts) > 1 {
				title = strings.TrimSpace(parts[1])
			}
		}
	}

	download.Status = "completed"
	download.Title = title
	download.Progress = 100
	db.Save(&download)
}

func StartCleanupJob(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		var downloads []Download
		cutoff := time.Now().Add(-24 * time.Hour)
		db.Where("created_at < ? AND status = ?", cutoff, "completed").Find(&downloads)

		for _, dl := range downloads {
			os.RemoveAll(dl.FilePath)
			zipPath := filepath.Join(DownloadDir, dl.JobID+".zip")
			os.Remove(zipPath)
			db.Delete(&dl)
		}
	}
}
