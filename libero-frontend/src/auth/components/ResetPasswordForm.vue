<template>
  <Form class="mt-8 space-y-6" @submit="handleResetPassword" :validation-schema="schema">
    <div class="space-y-4">
      <div class="relative pb-4">
        <label for="token" class="block text-sm font-medium text-gray-700 mb-1">Reset Token</label>
        <Field id="token" name="token" type="text"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Enter your reset token" />
        <ErrorMessage name="token" class="absolute text-xs text-red-600 bottom-0 left-0" />
      </div>
      <div class="relative pb-4">
        <label for="new-password" class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
        <Field id="new-password" name="newPassword" type="password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Enter new password (min. 8 characters)" />
        <ErrorMessage name="newPassword" class="absolute text-xs text-red-600 bottom-0 left-0" />
      </div>
      <div class="relative pb-4">
        <label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
        <Field id="confirm-password" name="confirmPassword" type="password"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Confirm your new password" />
        <ErrorMessage name="confirmPassword" class="absolute text-xs text-red-600 bottom-0 left-0" />
      </div>
    </div>

    <div>
      <button type="submit" :disabled="isLoading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500 disabled:opacity-50">
        <span v-if="isLoading">Resetting Password...</span>
        <span v-else>Reset Password</span>
      </button>
    </div>

    <!-- Success/Error Messages -->
    <div v-if="successMessage" class="mt-4 p-3 bg-green-100 text-green-700 border border-green-300 rounded-md text-sm">
      {{ successMessage }}
    </div>
    <div v-if="errorMessage" class="mt-4 p-3 bg-red-100 text-red-700 border border-red-300 rounded-md text-sm">
      {{ errorMessage }}
    </div>

    <!-- Back to login link -->
    <div class="text-center">
      <button @click="$emit('back-to-login')" type="button" class="text-sm text-amber-500 hover:text-amber-600">
        Back to login
      </button>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { resetPassword } from '@/services/api';
import { ref } from 'vue';

// Define validation schema
const schema = yup.object({
  token: yup.string().required('Reset token is required'),
  newPassword: yup.string().required('New password is required').min(8, 'Password must be at least 8 characters'),
  confirmPassword: yup.string()
    .required('Password confirmation is required')
    .oneOf([yup.ref('newPassword')], 'Passwords must match'),
});

// Define emits
const emit = defineEmits(['back-to-login', 'password-reset-success']);

// Reactive variables
const isLoading = ref(false);
const successMessage = ref<string | null>(null);
const errorMessage = ref<string | null>(null);

const handleResetPassword = async (values: any) => {
  isLoading.value = true;
  successMessage.value = null;
  errorMessage.value = null;

  try {
    const response = await resetPassword({
      token: values.token,
      new_password: values.newPassword,
      confirm_password: values.confirmPassword,
    });
    
    successMessage.value = response.message;
    
    // Redirect to login after successful reset
    setTimeout(() => {
      emit('password-reset-success');
    }, 2000);
    
  } catch (error: any) {
    console.error('Password reset failed:', error);
    errorMessage.value = error?.response?.data?.message || error?.message || 'Failed to reset password. Please try again.';
  } finally {
    isLoading.value = false;
  }
};
</script> 