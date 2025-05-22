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
    isLoadingStandings.value = true;
    standingsError.value = '';
    try {
      standings.value = await getStandings(competitionCode);
    } catch (error) {
      standingsError.value = 'Failed to load standings';
      console.error('Error fetching standings:', error);
    } finally {
      isLoadingStandings.value = false;
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