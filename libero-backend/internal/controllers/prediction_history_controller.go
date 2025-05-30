package controllers

import (
	"encoding/json"
	"fmt"
	"libero-backend/internal/middleware"
	"libero-backend/internal/models"
	"libero-backend/internal/service"
	"libero-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// PredictionHistoryController handles HTTP requests for prediction history
type PredictionHistoryController struct {
	predictionService service.PredictionHistoryService
}

// NewPredictionHistoryController creates a new prediction history controller instance
func NewPredictionHistoryController(predictionService service.PredictionHistoryService) *PredictionHistoryController {
	return &PredictionHistoryController{
		predictionService: predictionService,
	}
}

// CreatePrediction handles POST /api/predictions
func (c *PredictionHistoryController) CreatePrediction(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context (set by auth middleware)
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var request models.CreatePredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use service layer to create prediction
	prediction, err := c.predictionService.CreatePrediction(claims.UserID, &request)
	if err != nil {
		if err.Error() == "invalid input data" {
			http.Error(w, "Invalid input: home team, away team, leagues, and predicted result are required", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to save prediction", http.StatusInternalServerError)
		}
		return
	}

	// Return the created prediction
	utils.RespondWithJSON(w, http.StatusCreated, prediction.ToResponse())
}

// GetPredictions handles GET /api/predictions
func (c *PredictionHistoryController) GetPredictions(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 50

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get predictions from service
	predictions, total, err := c.predictionService.GetUserPredictions(claims.UserID, page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch predictions", http.StatusInternalServerError)
		return
	}

	// Debug logging
	fmt.Printf("Debug: Retrieved %d predictions from database\n", len(predictions))
	if len(predictions) > 0 {
		fmt.Printf("Debug: First prediction from DB: %+v\n", predictions[0])

		// Test: Return raw prediction to see what GORM loaded
		fmt.Printf("Debug: Raw prediction JSON would be: HomeTeam='%s', AwayTeam='%s'\n",
			predictions[0].HomeTeam, predictions[0].AwayTeam)
	}

	// Convert to response format (fixed to use ToResponse properly)
	var responses []models.PredictionHistoryResponse
	for _, prediction := range predictions {
		response := prediction.ToResponse()
		fmt.Printf("Debug: Converting prediction %d to response: %+v\n", prediction.ID, response)
		responses = append(responses, response)
	}

	// Debug logging for responses
	if len(responses) > 0 {
		fmt.Printf("Debug: First response object: %+v\n", responses[0])
		fmt.Printf("Debug: Response should have camelCase fields\n")
	}

	// Return converted responses, not raw predictions
	finalResponse := map[string]interface{}{
		"predictions": responses, // Use converted responses
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	fmt.Printf("Debug: Returning converted responses with camelCase fields\n")
	utils.RespondWithJSON(w, http.StatusOK, finalResponse)
}

// DeletePrediction handles DELETE /api/predictions/{id}
func (c *PredictionHistoryController) DeletePrediction(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse prediction ID from URL
	vars := mux.Vars(r)
	predictionIDStr := vars["id"]
	predictionID, err := strconv.ParseUint(predictionIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid prediction ID", http.StatusBadRequest)
		return
	}

	// Delete prediction using service
	if err := c.predictionService.DeletePrediction(uint(predictionID), claims.UserID); err != nil {
		http.Error(w, "Failed to delete prediction", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Prediction deleted successfully"})
}

// DeleteAllPredictions handles DELETE /api/predictions
func (c *PredictionHistoryController) DeleteAllPredictions(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Delete all predictions using service
	if err := c.predictionService.DeleteAllUserPredictions(claims.UserID); err != nil {
		http.Error(w, "Failed to delete predictions", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "All predictions deleted successfully"})
}

// GetPredictionStatistics handles GET /api/predictions/statistics
func (c *PredictionHistoryController) GetPredictionStatistics(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get statistics using service
	stats, err := c.predictionService.GetUserStatistics(claims.UserID)
	if err != nil {
		http.Error(w, "Failed to fetch statistics", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, stats)
}
