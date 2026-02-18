package ytdl

import "time"

const DownloadDir = "./downloads"

type Download struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	JobID      string    `json:"job_id" gorm:"uniqueIndex"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	Status     string    `json:"status"` // pending, downloading, completed, failed
	FilePath   string    `json:"file_path"`
	IsPlaylist bool      `json:"is_playlist"`
	BatchID    string    `json:"batch_id" gorm:"index"`
	CreatedAt  time.Time `json:"created_at"`
	Progress   int       `json:"progress"`
}

type BatchDownloadRequest struct {
	URLs []string `json:"urls" binding:"required"`

	Format      string `json:"format"`
	Quality     string `json:"quality"`
	AudioOnly   bool   `json:"audio_only"`
	AudioFormat string `json:"audio_format"`
	RemuxVideo  string `json:"remux_video"`
	FPS         string `json:"fps"`

	Playlist     bool `json:"playlist"`
	MaxDownloads int  `json:"max_downloads"`

	EmbedChapters  bool `json:"embed_chapters"`
	EmbedMetadata  bool `json:"embed_metadata"`
	EmbedThumbnail bool `json:"embed_thumbnail"`
	EmbedSubtitles bool `json:"embed_subtitles"`

	SubtitlesLang string `json:"subtitles_lang"`

	WriteInfoJSON    bool `json:"write_info_json"`
	WriteDescription bool `json:"write_description"`
	WriteThumbnail   bool `json:"write_thumbnail"`
	SplitChapters    bool `json:"split_chapters"`

	MaxFilesize string `json:"max_filesize"`
	MinFilesize string `json:"min_filesize"`

	RateLimit string `json:"rate_limit"`

	CookiesFromBrowser string `json:"cookies_from_browser"`
	Sponsorblock       bool   `json:"sponsorblock"`
	AudioMultistreams  bool   `json:"audio_multistreams"`
	OutputTemplate     string `json:"output_template"`
}

type DownloadRequest struct {
	// Basic
	URL string `json:"url" binding:"required"`

	// Format & Quality
	Format      string `json:"format"`
	Quality     string `json:"quality"`
	AudioOnly   bool   `json:"audio_only"`
	AudioFormat string `json:"audio_format"`
	RemuxVideo  string `json:"remux_video"`
	FPS         string `json:"fps"`

	// Playlist
	Playlist     bool `json:"playlist"`
	MaxDownloads int  `json:"max_downloads"`

	// Embedding
	EmbedChapters  bool `json:"embed_chapters"`
	EmbedMetadata  bool `json:"embed_metadata"`
	EmbedThumbnail bool `json:"embed_thumbnail"`
	EmbedSubtitles bool `json:"embed_subtitles"`

	// Subtitles
	SubtitlesLang string `json:"subtitles_lang"`

	// Additional Files
	WriteInfoJSON    bool `json:"write_info_json"`
	WriteDescription bool `json:"write_description"`
	WriteThumbnail   bool `json:"write_thumbnail"`
	SplitChapters    bool `json:"split_chapters"`

	// Size Limits
	MaxFilesize string `json:"max_filesize"`
	MinFilesize string `json:"min_filesize"`

	// Rate Limiting
	RateLimit string `json:"rate_limit"`

	// Advanced
	CookiesFromBrowser string `json:"cookies_from_browser"`
	Sponsorblock       bool   `json:"sponsorblock"`
	AudioMultistreams  bool   `json:"audio_multistreams"`
	OutputTemplate     string `json:"output_template"`
}
