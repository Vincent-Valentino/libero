<template>
  <div class="p-10 md:p-20 w-full flex flex-col">
    <div class="flex flex-col gap-2">
      <h1 class="text-4xl font-title mb-4">Upcoming Matches</h1>
      
      <!-- League Tabs -->
      <LeagueTabs @tab-change="handleTabChange" />
      <!-- Loading / Error / Empty -->
      <div v-if="isLoading" class="text-center py-10">Loading fixtures…</div>
      <div v-else-if="error" class="text-center py-10 text-red-500">Error: {{ error }}</div>
      <div v-else-if="displayLeagues.length === 0" class="text-center py-10 text-gray-500">No matches scheduled</div>
      <!-- Dynamic Sections: All or specific league buckets -->
      <div v-else-if="displayLeagues.length > 0" class="flex justify-start overflow-x-auto">
        <div class="flex gap-2 max-w-full">
          <div
            v-for="section in displayLeagues"
            :key="section.competition_name"
            class="w-56 flex-shrink-0 self-start"
          >
            <LeagueSection 
              :leagueName="section.competition_name"
              :leagueLogoPath="section.logo_url"
            >
              <!-- Empty state per section -->
              <div v-if="section.matches.length === 0" class="text-gray-500 italic p-4">
                No matches scheduled for {{ section.competition_name }}
              </div>
              <!-- Match list -->
              <div v-else class="space-y-2">
                <MatchCard
                  v-for="(match, idx) in section.matches"
                  :key="idx"
                  :class="{ 'mt-2': idx > 0 }"
                  :homeTeam="{ name: match.home_team_name, logo: isAllTab ? (match.home_logo_url || '') : '', score: match.home_score, acronym: '', colorClass: '' }"
                  :awayTeam="{ name: match.away_team_name, logo: isAllTab ? (match.away_logo_url || '') : '', score: match.away_score, acronym: '', colorClass: '' }"
                  :matchStarted="match.home_score != null && match.away_score != null"
                  :matchDate="formatTime(match.match_date)"
                  :matchStatus="''"
                  :stadium="match.venue || ''"
                />
              </div>
            </LeagueSection>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import LeagueTabs from './LeagueTabs.vue';
import LeagueSection from './LeagueSection.vue';
import MatchCard from './MatchCard.vue';
import { getTodaysFixtures, getFixturesSummary } from '@/services/api';
import type { FixturesSummaryDTO } from '@/services/api';

// Track the selected league tab
const selectedTab = ref<string>('all');

// Track if 'All Leagues' tab is active
const isAllTab = computed(() => selectedTab.value === 'all');

// Summary data for a specific league
const summary = ref<FixturesSummaryDTO | null>(null);
const allLeagues = ref<any[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

// Time formatting helper for WIB
function formatTime(isoString: string): string {
  const date = new Date(isoString);
  return date.toLocaleTimeString('en-GB', {
    hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'Asia/Jakarta'
  }) + ' WIB';
}

// Handle tab change from LeagueTabs component
function handleTabChange(tabId: string) {
  selectedTab.value = tabId;
}

// League codes to prefetch
const leagueCodes = ['pl','pd','sa','bl1','fl1','cl','el'];
// Map of competition code to summary
const summariesMap = ref<Record<string, typeof summary.value>>({});
// Fetch today's fixtures and all summaries on mount
onMounted(async () => {
  isLoading.value = true;
  try {
    allLeagues.value = await getTodaysFixtures();
    const promises = leagueCodes.map(code =>
      getFixturesSummary(code)
        .then(data => ({ code, data }))
        .catch(() => ({ code, data: null }))
    );
    const results = await Promise.all(promises);
    results.forEach(({ code, data }) => {
      summariesMap.value[code] = data;
    });
  } catch (e: any) {
    error.value = e.message;
  } finally {
    isLoading.value = false;
  }
});

// Update summary when tab changes, using prefetched data
watch(selectedTab, (code) => {
  if (code !== 'all') {
    summary.value = summariesMap.value[code] || null;
  }
});

// --- MOCK DATA FALLBACKS ---
const mockAllLeagues = [
  {
    competition_name: 'Premier League',
    logo_url: '/Premier League.svg',
    matches: [
      {
        match_date: '2025-05-25T16:00:00Z',
        home_team_name: 'Liverpool FC',
        away_team_name: 'Arsenal FC',
        home_logo_url: '/Liverpool FC.png',
        away_logo_url: '/Arsenal FC.png',
        home_score: null,
        away_score: null,
        venue: 'Anfield',
      },
      {
        match_date: '2025-05-26T18:30:00Z',
        home_team_name: 'Manchester City FC',
        away_team_name: 'Chelsea FC',
        home_logo_url: '/Manchester City FC.png',
        away_logo_url: '/Chelsea FC.png',
        home_score: null,
        away_score: null,
        venue: 'Etihad Stadium',
      },
    ],
  },
  {
    competition_name: 'La Liga',
    logo_url: '/LaLiga.svg',
    matches: [
      {
        match_date: '2025-05-27T20:00:00Z',
        home_team_name: 'Real Madrid',
        away_team_name: 'Barcelona',
        home_logo_url: '/Real Madrid.png',
        away_logo_url: '/Barcelona.png',
        home_score: null,
        away_score: null,
        venue: 'Santiago Bernabéu',
      },
    ],
  },
];

// Compute sections for display
const displayLeagues = computed(() => {
  // Use mock data if error or no data
  if (error.value || allLeagues.value.length === 0) {
    return mockAllLeagues;
  }
  if (selectedTab.value === 'all') {
    // Show each league's today's fixtures
    return allLeagues.value
      .filter((l: any) => l.matches && l.matches.length > 0)
      .map((l: any) => ({
        competition_name: l.competition_name,
        logo_url: l.logo_url,
        matches: (l.matches as any[]).slice(0, 4),
      }));
  }
  // Specific league tab: show Today, Tomorrow, Upcoming without team logos
  const summ = summariesMap.value[selectedTab.value];
  if (!summ) {
    return [];
  }
  const { today, tomorrow, upcoming } = summ;
  return [
    { competition_name: 'Today', logo_url: '', matches: today.slice(0, 4) },
    { competition_name: 'Tomorrow', logo_url: '', matches: tomorrow.slice(0, 4) },
    { competition_name: 'Upcoming', logo_url: '', matches: upcoming.slice(0, 4) },
  ];
});
</script>

<style scoped>
/* Add overflow handling for the containers */
.overflow-x-auto {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}

.overflow-x-auto::-webkit-scrollbar {
  display: none;  /* Chrome, Safari and Opera */
}
</style>