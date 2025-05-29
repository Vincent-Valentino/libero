<template>
  <div class="profile-page container mx-auto px-4 py-8 pt-32 max-w-6xl">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Team Following Dashboard</h1>
      <p class="text-gray-600">Follow your favorite teams and get predictions for their upcoming matches</p>
    </div>

    <div v-if="profileStore.isLoading" class="loading-indicator">
      <div class="flex items-center justify-center p-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Loading profile...</span>
      </div>
    </div>

    <div v-else-if="profileStore.error" class="error-message bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex items-center">
        <svg class="w-5 h-5 text-red-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <span class="text-red-700 font-medium">Error loading profile: {{ profileStore.error }}</span>
      </div>
    </div>

    <div v-else-if="profileStore.user" class="profile-content space-y-8">
      <!-- User Info Section -->
      <UserInfoCard 
        :user="profileStore.user" 
        :followed-teams-count="profileStore.followedTeams.length"
        class="shadow-lg"
      />

      <!-- Preferences Section -->
      <div class="bg-white rounded-lg shadow-lg">
        <PreferencesManager
          :followed-teams="profileStore.followedTeams"
          @remove-team="handleRemoveTeam"
          @add-team="handleAddTeam"
        />
      </div>
    </div>

    <div v-else class="no-profile bg-gray-50 border border-gray-200 rounded-lg p-8 text-center">
      <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Unable to Load Profile</h3>
      <p class="text-gray-500 mb-4">Could not load profile data. Please try again later.</p>
      <button 
        @click="profileStore.fetchProfile()" 
        class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
      >
        Retry
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useProfileStore } from '@/stores/profile';
import UserInfoCard from './components/UserInfoCard.vue';
import PreferencesManager from './components/PreferencesManager.vue';

const profileStore = useProfileStore();

// Fetch profile data when the component mounts
onMounted(() => {
  // Only fetch if user data isn't already loaded (or if forced refresh is needed)
  if (!profileStore.user) {
    profileStore.fetchProfile();
  }
});

// --- Event Handlers for Team Preferences ---
const handleRemoveTeam = (teamId: number) => {
  console.log('Removing team:', teamId);
  profileStore.removeTeam(teamId);
};

const handleAddTeam = (teamId: number) => {
  console.log('Adding team:', teamId);
  profileStore.addTeam(teamId);
};
</script>

<!-- Style block removed as Tailwind classes are used directly -->