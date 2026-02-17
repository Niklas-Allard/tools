package qrcode

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{}

func RegisterRoutes(rg *gin.RouterGroup, _ *gorm.DB) {
	h := &Handler{}
	rg.POST("/generate", h.GenerateQR)
}

func (h *Handler) GenerateQR(c *gin.Context) {
	var req QRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pngBytes, err := GeneratePNG(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=\"qrcode.png\"")
	c.Writer.Write(pngBytes)
}
