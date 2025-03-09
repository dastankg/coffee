package app

import (
	"coffee/configs"
	_ "coffee/docs"
	"coffee/internal/auth"
	"coffee/internal/coffee"
	"coffee/internal/notification"
	"coffee/internal/user"
	"coffee/pkg/db"
	httpSwagger "github.com/swaggo/http-swagger" // Add this import

	"coffee/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	userRepository := user.NewUserRepository(db)

	coffeeRepository := coffee.NewCoffeeRepository(db)

	authService := auth.NewAuthService(userRepository)

	coffee.NewCoffeeHandler(router, coffee.CoffeeHandlerDeps{
		CoffeeRepository: coffeeRepository,
		Config:           conf,
	})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	notification.NewNotificationHandler(router, notification.NotificationHandlerDeps{
		Config: conf,
	})
	router.Handle("/docs/", httpSwagger.WrapHandler)

	stack := middleware.Chain(middleware.CORS)
	return stack(router)
}
