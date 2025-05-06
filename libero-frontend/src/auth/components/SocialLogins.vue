<template>
  <div class="mt-6">
    <div class="relative">
      <div class="absolute inset-0 flex items-center">
        <div class="w-full border-t border-gray-300"></div>
      </div>
      <div class="relative flex justify-center text-sm">
        <span class="px-2 bg-white text-gray-500">
          Or continue with
        </span>
      </div>
    </div>

    <div class="mt-6 grid grid-cols-3 gap-3">
      <div>
        <button @click="handleSocialLogin('google')"
                class="w-full inline-flex cursor-pointer justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
          <span class="sr-only">Sign in with Google</span>
          <!-- Google Logo SVG -->
          <svg  xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
        </button>
      </div>

      <div>
        <button @click="handleSocialLogin('facebook')"
                class="w-full inline-flex cursor-pointer justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
          <span class="sr-only">Sign in with Facebook</span>
          <!-- TODO: Add Facebook Logo SVG -->
          <!-- Facebook Logo SVG -->
          <svg class="w-5 h-5" fill="#1877F2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path d="M22.676 0H1.324C.593 0 0 .593 0 1.324v21.352C0 23.407.593 24 1.324 24h11.494v-9.294H9.689v-3.621h3.129V8.41c0-3.1 1.893-4.785 4.659-4.785 1.325 0 2.464.097 2.796.141v3.24h-1.918c-1.504 0-1.795.715-1.795 1.763v2.309h3.587l-.467 3.621h-3.12V24h6.116c.73 0 1.323-.593 1.323-1.324V1.324C24 .593 23.407 0 22.676 0z"/>
          </svg>
        </button>
      </div>

      <div>
        <button @click="handleSocialLogin('github')"
                class="w-full inline-flex cursor-pointer justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
          <span class="sr-only">Sign in with GitHub</span>
          <!-- TODO: Add GitHub Logo SVG -->
          <!-- GitHub Logo SVG -->
          <svg class="w-5 h-5" fill="#1B1F23" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// import axios from 'axios'; // Import axios when ready
import { getGoogleLoginUrl, getFacebookLoginUrl, getGitHubLoginUrl } from '@/services/api';

const handleSocialLogin = (provider: 'google' | 'facebook' | 'github') => {
  console.log(`Attempting social login with ${provider}`);

  let redirectUrl: string;

  switch (provider) {
    case 'google':
      redirectUrl = getGoogleLoginUrl();
      break;
    case 'facebook':
      redirectUrl = getFacebookLoginUrl();
      break;
    case 'github':
      redirectUrl = getGitHubLoginUrl();
      break;
    default:
      console.error('Unknown social login provider:', provider);
      return;
  }

  // Redirect the browser to the backend OAuth endpoint
  // The backend will handle the redirect to the actual provider
  // Alternative method: Create and click a link
  window.location.href = redirectUrl;
};

</script>