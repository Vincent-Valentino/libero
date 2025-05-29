<template>
  <div class="preferences-manager p-6">
    <h2 class="text-2xl font-bold text-gray-900 mb-8">Team Following</h2>

    <!-- Section to Add New Teams -->
    <div class="add-team-section mb-8">
      <h3 class="text-xl font-semibold mb-4 text-gray-800">Follow New Teams</h3>
      <AddPreference @add-team="emitAddTeam" :followed-teams="followedTeams" />
    </div>

    <!-- Section to Display Followed Teams -->
    <div class="followed-teams-section mb-8">
      <h3 class="text-xl font-semibold mb-4 text-gray-800">Your Teams</h3>
      <div v-if="followedTeams && followedTeams.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="team in followedTeams" 
          :key="team.id"
          class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between">
            <div>
              <h4 class="font-semibold text-gray-900">{{ team.name }}</h4>
              <p class="text-sm text-gray-600">{{ getTeamLeague(team.name) }}</p>
            </div>
            <button 
              @click="emitRemoveTeam(team.id)"
              class="text-red-500 hover:text-red-700 p-2 rounded-lg hover:bg-red-50 transition-colors"
              title="Unfollow team"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
      <div v-else class="bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
        <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clip-rule="evenodd" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No Teams Followed</h3>
        <p class="text-gray-600">Start following your favorite teams to see their upcoming matches and predictions!</p>
      </div>
    </div>

    <!-- Upcoming Matches Section -->
    <div v-if="followedTeams && followedTeams.length > 0" class="upcoming-matches-section">
      <h3 class="text-xl font-semibold mb-6 text-gray-800">Upcoming Matches</h3>
      
      <!-- Loading State -->
      <div v-if="loadingMatches" class="text-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
        <p class="mt-3 text-gray-600">Loading upcoming matches...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="matchesError" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <p class="text-red-800">{{ matchesError }}</p>
        <button @click="fetchUpcomingMatches" class="mt-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors">
          Retry
        </button>
      </div>

      <!-- Matches Display -->
      <div v-else-if="upcomingMatches.length > 0" class="space-y-4">
        <div 
          v-for="match in upcomingMatches" 
          :key="`${match.home_team_name}-${match.away_team_name}-${match.match_date}`"
          class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center space-x-4">
              <div class="text-center">
                <h4 class="font-semibold text-gray-900">{{ match.home_team_name }}</h4>
                <p class="text-sm text-gray-600">Home</p>
              </div>
              <div class="text-2xl font-bold text-gray-400">VS</div>
              <div class="text-center">
                <h4 class="font-semibold text-gray-900">{{ match.away_team_name }}</h4>
                <p class="text-sm text-gray-600">Away</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-medium text-gray-900">{{ formatMatchDate(match.match_date) }}</p>
              <p class="text-sm text-gray-600">{{ formatMatchTime(match.match_date) }}</p>
            </div>
          </div>

          <!-- Prediction Section -->
          <div v-if="isTopLeagueMatch(match)" class="border-t pt-4">
            <div v-if="match.prediction" class="bg-gradient-to-r from-indigo-50 to-blue-50 rounded-lg p-4">
              <h5 class="font-semibold text-gray-900 mb-3">Match Prediction</h5>
              <div class="grid grid-cols-3 gap-4 text-center">
                <div class="text-center">
                  <p class="text-2xl font-bold text-indigo-600">{{ match.prediction.most_likely_home_score }}</p>
                  <p class="text-sm text-gray-600">{{ match.home_team_name }}</p>
                </div>
                <div class="text-center">
                  <p class="text-lg font-semibold text-gray-400">-</p>
                  <p class="text-sm text-gray-600">{{ getPredictionResult(match.prediction, match.home_team_name, match.away_team_name) }}</p>
                </div>
                <div class="text-center">
                  <p class="text-2xl font-bold text-indigo-600">{{ match.prediction.most_likely_away_score }}</p>
                  <p class="text-sm text-gray-600">{{ match.away_team_name }}</p>
                </div>
              </div>
              <div class="mt-3 grid grid-cols-3 gap-2 text-xs">
                <div class="bg-white rounded p-2 text-center">
                  <p class="font-semibold">Home Win</p>
                  <p class="text-green-600 font-bold">{{ (match.prediction.probabilities.home_win * 100).toFixed(0) }}%</p>
                </div>
                <div class="bg-white rounded p-2 text-center">
                  <p class="font-semibold">Draw</p>
                  <p class="text-gray-600 font-bold">{{ (match.prediction.probabilities.draw * 100).toFixed(0) }}%</p>
                </div>
                <div class="bg-white rounded p-2 text-center">
                  <p class="font-semibold">Away Win</p>
                  <p class="text-red-600 font-bold">{{ (match.prediction.probabilities.away_win * 100).toFixed(0) }}%</p>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-2">
              <button 
                @click="generatePrediction(match)"
                :disabled="loadingPredictions"
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 transition-colors"
              >
                Generate Prediction
              </button>
            </div>
          </div>
          <div v-else class="border-t pt-4">
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-center">
              <svg class="w-6 h-6 text-yellow-600 mx-auto mb-2" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              <p class="text-yellow-800 font-medium">Predictions not available</p>
              <p class="text-yellow-700 text-sm">Currently, predictions are only supported for Premier League, La Liga, Serie A, Bundesliga, and Ligue 1.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- No Matches State -->
      <div v-else class="bg-gray-50 border border-gray-200 rounded-lg p-8 text-center">
        <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
        </svg>
        <h4 class="text-lg font-medium text-gray-900 mb-2">No Upcoming Matches</h4>
        <p class="text-gray-600">No upcoming matches found for your followed teams in the next 7 days.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, watch, onMounted } from 'vue';
import { getTodaysFixtures, getUpcomingMatches, predictMatch, type FixtureMatchDTO, type PredictMatchResponse } from '@/services/api';
import AddPreference from './AddPreference.vue';

// Define interfaces
interface Team {
  id: number;
  name: string;
}

interface MatchWithPrediction extends FixtureMatchDTO {
  prediction?: PredictMatchResponse;
}

// Define props - only teams now
const props = defineProps<{
  followedTeams: Team[];
}>();

// Define emits - only team-related events
const emit = defineEmits<{
  (e: 'remove-team', id: number): void;
  (e: 'add-team', id: number): void;
}>();

// Reactive data
const upcomingMatches = ref<MatchWithPrediction[]>([]);
const loadingMatches = ref(false);
const matchesError = ref('');
const loadingPredictions = ref(false);

// Top 5 leagues mapping
const topLeagues = {
  'Premier League': 'E0',
  'La Liga': 'SP1',
  'Serie A': 'I1',
  'Bundesliga': 'D1',
  'Ligue 1': 'F1'
};

// Team to league mapping (simplified version)
const teamLeagueMap: Record<string, string> = {
  // Premier League
  'Arsenal': 'Premier League', 'Chelsea': 'Premier League', 'Liverpool': 'Premier League',
  'Manchester City': 'Premier League', 'Manchester United': 'Premier League', 'Tottenham': 'Premier League',
  'Newcastle': 'Premier League', 'Brighton': 'Premier League', 'Aston Villa': 'Premier League',
  'West Ham': 'Premier League', 'Crystal Palace': 'Premier League', 'Everton': 'Premier League',
  
  // La Liga
  'Barcelona': 'La Liga', 'Real Madrid': 'La Liga', 'Atletico Madrid': 'La Liga',
  'Athletic Bilbao': 'La Liga', 'Real Sociedad': 'La Liga', 'Villarreal': 'La Liga',
  'Valencia': 'La Liga', 'Sevilla': 'La Liga', 'Real Betis': 'La Liga',
  
  // Serie A
  'Juventus': 'Serie A', 'AC Milan': 'Serie A', 'Inter Milan': 'Serie A',
  'AS Roma': 'Serie A', 'Napoli': 'Serie A', 'Lazio': 'Serie A',
  'Atalanta': 'Serie A', 'Fiorentina': 'Serie A',
  
  // Bundesliga
  'Bayern Munich': 'Bundesliga', 'Borussia Dortmund': 'Bundesliga', 'RB Leipzig': 'Bundesliga',
  'Bayer Leverkusen': 'Bundesliga', 'Eintracht Frankfurt': 'Bundesliga', 'Wolfsburg': 'Bundesliga',
  
  // Ligue 1
  'Paris Saint-Germain': 'Ligue 1', 'Olympique Marseille': 'Ligue 1', 'AS Monaco': 'Ligue 1',
  'Lille': 'Ligue 1', 'Lyon': 'Ligue 1', 'Nice': 'Ligue 1'
};

// Methods
const getTeamLeague = (teamName: string): string => {
  return teamLeagueMap[teamName] || 'Unknown League';
};

const isTopLeagueMatch = (match: FixtureMatchDTO): boolean => {
  const homeLeague = getTeamLeague(match.home_team_name);
  const awayLeague = getTeamLeague(match.away_team_name);
  return Object.keys(topLeagues).includes(homeLeague) && Object.keys(topLeagues).includes(awayLeague);
};

const formatMatchDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    weekday: 'short', 
    month: 'short', 
    day: 'numeric' 
  });
};

const formatMatchTime = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', { 
    hour: '2-digit', 
    minute: '2-digit' 
  });
};

const getPredictionResult = (prediction: PredictMatchResponse, homeTeam: string, awayTeam: string): string => {
  const homeScore = prediction.most_likely_home_score;
  const awayScore = prediction.most_likely_away_score;
  
  if (homeScore === awayScore) {
    return 'Draw';
  } else if (homeScore > awayScore) {
    return `${homeTeam} Win`;
  } else {
    return `${awayTeam} Win`;
  }
};

const fetchUpcomingMatches = async () => {
  if (!props.followedTeams || props.followedTeams.length === 0) {
    upcomingMatches.value = [];
    return;
  }

  loadingMatches.value = true;
  matchesError.value = '';
  
  try {
    const fixtures = await getUpcomingMatches();
    const followedTeamNames = props.followedTeams.map(team => team.name);
    
    // Filter matches that include followed teams and are in the future
    const now = new Date();
    const relevantMatches = fixtures.flatMap(competition => 
      competition.matches.filter(match => {
        const matchDate = new Date(match.match_date);
        return matchDate > now && 
               (followedTeamNames.includes(match.home_team_name) || 
                followedTeamNames.includes(match.away_team_name));
      })
    );
    
    // Sort by date and limit to next 10 matches
    upcomingMatches.value = relevantMatches
      .sort((a, b) => new Date(a.match_date).getTime() - new Date(b.match_date).getTime())
      .slice(0, 10)
      .map(match => ({ ...match, prediction: undefined }));
      
  } catch (error: any) {
    console.error('Error fetching upcoming matches:', error);
    matchesError.value = 'Failed to load upcoming matches. Please try again.';
  } finally {
    loadingMatches.value = false;
  }
};

const generatePrediction = async (match: MatchWithPrediction) => {
  if (!isTopLeagueMatch(match)) return;
  
  loadingPredictions.value = true;
  
  try {
    const homeLeague = getTeamLeague(match.home_team_name);
    const leagueCode = topLeagues[homeLeague as keyof typeof topLeagues];
    
    const prediction = await predictMatch({
      league: leagueCode,
      home_team: match.home_team_name,
      away_team: match.away_team_name
    });
    
    // Update the specific match with prediction
    const matchIndex = upcomingMatches.value.findIndex(m => 
      m.home_team_name === match.home_team_name && 
      m.away_team_name === match.away_team_name && 
      m.match_date === match.match_date
    );
    
    if (matchIndex !== -1) {
      upcomingMatches.value[matchIndex].prediction = prediction;
    }
    
  } catch (error: any) {
    console.error('Error generating prediction:', error);
    // Could show individual match prediction error if needed
  } finally {
    loadingPredictions.value = false;
  }
};

// Emit functions
const emitRemoveTeam = (teamId: number) => {
  emit('remove-team', teamId);
};

const emitAddTeam = (teamId: number) => {
  emit('add-team', teamId);
};

// Watch for changes in followed teams
watch(() => props.followedTeams, () => {
  fetchUpcomingMatches();
}, { deep: true });

// Fetch matches on mount
onMounted(() => {
  fetchUpcomingMatches();
});
</script>

<!-- Style block removed as Tailwind classes are used directly -->