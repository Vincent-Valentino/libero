<template>
  <Form class="mt-8 space-y-6" @submit="handleLogin" :validation-schema="schema">
    <!-- Removed hidden input, vee-validate handles values -->
    <div class="space-y-4">
      <div class="relative">
        <label for="email-address" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
        <Field id="email-address" name="email" type="email" autocomplete="email"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="you@example.com" />
        <ErrorMessage name="email" class="absolute text-[0.2rem] font-roboto-condensed text-red-600 -bottom-4 left-0" />
      </div>
      <div class="relative">
        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
        <Field id="password" name="password" type="password" autocomplete="current-password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Password" />
        <ErrorMessage name="password" class="absolute text-[0.2rem] font-roboto-condensed text-red-600 -bottom-4 left-0" />
      </div>
    </div>

    <div class="flex items-center justify-end mt-6"> <!-- Changed justify-between to justify-end -->
      <!-- Removed remember me checkbox for simplicity -->
      <div class="text-sm">
        <button @click="$emit('forgot-password')" type="button" class="font-medium text-amber-500 hover:text-amber-600">
          Forgot your password?
        </button>
      </div>
    </div>

    <div>
      <button type="submit"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500">
        Log in
      </button>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { useRouter } from 'vue-router'; // Import router
import { useAuthStore } from '@/stores/auth'; // Import auth store
import { ref } from 'vue'; // Import ref for error message

// Define validation schema using Yup
const schema = yup.object({
  email: yup.string().required('Email is required').email('Must be a valid email address'),
  password: yup.string().required('Password is required'),
});

const router = useRouter();
const authStore = useAuthStore();
const loginError = ref<string | null>(null); // To display general login errors

// Define emits
defineEmits(['forgot-password']);

const handleLogin = async (values: any) => { // Removed setErrors from signature
  loginError.value = null; // Clear previous errors
  try {
    await authStore.login({ email: values.email,
 password: values.password });
    console.log('Login successful, redirecting...');
    // Redirect to home page after successful login
    router.push({ name: 'Home' });
  } catch (error: any) {
    console.error('Login failed in component:', error);
    // Display error message to the user
    // Check if the error message from the store/API is available
    const message = error?.response?.data?.message || error?.message || 'Invalid email or password.';

    // Check for specific known error messages from the backend
    if (typeof message === 'string') {
        if (message.toLowerCase().includes('inactive')) {
            loginError.value = 'Your account is inactive. Please contact support.';
        } else if (message.toLowerCase().includes('invalid email or password')) {
            // Standard message for invalid credentials (covers user not found / wrong password)
            loginError.value = 'Invalid email or password.';
            // Optionally set field errors using vee-validate's setErrors
            // setErrors({ email: ' ', password: ' ' }); // Mark both fields as potentially wrong
        }
         else {
            // General error for other issues (e.g., server error)
            loginError.value = 'Login failed. Please try again later.';
        }
    } else {
        loginError.value = 'An unexpected error occurred during login.';
    }
  }
};
</script>