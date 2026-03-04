package main

import (
	"devtools/api"
	"devtools/tools/ytdl"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("ytdl.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&ytdl.Download{})

	os.MkdirAll(ytdl.DownloadDir, 0755)
	go ytdl.StartCleanupJob(db)

	r := gin.Default()
	r.Use(cors.Default())

	// COOP/COEP Middleware – WICHTIG: vor Static!
	r.Use(func(c *gin.Context) {
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Resource-Policy", "cross-origin") // CDN
		c.Next()
	})

	api.SetupRoutes(r, db)

	// Serve frontend
	r.Static("/assets", "./frontend-vite/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend-vite/dist/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8") // Explizit
		c.File("./frontend-vite/dist/index.html")
	})

	r.Run(":8080")
}
