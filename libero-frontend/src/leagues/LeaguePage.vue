<template>
  <div class="league-page">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- League Header -->
      <div class="league-header flex items-center mb-8" v-if="currentLeague">
        <img :src="currentLeague.logo" :alt="currentLeague.name" class="w-16 h-16 mr-4">
        <h1 class="text-3xl font-bold" :style="{ color: currentLeague.themeColor }">{{ currentLeague.name }}</h1>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: League Table -->
        <div class="lg:col-span-2">
          <div v-if="isLoadingStandings" class="flex justify-center items-center h-64">
            <div class="animate-spin rounded-full h-8 w-8 border-4" :style="{ borderColor: `${currentLeague?.themeColor || '#000'} transparent` }"></div>
          </div>
          <div v-else-if="standingsError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> {{ standingsError }}</span>
          </div>
          <LeagueTable 
            v-else
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
            <div v-if="!isLoadingScorers && displayTopScorers.length > 0" class="space-y-4">
              <div v-for="scorer in displayTopScorers.slice(0, 5)" :key="scorer.id" class="flex items-center">
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
            <div v-else-if="scorersError && displayTopScorers.length === 0" class="text-red-600 text-sm p-4">
              {{ scorersError }}
            </div>
            <p v-else class="text-gray-500 text-center p-4">No scorer data available</p>
          </div>

          <!-- Upcoming Matches -->
          <div class="bg-white rounded-lg shadow p-4">
            <h3 class="text-lg font-semibold mb-4" :style="{ color: currentLeague?.themeColor }">Upcoming Matches</h3>
            <div v-if="!isLoadingFixtures && displayFixtures.length > 0" class="space-y-4">
              <div v-for="match in displayFixtures.slice(0, 3)" :key="match.match_date" class="text-sm">
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
            <div v-else-if="fixturesError && displayFixtures.length === 0" class="text-red-600 text-sm p-4">
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
  const code = (typeof route.params.code === 'string' ? route.params.code : 'PL').toUpperCase();
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
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-GB', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// --- MOCK DATA FALLBACKS ---
const mockTopScorers = [
  {
    id: 1,
    name: 'Erling Haaland',
    team: { name: 'Manchester City FC' },
    value: 32,
    photo: '/Erling Haaland.png',
  },
  {
    id: 2,
    name: 'Mohamed Salah',
    team: { name: 'Liverpool FC' },
    value: 28,
    photo: '/Mohamed Salah.png',
  },
  {
    id: 3,
    name: 'Jude Bellingham',
    team: { name: 'Real Madrid' },
    value: 25,
    photo: '/Jude Bellingham.png',
  },
  {
    id: 4,
    name: 'Kylian Mbappé',
    team: { name: 'Paris Saint-Germain' },
    value: 24,
    photo: '/Kylian Mbappe.png',
  },
  {
    id: 5,
    name: 'Vinícius Júnior',
    team: { name: 'Real Madrid' },
    value: 22,
    photo: '/Vinicius Junior.png',
  },
];

const mockUpcomingMatches = [
  {
    match_date: '2025-05-25T16:00:00Z',
    home_team_name: 'Liverpool FC',
    away_team_name: 'Arsenal FC',
  },
  {
    match_date: '2025-05-26T18:30:00Z',
    home_team_name: 'Manchester City FC',
    away_team_name: 'Chelsea FC',
  },
  {
    match_date: '2025-05-27T20:00:00Z',
    home_team_name: 'Newcastle United FC',
    away_team_name: 'Aston Villa FC',
  },
];

// Use mock data if API fails
const displayTopScorers = computed(() => {
  return scorersError.value ? mockTopScorers : topScorers.value;
});
const displayFixtures = computed(() => {
  return fixturesError.value ? mockUpcomingMatches : fixtures.value;
});

// Watch for route changes and load data
watch(
  () => route.params.code,
  async (newCode: string | string[]) => {
    try {
      const code = (typeof newCode === 'string' ? newCode : 'PL').toUpperCase();
      await store.fetchAllLeagueData(code);
    } catch (error: any) {
      console.error('Error fetching league data:', error);
    }
  },
  { immediate: true }
);
</script>