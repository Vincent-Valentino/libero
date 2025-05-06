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

export default apiClient; // Export the configured instance if needed elsewhere