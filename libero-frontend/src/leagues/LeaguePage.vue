<template>
  <div class="league-page">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!--// Watch for route changes and load data
watch(
  () => route.params.code,
  async (newCode: string | string[]) => {
    try {
      const code = (typeof newCode === 'string' ? newCode : 'PL').toUpperCase();
      console.log('Loading league data for:', code);
      await store.fetchAllLeagueData(code);
    } catch (error: any) {
      console.error('Error fetching league data:', error);
      // Error states are handled by the store actions
    }
  },
  { immediate: true });eader -->
      <div class="league-header flex items-center mb-8" v-if="currentLeague">
        <img :src="currentLeague.logo" :alt="currentLeague.name" class="w-16 h-16 mr-4">
        <h1 class="text-3xl font-bold" :style="{ color: currentLeague.themeColor }">{{ currentLeague.name }}</h1>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: League Table -->
        <div class="lg:col-span-2">
          <LeagueTable 
            :tableData="standings" 
            :themeColor="currentLeague?.themeColor || '#000'"
            :isLoading="isLoadingStandings"
            :error="standingsError"
          />
        </div>

        <!-- Right Column: Stats -->
        <div class="stats-column space-y-8">
          <!-- Top Scorers -->
          <div class="bg-white rounded-lg shadow p-4">
            <h3 class="text-lg font-semibold mb-4" :style="{ color: currentLeague?.themeColor }">Top Scorers</h3>
            <div v-if="!isLoadingScorers && !scorersError && topScorers.length > 0" class="space-y-4">
              <div v-for="scorer in topScorers.slice(0, 5)" :key="scorer.id" class="flex items-center">
                <img :src="scorer.photo" :alt="scorer.name" class="w-12 h-12 rounded-full mr-4 object-cover">
                <div>
                  <div class="font-medium">{{ scorer.name }}</div>
                  <div class="text-sm text-gray-500">{{ scorer.team.name }}</div>
                </div>
                <div class="ml-auto font-bold">{{ scorer.value }} goals</div>
              </div>
            </div>
            <div v-else-if="isLoadingScorers" class="flex justify-center p-4">
              <div class="animate-spin rounded-full h-5 w-5 border-2" :style="{ borderColor: `${currentLeague?.themeColor} transparent` }"></div>
            </div>
            <div v-else-if="scorersError" class="text-red-600 text-sm p-4">
              {{ scorersError }}
            </div>
            <p v-else class="text-gray-500 text-center p-4">No scorer data available</p>
          </div>

          <!-- Upcoming Matches -->
          <div class="bg-white rounded-lg shadow p-4">
            <h3 class="text-lg font-semibold mb-4" :style="{ color: currentLeague?.themeColor }">Upcoming Matches</h3>
            <div v-if="!isLoadingFixtures && !fixturesError && fixtures.length > 0" class="space-y-4">
              <div v-for="match in fixtures.slice(0, 3)" :key="match.match_date" class="text-sm">
                <div class="flex justify-between items-center">
                  <span>{{ match.home_team_name }}</span>
                  <span class="text-gray-500">vs</span>
                  <span>{{ match.away_team_name }}</span>
                </div>
                <div class="text-gray-500 text-xs mt-1">{{ formatDate(match.match_date) }}</div>
              </div>
            </div>
            <div v-else-if="isLoadingFixtures" class="flex justify-center p-4">
              <div class="animate-spin rounded-full h-5 w-5 border-2" :style="{ borderColor: `${currentLeague?.themeColor} transparent` }"></div>
            </div>
            <div v-else-if="fixturesError" class="text-red-600 text-sm p-4">
              {{ fixturesError }}
            </div>
            <p v-else class="text-gray-500 text-center">No upcoming matches</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useLeagueStore } from '@/stores/league';
import LeagueTable from './components/LeagueTable.vue';
import { leagueMetadata } from './mockData';

const route = useRoute();
const store = useLeagueStore();

// Get current league metadata based on route params
const currentLeague = computed(() => {
  const code = (route.params.code as string)?.toUpperCase() || 'PL';
  return leagueMetadata[code];
});

// Get data from store
const standings = computed(() => store.standings);
const topScorers = computed(() => store.topScorers);
const fixtures = computed(() => store.fixtures);

// Get loading states from store
const isLoadingStandings = computed(() => store.isLoadingStandings);
const isLoadingScorers = computed(() => store.isLoadingScorers);
const isLoadingFixtures = computed(() => store.isLoadingFixtures);

// Get error states from store
const standingsError = computed(() => store.standingsError);
const scorersError = computed(() => store.scorersError);
const fixturesError = computed(() => store.fixturesError);

// Format date helper function
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  });
};

// Watch for route changes and load data
watch(
  () => route.params.code,
  async (newCode) => {
    if (newCode) {
      await store.fetchAllLeagueData(newCode as string);
    }
  },
  { immediate: true }
);
</script>