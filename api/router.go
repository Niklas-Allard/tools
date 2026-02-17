package api

import (
	"devtools/tools/qrcode"
	"devtools/tools/ytdl"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes registers all tool API routes under /api.
// Add new tools here by calling their RegisterRoutes function.
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	ytdl.RegisterRoutes(api.Group("/ytdl"), db)
	qrcode.RegisterRoutes(api.Group("/qrcode"), db)
}
