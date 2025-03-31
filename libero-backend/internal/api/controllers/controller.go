package controllers

import (
	"libero-backend/internal/service"
)

// Service provides access to all service operations
type Controller struct {
	User *UserController
	Oauth *OAuthController
}

// New creates a new service instance with all services
func New(service *service.Service) *Controller {

	return &Controller{
		User: NewUserController(service.User),
		Oauth: NewOAuthController(service.OAuth),
	}
}