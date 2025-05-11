package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessToken  string = "access"
	RefreshToken string = "refresh"
)

type JWTData struct {
	Email     string
	ExpiresAt time.Time
	TokenType string
}

type JWT struct {
	AccessSecret  string
	RefreshSecret string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewJWT(accessSecret, refreshSecret string) *JWT {
	return &JWT{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
	}
}

func (j *JWT) CreateTokenPair(email string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	// Создаем access token
	accessToken, err := j.Create(JWTData{
		Email:     email,
		ExpiresAt: time.Now().Add(accessTTL),
		TokenType: AccessToken,
	}, j.AccessSecret)
	if err != nil {
		return nil, err
	}

	// Создаем refresh token
	refreshToken, err := j.Create(JWTData{
		Email:     email,
		ExpiresAt: time.Now().Add(refreshTTL),
		TokenType: RefreshToken,
	}, j.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (j *JWT) Create(data JWTData, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   data.ExpiresAt.Unix(),
		"type":  data.TokenType,
	})
	return t.SignedString([]byte(secret))
}

func (j *JWT) parse(token string, secret string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, nil
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	email, ok := claims["email"].(string)
	if !ok {
		return false, nil
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, nil
	}

	tokenType, ok := claims["type"].(string)
	if !ok {
		return false, nil
	}

	return t.Valid, &JWTData{
		Email:     email,
		ExpiresAt: time.Unix(int64(exp), 0),
		TokenType: tokenType,
	}
}

func (j *JWT) ParseAccessToken(token string) (bool, *JWTData) {
	return j.parse(token, j.AccessSecret)
}

func (j *JWT) ParseRefreshToken(token string) (bool, *JWTData) {
	return j.parse(token, j.RefreshSecret)
}

func (j *JWT) RefreshTokens(refreshToken string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	valid, data := j.ParseRefreshToken(refreshToken)
	if !valid || data == nil {
		return nil, jwt.ErrSignatureInvalid
	}

	if data.TokenType != RefreshToken {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return j.CreateTokenPair(data.Email, accessTTL, refreshTTL)
}
