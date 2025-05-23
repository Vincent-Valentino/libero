import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { LeagueTableRow, PlayerStat, FixtureMatchDTO } from '@/services/api';
import { getStandings, getTopScorers, getFixturesSummary } from '@/services/api';

export const useLeagueStore = defineStore('league', () => {
  // State
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
    console.log(`Fetching standings for ${competitionCode}...`);
    isLoadingStandings.value = true;
    standingsError.value = '';
    try {
      console.log('Making API call to get standings...');
      const data = await getStandings(competitionCode);
      console.log('Received standings data:', data);
      standings.value = data;
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
    isLoadingScorers.value = true;
    scorersError.value = '';
    try {
      topScorers.value = await getTopScorers(competitionCode);
    } catch (error) {
      scorersError.value = 'Failed to load top scorers';
      console.error('Error fetching top scorers:', error);
    } finally {
      isLoadingScorers.value = false;
    }
  };

  const fetchFixtures = async (competitionCode: string) => {
    isLoadingFixtures.value = true;
    fixturesError.value = '';
    try {
      const summary = await getFixturesSummary(competitionCode);
      fixtures.value = [...summary.today, ...summary.tomorrow, ...summary.upcoming];
    } catch (error) {
      fixturesError.value = 'Failed to load fixtures';
      console.error('Error fetching fixtures:', error);
    } finally {
      isLoadingFixtures.value = false;
    }
  };

  const fetchAllLeagueData = async (competitionCode: string) => {
    await Promise.all([
      fetchStandings(competitionCode),
      fetchTopScorers(competitionCode),
      fetchFixtures(competitionCode)
    ]);
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
    fetchAllLeagueData
  };
});