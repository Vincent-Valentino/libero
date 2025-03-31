import { createRouter, createWebHistory } from 'vue-router';
import Home from '../home/Home.vue';

// Import the Auth component
import Auth from '../auth/Auth.vue';

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


const routes = [
  { path: '/', name: 'Home', component: Home },
  // Use LeaguePage for the top 5 leagues
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
];

// Add the auth route
routes.push({ path: '/auth', name: 'Auth', component: Auth });

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;