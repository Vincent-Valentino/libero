package controllers

import (
	"libero-backend/config"
	"libero-backend/internal/repository"
	serv "libero-backend/internal/service"
)

// Service provides access to all service operations
type Controller struct {
	User              *UserController
	Oauth             *OAuthController
	SportsData        *SportsDataController
	Prediction        *PredictionController
	PredictionHistory *PredictionHistoryController
}

// New creates a new service instance with all services
func New(service *serv.Service, cfg *config.Config, repo *repository.Repository) *Controller {
	return &Controller{
		User:              NewUserController(service.User, service.Auth),
		Oauth:             NewOAuthController(service.OAuth, cfg),
		SportsData:        NewSportsDataController(service.ML, service.Fixtures, service.Football, repo.Cache),
		Prediction:        NewPredictionController(cfg),
		PredictionHistory: NewPredictionHistoryController(service.PredictionHistory),
	}
}
