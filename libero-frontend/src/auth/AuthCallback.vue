<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center p-8">
      <div class="mb-4">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      </div>
      <h2 class="text-2xl font-semibold mb-4">Processing login...</h2>
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
      // Call the single action to handle the token and fetch profile
      await authStore.handleAuthCallback(token);

      // After handleAuthCallback completes, check the store's state
      if (authStore.isAuthenticated && authStore.user) {
        console.log('OAuth Callback: handleAuthCallback successful. Redirecting to Profile.');
        // Redirect to the profile page for successful OAuth login
        router.push({ name: 'profile' });
      } else {
        // This case means handleAuthCallback finished but the user is not authenticated
        // (likely because the profile fetch within it failed and triggered logout)
        console.error('OAuth Callback: handleAuthCallback completed, but user is not authenticated. Profile fetch likely failed.');
        // Use the error message potentially set in the store by fetchUserProfile failure
        error.value = authStore.error || 'Login succeeded, but failed to load user details. Please try again.';
        // No need to call clearAuth here, as handleAuthCallback/fetchUserProfile should have handled it
        setTimeout(() => router.push({ name: 'Home' }), 4000); // Redirect home on error
      }
    } catch (err: any) {
      // Catch any unexpected errors during handleAuthCallback itself
      console.error('OAuth Callback: Unexpected error during handleAuthCallback:', err);
      error.value = err.message || 'An unexpected error occurred during login processing. Please try again.';
      // Ensure state is cleared if an unexpected error occurs
      authStore.logout(); // Use the logout action to clear state
      setTimeout(() => router.push({ name: 'Home' }), 4000); // Redirect home on error
    }
  } else {
    console.error('OAuth Callback: No token found in URL hash.');
    error.value = 'Authentication callback failed. No token received. Redirecting to login.';
    // Redirect back to login page if no token is present
    setTimeout(() => router.push({ name: 'Home' }), 3000); // Redirect home if no token
  }
});
</script>

<style scoped>
/* Add styles if needed */
</style>