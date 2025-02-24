package app

import (
	"coffee/configs"
	_ "coffee/docs"
	"coffee/internal/coffee"
	"coffee/pkg/db"
	httpSwagger "github.com/swaggo/http-swagger" // Add this import

	"coffee/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	coffeeRepository := coffee.NewCoffeeRepository(db)

	coffee.NewCoffeeHandler(router, coffee.CoffeeHandlerDeps{
		CoffeeRepository: coffeeRepository,
		Config:           conf,
	})

	router.Handle("/docs/", httpSwagger.WrapHandler)

	stack := middleware.Chain(middleware.CORS)
	return stack(router)
}
