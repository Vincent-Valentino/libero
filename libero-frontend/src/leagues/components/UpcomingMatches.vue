<!-- src/leagues/components/UpcomingMatches.vue -->
<template>
  <div class="upcoming-matches-section p-4">
    <h2 class="text-xl font-semibold mb-3" :style="{ color: themeColor }">Upcoming Matches</h2>
    <div class="bg-white shadow rounded p-4 flex space-x-6 overflow-x-auto pb-4 scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-gray-200">
      <div v-if="!matches || matches.length === 0" class="text-gray-500 italic px-3">
        No upcoming matches available.
      </div>
      <!-- Removed v-else from here -->
      <div v-for="(match, index) in matches" :key="match.id"
           class="min-w-[15rem]"
           :class="{ 'border-r border-gray-200 pr-6 mr-6': index < matches.length - 1 }">
        <div class="text-xs text-gray-500 mb-1">{{ formatDate(match.date) }} - {{ match.time }}</div>
        <div class="text-sm text-gray-600 mb-2">{{ match.venue }}</div>
        <div class="flex items-center justify-between mb-1">
          <div class="flex items-center space-x-2">
            <img
              :src="imgSrc(match.homeTeam.logo)"
              :alt="match.homeTeam.name"
              class="h-5 w-5 object-contain"
              @error="onImgError($event, match.homeTeam)"
            >
            <span class="font-medium text-sm">{{ match.homeTeam.name }}</span>
          </div>
          <span class="text-xs font-semibold" :style="{ color: themeColor }">vs</span>
        </div>
        <div class="flex items-center space-x-2">
           <img
             :src="imgSrc(match.awayTeam.logo)"
             :alt="match.awayTeam.name"
             class="h-5 w-5 object-contain"
             @error="onImgError($event, match.awayTeam)"
           >
           <span class="font-medium text-sm">{{ match.awayTeam.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
// Ensure the path to mockData is correct relative to this component
import type { Match } from '../mockData';
import { ref } from 'vue';

// Props are used in the template, so this declaration is correct
defineProps<{
  matches: Match[];
  themeColor: string;
}>();

// Helper function to format date
const formatDate = (dateString: string): string => {
  // Explicitly define options type and values as strings
  const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' };
  try {
    return new Date(dateString).toLocaleDateString(undefined, options);
  } catch (e) {
    console.error("Error formatting date:", dateString, e);
    return dateString; // Fallback to original string on error
  }
};

const fallbackTeamLogo = '/fallback-team.png';
const erroredTeams = ref<Record<string, boolean>>({});
function imgSrc(logo: string) {
  return erroredTeams.value[logo] ? fallbackTeamLogo : logo;
}
function onImgError(event: Event, team: any) {
  if (!erroredTeams.value[team.logo]) {
    erroredTeams.value[team.logo] = true;
    const target = event.target as HTMLImageElement | null;
    if (target) target.src = fallbackTeamLogo;
  }
}
</script>

<style scoped>
/* Custom scrollbar styling */
.scrollbar-thin {
  scrollbar-width: thin; /* Firefox */
  scrollbar-color: #9ca3af #e5e7eb; /* thumb track */
}
/* Webkit scrollbar styling */
.scrollbar-thumb-gray-400::-webkit-scrollbar {
  height: 8px; /* Height of horizontal scrollbar */
}
.scrollbar-thumb-gray-400::-webkit-scrollbar-thumb {
  background-color: #9ca3af; /* Tailwind gray-400 */
  border-radius: 4px;
}
.scrollbar-track-gray-200::-webkit-scrollbar-track {
  background-color: #e5e7eb; /* Tailwind gray-200 */
  border-radius: 4px;
}
/* Ensure container allows horizontal scrolling */
.overflow-x-auto {
  overflow-x: auto;
}
</style>