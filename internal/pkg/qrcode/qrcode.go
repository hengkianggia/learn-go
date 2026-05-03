package qrcode

import (
	"fmt"
	"os"
	"path/filepath"

	goqrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCodePNG(basePath, ticketCode string) (string, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create QR storage directory: %w", err)
	}

	filename := ticketCode + ".png"
	fullPath := filepath.Join(basePath, filename)

	if err := goqrcode.WriteFile(ticketCode, goqrcode.Medium, 256, fullPath); err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return fullPath, nil
}
