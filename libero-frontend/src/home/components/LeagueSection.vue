<template>
  <div class="bg-neutral-100 rounded-2xl p-3 pb-5 h-auto" :class="{ 'col-span-2': fullWidth }">
    <div class="flex items-center mb-1.5">
      <img :src="leagueLogoPath" :alt="leagueName" class="w-4 h-4 mr-1.5" />
      <h2 class="font-bold text-xs">{{ leagueName }}</h2>
    </div>
    
    <div :class="{ 'grid grid-cols-1 md:grid-cols-2 gap-2': fullWidth }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  leagueName: {
    type: String,
    required: true
  },
  leagueLogoPath: {
    type: String,
    required: true
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
});
</script>