<template>
  <Form class="mt-8 space-y-6" @submit="handleSignup" :validation-schema="schema">
    <div class="space-y-4">
      <div class="relative pb-4"> <!-- Added padding-bottom for error message space -->
        <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
        <Field id="name" name="name" type="text" autocomplete="name"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Enter your full name" />
        <ErrorMessage name="name" class="absolute text-xs text-red-600 bottom-0 left-0" />
      </div>
      <div class="relative pb-4"> <!-- Added padding-bottom for error message space -->
        <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
        <Field id="username" name="username" type="text" autocomplete="username"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Choose a username" />
        <ErrorMessage name="username" class="absolute text-xs text-red-600 bottom-0 left-0" /> <!-- Adjusted position -->
      </div>
      <div class="relative pb-4"> <!-- Added padding-bottom -->
        <label for="email-address-signup" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
        <Field id="email-address-signup" name="email" type="email" autocomplete="email"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="you@example.com" />
        <ErrorMessage name="email" class="absolute text-xs text-red-600 bottom-0 left-0" /> <!-- Adjusted position -->
      </div>
      <div class="relative pb-4"> <!-- Added padding-bottom -->
        <label for="password-signup" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
        <Field id="password-signup" name="password" type="password" autocomplete="new-password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Create a password (min. 8 characters)" />
        <ErrorMessage name="password" class="absolute text-xs text-red-600 bottom-0 left-0" /> <!-- Adjusted position -->
      </div>
      <div class="relative pb-4"> <!-- Added padding-bottom -->
        <label for="password-confirm" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
        <Field id="password-confirm" name="passwordConfirmation" type="password" autocomplete="new-password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Confirm your password" />
        <ErrorMessage name="passwordConfirmation" class="absolute text-xs text-red-600 bottom-0 left-0" /> <!-- Adjusted position -->
      </div>
    </div>

    <div class="mt-6">
      <button type="submit" :disabled="authStore.isLoading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500 disabled:opacity-50">
        <span v-if="authStore.isLoading">Signing up...</span>
        <span v-else>Sign up</span>
      </button>
    </div>
  </Form>
  <!-- Display Success/Error Messages -->
  <div v-if="signupSuccess" class="mt-4 p-3 bg-green-100 text-green-700 border border-green-300 rounded-md text-sm">
    Registration successful! You can now log in.
  </div>
  <div v-if="signupError" class="mt-4 p-3 bg-red-100 text-red-700 border border-red-300 rounded-md text-sm">
    {{ signupError }}
  </div>
</template>

<script setup lang="ts">
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { useAuthStore } from '@/stores/auth'; // Import auth store
import { registerUser } from '@/services/api'; // Import API function as fallback
import { ref } from 'vue'; // Import ref for status messages

// Define validation schema using Yup
const schema = yup.object({
  name: yup.string().required('Name is required'),
  username: yup.string().required('Username is required').min(3, 'Username must be at least 3 characters'),
  email: yup.string().required('Email is required').email('Must be a valid email address'),
  password: yup.string().required('Password is required').min(8, 'Password must be at least 8 characters'),
  passwordConfirmation: yup.string()
    .required('Password confirmation is required')
    .oneOf([yup.ref('password')], 'Passwords must match'), // Correctly check if it matches the 'password' field
});

const authStore = useAuthStore();
const signupError = ref<string | null>(null);
const signupSuccess = ref<boolean>(false);

// Define the type for the setErrors function from vee-validate context
type SetErrorsFunction = (errors: Record<string, string | undefined>) => void;

const handleSignup = async (values: Record<string, any>, { setErrors }: { setErrors: SetErrorsFunction }) => {
  signupError.value = null;
  signupSuccess.value = false;
  setErrors({}); // Clear previous vee-validate errors

  console.log('Signup attempt with values:', values);
  
  // Validate that we have the required fields
  if (!values.name || !values.username || !values.email || !values.password) {
    console.error('Missing required fields');
    signupError.value = 'Please fill in all required fields';
    return;
  }
  
  const userData = {
    name: values.name,
    username: values.username,
    email: values.email,
    password: values.password,
  };
  
  try {
    if (typeof authStore.register === 'function') {
      await authStore.register(userData);
    } else {
      // Fallback to direct API call
      console.log('Using direct API call');
      authStore.loading = true;
      try {
        await registerUser(userData);
      } finally {
        authStore.loading = false;
      }
    }
    
    console.log('Signup successful');
    signupSuccess.value = true;

  } catch (error: any) {
    console.error('Signup failed:', error);
    signupSuccess.value = false;
    const message = error?.response?.data?.message || error?.response?.data || error?.message || 'Registration failed. Please try again.';

    // Try to map common errors to specific fields using setErrors
    if (typeof message === 'string') {
        if (message.toLowerCase().includes('email') && (message.toLowerCase().includes('exists') || message.toLowerCase().includes('duplicate'))) {
          setErrors({ email: 'This email address is already registered.' });
        } else if (message.toLowerCase().includes('username') && (message.toLowerCase().includes('exists') || message.toLowerCase().includes('duplicate') || message.toLowerCase().includes('taken'))) {
          setErrors({ username: 'This username is already taken.' });
        } else {
          // Display general error if specific field mapping isn't clear
          signupError.value = message;
        }
    } else {
         // Handle non-string or complex errors if necessary
         signupError.value = 'An unexpected error occurred during registration.';
    }
  }
};
</script>