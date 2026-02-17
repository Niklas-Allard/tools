package qrcode

// QRRequest holds all parameters for generating a QR code.
type QRRequest struct {
	Content string `json:"content" binding:"required"`
	Size    int    `json:"size"`     // pixels, default 512
	Level   string `json:"level"`    // L, M, Q, H — default M
	FGColor string `json:"fg_color"` // hex, default #000000
	BGColor string `json:"bg_color"` // hex, default #ffffff
}

func (r *QRRequest) applyDefaults() {
	if r.Size <= 0 || r.Size > 4096 {
		r.Size = 512
	}
	if r.Level == "" {
		r.Level = "M"
	}
	if r.FGColor == "" {
		r.FGColor = "#000000"
	}
	if r.BGColor == "" {
		r.BGColor = "#ffffff"
	}
}
