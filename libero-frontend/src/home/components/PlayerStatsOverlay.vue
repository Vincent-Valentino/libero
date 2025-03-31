<template>
  <div class="absolute inset-0 bg-gradient-to-b from-amber-400 to-white text-black p-2 m-0 sm:p-3 md:p-4 opacity-0 group-hover:opacity-95 transition-opacity duration-300 overflow-y-auto flex flex-col z-10">
    <h3 class="font-roboto-condensed font-bold text-base sm:text-lg mb-1">{{ playerName }}</h3>
    
    <div class="text-2xs sm:text-sm font-roboto-condensed mb-2 sm:mb-3">
      <p>{{ age }} years • {{ position }}</p>
      <p>{{ team }}</p>
      <p>{{ nationality }}</p>
    </div>
    
    <div class="grid grid-cols-1 gap-x-1 sm:gap-x-2 text-2xs sm:text-xs font-roboto-condensed mt-2">
      <div class="flex justify-between">
        <span>Apps:</span>
        <span class="font-bold">{{ stats.appearances }}</span>
      </div>
      <div class="flex justify-between">
        <span>Goals:</span>
        <span class="font-bold">{{ stats.goals }}</span>
      </div>
      <div class="flex justify-between">
        <span>Assists:</span>
        <span class="font-bold">{{ stats.assists }}</span>
      </div>
      <div class="flex justify-between">
        <span>Key passes:</span>
        <span class="font-bold">{{ stats.keyPasses }}</span>
      </div>
      <div class="flex justify-between">
        <span>Dribbles/g:</span>
        <span class="font-bold">{{ stats.dribblesPerGame }}</span>
      </div>
      <div class="flex justify-between">
        <span>Aerial %:</span>
        <span class="font-bold">{{ stats.aerialPercentage }}%</span>
      </div>
      <div class="flex justify-between">
        <span>xG:</span>
        <span class="font-bold">{{ stats.xG }}</span>
      </div>
      <div class="flex justify-between">
        <span>xA:</span>
        <span class="font-bold">{{ stats.xA }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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
  playerName: {
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

<style scoped>
.from-gold {
  --tw-gradient-from: gold var(--tw-gradient-from-position);
  --tw-gradient-to: rgb(255 215 0 / 0) var(--tw-gradient-to-position);
  --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to);
}

/* Additional text size for very small screens */
.text-2xs {
  font-size: 0.65rem;
  line-height: 0.9rem;
}
</style>