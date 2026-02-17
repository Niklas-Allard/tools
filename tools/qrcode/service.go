package qrcode

import (
	"fmt"
	"image/color"
	"strconv"

	goqrcode "github.com/skip2/go-qrcode"
)

var levelMap = map[string]goqrcode.RecoveryLevel{
	"L": goqrcode.Low,
	"M": goqrcode.Medium,
	"Q": goqrcode.High,
	"H": goqrcode.Highest,
}

// GeneratePNG generates a QR code and returns it as a PNG byte slice.
func GeneratePNG(req QRRequest) ([]byte, error) {
	req.applyDefaults()

	level, ok := levelMap[req.Level]
	if !ok {
		level = goqrcode.Medium
	}

	qr, err := goqrcode.New(req.Content, level)
	if err != nil {
		return nil, fmt.Errorf("qr generation failed: %w", err)
	}

	fg, err := parseHexColor(req.FGColor)
	if err != nil {
		fg = color.Black
	}
	bg, err := parseHexColor(req.BGColor)
	if err != nil {
		bg = color.White
	}

	qr.ForegroundColor = fg
	qr.BackgroundColor = bg

	return qr.PNG(req.Size)
}

func parseHexColor(s string) (color.Color, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 6 {
		return nil, fmt.Errorf("invalid hex color: %s", s)
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return nil, err
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return nil, err
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return nil, err
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}
