import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface PredictionHistory {
  id: number;
  homeTeam: string;
  awayTeam: string;
  homeLeague: string;
  awayLeague: string;
  predictedHomeScore: number;
  predictedAwayScore: number;
  expectedHomeGoals: number;
  expectedAwayGoals: number;
  homeWinProbability: number;
  drawProbability: number;
  awayWinProbability: number;
  predictedResult: string;
  createdAt: string;
  userId?: number;
}

export interface CreatePredictionRequest {
  homeTeam: string;
  awayTeam: string;
  homeLeague: string;
  awayLeague: string;
  predictedHomeScore: number;
  predictedAwayScore: number;
  expectedHomeGoals: number;
  expectedAwayGoals: number;
  homeWinProbability: number;
  drawProbability: number;
  awayWinProbability: number;
  predictedResult: string;
}

export interface PredictionStatistics {
  total: number;
  homeWins: number;
  draws: number;
  awayWins: number;
  homeWinPercentage: number;
  drawPercentage: number;
  awayWinPercentage: number;
}

// API client setup
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

const apiClient = {
  async request(url: string, options: RequestInit = {}) {
    const token = localStorage.getItem('authToken');
    const headers = {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    };

    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  },

  get(url: string) {
    return this.request(url);
  },

  post(url: string, data: any) {
    return this.request(url, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  delete(url: string) {
    return this.request(url, {
      method: 'DELETE',
    });
  },
};

// Helper function to normalize prediction data from API (handles both snake_case and camelCase)
const normalizePrediction = (rawPrediction: any): PredictionHistory => {
  return {
    id: rawPrediction.id,
    homeTeam: rawPrediction.homeTeam || rawPrediction.home_team || '',
    awayTeam: rawPrediction.awayTeam || rawPrediction.away_team || '',
    homeLeague: rawPrediction.homeLeague || rawPrediction.home_league || '',
    awayLeague: rawPrediction.awayLeague || rawPrediction.away_league || '',
    predictedHomeScore: rawPrediction.predictedHomeScore ?? rawPrediction.predicted_home_score ?? 0,
    predictedAwayScore: rawPrediction.predictedAwayScore ?? rawPrediction.predicted_away_score ?? 0,
    expectedHomeGoals: rawPrediction.expectedHomeGoals ?? rawPrediction.expected_home_goals ?? 0,
    expectedAwayGoals: rawPrediction.expectedAwayGoals ?? rawPrediction.expected_away_goals ?? 0,
    homeWinProbability: rawPrediction.homeWinProbability ?? rawPrediction.home_win_probability ?? 0,
    drawProbability: rawPrediction.drawProbability ?? rawPrediction.draw_probability ?? 0,
    awayWinProbability: rawPrediction.awayWinProbability ?? rawPrediction.away_win_probability ?? 0,
    predictedResult: rawPrediction.predictedResult || rawPrediction.predicted_result || 'Unknown',
    createdAt: rawPrediction.createdAt || rawPrediction.created_at || new Date().toISOString(),
    userId: rawPrediction.userId || rawPrediction.user_id,
  };
};

export const usePredictionStore = defineStore('prediction', () => {
  const predictions = ref<PredictionHistory[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Fetch user's prediction history from API
  const fetchPredictions = async () => {
    isLoading.value = true;
    error.value = null;
    
    try {
      const token = localStorage.getItem('authToken');
      console.log('Debug: Token found:', token ? 'Yes (length: ' + token.length + ')' : 'No');
      
      if (!token) {
        // No token available, user is not authenticated
        predictions.value = [];
        error.value = 'Please log in to view your prediction history';
        return;
      }

      console.log('Debug: Making API request to /predictions');
      const response = await apiClient.get('/predictions');
      console.log('Debug: API response received with', response.predictions?.length || 0, 'predictions');
      if (response.predictions && response.predictions.length > 0) {
        console.log('Debug: First prediction fields:', Object.keys(response.predictions[0]));
        console.log('Debug: Sample prediction:', response.predictions[0]);
      }
      predictions.value = response.predictions?.map(normalizePrediction) || [];
    } catch (err: any) {
      console.error('Error fetching predictions:', err);
      
      // Check for 401 Unauthorized or similar auth errors
      if (err.message.includes('401') || err.message.includes('Authorization') || err.message.includes('Unauthorized')) {
        error.value = 'Please log in to view your prediction history';
      } else if (err.message.includes('403')) {
        error.value = 'Access denied. Please check your permissions.';
      } else {
        error.value = 'Failed to load prediction history. Please check your connection.';
      }
      
      // Fallback to localStorage for offline functionality
      const stored = localStorage.getItem('prediction_history');
      if (stored) {
        predictions.value = JSON.parse(stored).map(normalizePrediction);
      } else {
        predictions.value = [];
      }
    } finally {
      isLoading.value = false;
    }
  };

  // Save a new prediction to API
  const savePrediction = async (predictionData: CreatePredictionRequest) => {
    try {
      const response = await apiClient.post('/predictions', predictionData);
      
      // Add the new prediction to the beginning of the list
      predictions.value.unshift(normalizePrediction(response));
      
      // Also save to localStorage as backup
      localStorage.setItem('prediction_history', JSON.stringify(predictions.value));
      
      return response;
    } catch (err: any) {
      console.error('Error saving prediction:', err);
      error.value = 'Failed to save prediction';
      
      // Fallback to localStorage
      const newPrediction: PredictionHistory = {
        id: Date.now(), // Temporary ID generation
        ...predictionData,
        createdAt: new Date().toISOString(),
      };

      predictions.value.unshift(newPrediction);
      localStorage.setItem('prediction_history', JSON.stringify(predictions.value));
      
      throw err;
    }
  };

  // Delete a prediction from API
  const deletePrediction = async (id: number) => {
    try {
      await apiClient.delete(`/predictions/${id}`);
      predictions.value = predictions.value.filter(p => p.id !== id);
      
      // Update localStorage backup
      localStorage.setItem('prediction_history', JSON.stringify(predictions.value));
    } catch (err: any) {
      console.error('Error deleting prediction:', err);
      error.value = 'Failed to delete prediction';
      throw err;
    }
  };

  // Clear all predictions from API
  const clearAllPredictions = async () => {
    try {
      await apiClient.delete('/predictions');
      predictions.value = [];
      
      // Clear localStorage backup
      localStorage.removeItem('prediction_history');
    } catch (err: any) {
      console.error('Error clearing predictions:', err);
      error.value = 'Failed to clear predictions';
      throw err;
    }
  };

  // Get statistics from API or calculate locally
  const getStats = (): PredictionStatistics => {
    const total = predictions.value.length;
    
    if (total === 0) {
      return {
        total: 0,
        homeWins: 0,
        draws: 0,
        awayWins: 0,
        homeWinPercentage: 0,
        drawPercentage: 0,
        awayWinPercentage: 0,
      };
    }

    // Calculate stats based on predicted scores instead of text parsing
    const homeWins = predictions.value.filter(p => 
      p.predictedHomeScore != null && p.predictedAwayScore != null && 
      p.predictedHomeScore > p.predictedAwayScore
    ).length;
    
    const draws = predictions.value.filter(p => 
      p.predictedHomeScore != null && p.predictedAwayScore != null && 
      p.predictedHomeScore === p.predictedAwayScore
    ).length;
    
    const awayWins = predictions.value.filter(p => 
      p.predictedHomeScore != null && p.predictedAwayScore != null && 
      p.predictedHomeScore < p.predictedAwayScore
    ).length;

    return {
      total,
      homeWins,
      draws,
      awayWins,
      homeWinPercentage: total > 0 ? (homeWins / total) * 100 : 0,
      drawPercentage: total > 0 ? (draws / total) * 100 : 0,
      awayWinPercentage: total > 0 ? (awayWins / total) * 100 : 0,
    };
  };

  // Fetch statistics from API
  const fetchStatistics = async (): Promise<PredictionStatistics> => {
    try {
      return await apiClient.get('/predictions/statistics');
    } catch (err: any) {
      console.warn('Error fetching statistics from API, using local calculation:', err);
      return getStats();
    }
  };

  return {
    predictions,
    isLoading,
    error,
    fetchPredictions,
    savePrediction,
    deletePrediction,
    clearAllPredictions,
    getStats,
    fetchStatistics,
  };
}); 