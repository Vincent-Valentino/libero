<!-- src/leagues/components/PlayerShowdown.vue -->
<template>
  <div class="player-showdown-section p-4">
    <h2 class="text-xl font-semibold mb-3" :style="{ color: themeColor }">Player Showdown</h2>

    <!-- Category Tabs/Buttons -->
    <div class="flex space-x-2 mb-4 border-b" :style="{ borderColor: themeColor }">
      <button
        v-for="category in categories"
        :key="category.key"
        @click="selectedCategory = category.key"
        :class="[
          'py-2 px-4 text-sm font-medium focus:outline-none', // Base classes
          selectedCategory === category.key
            ? 'border-b-2 font-semibold' // Active classes
            : 'text-gray-600 hover:text-gray-800', // Inactive classes
        ]"
        :style="{
          borderBottomColor: selectedCategory === category.key ? themeColor : 'transparent',
          color: selectedCategory === category.key ? themeColor : ''
        }"
      >
        {{ category.label }}
      </button>
    </div>

    <!-- Player List for Selected Category -->
    <div v-if="selectedPlayers.length > 0" class="space-y-2">
       <div v-for="(player, index) in selectedPlayers" :key="player.id" class="flex items-center justify-between p-3 bg-gray-50 rounded">
         <div class="flex items-center space-x-3">
           <span class="font-semibold text-gray-500 w-5 text-right">{{ index + 1 }}.</span>
           <img :src="player.photo" :alt="player.name" class="h-8 w-8 rounded-full object-cover border">
           <div>
             <div class="font-medium text-sm">{{ player.name }}</div>
             <div class="text-xs text-gray-500">{{ player.team.name }}</div>
           </div>
         </div>
         <div class="font-bold text-lg" :style="{ color: themeColor }">{{ player.value }}</div>
       </div>
    </div>
     <div v-else class="text-gray-500 italic p-2">
        No player data available for this category.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineProps } from 'vue';
import type { PlayerStat } from '../mockData'; // Correct path assumed

// Define categories with explicit types and 'as const' for type safety
const categories = [
  { key: 'topScorers' as const, label: 'Top Scorers' },
  { key: 'topAssists' as const, label: 'Top Assists' },
  { key: 'mostCleanSheets' as const, label: 'Most Clean Sheets' },
];

// Derive the CategoryKey type from the categories array
type CategoryKey = typeof categories[number]['key'];

// Define props
const props = defineProps<{
  topScorers: PlayerStat[];
  topAssists: PlayerStat[];
  mostCleanSheets: PlayerStat[];
  themeColor: string;
}>();

// State for the selected category
const selectedCategory = ref<CategoryKey>(categories[0].key); // Default to the key of the first category

// Computed property to get the players for the selected category
const selectedPlayers = computed<PlayerStat[]>(() => {
  switch (selectedCategory.value) {
    case 'topScorers':
      return props.topScorers || [];
    case 'topAssists':
      return props.topAssists || [];
    case 'mostCleanSheets':
      return props.mostCleanSheets || [];
    default:
      // Should not happen with strict typing, but good practice
      const exhaustiveCheck: never = selectedCategory.value;
      console.warn(`Unhandled category: ${exhaustiveCheck}`);
      return [];
  }
});

</script>

<style scoped>
/* Add any specific styles if needed */
button {
  transition: all 0.2s ease-in-out;
  border-bottom-width: 2px; /* Ensure space for border even when transparent */
  border-bottom-style: solid;
  margin-bottom: -2px; /* Counteract the border space */
}
</style>