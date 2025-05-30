import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios';

// Backend URLs
const BACKEND_BASE_URL = 'http://localhost:8080'; // Only used for OAuth

// Create axios instance with custom config
const apiClient: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8080/api', // Update to match your backend URL
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // Increase timeout to 30 seconds
});

// Request Interceptor: Add auth token and debugging
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Fix common URL issues - API path corrections
    if (config.url && (config.url.startsWith('/api/') || config.url.startsWith('api/'))) {
      // Remove duplicate /api prefix since the baseURL already has it
      config.url = config.url.replace(/^\/?(api\/)/g, '');
      console.warn('Corrected duplicate API prefix in URL:', config.url);
    }
    
    console.log('API Request:', {
      url: config.url,
      method: config.method,
      headers: config.headers,
      data: config.data
    });
    return config;
  },
  (error: any): Promise<any> => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Response Interceptor: Add debugging and handle common errors
apiClient.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    console.log('API Response:', {
      url: response.config.url,
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error: any): Promise<any> => {
    // Handle network errors gracefully
    if (error.code === 'ECONNABORTED' || !error.response) {
      console.error('Network Error:', {
        url: error.config?.url || 'unknown',
        message: error.message,
        code: error.code
      });
    } else if (error.response) {
      console.error('API Error:', {
        url: error.config.url,
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers
      });
    } else if (error.request) {
      console.error('API Error: No response received', {
        url: error.config.url,
        request: error.request
      });
    } else {
      console.error('API Error:', error.message);
    }

    if (error.response && error.response.status === 401) {
      // Handle unauthorized access - e.g., clear token, redirect to login
      console.warn('Unauthorized access detected. Clearing token.');
      localStorage.removeItem('authToken');
      // Redirect to root page
      window.location.href = '/';
    }

    return Promise.reject(error);
  }
);

// Define interfaces for expected data structures
export interface LoginCredentials { // Add export keyword
  email: string;
  password: string;
}

interface LoginResponse {
  token: string;
}

export interface RegisterUserData {
  name: string;
  username: string;
  email: string;
  password: string;
}

// Response from registration endpoint (matches backend UserResponse)
export interface UserResponse {
  id: number;
  username: string;
  email: string;
  name?: string;
  role: string;
  created_at: string;
  updated_at: string;
}

// Define a basic user profile structure (adjust based on backend's UserResponse)
// Define interfaces for expected data structures from profile_feature_plan.md
interface Team {
  id: number;
  name: string;
  // Add other fields if needed later
}

interface Player {
  id: number;
  name: string;
  // Add other fields if needed later
}

// Add Competition interface
interface Competition {
  id: number;
  name: string;
  // Add other fields if needed later
}

interface UserPreferences {
  followed_teams: Team[];
  followed_players: Player[];
  followed_competitions: Competition[]; // Added competitions
}

// Updated UserProfile interface to include preferences
export interface UserProfile {
  id: number;
  username?: string; // Make username optional if backend doesn't always return it
  email: string;
  name?: string;
  role?: string; // Make role optional
  created_at?: string; // Make optional
  updated_at?: string; // Make optional
  preferences: UserPreferences; // Added preferences
}

// Interface for the PUT /api/users/preferences payload
export interface UpdatePreferencesPayload {
  add_teams?: number[];
  remove_teams?: number[];
  add_players?: number[];
  remove_players?: number[];
  add_competitions?: number[];    // Added competitions
  remove_competitions?: number[]; // Added competitions
}

// --- Fixtures DTOs ---

export interface FixtureMatchDTO {
  match_date: string;
  home_team_name: string;
  away_team_name: string;
  home_score?: number | null;
  away_score?: number | null;
  match_status: string;
  venue?: string;
  home_logo_url?: string;
  away_logo_url?: string;
}

export interface CompetitionFixturesDTO {
  competition_name: string;
  competition_code: string;
  logo_url?: string;
  matches: FixtureMatchDTO[];
}

// --- Fixtures Summary DTO ---
export interface FixturesSummaryDTO {
  competition_name: string;
  competition_code: string;
  logo_url?: string;
  today: FixtureMatchDTO[];
  tomorrow: FixtureMatchDTO[];
  upcoming: FixtureMatchDTO[];
}

// --- League Table & Player Stats DTOs ---
export interface LeagueTableRow {
  position: number;
  team: {
    id: number | string;
    name: string;
    logo: string;
  };
  played: number;
  won: number;
  drawn: number;
  lost: number;
  goalsFor: number;
  goalsAgainst: number;
  goalDifference: number;
  points: number;
}

export interface PlayerStat {
  id: number;
  name: string;
  team: {
    id: number;
    name: string;
    logo: string;
  };
  value: number; // goals, assists, or clean sheets
  photo: string;
}

// --- Authentication API Calls ---

/**
 * Logs in a user with email and password.
 * @param credentials - { email, password }
 * @returns Promise containing the login response (e.g., { token: "..." })
 */
export const loginUser = (credentials: LoginCredentials): Promise<LoginResponse> => {
  return apiClient.post<LoginResponse>('/auth/login', credentials)
    .then(response => response.data);
};

/**
 * Registers a new user.
 * @param userData - { username, email, password }
 * @returns Promise containing the registered user's profile (adjust type based on actual response)
 */
export const registerUser = (userData: RegisterUserData): Promise<UserResponse> => {
  console.log('[API] Registration request for user:', userData.username);
  
  // Assuming registration returns the user profile upon success
  return apiClient.post<UserResponse>('/auth/register', userData)
    .then(response => {
      console.log('[API] Registration successful for:', userData.username);
      return response.data;
    })
    .catch(error => {
      console.error('[API] Registration failed:', error.response?.data || error.message);
      throw error;
    });
};

/**
 * Requests a password reset token for the given email.
 * @param email - User's email address
 * @returns Promise containing the response message and token (for testing)
 */
export const requestPasswordReset = (email: string): Promise<{ message: string; token?: string }> => {
  return apiClient.post<{ message: string; token?: string }>('/auth/forgot-password', { email })
    .then(response => response.data);
};

/**
 * Resets user password using a reset token.
 * @param resetData - { token, new_password, confirm_password }
 * @returns Promise containing the success message
 */
export const resetPassword = (resetData: {
  token: string;
  new_password: string;
  confirm_password: string;
}): Promise<{ message: string }> => {
  return apiClient.post<{ message: string }>('/auth/reset-password', resetData)
    .then(response => response.data);
};

// --- User API Calls ---

/**
 * Fetches the profile of the currently authenticated user.
 * @returns Promise containing the user profile data
 */
export const getUserProfile = (): Promise<UserProfile> => {
  return apiClient.get<UserProfile>('/users/profile')
    .then(response => response.data);
};

/**
 * Updates the preferences of the currently authenticated user.
 * @param payload - Object containing arrays of IDs for teams/players to add/remove
 * @returns Promise<void | UserProfile> - Backend might return updated profile or just 204
 */
export const updateUserPreferences = (payload: UpdatePreferencesPayload): Promise<void | UserProfile> => {
  // Assuming the backend might return the updated profile or just a success status
  return apiClient.put<void | UserProfile>('/users/preferences', payload)
    .then(response => response.data); // Return data if available (e.g., updated profile)
};

// --- OAuth ---
// We don't call OAuth endpoints directly via JS/TS.
// Instead, we redirect the browser to the backend OAuth login URLs.

/**
 * Gets the backend URL to initiate Google OAuth login.
 * @returns {string}
 */
export const getGoogleLoginUrl = (): string => {
  // Note: We construct the full URL here, assuming the backend is at the root
  // Adjust if your setup is different (e.g., using API_BASE_URL)
  // Construct the full absolute URL to the backend endpoint
  return `${BACKEND_BASE_URL}/auth/google/login`; // <--- Changed
};

/**
 * Gets the backend URL to initiate Facebook OAuth login.
 * @returns {string}
 */
export const getFacebookLoginUrl = (): string => {
  return `${BACKEND_BASE_URL}/auth/facebook/login`; // <--- Changed
};

/**
 * Gets the backend URL to initiate GitHub OAuth login.
 * @returns {string}
 */
export const getGitHubLoginUrl = (): string => {
  return `${BACKEND_BASE_URL}/auth/github/login`; // <--- Changed
};


// --- Sports/Fixtures API Calls ---

/**
 * Fetches today's fixtures for relevant competitions.
 * @returns Promise containing an array of CompetitionFixturesDTO
 */
export const getTodaysFixtures = (): Promise<CompetitionFixturesDTO[]> => {
  return apiClient.get<CompetitionFixturesDTO[]>('/sports/fixtures/today')
    .then(response => response.data);
};

/**
 * Fetches summary of fixtures (today, tomorrow, upcoming) for a specific competition.
 * @param competitionCode - competition code (e.g., 'PL', 'CL')
 * @returns Promise containing the fixtures summary DTO
 */
export const getFixturesSummary = async (competitionCode: string): Promise<FixturesSummaryDTO> => {
  try {
    // Make the API request with correct path
    const response = await apiClient.get<FixturesSummaryDTO>(`/sports/fixtures/summary?competition=${competitionCode}`);
    return response.data;
  } catch (error: any) {
    console.error(`Error fetching fixtures summary for ${competitionCode}:`, error);
    
    // Create an empty response structure if the API fails
    // This allows the UI to render properly without errors
    const emptyResponse: FixturesSummaryDTO = {
      competition_name: "",
      competition_code: competitionCode,
      logo_url: "",
      today: [],
      tomorrow: [],
      upcoming: []
    };
    
    if (error.response && error.response.status === 404) {
      console.warn(`No fixtures found for ${competitionCode}`);
      return emptyResponse;
    }
    
    // For other errors, throw to let the store handle them
    throw error;
  }
};


// interface StandingsResponse {
//   standings: Array<{
//     stage: string;
//     table: Array<{
//       position: number;
//       team: {
//         id: number;
//         name: string;
//         crest: string;
//       };
//       playedGames: number;
//       won: number;
//       draw: number;
//       lost: number;
//       points: number;
//       goalsFor: number;
//       goalsAgainst: number;
//       goalDifference: number;
//     }>;
//   }>;
// }

// export const getStandings = async (competitionCode: string): Promise<LeagueTableRow[]> => {
//   const response = await apiClient.get<StandingsResponse>(`/standings?competition=${competitionCode}`);
//   // Find the total standings (usually the first group for league format)
//   const totalStandings = response.data.standings.find(s => s.stage === 'REGULAR_SEASON') || response.data.standings[0];
  
//   return totalStandings.table.map(row => ({
//     position: row.position,
//     team: {
//       id: row.team.id,
//       name: row.team.name,
//       logo: row.team.crest
//     },
//     played: row.playedGames,
//     won: row.won,
//     drawn: row.draw,
//     lost: row.lost,
//     goalDifference: row.goalDifference,
//     points: row.points
//   }));
// };


// In api.ts, for clarity, rename this:
interface BackendCompetitionStandingsDTO {
  competition_name: string;
  competition_code: string;
  season: number;
  standings: Array<{
      position: number;
      team_name: string;
      team_crest: string;
      played: number;  // Fixed: backend returns "played", not "played_games"
      won: number;
      drawn: number;  // Fixed: backend returns "drawn", not "draw"
      lost: number;
      goals_for: number;
      goals_against: number;
      goal_difference: number;
      points: number;
  }>;
}

export const getStandings = async (competitionCode: string): Promise<LeagueTableRow[]> => {
  console.log(`[API] Getting standings for ${competitionCode}`);
  // Expecting backend to return CompetitionStandingsDTO
  const response = await apiClient.get<BackendCompetitionStandingsDTO>(`/standings?competition=${competitionCode}`);
  
  // Show the full API response in a readable way for debugging
  console.log('[API] Raw standings response:', response.data);

  // Check if the response has the expected structure
  if (!response.data.standings || !Array.isArray(response.data.standings)) {
    console.error('[API] Invalid standings data structure:', response.data);
    return [];
  }

  console.log('[API] Processing standings data:', response.data.standings);

  // Transform backend data to match frontend interface
  return response.data.standings.map(row => ({
    position: row.position,
    team: {
      id: String(row.position), // Converting position to string since team.id can be string | number
      name: row.team_name,
      logo: row.team_crest,
    },
    played: row.played,  // Fixed: match corrected backend field
    won: row.won,
    drawn: row.drawn,  // Fixed: match corrected backend field
    lost: row.lost,
    goalsFor: row.goals_for,
    goalsAgainst: row.goals_against,
    goalDifference: row.goal_difference,
    points: row.points,
  }));
};

export const getTopScorers = async (competitionCode: string): Promise<PlayerStat[]> => {
  try {
    // Backend returns CompetitionScorersDTO with a 'scorers' array of flat objects
    const response = await apiClient.get<{ scorers: Array<{
      player_name: string;
      team_name: string;
      team_crest: string;
      goals: number;
      assists: number;
      penalties: number;
    }> }>(`/topscorers?competition=${competitionCode}`);

    const scorers = response.data.scorers || [];
    // If scorers is empty or not an array, return an empty array
    if (!Array.isArray(scorers) || scorers.length === 0) return [];

    return scorers.map((scorer, idx) => ({
      id: idx + 1, // Use index as fallback id
      name: scorer.player_name,
      team: {
        id: idx + 1, // No team id from backend, use index
        name: scorer.team_name,
        logo: scorer.team_crest
      },
      value: scorer.goals,
      photo: `/public/${scorer.player_name}.png` // Assuming player photos are stored with player names
    }));
  } catch (error) {
    console.error(`Error fetching top scorers for ${competitionCode}:`, error);
    // Return empty array to allow UI to fallback to mock data
    return [];
  }
};

// --- Prediction DTOs ---
export interface PredictMatchRequest {
  league: string;
  home_team: string;
  away_team: string;
}

export interface PredictMatchResponse {
  prediction: number;
  probabilities: {
    home_win: number;
    draw: number;
    away_win: number;
  };
  expected_home_goals: number;
  expected_away_goals: number;
  most_likely_home_score: number;
  most_likely_away_score: number;
}

// --- Top Scorers DTO ---
export interface TopScorerDTO {
  id: number;
  name: string;
  team: {
    id: number;
    name: string;
    logo: string;
  };
  value: number; // goals, assists, or clean sheets
}

export default apiClient;

// --- Match Prediction API Functions ---

/**
 * Predict the outcome of a football match
 * @param request - match prediction request with league, home_team, and away_team
 * @returns Promise containing the prediction response with scores and probabilities
 */
export const predictMatch = async (request: PredictMatchRequest): Promise<PredictMatchResponse> => {
  try {
    console.log('[API] Predicting match:', request);
    const response = await apiClient.post<PredictMatchResponse>('/predict/match', request);
    console.log('[API] Prediction response:', response.data);
    return response.data;
  } catch (error: any) {
    console.error('Error predicting match:', error);
    throw error;
  }
};

/**
 * Get list of available teams from the ML service
 * @returns Promise containing array of team names
 */
export const getAvailableTeams = async (): Promise<string[]> => {
  try {
    console.log('[API] Fetching available teams');
    const response = await apiClient.get<{ teams: string[] }>('/predict/teams');
    console.log('[API] Available teams:', response.data.teams);
    return response.data.teams;
  } catch (error: any) {
    console.error('Error fetching teams:', error);
    throw error;
  }
};

/**
 * Get list of available leagues from the ML service
 * @returns Promise containing array of league names
 */
export const getAvailableLeagues = async (): Promise<string[]> => {
  try {
    console.log('[API] Fetching available leagues');
    const response = await apiClient.get<{ leagues: string[] }>('/predict/leagues');
    console.log('[API] Available leagues:', response.data.leagues);
    return response.data.leagues;
  } catch (error: any) {
    console.error('Error fetching leagues:', error);
    throw error;
  }
};

// Get upcoming matches for team following
export const getUpcomingMatches = async (): Promise<CompetitionFixturesDTO[]> => {
  try {
    const response = await apiClient.get('/matches/upcoming');
    return response.data;
  } catch (error: any) {
    console.error('Error fetching upcoming matches:', error);
    throw new Error(`Failed to fetch upcoming matches: ${error.response?.data?.message || error.message}`);
  }
};
