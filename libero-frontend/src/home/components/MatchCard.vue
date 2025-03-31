<template>
  <div class="bg-white rounded-lg py-4 px-2 shadow-sm hover:shadow-md transition-shadow">
    <div class="flex justify-between items-center">
      <div class="flex-1 flex justify-start">
        <div 
          v-if="homeTeam.logo" 
          class="size-10 rounded-full flex items-center justify-center overflow-hidden ml-auto mr-3">
          <img :src="homeTeam.logo" :alt="homeTeam.name" class="w-full h-full object-contain" />
        </div>
        <div 
          v-else
          class="w-7 h-7 rounded-full flex items-center justify-center" 
          :class="homeTeam.colorClass">
          <span 
            class="text-xs font-bold" 
            :class="homeTeam.textColorClass">{{ homeTeam.acronym }}</span>
        </div>
      </div>
      
      <div class="text-center flex-1">
        <div v-if="matchStarted" class="font-bold text-sm">{{ homeTeam.score }} - {{ awayTeam.score }}</div>
        <div v-else class="font-medium text-xs">{{ matchDate }}</div>
        <div class="text-xs text-gray-500">{{ matchStatus }}</div>
      </div>
      
      <div class="flex-1 flex justify-end">
        <div 
          v-if="awayTeam.logo" 
          class="size-10 rounded-full flex items-center justify-center overflow-hidden mr-auto ml-3">
          <img :src="awayTeam.logo" :alt="awayTeam.name" class="w-full h-full object-contain" />
        </div>
        <div 
          v-else
          class="w-7 h-7 rounded-full flex items-center justify-center" 
          :class="awayTeam.colorClass">
          <span 
            class="text-xs font-bold" 
            :class="awayTeam.textColorClass">{{ awayTeam.acronym }}</span>
        </div>
      </div>
    </div>
    <div class="text-xs text-gray-500 text-center mt-1 truncate">{{ stadium }}</div>
  </div>
</template>

<script setup lang="ts">
interface Team {
  name: string;
  acronym: string;
  colorClass: string;
  textColorClass?: string;
  score?: number;
  logo?: string;
}

interface MatchCardProps {
  homeTeam: Team;
  awayTeam: Team;
  matchStarted: boolean;
  matchDate?: string;
  matchStatus: string;
  stadium: string;
}

defineProps<MatchCardProps>();
</script>

<style scoped>
/* Optional hover effect for better interaction */
.bg-white {
  transition: transform 0.15s ease-in-out;
}

.bg-white:hover {
  transform: translateY(-1px);
}
</style>