<template>
  <Form class="mt-8 space-y-6" @submit="handleForgotPassword" :validation-schema="schema">
    <div class="space-y-4">
      <div class="relative">
        <label for="email-forgot" class="block text-sm font-medium text-gray-700 mb-1">Email address</label>
        <Field id="email-forgot" name="email" type="email" autocomplete="email"
               class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-amber-500 focus:border-amber-500 focus:z-10 sm:text-sm"
               placeholder="Enter your email address" />
        <ErrorMessage name="email" class="absolute text-xs text-red-600 -bottom-4 left-0" />
      </div>
    </div>

    <div>
      <button type="submit" :disabled="isLoading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-amber-500 hover:bg-amber-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500 disabled:opacity-50">
        <span v-if="isLoading">Sending...</span>
        <span v-else>Send Reset Link</span>
      </button>
    </div>

    <!-- Success/Error Messages -->
    <div v-if="successMessage" class="mt-4 p-3 bg-green-100 text-green-700 border border-green-300 rounded-md text-sm">
      {{ successMessage }}
    </div>
    
    <!-- Reset Token Display (For Testing) -->
    <div v-if="resetToken" class="mt-4 p-4 bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg">
      <div class="text-center">
        <div class="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-2">
          For Testing Only
        </div>
        <div class="text-sm text-gray-700 mb-2">Reset Token:</div>
        <div class="font-mono text-xs bg-white p-2 border rounded break-all select-all cursor-pointer" 
             @click="copyToClipboard"
             title="Click to copy">
          {{ resetToken }}
        </div>
        <div class="text-xs text-gray-500 mt-1">Click to copy</div>
      </div>
    </div>
    
    <div v-if="errorMessage" class="mt-4 p-3 bg-red-100 text-red-700 border border-red-300 rounded-md text-sm">
      {{ errorMessage }}
    </div>

    <!-- Back to login link -->
    <div class="text-center space-y-2">
      <button @click="$emit('go-to-reset')" type="button" class="text-sm text-amber-500 hover:text-amber-600">
        Already have a reset token?
      </button>
      <br>
      <button @click="$emit('back-to-login')" type="button" class="text-sm text-amber-500 hover:text-amber-600">
        Back to login
      </button>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';
import { requestPasswordReset } from '@/services/api';
import { ref } from 'vue';

// Define validation schema
const schema = yup.object({
  email: yup.string().required('Email is required').email('Must be a valid email address'),
});

// Define emits
defineEmits(['back-to-login', 'go-to-reset']);

// Reactive variables
const isLoading = ref(false);
const successMessage = ref<string | null>(null);
const errorMessage = ref<string | null>(null);
const resetToken = ref<string | null>(null);

const handleForgotPassword = async (values: any) => {
  isLoading.value = true;
  successMessage.value = null;
  errorMessage.value = null;
  resetToken.value = null;

  try {
    const response = await requestPasswordReset(values.email);
    successMessage.value = response.message;
    // For testing purposes, show the token
    if (response.token) {
      resetToken.value = response.token;
    }
  } catch (error: any) {
    console.error('Password reset request failed:', error);
    errorMessage.value = error?.response?.data?.message || error?.message || 'Failed to send reset email. Please try again.';
  } finally {
    isLoading.value = false;
  }
};

const copyToClipboard = async () => {
  if (resetToken.value) {
    try {
      await navigator.clipboard.writeText(resetToken.value);
      // Could add a toast notification here if you have one
      console.log('Reset token copied to clipboard');
    } catch (err) {
      console.error('Failed to copy text: ', err);
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = resetToken.value;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    }
  }
};
</script> 