package db

import (
	"coffee/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.DATABASE_URL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &Db{db}
}
