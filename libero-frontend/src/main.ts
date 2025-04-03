import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia' // Import Pinia
import router from './router' // Import the router
import { useAuthStore } from './stores/auth' // Import the auth store

const app = createApp(App)
const pinia = createPinia() // Create Pinia instance

app.use(pinia) // Use Pinia
app.use(router) // Use the router

app.mount('#app') // Mount the app first

// Initialize authentication state after app is mounted
// This ensures Pinia is ready
const authStore = useAuthStore()
authStore.initializeAuth().catch(error => {
  console.error("Failed to initialize authentication:", error);
  // Handle initialization error if needed
});
