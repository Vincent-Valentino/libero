<template>
      <div class="container mx-auto p-8">
        <h1 class="text-3xl font-bold mb-6">User Profile</h1>
        <div v-if="authStore.isLoading" class="text-center">
          Loading profile...
        </div>
        <div v-else-if="authStore.user" class="bg-white p-6 rounded shadow">
          <p><strong class="font-semibold">ID:</strong> {{ authStore.user.id ?? 'N/A' }}</p>
          <p><strong class="font-semibold">Username:</strong> {{ authStore.user.username ?? 'N/A' }}</p>
          <p><strong class="font-semibold">Email:</strong> {{ authStore.user.email ?? 'N/A' }}</p>
          <p><strong class="font-semibold">Name:</strong> {{ authStore.user.name || 'N/A' }}</p>
          <p><strong class="font-semibold">Role:</strong> {{ authStore.user.role ?? 'N/A' }}</p>
          <p><strong class="font-semibold">Joined:</strong> {{ authStore.user.created_at ? new Date(authStore.user.created_at).toLocaleDateString() : 'N/A' }}</p>
        </div>
        <div v-else class="text-red-500">
          Could not load user profile. Please try logging in again.
        </div>
      </div>
    </template>

    <script setup lang="ts">
    import { useAuthStore } from '@/stores/auth';
    import { onMounted } from 'vue';

    const authStore = useAuthStore();

    // Fetch profile if not already loaded (e.g., on page refresh)
    onMounted(() => {
      if (authStore.isAuthenticated && !authStore.user) {
        authStore.fetchUserProfile(); // Corrected action name
      }
    });
    </script>

    <style scoped>
    /* Add any specific styles if needed */
    </style>