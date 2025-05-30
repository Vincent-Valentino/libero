package service

import (
	"errors"
	"libero-backend/internal/models"
	"libero-backend/internal/repository"
)

// Error definitions for prediction history service
var (
	ErrInvalidInput = errors.New("invalid input data")
)

// PredictionHistoryService defines the interface for prediction history business logic
type PredictionHistoryService interface {
	CreatePrediction(userID uint, request *models.CreatePredictionRequest) (*models.PredictionHistory, error)
	GetUserPredictions(userID uint, page, limit int) ([]models.PredictionHistory, int64, error)
	DeletePrediction(predictionID, userID uint) error
	DeleteAllUserPredictions(userID uint) error
	GetUserStatistics(userID uint) (*models.PredictionStatistics, error)
}

// predictionHistoryService implements the PredictionHistoryService interface
type predictionHistoryService struct {
	predictionRepo repository.PredictionHistoryRepository
}

// NewPredictionHistoryService creates a new prediction history service instance
func NewPredictionHistoryService(predictionRepo repository.PredictionHistoryRepository) PredictionHistoryService {
	return &predictionHistoryService{
		predictionRepo: predictionRepo,
	}
}

// CreatePrediction creates a new prediction record
func (s *predictionHistoryService) CreatePrediction(userID uint, request *models.CreatePredictionRequest) (*models.PredictionHistory, error) {
	// Normalize request to handle both camelCase and snake_case
	request.Normalize()

	// Validate request
	if request.HomeTeam == "" || request.AwayTeam == "" ||
		request.HomeLeague == "" || request.AwayLeague == "" ||
		request.PredictedResult == "" {
		return nil, ErrInvalidInput
	}

	// Create prediction model
	prediction := &models.PredictionHistory{
		UserID:             userID,
		HomeTeam:           request.HomeTeam,
		AwayTeam:           request.AwayTeam,
		HomeLeague:         request.HomeLeague,
		AwayLeague:         request.AwayLeague,
		PredictedHomeScore: request.PredictedHomeScore,
		PredictedAwayScore: request.PredictedAwayScore,
		ExpectedHomeGoals:  request.ExpectedHomeGoals,
		ExpectedAwayGoals:  request.ExpectedAwayGoals,
		HomeWinProbability: request.HomeWinProbability,
		DrawProbability:    request.DrawProbability,
		AwayWinProbability: request.AwayWinProbability,
		PredictedResult:    request.PredictedResult,
	}

	// Save to database
	if err := s.predictionRepo.Create(prediction); err != nil {
		return nil, err
	}

	return prediction, nil
}

// GetUserPredictions retrieves predictions for a specific user with pagination
func (s *predictionHistoryService) GetUserPredictions(userID uint, page, limit int) ([]models.PredictionHistory, int64, error) {
	return s.predictionRepo.FindByUserID(userID, page, limit)
}

// DeletePrediction removes a specific prediction (only if it belongs to the user)
func (s *predictionHistoryService) DeletePrediction(predictionID, userID uint) error {
	return s.predictionRepo.Delete(predictionID, userID)
}

// DeleteAllUserPredictions removes all predictions for a specific user
func (s *predictionHistoryService) DeleteAllUserPredictions(userID uint) error {
	return s.predictionRepo.DeleteAllByUserID(userID)
}

// GetUserStatistics calculates and returns prediction statistics for a user
func (s *predictionHistoryService) GetUserStatistics(userID uint) (*models.PredictionStatistics, error) {
	return s.predictionRepo.GetStatistics(userID)
}
