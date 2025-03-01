package coffee

import "mime/multipart"

type CoffeeCreateRequest struct {
	Name        string          `json:"name" validate:"required,max=50"`
	Slug        string          `json:"slug" validate:"required,max=50"`
	Price       float64         `json:"price" validate:"required,gt=0"`
	Description string          `json:"description" validate:"required"`
	Dollar      float64         `json:"dollar" validate:"required,gt=0"`
	Ruble       float64         `json:"ruble" validate:"required,gt=0"`
	Image       *multipart.Form `json:"image" validate:"required"`
	FlagIcon    *multipart.Form `json:"flag_icon" validate:"required"`
}

type CoffeeGetAllResponse struct {
	Coffee []Coffee `json:"coffee"`
	Count  int64    `json:"count"`
}
type CoffeeGetResponse struct {
	Coffee Coffee `json:"coffee"`
}
type CoffeeDeleteResponse struct {
	Message string `json:"message"`
}

type CoffeeUpdateRequest struct {
	Name        string          `json:"name" validate:"required,max=50"`
	Slug        string          `json:"slug" validate:"required,max=50"`
	Price       float64         `json:"price" validate:"required,gt=0"`
	Description string          `json:"description" validate:"required"`
	Dollar      float64         `json:"dollar" validate:"required,gt=0"`
	Ruble       float64         `json:"ruble" validate:"required,gt=0"`
	Image       *multipart.Form `json:"image" validate:"required"`
	FlagIcon    *multipart.Form `json:"flag_icon" validate:"required"`
}
