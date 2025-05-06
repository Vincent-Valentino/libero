<template>
  <div class="add-preference space-y-4">
    <!-- Add Team Section -->
    <div class="add-team-form flex items-end space-x-2">
      <div class="flex-grow">
        <label for="team-search" class="block text-sm font-medium text-gray-700 mb-1">Search Team</label>
        <input
          type="text"
          id="team-search"
          v-model="teamSearchQuery"
          placeholder="Enter team name..."
          class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
        />
        <!-- In a real app, search results would appear here -->
      </div>
      <button
        @click="handleAddTeam"
        :disabled="!teamSearchQuery"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Add Team
      </button>
    </div>

    <!-- Add Player Section -->
    <div class="add-player-form flex items-end space-x-2">
       <div class="flex-grow">
        <label for="player-search" class="block text-sm font-medium text-gray-700 mb-1">Search Player</label>
        <input
          type="text"
          id="player-search"
          v-model="playerSearchQuery"
          placeholder="Enter player name..."
          class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
        />
         <!-- In a real app, search results would appear here -->
      </div>
      <button
        @click="handleAddPlayer"
        :disabled="!playerSearchQuery"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Add Player
      </button>
    </div>

    <!-- Add Competition Section -->
    <div class="add-competition-form flex items-end space-x-2">
       <div class="flex-grow">
        <label for="competition-search" class="block text-sm font-medium text-gray-700 mb-1">Search Competition</label>
        <input
          type="text"
          id="competition-search"
          v-model="competitionSearchQuery"
          placeholder="Enter competition name..."
          class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
        />
         <!-- In a real app, search results would appear here -->
      </div>
      <button
        @click="handleAddCompetition"
        :disabled="!competitionSearchQuery"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Add Competition
      </button>
    </div>

     <!-- Basic feedback message -->
     <p v-if="feedbackMessage" class="text-sm text-green-600 mt-2">{{ feedbackMessage }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, defineEmits } from 'vue';

const teamSearchQuery = ref('');
const playerSearchQuery = ref('');
const competitionSearchQuery = ref(''); // Added competition query
const feedbackMessage = ref(''); // For simple user feedback

// Define emits
const emit = defineEmits<{
  (e: 'add-team', id: number): void;
  (e: 'add-player', id: number): void;
  (e: 'add-competition', id: number): void; // Added competition emit
}>();

const handleAddTeam = () => {
  if (!teamSearchQuery.value) return;
  // In a real app, you would get the actual team ID from search results
  const placeholderTeamId = Math.floor(Math.random() * 1000); // Placeholder
  console.log('Emitting add-team with placeholder ID:', placeholderTeamId);
  emit('add-team', placeholderTeamId);
  feedbackMessage.value = `Follow request sent for team: ${teamSearchQuery.value}`;
  teamSearchQuery.value = ''; // Clear input
  setTimeout(() => feedbackMessage.value = '', 3000); // Clear feedback after 3s
};

const handleAddPlayer = () => {
  if (!playerSearchQuery.value) return;
  // In a real app, you would get the actual player ID from search results
  const placeholderPlayerId = Math.floor(Math.random() * 10000); // Placeholder
  console.log('Emitting add-player with placeholder ID:', placeholderPlayerId);
  emit('add-player', placeholderPlayerId);
  feedbackMessage.value = `Follow request sent for player: ${playerSearchQuery.value}`;
  playerSearchQuery.value = ''; // Clear input
  setTimeout(() => feedbackMessage.value = '', 3000); // Clear feedback after 3s
};

const handleAddCompetition = () => {
  if (!competitionSearchQuery.value) return;
  // In a real app, you would get the actual competition ID from search results
  const placeholderCompetitionId = Math.floor(Math.random() * 100000); // Placeholder
  console.log('Emitting add-competition with placeholder ID:', placeholderCompetitionId);
  emit('add-competition', placeholderCompetitionId);
  feedbackMessage.value = `Follow request sent for competition: ${competitionSearchQuery.value}`;
  competitionSearchQuery.value = ''; // Clear input
  setTimeout(() => feedbackMessage.value = '', 3000); // Clear feedback after 3s
};

</script>

<!-- No scoped styles needed as Tailwind is used -->