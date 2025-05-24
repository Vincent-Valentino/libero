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
          'py-2 px-4 text-sm font-medium focus:outline-none',
          selectedCategory === category.key
            ? 'border-b-2 font-semibold'
            : 'text-gray-600 hover:text-gray-800'
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
      <div
        v-for="(player, index) in selectedPlayers"
        :key="player.id"
        class="flex items-center justify-between p-3 bg-gray-50 rounded"
      >
        <div class="flex items-center space-x-3">
          <span class="font-semibold text-gray-500 w-5 text-right">{{ index + 1 }}.</span>
          <img
            v-if="player.photo"
            loading="lazy"
            :src="player.photo"
            :alt="player.name"
            class="h-8 w-8 rounded-full object-cover border"
          >
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
import type { PlayerStat } from '../mockData'; // Update if needed

// Define categories
const categories = [
  { key: 'topScorers' as const, label: 'Top Scorers' },
  { key: 'topAssists' as const, label: 'Top Assists' },
  { key: 'mostCleanSheets' as const, label: 'Most Clean Sheets' },
];

type CategoryKey = typeof categories[number]['key'];

// Props
const props = defineProps<{
  topScorers: PlayerStat[];
  topAssists: PlayerStat[];
  mostCleanSheets: PlayerStat[];
  themeColor: string;
}>();

// State
const selectedCategory = ref<CategoryKey>('topScorers');

// Dynamic selected players based on category
const selectedPlayers = computed<PlayerStat[]>(() => {
  switch (selectedCategory.value) {
    case 'topAssists':
      return props.topAssists || [];
    case 'mostCleanSheets':
      return props.mostCleanSheets || [];
    case 'topScorers':
    default:
      return props.topScorers || [];
  }
});
</script>

<style scoped>
button {
  transition: all 0.2s ease-in-out;
  border-bottom-width: 2px;
  border-bottom-style: solid;
  margin-bottom: -2px;
}
</style>
