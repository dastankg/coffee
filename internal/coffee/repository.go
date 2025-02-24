package coffee

import "coffee/pkg/db"

type CoffeeRepository struct {
	Database *db.Db
}

func NewCoffeeRepository(db *db.Db) *CoffeeRepository {
	return &CoffeeRepository{
		Database: db,
	}
}

func (repo *CoffeeRepository) CreateCoffee(coffee *Coffee) (*Coffee, error) {
	result := repo.Database.DB.Create(coffee)
	if result.Error != nil {
		panic(result.Error)
	}
	return coffee, nil
}
