package user

import "coffee/pkg/db"

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (repo *UserRepository) CreateUser(user *User) (*User, error) {
	result := repo.database.DB.Create(user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := repo.database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
