<template>
  <div class="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-white rounded-xl shadow-lg z-10">
      <div class="text-center">
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
          {{ currentTitle }}
        </h2>
        <p v-if="currentView === 'login' || currentView === 'signup'" class="mt-2 text-sm text-gray-600">
          Or
          <button @click="toggleForm" class="font-medium text-amber-500 hover:text-amber-600">
            {{ isLogin ? 'create an account' : 'log in to your account' }}
          </button>
        </p>
      </div>

      <!-- Dynamic Component for Different Forms -->
      <component 
        :is="currentFormComponent" 
        @forgot-password="showForgotPassword"
        @go-to-reset="showResetPassword"
        @back-to-login="showLogin"
        @password-reset-success="showLogin" />

      <!-- Social Logins (only for login/signup) -->
      <SocialLogins v-if="currentView === 'login' || currentView === 'signup'" />

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import LoginForm from './components/LoginForm.vue';
import SignupForm from './components/SignupForm.vue';
import ForgotPasswordForm from './components/ForgotPasswordForm.vue';
import ResetPasswordForm from './components/ResetPasswordForm.vue';
import SocialLogins from './components/SocialLogins.vue';

type ViewType = 'login' | 'signup' | 'forgot-password' | 'reset-password';

const currentView = ref<ViewType>('login');
const isLogin = ref(true);

const toggleForm = () => {
  if (currentView.value === 'login') {
    currentView.value = 'signup';
    isLogin.value = false;
  } else {
    currentView.value = 'login';
    isLogin.value = true;
  }
};

const showLogin = () => {
  currentView.value = 'login';
  isLogin.value = true;
};

const showForgotPassword = () => {
  currentView.value = 'forgot-password';
};

const showResetPassword = () => {
  currentView.value = 'reset-password';
};

const currentFormComponent = computed(() => {
  switch (currentView.value) {
    case 'login':
      return LoginForm;
    case 'signup':
      return SignupForm;
    case 'forgot-password':
      return ForgotPasswordForm;
    case 'reset-password':
      return ResetPasswordForm;
    default:
      return LoginForm;
  }
});

const currentTitle = computed(() => {
  switch (currentView.value) {
    case 'login':
      return 'Log in to your account';
    case 'signup':
      return 'Create an account';
    case 'forgot-password':
      return 'Reset your password';
    case 'reset-password':
      return 'Set new password';
    default:
      return 'Log in to your account';
  }
});

// TODO: Add Axios logic for form submissions
</script>

<style scoped>
/* Add any specific styles for Auth.vue if needed */
</style>