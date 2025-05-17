package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"libero-backend/config"

	"libero-backend/internal/api/routes"
	"libero-backend/internal/repository"
	"libero-backend/internal/scheduler"
	"libero-backend/internal/service"
)

// App encapsulates the application components.
type App struct {
	Router        *mux.Router
	Config        *config.Config
	Repository    *repository.Repository
	Service       *service.Service
	Scheduler     *scheduler.Scheduler
	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
	// Controllers are now instantiated within SetupRoutes
}

// Initialize the application with all dependencies.
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
	routes.SetupRoutes(app.Router, app.Service, app.Config) // Pass config as well
	
	// Setup cache cleanup
	app.cleanupCtx, app.cleanupCancel = context.WithCancel(context.Background())
	go app.startCacheCleanup()
	
	// Initialize and start scheduler
	app.Scheduler = scheduler.New(app.Service.Fixtures)
	app.Scheduler.Start()

	return app
}

// startCacheCleanup runs a background goroutine to clean expired cache entries.
func (a *App) startCacheCleanup() {
	ticker := time.NewTicker(15 * time.Minute) // Clean every 15 minutes
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if err := a.Repository.Cache.CleanExpiredCache(); err != nil {
				log.Printf("Error cleaning cache: %v", err)
			} else {
				log.Println("Cache cleanup completed")
			}
		case <-a.cleanupCtx.Done():
			log.Println("Cache cleanup terminated")
			return
		}
	}
}

// Run starts the HTTP server.
func (a *App) Run() {
	addr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)
	log.Printf("Server starting on %s", addr)
	
	// Create a server with timeout settings
	server := &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Shutdown gracefully terminates the app.
func (a *App) Shutdown() {
	// Stop the scheduler
	if a.Scheduler != nil {
		a.Scheduler.Stop()
	}
	
	// Cancel the cleanup goroutine
	if a.cleanupCancel != nil {
		a.cleanupCancel()
	}
}