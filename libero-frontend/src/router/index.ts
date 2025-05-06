import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'; // Use type-only import
import Home from '../home/Home.vue';

// --- Auth Components ---
// Import the Auth component
import Auth from '../auth/Auth.vue';
// Placeholder for OAuth callback handler
const AuthCallback = () => import('../auth/AuthCallback.vue'); // Use dynamic import

// Import the LeaguePage component
import LeaguePage from '../leagues/LeaguePage.vue';

// Placeholder components for league pages - we'll create these next
// Placeholders for leagues we are not implementing yet
const UCL = { template: '<div>UCL Page</div>' };
const UEL = { template: '<div>UEL Page</div>' };
// Add placeholders for other links if needed
const Player = { template: '<div>Player Page</div>' };
const Team = { template: '<div>Team Page</div>' };
const Nations = { template: '<div>Nations Page</div>' };
const Awards = { template: '<div>Awards Page</div>' };


// --- Pinia Store ---
// Import outside router setup to be used in navigation guard
import { useAuthStore } from '@/stores/auth';


const routes: Array<RouteRecordRaw> = [ // Added type annotation
  { path: '/', name: 'Home', component: Home },
  // Auth routes (public)
  { path: '/auth', name: 'Auth', component: Auth, meta: { guestOnly: true } }, // Redirect if logged in
  { path: '/auth/callback', name: 'AuthCallback', component: AuthCallback }, // Handles OAuth redirect
  // Protected routes
  // Existing routes (assuming public for now, add meta: { requiresAuth: true } if needed)
  // Use LeaguePage for the top 5 leaguesx`
  { path: '/premier-league', name: 'PremierLeague', component: LeaguePage },
  { path: '/la-liga', name: 'LaLiga', component: LeaguePage },
  { path: '/serie-a', name: 'SerieA', component: LeaguePage },
  { path: '/bundesliga', name: 'Bundesliga', component: LeaguePage },
  { path: '/ligue-1', name: 'Ligue1', component: LeaguePage },
  // Keep placeholders for others
  { path: '/ucl', name: 'UCL', component: UCL },
  { path: '/uel', name: 'UEL', component: UEL },
  { path: '/player', name: 'Player', component: Player },
  { path: '/team', name: 'Team', component: Team },
  { path: '/nations', name: 'Nations', component: Nations },
  { path: '/awards', name: 'Awards', component: Awards },
  // Profile route (protected)
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/profile/ProfilePage.vue'), // Lazy load
    meta: { requiresAuth: true } // Ensure user is logged in
  },
];

// Catch-all route for 404 errors - MUST BE LAST
const catchAllRoute: RouteRecordRaw = {
  path: '/:pathMatch(.*)*', // Matches any path not matched above
  name: 'NotFound',
  redirect: { name: 'Home' } // Redirect to the root page
};

routes.push(catchAllRoute); // Add the catch-all route to the array

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation Guard
router.beforeEach((to, _, next) => {
  const authStore = useAuthStore(); // Get store instance
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const guestOnly = to.matched.some(record => record.meta.guestOnly);
  const isAuthenticated = authStore.isAuthenticated;

  // Ensure Pinia store is initialized (especially user profile) if token exists
  // Note: initializeAuth is likely already called in main.ts, but this adds safety
  // if (!authStore.user && isAuthenticated) {
  //   await authStore.fetchProfile(); // Consider potential race conditions/performance
  // }

  if (requiresAuth && !isAuthenticated) {
    // Redirect to login page if trying to access protected route without auth
    console.log(`Navigation Guard: Route ${to.path} requires auth. Redirecting to /.`);
    next({ name: 'Home' }); // <--- Changed from 'Auth'
  } else if (guestOnly && isAuthenticated) {
    // Redirect to home/profile if trying to access login/register page while authenticated
    console.log(`Navigation Guard: Route ${to.path} is guest only. Redirecting to /.`);
    next({ name: 'Home' }); // Or 'Profile'
  } else {
    // Allow navigation
    next();
  }
});

export default router;