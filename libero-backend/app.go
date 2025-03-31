package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"libero-backend/config"
	// "libero-backend/internal/api/controllers" // Ensure this is removed
	"libero-backend/internal/api/routes"
	"libero-backend/internal/repository"
	"libero-backend/internal/service"
)

// App encapsulates the application components
type App struct {
	Router     *mux.Router
	Config     *config.Config
	Repository *repository.Repository
	Service    *service.Service
	// Controllers are now instantiated within SetupRoutes
}

// Initialize the application with all dependencies
func initApp() *App {
	app := &App{}

	// Initialize configuration
	app.Config = config.New()

	// Initialize router with CORS
	app.Router = mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	app.Router.Use(c.Handler)

	// Initialize database and repository
	db := config.InitDB(app.Config)
	app.Repository = repository.New(db)

	// Initialize services (which includes AuthService with JWT config)
	app.Service = service.New(app.Repository)

	// Setup routes
	// Pass the main service struct to SetupRoutes
	routes.SetupRoutes(app.Router, app.Service)

	return app
}