import { defineStore } from 'pinia';
// Import API functions and types
import { loginUser, registerUser, getUserProfile, type LoginCredentials, type RegisterUserData, type UserProfile } from '@/services/api'; // Use @ alias for cleaner imports

// User interface is now imported as UserProfile

// Define the state structure
interface AuthState {
  isAuthenticated: boolean;
  user: UserProfile | null; // Use the imported UserProfile type
  token: string | null;
  loading: boolean;
  error: string | null;
}

// Define the store
export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    isAuthenticated: !!localStorage.getItem('authToken'), // Initialize based on stored token
    user: null,
    token: localStorage.getItem('authToken'), // Load token from localStorage
    loading: false,
    error: null,
  }),

  getters: {
    // Example getter: Check if user is authenticated
    isLoggedIn(): boolean {
      return this.isAuthenticated && !!this.token;
    },
    // Example getter: Get user profile
    getUserProfile(): UserProfile | null { // Update return type
        return this.user;
    },
    // Example getter: Get loading state
    isLoading(): boolean {
        return this.loading;
    },
    // Example getter: Get error message
    getError(): string | null {
        return this.error;
    }
  },

  actions: {
    // Action to handle user login
    // Action to handle user login using the api service
    async login(credentials: LoginCredentials) { // Use imported LoginCredentials type
      this.loading = true;
      this.error = null;
      try {
        // Call the loginUser function from api.ts
        const loginResponse = await loginUser(credentials);
        const token = loginResponse.token;

        // Store the token
        this.token = token;
        localStorage.setItem('authToken', token); // Store token securely
        this.isAuthenticated = true; // Mark as authenticated

        // After successful login and token storage, fetch the user profile
        await this.fetchUserProfile(); // Fetch profile using the new token

        // If fetchUserProfile fails, the error state will be set there.
        // If fetchUserProfile succeeds, this.user will be populated.

      } catch (err: any) {
        // Error handling delegated to api.ts interceptor for 401,
        // but catch other potential errors (network, validation, etc.)
        this.error = err.response?.data?.message || err.message || 'Login failed';
        this.isAuthenticated = false;
        this.user = null;
        this.token = null;
        localStorage.removeItem('authToken');
      } finally {
        this.loading = false;
      }
    },

    // Action to fetch user profile (e.g., after initial login or page load)
    // Action to fetch user profile using the api service
    async fetchUserProfile() {
      // Token check is useful to prevent unnecessary API calls
      if (!this.token) {
        // This case might happen if initializeAuth is called before a token is set (e.g., first visit)
        // Or if the token was cleared due to an error elsewhere.
        // console.warn('Attempted to fetch profile without a token.');
        // No need to set an error here unless it's unexpected.
        // Ensure state is clean if no token exists.
        this.isAuthenticated = false;
        this.user = null;
        return;
      }

      this.loading = true;
      this.error = null;
      try {
        // Call the getUserProfile function from api.ts
        // Auth header is handled by the interceptor in api.ts
        const userProfile = await getUserProfile();
        this.user = userProfile;
        this.isAuthenticated = true; // Re-affirm authentication status
      } catch (err: any) {
        // The api.ts interceptor handles 401 by clearing the token and redirecting.
        // The store's state needs to be updated accordingly.
        // The interceptor might have already called logout logic indirectly (via page reload).
        // We still catch other errors (e.g., network, server errors 5xx).
        this.error = err.response?.data?.message || err.message || 'Failed to fetch user profile';
        // If an error occurs fetching the profile, assume authentication is compromised
        this.logout(); // Clean up store state on profile fetch failure
      } finally {
        this.loading = false;
      }
    },

    // Action to handle user logout
    logout() {
      this.isAuthenticated = false;
      this.user = null;
      this.token = null;
      this.error = null;
      localStorage.removeItem('authToken'); // Remove token from storage
      // Optionally redirect to login page or perform other cleanup
      // Example: router.push('/login');
    },

    // Action to initialize the store, e.g., on app startup
    async initializeAuth() {
        if (this.token && !this.user) {
            // If token exists but user profile is not loaded, fetch it
            await this.fetchUserProfile();
        }
        // TODO: Implement token refresh logic if needed
        // Check token expiry, refresh if necessary using a refresh token
    },

    // Action to handle authentication callback (e.g., for OAuth)
    // Action to handle authentication callback (e.g., for OAuth)
    async handleAuthCallback(token: string) { // Make async to await fetchUserProfile
        this.token = token;
        this.isAuthenticated = true; // Assume authenticated for now
        localStorage.setItem('authToken', token);
        // Fetch user profile immediately after obtaining token
        await this.fetchUserProfile(); // Await the profile fetch
        // If fetchUserProfile fails, it will handle logout/error state.
    },

    // Action to clear errors
    clearError() {
        this.error = null;
    },

    // Action to register a new user
    async register(userData: RegisterUserData) {
      this.loading = true;
      this.error = null;
      try {
        // Call the registerUser function from api.ts
        // Registration returns user profile but doesn't automatically log in
        await registerUser(userData);
        
        // Registration successful - user needs to login separately
        // Don't set authentication state here since registration != login
        
      } catch (err: any) {
        // Error handling for registration errors
        this.error = err.response?.data?.message || err.message || 'Registration failed';
        throw err; // Re-throw so the component can handle specific errors
      } finally {
        this.loading = false;
      }
    },
  },
});