package qr

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"os"
	"path/filepath"
)

type SimpleQRCode struct {
	Content string
	Size    int
}

func (code *SimpleQRCode) Generate() ([]byte, error) {
	qrCode, err := qrcode.Encode(code.Content, qrcode.Medium, code.Size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

func (code *SimpleQRCode) SaveToFile(imagePath string) (string, error) {
	qr, err := qrcode.New(code.Content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать QR-код: %w", err)
	}
	if err := os.MkdirAll(imagePath, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории: %w", err)
	}
	img := qr.Image(code.Size)
	filename := fmt.Sprintf("%s.jpg", uuid.New().String())
	fullPath := filepath.Join(imagePath, filename)
	dest, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer dest.Close()
	opt := jpeg.Options{
		Quality: 90,
	}
	if err := jpeg.Encode(dest, img, &opt); err != nil {
		return "", fmt.Errorf("ошибка кодирования JPEG: %w", err)
	}

	return filename, nil

}
