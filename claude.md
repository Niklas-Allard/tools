🚀 YouTube Downloader Pro - Claude AI Prompts
🎯 Project Description

An ultimate YouTube Downloader web app with Gin (Go backend) + Vue.js frontend. Features:

    25+ yt-dlp options (format, quality, chapters, metadata, SponsorBlock, etc.)

    Auto-cleanup after 24h

    ZIP download for playlists/folders

    Proxmox-ready LXC container

    ARM Windows compatible (./yt-dlp.exe)

📁 Project Structure

text
yt-downloader/
├── main.go                 ← Backend (Gin + SQLite + yt-dlp)
├── yt-dlp.exe             ← yt-dlp binary
├── ytdl.db                ← SQLite DB
├── downloads/             ← Temporary downloads
├── frontend/
│   ├── dist/              ← Vue build
│   └── src/App.vue        ← Advanced UI
└── go.mod


