package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" example:"John" gorm:"size:50;not null"`
	Email    string `json:"email" example:"<EMAIL>" gorm:"size:50;unique;index;not null"`
	Password string `json:"password" example:"<PASSWORD>" gorm:"size:50;not null"`
}
