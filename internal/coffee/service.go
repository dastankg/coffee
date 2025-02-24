package coffee

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (handler *CoffeeHandler) parseNumericValues(r *http.Request) (price, dollar, ruble float64, err error) {
	price, err = strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("некорректная цена: %w", err)
	}

	dollar, err = strconv.ParseFloat(r.FormValue("dollar"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("некорректное значение в долларах: %w", err)
	}

	ruble, err = strconv.ParseFloat(r.FormValue("ruble"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("некорректное значение в рублях: %w", err)
	}

	return price, dollar, ruble, nil
}

func (handler *CoffeeHandler) saveFile(r *http.Request, fieldName string) (string, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("ошибка получения файла %s: %w", fieldName, err)
	}
	defer file.Close()

	// Создаем директорию, если не существует
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории: %w", err)
	}

	// Получаем расширение исходного файла
	ext := filepath.Ext(fileHeader.Filename)

	// Генерируем UUID для имени файла
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Формируем полный путь к файлу
	fullPath := filepath.Join(uploadDir, filename)

	dest, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return "", fmt.Errorf("ошибка копирования файла: %w", err)
	}

	return filename, nil
}
