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
	//router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
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
