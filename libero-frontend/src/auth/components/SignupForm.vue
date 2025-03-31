<template>
  <Form class="mt-8 space-y-6" @submit="handleSignup" :validation-schema="schema">
    <div class="space-y-4">
      <div class="relative pb-4"> <!-- Added padding-bottom for error message space -->
        <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
        <Field id="name" name="name" type="text" autocomplete="name"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Your Name" />
        <ErrorMessage name="name" class="absolute text-xs text-red-600 bottom-0 left-0" /> <!-- Adjusted position -->
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
      <button type="submit"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500">
        Sign up
      </button>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';

// Define validation schema using Yup
const schema = yup.object({
  name: yup.string().required('Name is required'),
  email: yup.string().required('Email is required').email('Must be a valid email address'),
  password: yup.string().required('Password is required').min(8, 'Password must be at least 8 characters'),
  passwordConfirmation: yup.string()
    .required('Password confirmation is required')
    .oneOf([yup.ref('password')], 'Passwords must match'), // Correctly check if it matches the 'password' field
});

// Define the type for the setErrors function if needed, though often inferred
// interface SetErrorsFunction { (errors: Record<string, string>): void; }

const handleSignup = async (values: Record<string, any>, { setErrors: _setErrors }: { setErrors: (errors: Record<string, string>) => void }) => {
  // values object contains the validated form data
  console.log('Signup attempt with validated values:', values);
  // TODO: Implement Axios POST request for signup
  // try {
  //   // Example: Replace with your actual API call
  //   // const response = await axios.post('/api/signup', {
  //   //   name: values.name,
  //   //   email: values.email,
  //   //   password: values.password,
  //   // });
  //   // console.log('Signup successful:', response.data);
  //   // Handle successful signup (e.g., redirect to login or dashboard)
  // } catch (error: any) { // Catch specific error type if possible
  //   console.error('Signup failed:', error);
  //   // Handle signup error (e.g., show error message like email already exists)
  //   // Example: Check error response and set specific field errors
  //   if (error.response && error.response.data && error.response.data.errors) {
  //      setErrors(error.response.data.errors); // Assuming API returns errors in { field: message } format
  //   } else {
  //      setErrors({ name: 'An unexpected error occurred.' }); // Generic error for the 'name' field or a general form error display
  //   }
  // }
};
</script>