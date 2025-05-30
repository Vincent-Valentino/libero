package repository

import (
	"libero-backend/internal/models"

	"gorm.io/gorm"
)

// PredictionHistoryRepository defines the interface for prediction history data operations
type PredictionHistoryRepository interface {
	Create(prediction *models.PredictionHistory) error
	FindByUserID(userID uint, page, limit int) ([]models.PredictionHistory, int64, error)
	FindByID(id uint) (*models.PredictionHistory, error)
	Delete(id uint, userID uint) error
	DeleteAllByUserID(userID uint) error
	GetStatistics(userID uint) (*models.PredictionStatistics, error)
}

// predictionHistoryRepository implements the PredictionHistoryRepository interface
type predictionHistoryRepository struct {
	db *gorm.DB
}

// NewPredictionHistoryRepository creates a new prediction history repository instance
func NewPredictionHistoryRepository(db *gorm.DB) PredictionHistoryRepository {
	return &predictionHistoryRepository{db: db}
}

// Create adds a new prediction to the database
func (r *predictionHistoryRepository) Create(prediction *models.PredictionHistory) error {
	return r.db.Create(prediction).Error
}

// FindByUserID retrieves predictions for a specific user with pagination
func (r *predictionHistoryRepository) FindByUserID(userID uint, page, limit int) ([]models.PredictionHistory, int64, error) {
	var predictions []models.PredictionHistory
	var count int64

	offset := (page - 1) * limit

	// Get total count for the user
	if err := r.db.Model(&models.PredictionHistory{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get predictions for the current page, ordered by creation date (newest first)
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&predictions).Error; err != nil {
		return nil, 0, err
	}

	return predictions, count, nil
}

// FindByID retrieves a prediction by ID
func (r *predictionHistoryRepository) FindByID(id uint) (*models.PredictionHistory, error) {
	var prediction models.PredictionHistory
	err := r.db.First(&prediction, id).Error
	if err != nil {
		return nil, err
	}
	return &prediction, nil
}

// Delete removes a prediction from the database (only if it belongs to the user)
func (r *predictionHistoryRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PredictionHistory{}).Error
}

// DeleteAllByUserID removes all predictions for a specific user
func (r *predictionHistoryRepository) DeleteAllByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PredictionHistory{}).Error
}

// GetStatistics calculates and returns prediction statistics for a user
func (r *predictionHistoryRepository) GetStatistics(userID uint) (*models.PredictionStatistics, error) {
	var stats models.PredictionStatistics

	// Count total predictions
	var total int64
	if err := r.db.Model(&models.PredictionHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}

	stats.Total = int(total)
	if stats.Total == 0 {
		return &stats, nil
	}

	// Count home wins (predictions where home team wins)
	var homeWins int64
	if err := r.db.Model(&models.PredictionHistory{}).
		Where("user_id = ? AND predicted_home_score > predicted_away_score", userID).
		Count(&homeWins).Error; err != nil {
		return nil, err
	}
	stats.HomeWins = int(homeWins)

	// Count draws
	var draws int64
	if err := r.db.Model(&models.PredictionHistory{}).
		Where("user_id = ? AND predicted_home_score = predicted_away_score", userID).
		Count(&draws).Error; err != nil {
		return nil, err
	}
	stats.Draws = int(draws)

	// Count away wins (predictions where away team wins)
	var awayWins int64
	if err := r.db.Model(&models.PredictionHistory{}).
		Where("user_id = ? AND predicted_home_score < predicted_away_score", userID).
		Count(&awayWins).Error; err != nil {
		return nil, err
	}
	stats.AwayWins = int(awayWins)

	// Calculate percentages
	totalFloat := float64(stats.Total)
	stats.HomeWinPercentage = (float64(stats.HomeWins) / totalFloat) * 100
	stats.DrawPercentage = (float64(stats.Draws) / totalFloat) * 100
	stats.AwayWinPercentage = (float64(stats.AwayWins) / totalFloat) * 100

	return &stats, nil
}
