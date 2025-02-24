package coffee

import (
	"coffee/configs"
	"coffee/pkg/res"
	"net/http"
)

type CoffeeHandler struct {
	CoffeeRepository *CoffeeRepository
}

type CoffeeHandlerDeps struct {
	CoffeeRepository *CoffeeRepository
	Config           *configs.Config
}

func NewCoffeeHandler(router *http.ServeMux, deps CoffeeHandlerDeps) {
	handler := &CoffeeHandler{
		CoffeeRepository: deps.CoffeeRepository,
	}
	router.HandleFunc("POST /coffee/create", handler.CreateCoffee())
}

const (
	maxFileSize = 10 << 20  // 10 MB
	uploadDir   = "uploads" // директория для загрузки файлов
)

// CreateCoffee ... Create Coffee
// @Summary Create Coffee
// @Description Create coffee
// @Tags Coffee CRUD
// @Accept json
// @Produce json
// @Param name formData string true "Название кофе"
// @Param slug formData string true "URL-friendly идентификатор"
// @Param price formData number true "Цена кофе"
// @Param description formData string true "Описание кофе"
// @Param dollar formData number true "Цена в долларах"
// @Param ruble formData number true "Цена в рублях"
// @Param image formData file true "Изображение кофе"
// @Param flagIcon formData file true "Иконка флага страны происхождения"
// @Success 201 {object} Coffee
// @Router /coffee/create [post]  // <-- Изменили путь здесь
func (handler *CoffeeHandler) CreateCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Загрузка и сохранение изображений
		imagePath, err := handler.saveFile(r, "image")
		if err != nil {
			http.Error(w, "Ошибка при сохранении изображения: "+err.Error(), http.StatusBadRequest)
			return
		}

		flagIconPath, err := handler.saveFile(r, "flagIcon")
		if err != nil {
			http.Error(w, "Ошибка при сохранении иконки флага: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Парсинг числовых значений
		price, dollar, ruble, err := handler.parseNumericValues(r)
		if err != nil {
			http.Error(w, "Ошибка в числовых значениях: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Создание объекта кофе
		coffee := NewCoffee(
			r.FormValue("name"),
			r.FormValue("slug"),
			price,
			r.FormValue("description"),
			dollar,
			ruble,
			imagePath,
			flagIconPath,
		)

		// Сохранение в базу данных
		createdCoffee, err := handler.CoffeeRepository.CreateCoffee(coffee)
		if err != nil {
			http.Error(w, "Ошибка при создании записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, createdCoffee, http.StatusCreated)
	}
}
