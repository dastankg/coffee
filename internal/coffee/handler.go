package coffee

import (
	"coffee/configs"
	"coffee/pkg/middleware"
	"coffee/pkg/qr"
	"coffee/pkg/res"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	router.Handle("POST /coffees", middleware.IsAuthed(handler.CreateCoffee(), deps.Config))
	router.HandleFunc("GET /coffees", handler.GetAllCoffee())
	router.HandleFunc("GET /coffees/{slug}", handler.GetCoffee())
	router.HandleFunc("GET /coffees/static/images/{dir}/{filename}", handler.GetCoffeeImage())
	router.Handle("DELETE /coffees/{slug}", middleware.IsAuthed(handler.DeleteCoffee(), deps.Config))
	router.Handle("PUT /coffees/{slug}", middleware.IsAuthed(handler.UpdateCoffee(), deps.Config))
}

const (
	maxFileSize = 10 << 20 // 10 MB
	uploadDir   = "static/images"
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
// @Router /coffees [post]
func (handler *CoffeeHandler) CreateCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
			return
		}

		imagePath, err := handler.saveFile(r, "image", uploadDir+"/products")
		if err != nil {
			http.Error(w, "Ошибка при сохранении изображения: "+err.Error(), http.StatusBadRequest)
			return
		}

		flagIconPath, err := handler.saveFile(r, "flagIcon", uploadDir+"/flagsIcon")
		if err != nil {
			http.Error(w, "Ошибка при сохранении иконки флага: "+err.Error(), http.StatusBadRequest)
			return
		}

		price, dollar, ruble, err := handler.parseNumericValues(r)
		if err != nil {
			http.Error(w, "Ошибка в числовых значениях: "+err.Error(), http.StatusBadRequest)
			return
		}
		qrCode := qr.SimpleQRCode{Content: "http://139.59.2.151:8081/coffee/coffee/" + r.FormValue("slug"), Size: 256}

		qrImage, err := qrCode.SaveToFile(uploadDir + "/qr")
		if err != nil {
			http.Error(w, "Qr code create error", http.StatusBadRequest)
		}
		coffee := NewCoffee(
			r.FormValue("name"),
			r.FormValue("slug"),
			price,
			r.FormValue("description"),
			dollar,
			ruble,
			imagePath,
			flagIconPath,
			qrImage,
		)

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
// @Router /coffees [get]
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
// @Param slug path string true "slug кофе"
// @Success 200 {object} nil "Успешное удаление"
// @Failure 400 {string} string "Неверный ID"
// @Failure 401 {string} string "Unauthorized"
// @Router /coffees/{slug} [delete]
// ]
func (handler *CoffeeHandler) DeleteCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		coffee, err := handler.CoffeeRepository.GetBySlug(slug)
		_ = handler.deleteFile(coffee.Image)
		_ = handler.deleteFile(coffee.FlagIcon)
		_ = handler.deleteFile(coffee.QrImage)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}
		err = handler.CoffeeRepository.Delete(slug)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		res.Json(w, CoffeeDeleteResponse{
			Message: "Товар удален",
		}, http.StatusOK)
	}
}

// Update ... Обновление информации о кофе
// @Summary Обновление кофе
// @Description Обновляет информацию о кофе по указанному ID-
// @Tags Coffee
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer токен авторизации" default(Bearer <token>)
// @Param slug path string true "slug кофе"
// @Param name formData string false "Название кофе"
// @Param slug formData string false "URL-friendly идентификатор"
// @Param price formData number false "Цена кофе"
// @Param description formData string false "Описание кофе"
// @Param dollar formData number false "Цена в долларах"
// @Param ruble formData number false "Цена в рублях"
// @Param image formData file false "Изображение кофе"
// @Param flagIcon formData file false "Иконка флага страны происхождения"
// @Success 200 {object} Coffee "Обновленная информация о кофе"
// @Failure 400 {string} string "Ошибка в запросе или неверный ID"
// @Failure 401 {string} string "Unauthorized"
// @Router /coffees/{slug} [put]
func (handler *CoffeeHandler) UpdateCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		existingCoffee, err := handler.CoffeeRepository.GetBySlug(slug)
		if err != nil {
			http.Error(w, "coffee not found", http.StatusBadRequest)
			return
		}

		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			name = existingCoffee.Name
		}

		slug = r.FormValue("slug")
		if slug == "" {
			slug = existingCoffee.Slug
		}

		description := r.FormValue("description")
		if description == "" {
			description = existingCoffee.Description
		}

		price, dollar, ruble, err := handler.parseNumericValues(r)
		if err != nil {
			price = existingCoffee.Price
			dollar = existingCoffee.Dollar
			ruble = existingCoffee.Ruble
		}

		imagePath := existingCoffee.Image
		if _, fileHeader, _ := r.FormFile("image"); fileHeader != nil {
			newImagePath, err := handler.saveFile(r, "image", uploadDir+"/products")
			if err == nil {
				_ = handler.deleteFile(filepath.Join(uploadDir+"/products", filepath.Base(imagePath)))
				imagePath = newImagePath
			}
		}

		flagIconPath := existingCoffee.FlagIcon
		if _, fileHeader, _ := r.FormFile("flagIcon"); fileHeader != nil {
			newFlagIconPath, err := handler.saveFile(r, "flagIcon", uploadDir+"/flagsIcon")
			if err == nil {
				_ = handler.deleteFile(filepath.Join(uploadDir+"/flagsIcon", filepath.Base(flagIconPath)))

				flagIconPath = newFlagIconPath
			}
		}

		updatedCoffee, err := handler.CoffeeRepository.Update(&Coffee{
			Name:        name,
			Slug:        slug,
			Price:       price,
			Description: description,
			Dollar:      dollar,
			Ruble:       ruble,
			Image:       imagePath,
			FlagIcon:    flagIconPath,
		})

		if err != nil {
			http.Error(w, "failed to update coffee: "+err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, updatedCoffee, http.StatusOK)
	}
}

// @Summary Получение кофе
// @Description Возвращает кофе
// @Tags Coffee
// @Accept json
// @Produce json
// @Param slug path string true "slug кофе"
// @Success 200 {object} CoffeeGetResponse "кофе"
// @Failure 400 {string} string "Неверные параметры"
// @Router /coffees/{slug} [get]
func (handler *CoffeeHandler) GetCoffee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		coffee, err := handler.CoffeeRepository.GetBySlug(slug)
		if err != nil {
			http.Error(w, "coffee not found", http.StatusNotFound)
			return
		}
		result := CoffeeGetResponse{
			Coffee: *coffee,
		}
		res.Json(w, result, http.StatusOK)
	}
}

func (handler *CoffeeHandler) GetCoffeeImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")
		dir := r.PathValue("dir")
		imagePath := path.Join(uploadDir+"/"+dir, filename)
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		contentType := "image/jpeg"
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", "inline")
		http.ServeFile(w, r, imagePath)
	}
}
