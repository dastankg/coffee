package auth

import (
	"coffee/configs"
	"net/http"
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
	//handler := &AuthHandler{
	//	Config:      deps.Config,
	//	AuthService: deps.AuthService,
	//}
	//router.HandleFunc("POST /auth/login", handler.Login())
	//router.HandleFunc("POST /auth/register", handler.Register())
}
