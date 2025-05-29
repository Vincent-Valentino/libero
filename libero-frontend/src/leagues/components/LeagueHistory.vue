<!-- src/leagues/components/LeagueHistory.vue -->
<template>
  <div class="league-history-section p-4">
    <h2 class="text-xl font-semibold mb-3" :style="{ color: themeColor }">League History</h2>
    <div v-if="!history || history.length === 0" class="text-gray-500 italic p-2">
      No history data available.
    </div>
    <div v-else class="space-y-4">
      <div v-for="item in history" :key="item.season" class="p-4 bg-gray-50 rounded shadow-sm">
        <h3 class="font-semibold text-lg mb-3">{{ item.season }}</h3>
        <div class="space-y-2 text-sm">
          <div v-if="item.winner" class="flex items-center space-x-2">
            <span class="font-medium w-24 inline-block">Winner:</span>
            <img :src="item.winner.logo" :alt="item.winner.name" class="h-5 w-5 object-contain inline-block">
            <span>{{ item.winner.name }}</span>
          </div>
          <div v-if="item.topScorer" class="flex items-center space-x-2">
            <span class="font-medium w-24 inline-block">Top Scorer:</span>
            <span>{{ item.topScorer.name }} ({{ item.topScorer.value }} goals)</span>
          </div>
           <div v-if="item.bestPlayer" class="flex items-center space-x-2">
            <span class="font-medium w-24 inline-block">Best Player:</span>
            <img :src="item.bestPlayer.photo" :alt="item.bestPlayer.name" class="h-5 w-5 rounded-full object-cover inline-block">
            <span>{{ item.bestPlayer.name }}</span>
          </div>
          <!-- Add more history details like Team of the Season if available -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import type { LeagueHistoryItem } from "../mockData"; // Ensure the path to mockData is correct relative to this component

defineProps<{
  history: LeagueHistoryItem[];
  themeColor: string;
}>();
</script>

<style scoped>
/* Add minor styling adjustments if needed */
</style>