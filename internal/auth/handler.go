package auth

import (
	"coffee/configs"
	"coffee/pkg/jwt"
	"coffee/pkg/req"
	"coffee/pkg/res"
	"net/http"
	"time"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/refresh", handler.Refresh())
}

// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse "Успешная регистрация"
// @Router /auth/register [post]
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Name, body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)
		tokens, err := jwtService.CreateTokenPair(
			email,
			15*time.Minute, // access token на 15 минут
			24*7*time.Hour, // refresh token на 7 дней
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}

		res.Json(w, data, 201)
	}
}

// @Summary Login нового пользователя
// @Description Login пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для регистрации"
// @Success 201 {object} LoginResponse "Успешная регистрация"
// @Router /auth/login [post]
func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)
		tokens, err := jwtService.CreateTokenPair(
			email,
			15*time.Minute, // access token на 15 минут
			24*7*time.Hour, // refresh token на 7 дней
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		res.Json(w, data, 201)
	}
}

// @Summary Обновление токена доступа
// @Description Обновляет access token используя refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh токен"
// @Success 200 {object} RefreshResponse "Новая пара токенов"
// @Router /auth/refresh [post]
func (handler *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RefreshRequest](&w, r)
		if err != nil {
			return
		}

		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)

		// Проверяем и парсим refresh токен
		flag, claims := jwtService.ParseRefreshToken(body.RefreshToken)
		if !flag {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Вычисляем оставшееся время жизни refresh токена
		expirationTime := time.Unix(claims.ExpiresAt.Unix(), 0)
		remainingTime := time.Until(expirationTime)

		if remainingTime <= 0 {
			http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
			return
		}

		// Генерируем только новый access токен
		accessToken, err := jwtService.Create(jwt.JWTData{
			Email:     claims.Email,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			TokenType: "access",
		}, jwtService.AccessSecret)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RefreshResponse{
			AccessToken: accessToken,
		}

		res.Json(w, data, http.StatusOK)
	}
}
