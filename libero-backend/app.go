package main

import (
	"context"
	"fmt"
	"libero-backend/config"
	"libero-backend/internal/repository"
	"libero-backend/internal/routes"
	"libero-backend/internal/scheduler"
	"libero-backend/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// App encapsulates app-wide dependencies and configuration
type App struct {
	Router        *mux.Router
	Config        *config.Config
	Repository    *repository.Repository
	Service       *service.Service
	Scheduler     *scheduler.Scheduler
	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
}

// Initialize the application with all dependencies.
func initApp() *App {
	app := &App{}

	// Initialize configuration
	app.Config = config.New()

	// Initialize database and repository
	db := config.InitDB(app.Config)
	app.Repository = repository.New(db)

	// Initialize services
	app.Service = service.New(app.Repository)

	// Initialize router
	app.Router = mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(app.Router, app.Service, app.Config, app.Repository)

	// Setup cache cleanup
	app.cleanupCtx, app.cleanupCancel = context.WithCancel(context.Background())
	go app.startCacheCleanup()

	// Initialize and start scheduler
	app.Scheduler = scheduler.New(app.Service.Fixtures)
	app.Scheduler.Start()

	return app
}

// Run starts the HTTP server
func (a *App) Run() error {
	port := a.Config.Server.Port
	host := a.Config.Server.Host

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: a.Router,
	}

	// Channel to listen for errors coming from the server.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		fmt.Printf("Server is running on %s:%d\n", host, port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %v", err)

	case sig := <-shutdown:
		fmt.Printf("Starting shutdown... Signal: %v\n", sig)

		// Stop the cache cleanup
		if a.cleanupCancel != nil {
			a.cleanupCancel()
		}

		// Stop the scheduler
		if a.Scheduler != nil {
			a.Scheduler.Stop()
		}

		// Create context for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			// If graceful shutdown fails, force close
			if err := srv.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %v", err)
			}
		}
	}

	return nil
}

// Shutdown performs cleanup of application resources
func (a *App) Shutdown() {
	// Stop the cache cleanup
	if a.cleanupCancel != nil {
		a.cleanupCancel()
	}

	// Stop the scheduler
	if a.Scheduler != nil {
		a.Scheduler.Stop()
	}

	fmt.Println("Application cleanup completed")
}

// startCacheCleanup runs periodic cache cleanup
func (a *App) startCacheCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.Repository.Cache.CleanExpiredCache(); err != nil {
				fmt.Printf("Error cleaning cache: %v\n", err)
			}
		case <-a.cleanupCtx.Done():
			return
		}
	}
}
