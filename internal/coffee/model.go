package coffee

import "time"

type Coffee struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string  `json:"name" example:"Espresso" gorm:"size:50;not null"`
	Slug        string  `json:"slug" example:"espresso" gorm:"size:50;unique;index;not null"`
	Price       float64 `json:"price" example:"4.99" gorm:"type:decimal(20,2);not null"`
	Description string  `json:"description" example:"Strong Italian coffee" gorm:"type:text;not null"`
	Dollar      float64 `json:"dollar" example:"1.0" gorm:"type:decimal(20,2);not null"`
	Ruble       float64 `json:"ruble" example:"89.5" gorm:"type:decimal(20,2);not null"`
	Image       string  `json:"image" example:"espresso.jpg" gorm:"type:varchar(500);not null"`
	FlagIcon    string  `json:"flag_icon" example:"italy.png" gorm:"type:varchar(500);not null"`
	QrImage     string  `json:"qrImage" example:"espresso.png" gorm:"type:varchar(500);not null"`
}

func NewCoffee(name string, coffeeSlug string, price float64, Description string, dollar, ruble float64, image, flagIcon, qrImage string) *Coffee {
	return &Coffee{
		Name:        name,
		Slug:        coffeeSlug,
		Price:       price,
		Description: Description,
		Dollar:      dollar,
		Ruble:       ruble,
		Image:       image,
		FlagIcon:    flagIcon,
		QrImage:     qrImage,
	}

}
