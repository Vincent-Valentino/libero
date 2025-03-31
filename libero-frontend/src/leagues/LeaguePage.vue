<template>
  <div class="bg-gray-100 p-4 min-h-screen">
    <h1 class="text-3xl font-bold mb-6 md:mb-10 flex items-center" :style="{ color: leagueData?.themeColor || '#cccccc', borderBottom: `3px solid ${leagueData?.themeColor || '#cccccc'}`, paddingBottom: '0.5rem' }">
      <img v-if="leagueData?.logo" :src="leagueData.logo" :alt="leagueData.name" class="h-8 mr-3">
      {{ leagueData?.name || 'League' }}
    </h1>

    <!-- Leagues Navigation -->
    <div class="flex overflow-x-auto p-4 mb-4 gap-4 mt-15 bg-white w-full">
      <router-link 
        v-for="(league, id) in allLeaguesData" 
        :key="id" 
        :to="`/${id}`" 
        class="flex flex-col items-center min-w-[80px] cursor-pointer transition-transform hover:scale-105"
        :class="{ 'opacity-100 font-bold': leagueId === id, 'opacity-70': leagueId !== id }"
      >
        <img :src="league.logo" :alt="league.name" class="h-10 w-10 object-contain mb-1">
        <span class="text-xs text-center" :style="{ color: leagueId === id ? league.themeColor : 'inherit' }">{{ league.name }}</span>
      </router-link>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center text-gray-500 mt-20 md:mt-40">
      Loading league data...
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center text-red-500 mt-20 md:mt-40 p-4 bg-red-100 rounded border border-red-300">
      Error loading league data: {{ error }}
    </div>

    <!-- Content Area: Only show if NOT loading, NO error, and data EXISTS -->
    <div v-else-if="leagueData" class=" mb-10 flex flex-col md:grid md:grid-cols-3 md:grid-rows-3 md:gap-6">
      <!-- Mobile view - Order components differently -->
      <div class="block md:hidden space-y-5">
        <!-- League Table - Most important content first on mobile -->
        <div class="bg-white shadow rounded overflow-hidden min-h-[400px]">
          <LeagueTable
            :tableData="leagueData.table"
            :themeColor="leagueData.themeColor"
          />
        </div>
        
        <!-- Upcoming Matches - Second most important -->
        <div class="bg-white shadow rounded overflow-hidden">
          <UpcomingMatches
            :matches="leagueData.upcomingMatches"
            :themeColor="leagueData.themeColor"
          />
        </div>
  
        <!-- Player Showdown -->
        <div class="bg-white shadow rounded overflow-hidden">
          <PlayerShowdown
            :topScorers="leagueData.topScorers"
            :topAssists="leagueData.topAssists"
            :mostCleanSheets="leagueData.mostCleanSheets"
            :themeColor="leagueData.themeColor"
          />
        </div>
  
        <!-- League History -->
        <div class="bg-white shadow rounded overflow-hidden">
          <LeagueHistory
            :history="leagueData.history"
            :themeColor="leagueData.themeColor"
          />
        </div>
      </div>

      <!-- Desktop view - Grid layout -->
      <div class="hidden md:block md:row-span-3 h-full">
        <div class="bg-white shadow rounded overflow-hidden h-full">
          <LeagueTable
            :tableData="leagueData.table"
            :themeColor="leagueData.themeColor"
          />
        </div>
      </div>
      
      <div class="hidden md:block md:col-span-2 md:row-span-2">
        <div class="bg-white shadow rounded overflow-hidden h-full">
          <UpcomingMatches
            :matches="leagueData.upcomingMatches"
            :themeColor="leagueData.themeColor"
          />
        </div>
      </div>

      <div class="hidden md:block md:col-span-1">
        <div class="bg-white shadow rounded overflow-hidden h-full">
          <PlayerShowdown
            :topScorers="leagueData.topScorers"
            :topAssists="leagueData.topAssists"
            :mostCleanSheets="leagueData.mostCleanSheets"
            :themeColor="leagueData.themeColor"
          />
        </div>
      </div>

      <div class="hidden md:block md:col-span-1">
        <div class="bg-white shadow rounded overflow-hidden h-full">
          <LeagueHistory
            :history="leagueData.history"
            :themeColor="leagueData.themeColor"
          />
        </div>
      </div>
    </div>

    <!-- No Data State: Only show if NOT loading, NO error, and NO data -->
    <div v-else class="text-center text-gray-500 mt-20 md:mt-40">
       Data not available for this league.
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { allLeaguesData, type LeagueData } from './mockData';

// Import the components
import UpcomingMatches from './components/UpcomingMatches.vue';
import PlayerShowdown from './components/PlayerShowdown.vue';
import LeagueTable from './components/LeagueTable.vue';
import LeagueHistory from './components/LeagueHistory.vue';

const route = useRoute();
// Ensure leagueId updates reactively when the route changes
const leagueId = computed(() => {
    // Extract the last part of the path as the league identifier
    const pathSegments = route.path.split('/');
    return pathSegments[pathSegments.length - 1] || '';
});

const leagueData = ref<LeagueData | null>(null);

const isLoading = ref(false);
const error = ref<string | null>(null);

const loadLeagueData = (id: string) => {
  isLoading.value = true;
  error.value = null;
  leagueData.value = null; // Clear previous data

  // Simulate async fetch (replace with actual API call later)
  setTimeout(() => {
    try {
      if (id && allLeaguesData[id]) {
        leagueData.value = allLeaguesData[id];
        console.log(`Loaded data for: ${id}`);
      } else {
        console.warn(`No data found for league ID: ${id}`);
        // Keep leagueData null, the 'No Data' state will show
      }
    } catch (err) {
      console.error("Error loading league data:", err);
      error.value = err instanceof Error ? err.message : 'An unknown error occurred';
    } finally {
      isLoading.value = false;
    }
  }, 500); // Simulate 500ms network delay
};

// Load data when component mounts
onMounted(() => {
  loadLeagueData(leagueId.value);
});

// Watch for changes in leagueId (when navigating between league pages)
watch(leagueId, (newId) => {
  loadLeagueData(newId);
});

</script>