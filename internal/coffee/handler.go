package coffee

import (
	"coffee/configs"
	"coffee/pkg/middleware"
	"coffee/pkg/res"
	"net/http"
	"strconv"
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
	router.Handle("POST /coffee/create", middleware.IsAuthed(handler.CreateCoffee(), deps.Config))
	router.HandleFunc("POST /coffee/coffees", handler.GetAllCoffee())
	router.Handle("POST /coffee/delete/{id}", middleware.IsAuthed(handler.DeleteCoffee(), deps.Config))
}

const (
	maxFileSize = 10 << 20  // 10 MB
	uploadDir   = "uploads" // директория для загрузки файлов
)

// CreateCoffee ... Create Coffee
// @Summary Create Coffee
// @Description Create coffee
// @Tags Coffee
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer токен авторизации" default(Bearer <token>)
// @Param name formData string true "Название кофе"
// @Param slug formData string true "URL-friendly идентификатор"
// @Param price formData number true "Цена кофе"
// @Param description formData string true "Описание кофе"
// @Param dollar formData number true "Цена в долларах"
// @Param ruble formData number true "Цена в рублях"
// @Param image formData file true "Изображение кофе"
// @Param flagIcon formData file true "Иконка флага страны происхождения"
// @Success 201 {object} Coffee
// @Failure 401 {string} string "Unauthorized"
// @Router /coffee/create [post]
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

// @Summary Получение списка кофе
// @Description Возвращает список кофе с пагинацией
// @Tags Coffee
// @Accept json
// @Produce json
// @Param limit query int true "Количество записей на странице"
// @Param offset query int true "Смещение от начала списка"
// @Success 200 {object} CoffeeGetAllResponse "Список кофе и общее количество"
// @Failure 400 {string} string "Неверные параметры пагинации"
// @Router /coffee/coffees [get]
func (handler *CoffeeHandler) GetAllCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		coffees := handler.CoffeeRepository.GetAllCoffee(limit, offset)
		count := handler.CoffeeRepository.Count()
		resultat := CoffeeGetAllResponse{
			Coffee: coffees,
			Count:  count,
		}
		res.Json(w, resultat, http.StatusOK)

	}
}

// DeleteCoffee ... Удаление кофе
// @Summary Удаление кофе
// @Description Удаляет кофе по указанному ID
// @Tags Coffee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer токен авторизации" default(Bearer <token>)
// @Param id path int true "ID кофе"
// @Success 200 {object} nil "Успешное удаление"
// @Failure 400 {string} string "Неверный ID"
// @Failure 401 {string} string "Unauthorized"
// @Router /coffee/delete/{id} [post]
func (handler *CoffeeHandler) DeleteCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		_, err = handler.CoffeeRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}
		err = handler.CoffeeRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		res.Json(w, CoffeeDeleteResponse{
			Message: "Товар удален",
		}, http.StatusOK)
	}
}
