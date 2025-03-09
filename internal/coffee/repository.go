package coffee

import (
	"coffee/pkg/db"
)

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

func (repo *CoffeeRepository) GetAllCoffee(limit, offset int) []Coffee {
	var coffees []Coffee
	repo.Database.Table("coffees").Limit(limit).Offset(offset).Scan(&coffees)
	return coffees
}

func (repo *CoffeeRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("coffees").
		Count(&count)
	return count
}

func (repo *CoffeeRepository) Delete(slug string) error {
	result := repo.Database.DB.Unscoped().Where("slug = ?", slug).Delete(&Coffee{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CoffeeRepository) GetBySlug(slug string) (*Coffee, error) {
	var coffee Coffee
	result := repo.Database.DB.Where("slug = ?", slug).First(&coffee)
	if result.Error != nil {
		return nil, result.Error
	}
	return &coffee, nil
}

func (repo *CoffeeRepository) Update(coffee *Coffee) (*Coffee, error) {
	result := repo.Database.DB.Model(&Coffee{}).
		Where("slug = ?", coffee.Slug).
		Updates(coffee)
	if result.Error != nil {
		return nil, result.Error
	}
	return coffee, nil
}
