<template>
  <div class="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-white rounded-xl shadow-lg z-10">
      <div class="text-center">
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
          {{ isLogin ? 'Log in to your account' : 'Create an account' }}
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          Or
          <button @click="toggleForm" class="font-medium text-amber-500 hover:text-amber-600">
            {{ isLogin ? 'create an account' : 'log in to your account' }}
          </button>
        </p>
      </div>

      <!-- Dynamic Component for Login/Signup Form -->
      <component :is="currentFormComponent" />

      <!-- Social Logins -->
      <SocialLogins />

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import LoginForm from './components/LoginForm.vue';
import SignupForm from './components/SignupForm.vue';
import SocialLogins from './components/SocialLogins.vue';

const isLogin = ref(true);

const toggleForm = () => {
  isLogin.value = !isLogin.value;
};

const currentFormComponent = computed(() => {
  return isLogin.value ? LoginForm : SignupForm;
});

// TODO: Add Axios logic for form submissions
</script>

<style scoped>
/* Add any specific styles for Auth.vue if needed */
</style>