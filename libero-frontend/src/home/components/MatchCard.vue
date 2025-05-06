<template>
  <div class="bg-white rounded-lg py-4 px-2 shadow-sm hover:shadow-md transition-shadow">
    <div class="flex items-center">
      <!-- Home zone: fixed width for alignment -->
      <div class="w-16 flex flex-col items-center">
        <div
          v-if="homeTeam.logo"
          class="w-10 h-10 rounded-full overflow-hidden"
        >
          <img :src="homeTeam.logo" :alt="homeTeam.name" class="w-full h-full object-contain" />
        </div>
        <div
          v-else
          class="w-7 h-7 rounded-full flex items-center justify-center"
          :class="homeTeam.colorClass"
        >
          <span
            class="text-xs font-bold"
            :class="homeTeam.textColorClass"
          >{{ homeTeam.acronym }}</span>
        </div>
        <div class="w-full text-[10px] text-gray-800 truncate mt-1 text-center">{{ homeTeam.name }}</div>
      </div>
      
      <!-- Center zone: flex-grow -->
      <div class="flex-1 flex flex-col items-center text-center">
        <div v-if="matchStarted" class="font-bold text-sm">{{ homeTeam.score }} - {{ awayTeam.score }}</div>
        <div v-else class="font-medium text-xs">{{ matchDate }}</div>
        <div class="text-xs text-gray-500">{{ matchStatus }}</div>
      </div>
      
      <!-- Away zone: fixed width for alignment -->
      <div class="w-16 flex flex-col items-center">
        <div 
          v-if="awayTeam.logo" 
          class="w-10 h-10 rounded-full overflow-hidden">
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
        <div class="w-full text-[10px] text-gray-800 truncate mt-1 text-center">{{ awayTeam.name }}</div>
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