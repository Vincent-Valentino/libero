import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import * as api from '@/services/api'; // Import API service functions
import { AxiosError } from 'axios'; // Import AxiosError for type checking

// Define the structure for the user profile state
// Align this with the UserProfile interface in api.ts
interface UserState {
  id: number | null;
  username: string | null;
  email: string | null;
  name?: string | null;
  role: string | null;
  // Add other relevant fields from UserProfile if needed
}

// Define the structure for the auth state
interface AuthState {
  token: string | null;
  user: UserState | null;
  status: 'idle' | 'loading' | 'success' | 'error'; // For tracking async operations
}

export const useAuthStore = defineStore('auth', () => {
  // --- State ---
  const token = ref<string | null>(localStorage.getItem('authToken') || null);
  const user = ref<UserState | null>(null); // User profile info
  const status = ref<AuthState['status']>('idle'); // Track loading/error states

  // --- Getters (Computed Properties) ---
  const isAuthenticated = computed<boolean>(() => !!token.value);
  const isLoading = computed<boolean>(() => status.value === 'loading');
  const authUser = computed<UserState | null>(() => user.value);

  // --- Actions ---

  /**
   * Sets the authentication token and stores it.
   * @param newToken The JWT token string.
   */
  function setToken(newToken: string | null) {
    token.value = newToken;
    if (newToken) {
      localStorage.setItem('authToken', newToken);
    } else {
      localStorage.removeItem('authToken');
    }
  }

  /**
   * Clears authentication state (token and user).
   */
  function clearAuth() {
    setToken(null);
    user.value = null;
    status.value = 'idle';
  }

  /**
   * Attempts to log in the user with credentials.
   * @param credentials User's email and password.
   */
  async function login(credentials: { email: string; password: string }): Promise<void> {
    status.value = 'loading';
    try {
      // Pass credentials as a single object argument
      const response = await api.loginUser({ email: credentials.email, password: credentials.password });
      setToken(response.token);
      await fetchProfile(); // Fetch profile after successful login
      status.value = 'success';
    } catch (error: unknown) { // Explicitly type error as unknown
      console.error('Login failed:', error);
      clearAuth(); // Clear any partial state on failure
      status.value = 'error';
      throw error; // Re-throw error to be handled by the component
    }
  }

  /**
   * Registers a new user.
   * @param userData User registration data.
   */
  async function register(userData: { username: string; email: string; password: string }): Promise<void> {
    status.value = 'loading';
    try {
      // Assuming registration directly returns the user profile (adjust if it returns something else)
      const registeredUser = await api.registerUser(userData);
      // Optionally log the user in immediately after registration
      // await login({ email: userData.email, password: userData.password });
      status.value = 'success'; // Or 'idle' if not logging in automatically
      // Note: Registration might not return a token directly.
      // The user might need to log in separately after registering.
    } catch (error) {
      console.error('Registration failed:', error);
      // No need to check error type here unless we need specific handling
      status.value = 'error';
      throw error; // Re-throw error
    }
  }

  /**
   * Logs out the current user.
   */
  function logout(): void {
    clearAuth();
    // Optionally redirect to login page or home page via router
    // import router from '@/router'; // Be careful with imports inside store actions
    // router.push('/login');
    console.log('User logged out.');
  }

  /**
   * Fetches the current user's profile from the backend.
   */
  async function fetchProfile(): Promise<void> {
    if (!token.value) {
      console.warn('Cannot fetch profile without a token.');
      return;
    }
    status.value = 'loading';
    try {
      const profileData = await api.getUserProfile();
      user.value = { // Map backend response to UserState
        id: profileData.id,
        username: profileData.username,
        email: profileData.email,
        name: profileData.name,
        role: profileData.role,
      };
      status.value = 'success';
    } catch (error: unknown) { // Explicitly type error as unknown
      console.error('Failed to fetch user profile:', error);
      // If profile fetch fails (e.g., token expired), clear auth state
      // Check if it's an AxiosError before accessing response
      if (error instanceof AxiosError && error.response && error.response.status === 401) {
          clearAuth();
      }
      status.value = 'error';
      // Don't re-throw here unless necessary, as it might break initial load checks
    }
  }

  /**
   * Initializes the store, typically on app load.
   * Checks for an existing token and fetches the user profile if found.
   */
  async function initializeAuth(): Promise<void> {
      if (token.value && !user.value) { // Only fetch if token exists but user data is missing
          console.log('Initializing auth: Found token, fetching profile...');
          await fetchProfile();
      } else if (!token.value) {
          console.log('Initializing auth: No token found.');
          status.value = 'idle';
      } else {
          console.log('Initializing auth: User data already present.');
          status.value = 'success'; // Assume success if user data exists
      }
  }


  // --- Return state, getters, and actions ---
  return {
    token,
    user,
    status,
    isAuthenticated,
    isLoading,
    authUser,
    login,
    register,
    logout,
    fetchProfile,
    setToken, // Expose if needed externally (e.g., for OAuth callback)
    initializeAuth,
    clearAuth,
  };
});