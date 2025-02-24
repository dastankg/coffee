package auth

import (
	"coffee/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(name, email, password string) (string, error) {
	existedUser, err := service.UserRepository.GetByEmail(email)
	if existedUser != nil {
		return "", err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	_, err = service.UserRepository.CreateUser(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
