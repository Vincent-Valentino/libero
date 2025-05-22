import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios';

// Define interfaces for expected data structures
export interface LoginCredentials { // Add export keyword
  email: string;
  password: string;
}

interface LoginResponse {
  token: string;
}

interface RegisterUserData {
  username: string;
  email: string;
  password: string;
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

// Base URL for the backend API
// TODO: Make this configurable via environment variables (.env)
// Use relative path for API calls; Vite proxy will handle forwarding
const API_BASE_URL: string = '/api';
// Explicit base URL for backend, needed for constructing full OAuth URLs
const BACKEND_BASE_URL: string = 'http://localhost:8080'; // <-- Added

const apiClient: AxiosInstance = axios.create({
  baseURL: `${BACKEND_BASE_URL}${API_BASE_URL}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Add JWT token to headers if available
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token: string | null = localStorage.getItem('authToken'); // Or sessionStorage
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error: any): Promise<any> => {
    return Promise.reject(error);
  }
);

// Response Interceptor: Handle common errors (e.g., 401 Unauthorized)
apiClient.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    // Any status code that lie within the range of 2xx cause this function to trigger
    return response;
  },
  (error: any): Promise<any> => {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    console.error('API Error:', error.response?.data || error.message);

    if (error.response && error.response.status === 401) {
      // Handle unauthorized access - e.g., clear token, redirect to login
      console.warn('Unauthorized access detected. Clearing token.');
      localStorage.removeItem('authToken'); // Or sessionStorage
      // Redirect to root page
      window.location.href = '/'; // <--- Changed & Uncommented
      // Optionally, could use router.push('/') if router instance is available here
    }

    // Return the error so it can be caught by the calling code
    return Promise.reject(error);
  }
);

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
export const registerUser = (userData: RegisterUserData): Promise<UserProfile> => {
  // Assuming registration returns the user profile upon success
  return apiClient.post<UserProfile>('/auth/register', userData)
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
export const getFixturesSummary = (competitionCode: string): Promise<FixturesSummaryDTO> => {
  return apiClient.get<FixturesSummaryDTO>(`/sports/fixtures/summary?competition=${competitionCode}`)
    .then(response => response.data);
};

// --- Football Data API Types ---
export interface FootballPlayer {
  id: number;
  name: string;
  firstName?: string;
  lastName?: string;
  dateOfBirth?: string;
  nationality?: string;
  position?: string;
  photo?: string;
}

export interface FootballTeam {
  id: number;
  name: string;
  shortName?: string;
  crest?: string;
}

export interface TopScorer {
  player: FootballPlayer;
  team: FootballTeam;
  playedMatches?: number;
  goals: number;
  assists?: number;
  penalties?: number;
}

export interface TopScorersResponse {
  count: number;
  filters?: any;
  competition?: any;
  season?: any;
  scorers: TopScorer[];
}

// --- Football Data API Client ---

// Create a separate axios instance for football-data.org API
const footballApiClient: AxiosInstance = axios.create({
  baseURL: 'https://api.football-data.org/v4',
  headers: {
    'Content-Type': 'application/json',
    // Use the API key from environment variables
    'X-Auth-Token': import.meta.env.THIRD_PARTY_FOOTBALL_API_KEY
  },
  timeout: 10000 // 10 second timeout
});

// Error handling for football API
footballApiClient.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    return response;
  },
  (error: any): Promise<any> => {
    // Log detailed error information
    console.error('Football API Error:', error.response?.data || error.message);
    
    // Handle common errors
    if (error.message && error.message.includes('Network Error')) {
      console.error(`
        CORS ERROR DETECTED: The football-data.org API has CORS restrictions.
        
        Possible solutions:
        1. Use a backend proxy server to make API requests
        2. Try a browser extension that disables CORS for development
        3. Make sure your API key is valid
        
        For more information, see: https://www.football-data.org/documentation/quickstart
      `);
    }
    
    if (error.response && error.response.status === 403) {
      console.error(`
        API KEY ERROR: Your API key was rejected by football-data.org
        
        Make sure:
        1. You have a valid API key in your .env file
        2. Your VITE_THIRD_PARTY_FOOTBALL_API_KEY environment variable is set correctly
        3. Your API subscription is active
        
        For more information, see: https://www.football-data.org/documentation/quickstart
      `);
    }
    
    return Promise.reject(error);
  }
);

// --- Football Data API Calls ---

/**
 * Fetches top scorers for a competition from football-data.org
 * @param competitionCode - Competition code (e.g., 'PL' for Premier League)
 * @param limit - Number of top scorers to retrieve (default: 7)
 * @returns Promise containing the top scorers data
 */
export const getTopScorers = (competitionCode: string, limit: number = 7): Promise<TopScorersResponse> => {
  return footballApiClient.get<TopScorersResponse>(`/competitions/${competitionCode}/scorers?limit=${limit}`)
    .then(response => response.data)
    .catch(error => {
      console.error(`Failed to fetch top scorers for ${competitionCode}:`, error);
      // Return empty result structure on error
      return {
        count: 0,
        scorers: []
      };
    });
};

/**
 * Fetches details for a specific player from football-data.org
 * @param playerId - ID of the player
 * @returns Promise containing the player details
 */
export const getPlayerDetails = (playerId: number): Promise<FootballPlayer> => {
  return footballApiClient.get<FootballPlayer>(`/persons/${playerId}`)
    .then(response => response.data)
    .catch(error => {
      console.error(`Failed to fetch player details for ID ${playerId}:`, error);
      throw error;
    });
};

/**
 * Fetches details for a specific team from football-data.org
 * @param teamId - ID of the team
 * @returns Promise containing the team details
 */
export const getTeamDetails = (teamId: number): Promise<FootballTeam> => {
  return footballApiClient.get<FootballTeam>(`/teams/${teamId}`)
    .then(response => response.data)
    .catch(error => {
      console.error(`Failed to fetch team details for ID ${teamId}:`, error);
      throw error;
    });
};

export default apiClient; // Export the configured instance if needed elsewhere

// Also export the football API client for direct use if needed
export { footballApiClient };