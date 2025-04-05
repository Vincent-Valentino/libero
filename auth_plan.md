# Authentication and Routing Configuration Plan

## Goal

Configure the routing and authentication flow for the `libero-frontend` and `libero-backend` application. This includes:
1.  Fixing the Google OAuth login flow (currently showing a blank page).
2.  Implementing frontend 404 handling: Redirect invalid frontend routes to `/`.
3.  Implementing authentication required handling: Redirect unauthenticated users trying to access protected routes or receiving a 401 API response to `/`.
4.  Implementing post-login handling: Redirect users to `/` after successful login/OAuth callback.
5.  Ensuring backend 404 API errors are handled locally in the component without a redirect.

## Backend Base URL

The backend server is assumed to be running at: `http://localhost:8080`

## Required Changes

### 1. Modify `libero-frontend/src/services/api.ts`

*   **Update OAuth URL functions:** Prepend the full backend base URL (`http://localhost:8080`) to the relative paths (`/auth/google/login`, etc.).
*   **Update 401 Interceptor:** Ensure the response interceptor redirects the user to the root (`/`) upon receiving a 401 Unauthorized status from the backend.

```typescript
// libero-frontend/src/services/api.ts
import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios';

// ... (interfaces) ...

// Base URL for the backend API (for proxied requests)
const API_BASE_URL: string = '/api';
// Explicit base URL for backend, needed for constructing full OAuth URLs
const BACKEND_BASE_URL: string = 'http://localhost:8080'; // <-- Added

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Add JWT token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const token: string | null = localStorage.getItem('authToken');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error: any): Promise<any> => {
    return Promise.reject(error);
  }
);

// Response Interceptor: Handle common errors
apiClient.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    return response;
  },
  (error: any): Promise<any> => {
    console.error('API Error:', error.response?.data || error.message);

    if (error.response && error.response.status === 401) {
      console.warn('Unauthorized access detected. Clearing token and redirecting to root.');
      localStorage.removeItem('authToken');
      // Redirect to root page
      window.location.href = '/'; // <--- Changed & Uncommented
    }
    // 404s will naturally reject and be handled by calling code

    return Promise.reject(error);
  }
);

// ... (loginUser, registerUser, getUserProfile) ...

// --- OAuth ---

export const getGoogleLoginUrl = (): string => {
  // Construct the full absolute URL to the backend endpoint
  return `${BACKEND_BASE_URL}/auth/google/login`; // <--- Changed
};

export const getFacebookLoginUrl = (): string => {
  return `${BACKEND_BASE_URL}/auth/facebook/login`; // <--- Changed
};

export const getGitHubLoginUrl = (): string => {
  return `${BACKEND_BASE_URL}/auth/github/login`; // <--- Changed
};

export default apiClient;
```

### 2. Modify `libero-frontend/src/router/index.ts`

*   **Add Catch-all Route:** Add a route definition (`path: '/:pathMatch(.*)*'`) at the end of the `routes` array that redirects any unmatched paths to the Home route (`/`).
*   **Update Navigation Guard:** Modify the `router.beforeEach` guard to redirect to `{ name: 'Home' }` (path `/`) when `requiresAuth` is true but the user is not authenticated.

```typescript
// libero-frontend/src/router/index.ts
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import Home from '../home/Home.vue';
import Auth from '../auth/Auth.vue';
const AuthCallback = { template: '<div>Processing login...</div>' }; // Placeholder
const UserProfile = { template: '<div>User Profile Page (Protected)</div>' };
import LeaguePage from '../leagues/LeaguePage.vue';
// ... (other imports) ...
import { useAuthStore } from '@/stores/auth';

const routes: Array<RouteRecordRaw> = [
  { path: '/', name: 'Home', component: Home },
  { path: '/auth', name: 'Auth', component: Auth, meta: { guestOnly: true } },
  { path: '/auth/callback', name: 'AuthCallback', component: AuthCallback }, // Needs implementation
  { path: '/profile', name: 'Profile', component: UserProfile, meta: { requiresAuth: true } },
  // ... (league routes, etc.) ...

  // Catch-all route for 404 errors - MUST BE LAST
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: { name: 'Home' } // Redirect to the root page
  } // <--- Added
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation Guard
router.beforeEach((to, _, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const guestOnly = to.matched.some(record => record.meta.guestOnly);
  const isAuthenticated = authStore.isAuthenticated;

  if (requiresAuth && !isAuthenticated) {
    console.log(`Navigation Guard: Route ${to.path} requires auth. Redirecting to /.`);
    next({ name: 'Home' }); // <--- Changed from 'Auth'
  } else if (guestOnly && isAuthenticated) {
    console.log(`Navigation Guard: Route ${to.path} is guest only. Redirecting to /.`);
    next({ name: 'Home' });
  } else {
    next();
  }
});

export default router;
```

### 3. Implement `libero-frontend/src/auth/AuthCallback.vue`

This component needs implementation to handle the redirect back from the backend after successful OAuth.

*   **Logic:**
    1.  Check `window.location.hash` for `#token=...`.
    2.  Extract the token value.
    3.  Use the `useAuthStore` to store the token (e.g., in `localStorage`) and update the application's auth state.
    4.  Use the router (`useRouter`) to navigate the user to the root (`/`) via `router.push({ name: 'Home' })`.
    5.  Handle cases where the token is missing or invalid.

*   **Example Structure (Conceptual):**

```vue
<template>
  <div>Processing login...</div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth'; // Assuming store path

const router = useRouter();
const authStore = useAuthStore();

onMounted(() => {
  const hash = window.location.hash;
  if (hash.startsWith('#token=')) {
    const token = hash.substring('#token='.length);
    if (token) {
      console.log('AuthCallback: Token received');
      // TODO: Implement this method in the auth store
      // It should store the token, update state, maybe fetch profile, then redirect
      authStore.handleTokenCallback(token)
        .then(() => {
          console.log('AuthCallback: Redirecting to Home');
          router.push({ name: 'Home' });
        })
        .catch(error => {
          console.error('AuthCallback: Error handling token:', error);
          // Redirect to home even on error, or maybe a specific error page/login
          router.push({ name: 'Home' });
        });
    } else {
      console.error('AuthCallback: Token fragment found but token is empty.');
      router.push({ name: 'Home' }); // Redirect if token is empty
    }
  } else {
    console.error('AuthCallback: No token fragment found in URL hash.');
    router.push({ name: 'Home' }); // Redirect if no token
  }
});
</script>
```

## Implementation Flow Diagram

```mermaid
graph TD
    subgraph Frontend Changes
        A[SocialLogins.vue] -- calls --> B(api.ts: getGoogleLoginUrl);
        B -- returns --> C{Full Backend URL\nhttp://.../auth/google/login};
        A -- window.location.href --> D[Browser Navigates to Backend];

        E[api.ts: Interceptor] -- detects 401 --> F[Clear Token];
        F -- window.location.href --> G[Redirect to /];

        H[router/index.ts] -- defines --> I(Catch-all Route :pathMatch);
        I -- redirects to --> G;

        J[router/index.ts] -- Navigation Guard --> K{Requires Auth?};
        K -- Yes & Not Authenticated --> G;

        L[Backend Redirects After OAuth] -- to --> M(FE: /auth/callback#token=...);
        N[AuthCallback.vue] -- reads --> O(window.location.hash);
        N -- extracts --> P[Token];
        N -- calls --> Q(authStore.handleTokenCallback);
        Q -- stores --> R[localStorage];
        Q -- router.push --> G;
    end

    D --> Backend;
    Backend --> L;

    style G fill:#f9f,stroke:#333,stroke-width:2px