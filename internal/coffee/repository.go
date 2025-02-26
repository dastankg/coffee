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

func (repo *CoffeeRepository) GetAllCoffee(limit, offset int) []Coffee {
	var coffees []Coffee
	repo.Database.Table("coffees").Where("deleted_at is null").Limit(limit).Offset(offset).Scan(&coffees)
	return coffees
}

func (repo *CoffeeRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("coffees").
		Where("deleted_at is null").
		Count(&count)
	return count
}

func (repo *CoffeeRepository) Delete(id uint) error {
	result := repo.Database.DB.Unscoped().Delete(&Coffee{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CoffeeRepository) GetById(id uint) (*Coffee, error) {
	var coffee Coffee
	result := repo.Database.DB.First(&coffee, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &coffee, nil
}
