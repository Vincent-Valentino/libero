<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center p-8">
      <h2 class="text-2xl font-semibold mb-4">Processing login...</h2>
      <!-- Optional: Add a loading spinner -->
      <p v-if="error" class="text-red-500">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();
const error = ref<string | null>(null);

onMounted(async () => {
  const hash = window.location.hash;
  let token: string | null = null;

  if (hash.startsWith('#token=')) {
    token = hash.substring('#token='.length);
  }

  // Clear the token from the URL hash regardless of success/failure
  // Use replaceState to avoid adding a new entry to the browser history
  history.replaceState(null, '', window.location.pathname + window.location.search);


  if (token) {
    console.log('OAuth Callback: Token found in hash.');
    try {
      // Store the token using the auth store action
      authStore.setToken(token);

      // Fetch the user profile using the new token
      await authStore.fetchProfile();

      if (authStore.isAuthenticated && authStore.user) {
        console.log('OAuth Callback: Profile fetched successfully. Redirecting to profile.');
        // Redirect to the user's profile page or dashboard
        router.push({ name: 'Profile' }); // Or 'Home' or a saved redirect path
      } else {
        // This case might happen if fetchProfile fails despite having a token
        console.error('OAuth Callback: Failed to fetch profile after setting token.');
        error.value = 'Login successful, but failed to retrieve user details. Please try logging in again.';
        authStore.clearAuth(); // Clear potentially invalid token
        // Redirect to login after a delay or keep the error message displayed
        setTimeout(() => router.push({ name: 'Auth' }), 4000);
      }
    } catch (err) {
      console.error('OAuth Callback: Error processing token or fetching profile:', err);
      error.value = 'An error occurred during the final login step. Please try again.';
      authStore.clearAuth();
      setTimeout(() => router.push({ name: 'Auth' }), 4000);
    }
  } else {
    console.error('OAuth Callback: No token found in URL hash.');
    error.value = 'Authentication callback failed. No token received. Redirecting to login.';
    // Redirect back to login page if no token is present
    setTimeout(() => router.push({ name: 'Auth' }), 3000);
  }
});
</script>

<style scoped>
/* Add styles if needed */
</style>