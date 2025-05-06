<template>
  <div class="profile-page container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">User Profile</h1>

    <div v-if="profileStore.isLoading" class="loading-indicator">
      Loading profile...
    </div>

    <div v-else-if="profileStore.error" class="error-message text-red-600">
      Error loading profile: {{ profileStore.error }}
    </div>

    <div v-else-if="profileStore.user" class="profile-content">
      <!-- User Info Section -->
      <UserInfoCard :user="profileStore.user" class="mb-6" />

      <!-- Preferences Section -->
      <PreferencesManager
        :followed-teams="profileStore.followedTeams"
        :followed-players="profileStore.followedPlayers"
        :followed-competitions="profileStore.followedCompetitions"
        @remove-team="handleRemoveTeam"
        @remove-player="handleRemovePlayer"
        @remove-competition="handleRemoveCompetition"
        @add-team="handleAddTeam"
        @add-player="handleAddPlayer"
        @add-competition="handleAddCompetition" 
      />

    </div>

    <div v-else class="no-profile text-gray-500">
      Could not load profile data. Please try again later.
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useProfileStore } from '@/stores/profile';
import UserInfoCard from './components/UserInfoCard.vue'; // Placeholder import
import PreferencesManager from './components/PreferencesManager.vue'; // Placeholder import

const profileStore = useProfileStore();

// Fetch profile data when the component mounts
onMounted(() => {
  // Only fetch if user data isn't already loaded (or if forced refresh is needed)
  if (!profileStore.user) {
    profileStore.fetchProfile();
  }
});

// --- Event Handlers for Preferences ---
// These will call the store actions

const handleRemoveTeam = (teamId: number) => {
  console.log('Removing team:', teamId);
  profileStore.removeTeam(teamId);
};

const handleRemovePlayer = (playerId: number) => {
  console.log('Removing player:', playerId);
  profileStore.removePlayer(playerId);
};

const handleAddTeam = (teamId: number) => {
  // In a real scenario, this might come from AddPreference component
  console.log('Adding team:', teamId);
  profileStore.addTeam(teamId);
};

const handleAddPlayer = (playerId: number) => {
  // In a real scenario, this might come from AddPreference component
  console.log('Adding player:', playerId);
  profileStore.addPlayer(playerId);
};

const handleAddCompetition = (competitionId: number) => {
  console.log('Adding competition:', competitionId);
  profileStore.addCompetition(competitionId);
};

const handleRemoveCompetition = (competitionId: number) => {
  console.log('Removing competition:', competitionId);
  profileStore.removeCompetition(competitionId);
};

</script>

<!-- Style block removed as Tailwind classes are used directly -->