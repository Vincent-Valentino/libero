<template>
  <Form class="mt-8 space-y-6" @submit="handleLogin" :validation-schema="schema">
    <!-- Removed hidden input, vee-validate handles values -->
    <div class="space-y-4">
      <div class="relative">
        <label for="email-address" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
        <Field id="email-address" name="email" type="email" autocomplete="email"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="you@example.com" />
        <ErrorMessage name="email" class="absolute text-xs text-red-600 -bottom-4 left-0" />
      </div>
      <div class="relative">
        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
        <Field id="password" name="password" type="password" autocomplete="current-password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Password" />
        <ErrorMessage name="password" class="absolute text-xs text-red-600 -bottom-4 left-0" />
      </div>
    </div>

    <div class="flex items-center justify-end mt-6"> <!-- Changed justify-between to justify-end -->
      <!-- Removed remember me checkbox for simplicity -->
      <div class="text-sm">
        <a href="#" class="font-medium text-amber-500 hover:text-amber-600">
          Forgot your password?
        </a>
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

// Define validation schema using Yup
const schema = yup.object({
  email: yup.string().required('Email is required').email('Must be a valid email address'),
  password: yup.string().required('Password is required'),
});

const handleLogin = async (values: any, { setErrors }: any) => {
  console.log('Login attempt with validated values:', values);
  // TODO: Implement Axios POST request for login
  // try {
  //   const response = await axios.post('/api/login', { // Replace with your actual API endpoint
  //     email: values.email,
  //     password: values.password.,
  //   });
  //   console.log('Login successful:', response.data);
  //   // Handle successful login (e.g., store token, redirect)
  // } catch (error) {
  //   console.error('Login failed:', error);
  //   // Handle login error (e.g., show error message)
  //   // Example: setErrors({ email: 'Invalid credentials', password: ' ' }); // Set form-level or field-level errors
  // }
};
</script>