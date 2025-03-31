<template>
  <div class="bg-white rounded-sm shadow-md flex flex-col items-center justify-center gap-2 sm:gap-3 group relative cursor-pointer transition-all hover:-translate-y-1">
    <PlayerStatsOverlay 
        :player-name="name"
        :age="age"
        :position="position"
        :team="team"
        :nationality="nationality"
        :stats="stats"
      />
    <div class="relative w-full">
      <img :src="imagePath" :alt="name" class="w-full h-auto object-contain rounded-xl p-2 sm:p-3 md:p-5" />
    </div>
    <p class="text-center font-roboto-condensed font-black text-sm sm:text-base truncate w-full p-2 sm:p-3 md:p-5 mt-auto">{{ name }}</p>
  </div>
</template>

<script setup lang="ts">
import PlayerStatsOverlay from './PlayerStatsOverlay.vue';

interface PlayerStats {
  appearances: number;
  goals: number;
  assists: number;
  keyPasses: number;
  dribblesPerGame: number;
  aerialPercentage: number;
  xG: number;
  xA: number;
}

defineProps({
  name: {
    type: String,
    required: true
  },
  imagePath: {
    type: String,
    required: true
  },
  age: {
    type: Number,
    required: true
  },
  position: {
    type: String,
    required: true
  },
  team: {
    type: String,
    required: true
  },
  nationality: {
    type: String,
    required: true
  },
  stats: {
    type: Object as () => PlayerStats,
    required: true
  }
});
</script>