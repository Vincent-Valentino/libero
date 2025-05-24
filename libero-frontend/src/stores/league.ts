import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { LeagueTableRow, PlayerStat, FixtureMatchDTO } from '@/services/api';
import { getStandings, getTopScorers, getFixturesSummary } from '@/services/api';

export const useLeagueStore = defineStore('league', () => {
  // Cache by competition code
  const standingsCache = ref<Record<string, LeagueTableRow[]>>({});
  const topScorersCache = ref<Record<string, PlayerStat[]>>({});
  const fixturesCache = ref<Record<string, FixtureMatchDTO[]>>({});
  
  // Current data
  const standings = ref<LeagueTableRow[]>([]);
  const topScorers = ref<PlayerStat[]>([]);
  const fixtures = ref<FixtureMatchDTO[]>([]);
  
  // Loading states
  const isLoadingStandings = ref(false);
  const isLoadingScorers = ref(false);
  const isLoadingFixtures = ref(false);
  
  // Error states
  const standingsError = ref<string>('');
  const scorersError = ref<string>('');
  const fixturesError = ref<string>('');

  // Actions
  const fetchStandings = async (competitionCode: string) => {
    // If we already have data for this competition, use it
    if (standingsCache.value[competitionCode] && standingsCache.value[competitionCode].length > 0) {
      console.log(`Using cached standings for ${competitionCode}`);
      standings.value = standingsCache.value[competitionCode];
      return;
    }

    console.log(`Fetching standings for ${competitionCode}...`);
    isLoadingStandings.value = true;
    standingsError.value = '';
    try {
      console.log('Making API call to get standings...');
      const data = await getStandings(competitionCode);
      console.log('Received standings data:', data);
      standings.value = data;
      // Only cache if data is non-empty
      if (Array.isArray(data) && data.length > 0) {
        standingsCache.value[competitionCode] = data;
      }
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || error.message || 'Failed to load standings';
      standingsError.value = errorMessage;
      console.error('Error fetching standings:', {
        code: competitionCode,
        error: error,
        response: error.response?.data
      });
    } finally {
      isLoadingStandings.value = false;
      console.log('Standings fetch complete. HasData:', standings.value.length > 0);
    }
  };
  const fetchTopScorers = async (competitionCode: string) => {
    // If we already have data for this competition, use it
    if (topScorersCache.value[competitionCode] && topScorersCache.value[competitionCode].length > 0) {
      console.log(`Using cached top scorers for ${competitionCode}`);
      topScorers.value = topScorersCache.value[competitionCode];
      return;
    }

    isLoadingScorers.value = true;
    scorersError.value = '';
    try {
      const data = await getTopScorers(competitionCode);
      topScorers.value = data;
      // Only cache if data is non-empty
      if (Array.isArray(data) && data.length > 0) {
        topScorersCache.value[competitionCode] = data;
      }
    } catch (error: any) {
      // Use mock data if available on network error (especially for EL)
      const errorMessage = error.response?.data?.message || error.message || 'Failed to load top scorers';
      scorersError.value = errorMessage;
      console.warn(`Error fetching top scorers for ${competitionCode}:`, error);
      // Do not cache empty/error result
    } finally {
      isLoadingScorers.value = false;
    }
  };

  const fetchFixtures = async (competitionCode: string) => {
    // If we already have data for this competition, use it
    if (fixturesCache.value[competitionCode] && fixturesCache.value[competitionCode].length > 0) {
      console.log(`Using cached fixtures for ${competitionCode}`);
      fixtures.value = fixturesCache.value[competitionCode];
      return;
    }

    isLoadingFixtures.value = true;
    fixturesError.value = '';
    try {
      const summary = await getFixturesSummary(competitionCode);
      const allFixtures = [...summary.today, ...summary.tomorrow, ...summary.upcoming]; 
      fixtures.value = allFixtures;
      // Only cache if data is non-empty
      if (Array.isArray(allFixtures) && allFixtures.length > 0) {
        fixturesCache.value[competitionCode] = allFixtures;
      }
    } catch (error: any) {
      // Handle network errors gracefully
      const errorMessage = error.response?.data?.message || error.message || 'Failed to load fixtures';
      fixturesError.value = errorMessage; 
      console.warn(`Error fetching fixtures for ${competitionCode}:`, error);
      // Do not cache empty/error result
    } finally {
      isLoadingFixtures.value = false;
    }
  };
  const fetchAllLeagueData = async (competitionCode: string, forceRefresh = false) => {
    // If forcing refresh, clear the cache for this competition
    if (forceRefresh) {
      clearCacheForCompetition(competitionCode);
    }
    
    // Load data in sequence to avoid overwhelming the backend
    try {
      await fetchStandings(competitionCode);
      await fetchTopScorers(competitionCode);
      await fetchFixtures(competitionCode);
    } catch (error) {
      console.error("Error fetching all league data:", error);
      // Individual error states are already set by the respective functions
    }
  };
  
  // New utility functions for cache management
  const clearAllCache = () => {
    standingsCache.value = {};
    topScorersCache.value = {};
    fixturesCache.value = {};
  };
  
  const clearCacheForCompetition = (competitionCode: string) => {
    delete standingsCache.value[competitionCode];
    delete topScorersCache.value[competitionCode];
    delete fixturesCache.value[competitionCode];
  };

  return {
    // State
    standings,
    topScorers,
    fixtures,
    // Loading states
    isLoadingStandings,
    isLoadingScorers,
    isLoadingFixtures,
    // Error states
    standingsError,
    scorersError,
    fixturesError,
    // Actions
    fetchStandings,
    fetchTopScorers,
    fetchFixtures,
    fetchAllLeagueData,
    // Cache management
    clearAllCache,
    clearCacheForCompetition
  };
});