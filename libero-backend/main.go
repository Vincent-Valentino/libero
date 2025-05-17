package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize the application
	app := initApp()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start the server in a goroutine
	go app.Run()
	
	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down server...")
	
	// Clean up resources
	app.Shutdown()
	log.Println("Server shut down complete")
}