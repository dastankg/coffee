package main

import (
	"coffee/internal/app"
	"net/http"
)

// @title People-Base API documentation
// @version 1.0.1
// @host 139.59.2.151:8081
// @BasePath
func main() {
	app := app.App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
