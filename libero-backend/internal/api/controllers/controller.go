package controllers

import (
	"libero-backend/internal/service"
	"libero-backend/config" // Added config import
)

// Service provides access to all service operations
type Controller struct {
	User *UserController
	Oauth *OAuthController
}

// New creates a new service instance with all services
func New(service *service.Service, cfg *config.Config) *Controller { // Added cfg argument

	return &Controller{
		User: NewUserController(service.User, service.Auth),
		Oauth: NewOAuthController(service.OAuth, cfg), // Pass cfg here
	}
}