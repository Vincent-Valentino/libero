<template>
  <div class="max-w-6xl mx-auto p-6">
    <div class="text-center mb-12">
      <h2 class="text-4xl font-bold text-gray-900 mb-4">Match Predictor</h2>
      <p class="text-lg text-gray-600 max-w-2xl mx-auto">Advanced football match prediction using Poisson regression models and historical data analysis</p>
    </div>
    
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-6 text-gray-600 text-lg">Loading teams and leagues...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6 mb-8">
      <div class="flex items-center">
        <svg class="w-6 h-6 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <div>
          <h3 class="text-lg font-medium text-red-800">Error Loading Data</h3>
          <p class="text-red-700">{{ error }}</p>
        </div>
      </div>
      <button @click="loadData" class="mt-4 px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium">
        Retry Loading
      </button>
    </div>

    <!-- Main Predictor -->
    <div v-if="!loading && !error" class="space-y-10">
      
      <!-- Team Selection Section -->
      <div class="bg-white rounded-2xl shadow-xl border border-gray-100 p-8">
        <h3 class="text-2xl font-bold mb-8 text-center text-gray-900">Team Selection</h3>
        
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-10 items-start">
          <!-- Home Team Selector -->
          <div class="team-selector-section">
            <div class="text-center mb-6">
              <h4 class="text-xl font-semibold text-gray-800 mb-3">Home Team</h4>
              <div v-if="selectedHomeTeam" class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-xl p-4 border border-blue-200">
                <div class="font-bold text-blue-900 text-lg">{{ selectedHomeTeam.name }}</div>
                <div class="text-sm text-blue-700 font-medium">{{ getLeagueDisplayName(selectedHomeTeam.league) }}</div>
              </div>
              <div v-else class="bg-gray-50 rounded-xl p-4 border-2 border-dashed border-gray-300">
                <div class="text-gray-500 font-medium">Select home team</div>
              </div>
            </div>
            
            <!-- Home Team League Tabs -->
            <div class="league-tabs mb-6">
              <div class="flex flex-wrap gap-2 justify-center">
                <button
                  v-for="league in availableLeagues"
                  :key="league"
                  @click="homeTeamSelectedLeague = league"
                  :class="[
                    'px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
                    homeTeamSelectedLeague === league
                      ? 'bg-indigo-600 text-white shadow-lg'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  ]"
                >
                  {{ getLeagueDisplayName(league) }}
                </button>
              </div>
            </div>
            
            <!-- Home Team Selection -->
            <div class="team-list max-h-80 overflow-y-auto border rounded-xl bg-gray-50">
              <button
                v-for="team in getTeamsForLeague(homeTeamSelectedLeague)"
                :key="team.name"
                @click="selectHomeTeam(team)"
                :class="[
                  'w-full text-left px-4 py-3 border-b border-gray-200 hover:bg-blue-50 transition-colors last:border-b-0',
                  selectedHomeTeam?.name === team.name ? 'bg-blue-100 border-blue-300 font-semibold text-blue-900' : 'text-gray-700'
                ]"
              >
                <div class="font-medium">{{ team.name }}</div>
              </button>
            </div>
          </div>

          <!-- Prediction Section -->
          <div class="text-center">
            <div class="text-5xl font-light text-gray-400 mb-6">VS</div>
            
            <!-- Match Context -->
            <div v-if="selectedHomeTeam && selectedAwayTeam" class="bg-gray-50 rounded-xl p-5 mb-6 border">
              <div class="text-sm text-gray-600 mb-2 font-medium">Match Context</div>
              <div class="font-semibold text-gray-900">{{ getMatchContext() }}</div>
            </div>

            <!-- Prediction Display -->
            <div v-if="prediction" class="bg-gradient-to-br from-indigo-50 via-blue-50 to-purple-50 rounded-2xl p-6 border border-indigo-200">
              <div class="text-xl font-bold text-gray-900 mb-4">Prediction Results</div>
              <div class="text-5xl font-bold text-indigo-700 mb-4">
                {{ prediction.most_likely_home_score }} - {{ prediction.most_likely_away_score }}
              </div>
              <div class="space-y-2 mb-6">
                <div class="text-gray-700">Expected Goals: {{ prediction.expected_home_goals.toFixed(1) }} - {{ prediction.expected_away_goals.toFixed(1) }}</div>
                <div class="text-indigo-700 font-bold text-lg">{{ getPredictionResult() }}</div>
                <div v-if="savingToHistory" class="text-sm text-blue-600 font-medium flex items-center justify-center gap-1 mt-3">
                  <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  Saving to your prediction history...
                </div>
                <div v-else-if="savedToHistory" class="text-sm text-green-600 font-medium flex items-center justify-center gap-1 mt-3">
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Automatically saved to your prediction history
                </div>
                <div v-else class="text-sm text-gray-500 font-medium flex items-center justify-center gap-1 mt-3">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                  Auto-saves to prediction history
                </div>
              </div>
              <div class="grid grid-cols-3 gap-3 mb-4">
                <div class="bg-white rounded-lg p-4 border border-green-200">
                  <div class="font-semibold text-gray-800">Home Win</div>
                  <div class="text-2xl font-bold text-green-600">{{ (prediction.probabilities.home_win * 100).toFixed(0) }}%</div>
                </div>
                <div class="bg-white rounded-lg p-4 border border-gray-200">
                  <div class="font-semibold text-gray-800">Draw</div>
                  <div class="text-2xl font-bold text-gray-600">{{ (prediction.probabilities.draw * 100).toFixed(0) }}%</div>
                </div>
                <div class="bg-white rounded-lg p-4 border border-red-200">
                  <div class="font-semibold text-gray-800">Away Win</div>
                  <div class="text-2xl font-bold text-red-600">{{ (prediction.probabilities.away_win * 100).toFixed(0) }}%</div>
                </div>
              </div>
              
              <!-- Action Buttons -->
              <div class="flex gap-3 justify-center">
                <router-link 
                  to="/profile" 
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors flex items-center gap-2"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  View History
                </router-link>
              </div>
            </div>

            <!-- Predict Button -->
            <div class="mt-6">
              <button 
                @click="generatePrediction"
                :disabled="!selectedHomeTeam || !selectedAwayTeam || predicting"
                class="px-8 py-4 bg-indigo-600 text-white rounded-xl font-semibold text-lg hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl"
              >
                <span v-if="predicting">Generating Prediction...</span>
                <span v-else>Generate Prediction</span>
              </button>
            </div>
          </div>

          <!-- Away Team Selector -->
          <div class="team-selector-section">
            <div class="text-center mb-6">
              <h4 class="text-xl font-semibold text-gray-800 mb-3">Away Team</h4>
              <div v-if="selectedAwayTeam" class="bg-gradient-to-r from-red-50 to-red-100 rounded-xl p-4 border border-red-200">
                <div class="font-bold text-red-900 text-lg">{{ selectedAwayTeam.name }}</div>
                <div class="text-sm text-red-700 font-medium">{{ getLeagueDisplayName(selectedAwayTeam.league) }}</div>
              </div>
              <div v-else class="bg-gray-50 rounded-xl p-4 border-2 border-dashed border-gray-300">
                <div class="text-gray-500 font-medium">Select away team</div>
              </div>
            </div>
            
            <!-- Away Team League Tabs -->
            <div class="league-tabs mb-6">
              <div class="flex flex-wrap gap-2 justify-center">
                <button
                  v-for="league in availableLeagues"
                  :key="league"
                  @click="awayTeamSelectedLeague = league"
                  :class="[
                    'px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
                    awayTeamSelectedLeague === league
                      ? 'bg-red-600 text-white shadow-lg'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  ]"
                >
                  {{ getLeagueDisplayName(league) }}
                </button>
              </div>
            </div>
            
            <!-- Away Team Selection -->
            <div class="team-list max-h-80 overflow-y-auto border rounded-xl bg-gray-50">
              <button
                v-for="team in getTeamsForLeague(awayTeamSelectedLeague)"
                :key="team.name"
                @click="selectAwayTeam(team)"
                :class="[
                  'w-full text-left px-4 py-3 border-b border-gray-200 hover:bg-red-50 transition-colors last:border-b-0',
                  selectedAwayTeam?.name === team.name ? 'bg-red-100 border-red-300 font-semibold text-red-900' : 'text-gray-700'
                ]"
              >
                <div class="font-medium">{{ team.name }}</div>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Prediction Error -->
      <div v-if="predictionError" class="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
        <svg class="w-8 h-8 text-red-500 mx-auto mb-3" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <h3 class="text-lg font-semibold text-red-800 mb-2">Prediction Error</h3>
        <p class="text-red-700">{{ predictionError }}</p>
      </div>

      <!-- Featured Matchups -->
      <div class="bg-gradient-to-r from-gray-50 to-gray-100 rounded-2xl p-8 border border-gray-200">
        <h3 class="text-2xl font-bold mb-6 text-center text-gray-900">Featured Matchups</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <button
            v-for="suggestion in featuredMatchups"
            :key="suggestion.id"
            @click="loadSuggestion(suggestion)"
            class="p-5 bg-white rounded-xl border border-gray-200 hover:border-indigo-300 hover:shadow-lg transition-all duration-200 text-left"
          >
            <div class="font-bold text-gray-900 mb-1">{{ suggestion.home }} vs {{ suggestion.away }}</div>
            <div class="text-sm text-gray-600">{{ suggestion.description }}</div>
            <div class="text-xs text-gray-500 mt-2">{{ getLeagueDisplayName(suggestion.homeLeague) }}{{ suggestion.homeLeague !== suggestion.awayLeague ? ` vs ${getLeagueDisplayName(suggestion.awayLeague)}` : '' }}</div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { predictMatch, getAvailableTeams, getAvailableLeagues, type PredictMatchResponse } from '@/services/api';
import { usePredictionStore } from '@/stores/prediction';

// Types
interface Team {
  name: string;
  league: string;
}

interface MatchSuggestion {
  id: string;
  home: string;
  away: string;
  homeLeague: string;
  awayLeague: string;
  description: string;
}

// Reactive data
const loading = ref(true);
const error = ref('');
const predicting = ref(false);
const predictionError = ref('');
const savingToHistory = ref(false);
const savedToHistory = ref(false);

const allTeams = ref<Team[]>([]);
const availableLeagues = ref<string[]>([]);

const homeTeamSelectedLeague = ref('');
const awayTeamSelectedLeague = ref('');

const selectedHomeTeam = ref<Team | null>(null);
const selectedAwayTeam = ref<Team | null>(null);

const prediction = ref<PredictMatchResponse | null>(null);

const predictionStore = usePredictionStore();

// League display name mapping
const leagueDisplayNames: Record<string, string> = {
  'E0': 'Premier League',
  'SP1': 'La Liga',
  'I1': 'Serie A',
  'D1': 'Bundesliga',
  'F1': 'Ligue 1',
  'Premier League': 'Premier League',
  'La Liga': 'La Liga',
  'Serie A': 'Serie A',
  'Bundesliga': 'Bundesliga',
  'Ligue 1': 'Ligue 1'
};

// Featured matchup suggestions
const featuredMatchups: MatchSuggestion[] = [
  { id: '1', home: 'Barcelona', away: 'Real Madrid', homeLeague: 'SP1', awayLeague: 'SP1', description: 'El Clasico' },
  { id: '2', home: 'Manchester United', away: 'Manchester City', homeLeague: 'E0', awayLeague: 'E0', description: 'Manchester Derby' },
  { id: '3', home: 'Liverpool', away: 'Arsenal', homeLeague: 'E0', awayLeague: 'E0', description: 'Premier League Classic' },
  { id: '4', home: 'Milan', away: 'Inter', homeLeague: 'I1', awayLeague: 'I1', description: 'Derby della Madonnina' },
  { id: '5', home: 'Bayern Munich', away: 'Borussia Dortmund', homeLeague: 'D1', awayLeague: 'D1', description: 'Der Klassiker' },
  { id: '6', home: 'PSG', away: 'Marseille', homeLeague: 'F1', awayLeague: 'F1', description: 'Le Classique' },
  { id: '7', home: 'Juventus', away: 'Roma', homeLeague: 'I1', awayLeague: 'I1', description: 'North vs South Italy' },
  { id: '8', home: 'Chelsea', away: 'Tottenham', homeLeague: 'E0', awayLeague: 'E0', description: 'London Derby' },
  { id: '9', home: 'Atletico Madrid', away: 'Real Madrid', homeLeague: 'SP1', awayLeague: 'SP1', description: 'Madrid Derby' }
];

// Computed properties
const getTeamsForLeague = (league: string) => {
  return allTeams.value.filter(team => team.league === league);
};

// Methods
const getLeagueDisplayName = (leagueCode: string): string => {
  return leagueDisplayNames[leagueCode] || leagueCode;
};

const getMatchContext = (): string => {
  if (!selectedHomeTeam.value || !selectedAwayTeam.value) return '';
  
  const homeLeague = getLeagueDisplayName(selectedHomeTeam.value.league);
  const awayLeague = getLeagueDisplayName(selectedAwayTeam.value.league);
  
  if (selectedHomeTeam.value.league === selectedAwayTeam.value.league) {
    return `${homeLeague} Match`;
  } else {
    return `Cross-League: ${homeLeague} vs ${awayLeague}`;
  }
};

// Fixed draw detection bug
const getPredictionResult = (): string => {
  if (!prediction.value || !selectedHomeTeam.value || !selectedAwayTeam.value) return '';
  
  // Check the actual predicted scores for draw detection
  const homeScore = prediction.value.most_likely_home_score;
  const awayScore = prediction.value.most_likely_away_score;
  
  if (homeScore === awayScore) {
    return 'Draw';
  } else if (homeScore > awayScore) {
    return `${selectedHomeTeam.value.name} Victory`;
  } else {
    return `${selectedAwayTeam.value.name} Victory`;
  }
};

const selectHomeTeam = (team: Team) => {
  selectedHomeTeam.value = team;
  clearPrediction();
};

const selectAwayTeam = (team: Team) => {
  selectedAwayTeam.value = team;
  clearPrediction();
};

const loadSuggestion = (suggestion: MatchSuggestion) => {
  const homeTeam = allTeams.value.find(t => t.name === suggestion.home && t.league === suggestion.homeLeague);
  const awayTeam = allTeams.value.find(t => t.name === suggestion.away && t.league === suggestion.awayLeague);
  
  if (homeTeam && awayTeam) {
    selectedHomeTeam.value = homeTeam;
    selectedAwayTeam.value = awayTeam;
    homeTeamSelectedLeague.value = homeTeam.league;
    awayTeamSelectedLeague.value = awayTeam.league;
    clearPrediction();
  }
};

const loadData = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    // Load teams and leagues from the backend
    const [teams, leagues] = await Promise.all([
      getAvailableTeams(),
      getAvailableLeagues()
    ]);
    
    // Build team-league mapping from teams data
    allTeams.value = buildTeamLeagueMapping(teams);
    availableLeagues.value = leagues.sort();
    
    // Set default league tabs
    if (leagues.length > 0) {
      homeTeamSelectedLeague.value = leagues[0];
      awayTeamSelectedLeague.value = leagues[0];
    }
    
    console.log('Loaded teams:', teams.length);
    console.log('Loaded leagues:', leagues.length);
    console.log('Team-league mapping:', allTeams.value.length);
    
  } catch (err: any) {
    console.error('Error loading data:', err);
    error.value = 'Failed to load teams and leagues. Please check that the backend services are running.';
  } finally {
    loading.value = false;
  }
};

const buildTeamLeagueMapping = (teams: string[]): Team[] => {
  // Extended team-league mapping with more comprehensive team coverage
  const teamLeagueMap: Record<string, string> = {
    // Premier League teams (E0)
    'Arsenal': 'E0', 'Aston Villa': 'E0', 'Brighton': 'E0', 'Burnley': 'E0',
    'Chelsea': 'E0', 'Crystal Palace': 'E0', 'Everton': 'E0', 'Fulham': 'E0',
    'Liverpool': 'E0', 'Luton Town': 'E0', 'Manchester City': 'E0', 'Manchester United': 'E0',
    'Newcastle': 'E0', 'Nottingham Forest': 'E0', 'Sheffield United': 'E0', 'Tottenham': 'E0',
    'West Ham': 'E0', 'Wolves': 'E0', 'Bournemouth': 'E0', 'Brentford': 'E0',
    'Leicester': 'E0', 'Leeds': 'E0', 'Watford': 'E0', 'Norwich': 'E0',
    
    // La Liga teams (SP1)
    'Barcelona': 'SP1', 'Real Madrid': 'SP1', 'Atletico Madrid': 'SP1', 'Athletic Bilbao': 'SP1',
    'Real Sociedad': 'SP1', 'Villarreal': 'SP1', 'Valencia': 'SP1', 'Sevilla': 'SP1',
    'Real Betis': 'SP1', 'Osasuna': 'SP1', 'Getafe': 'SP1', 'Las Palmas': 'SP1',
    'Girona': 'SP1', 'Alaves': 'SP1', 'Celta Vigo': 'SP1', 'Mallorca': 'SP1',
    'Cadiz': 'SP1', 'Granada': 'SP1', 'Almeria': 'SP1', 'Rayo Vallecano': 'SP1',
    'Ath Madrid': 'SP1', 'Ath Bilbao': 'SP1', 'Espanol': 'SP1',
    
    // Serie A teams (I1)
    'Juventus': 'I1', 'AC Milan': 'I1', 'Inter Milan': 'I1', 'AS Roma': 'I1',
    'Napoli': 'I1', 'Lazio': 'I1', 'Atalanta': 'I1', 'Fiorentina': 'I1',
    'Bologna': 'I1', 'Torino': 'I1', 'Udinese': 'I1', 'Monza': 'I1',
    'Genoa': 'I1', 'Lecce': 'I1', 'Cagliari': 'I1', 'Hellas Verona': 'I1',
    'Frosinone': 'I1', 'Sassuolo': 'I1', 'Empoli': 'I1', 'Salernitana': 'I1',
    'Inter': 'I1', 'Milan': 'I1', 'Roma': 'I1', 'Verona': 'I1',
    
    // Bundesliga teams (D1)
    'Bayern Munich': 'D1', 'Borussia Dortmund': 'D1', 'RB Leipzig': 'D1', 'Union Berlin': 'D1',
    'SC Freiburg': 'D1', 'Bayer Leverkusen': 'D1', 'Eintracht Frankfurt': 'D1', 'Wolfsburg': 'D1',
    'Borussia Monchengladbach': 'D1', 'Mainz': 'D1', 'Augsburg': 'D1', 'Werder Bremen': 'D1',
    'Hoffenheim': 'D1', 'FC Koln': 'D1', 'VfL Bochum': 'D1', 'Heidenheim': 'D1',
    'VfB Stuttgart': 'D1', 'Darmstadt': 'D1', 'Bayern': 'D1', 'Dortmund': 'D1',
    'Leipzig': 'D1', 'Leverkusen': 'D1', 'Frankfurt': 'D1', 'Gladbach': 'D1',
    'Stuttgart': 'D1', 'Freiburg': 'D1', 'Koln': 'D1', 'Bochum': 'D1',
    
    // Ligue 1 teams (F1)
    'Paris Saint-Germain': 'F1', 'Olympique Marseille': 'F1', 'AS Monaco': 'F1', 'Lille': 'F1',
    'Lyon': 'F1', 'Nice': 'F1', 'Lens': 'F1', 'Rennes': 'F1',
    'Montpellier': 'F1', 'Strasbourg': 'F1', 'Nantes': 'F1', 'Reims': 'F1',
    'Toulouse': 'F1', 'Le Havre': 'F1', 'Brest': 'F1', 'Clermont': 'F1',
    'Lorient': 'F1', 'Metz': 'F1', 'Angers': 'F1', 'Monaco': 'F1',
    'Marseille': 'F1', 'PSG': 'F1', 'Paris SG': 'F1', 'St Etienne': 'F1'
  };
  
  // Map teams to their leagues, filtering out unknown teams and providing debugging
  const mappedTeams = teams.map(teamName => {
    const league = teamLeagueMap[teamName];
    if (!league) {
      console.warn(`Unknown team: "${teamName}" - this team will be excluded from the interface`);
      return null;
    }
    return {
      name: teamName,
      league: league
    };
  }).filter((team): team is Team => team !== null);
  
  // Group teams by league for debugging
  const teamsByLeague = mappedTeams.reduce((acc, team) => {
    if (!acc[team.league]) acc[team.league] = [];
    acc[team.league].push(team.name);
    return acc;
  }, {} as Record<string, string[]>);
  
  console.log('Teams by league:');
  Object.entries(teamsByLeague).forEach(([league, teams]) => {
    console.log(`${getLeagueDisplayName(league)}: ${teams.length} teams`);
  });
  
  return mappedTeams;
};

const clearPrediction = () => {
  prediction.value = null;
  predictionError.value = '';
  savedToHistory.value = false;
  savingToHistory.value = false;
};

const generatePrediction = async () => {
  if (!selectedHomeTeam.value || !selectedAwayTeam.value) {
    return;
  }
  
  predicting.value = true;
  predictionError.value = '';
  savedToHistory.value = false;
  
  try {
    // For cross-league matches, we'll use the home team's league as the context
    const league = selectedHomeTeam.value.league;
    
    console.log('Making prediction:', {
      league: league,
      home_team: selectedHomeTeam.value.name,
      away_team: selectedAwayTeam.value.name
    });
    
    const result = await predictMatch({
      league: league,
      home_team: selectedHomeTeam.value.name,
      away_team: selectedAwayTeam.value.name
    });
    
    prediction.value = result;
    console.log('Prediction result:', result);
    
    // Automatically save to history after successful prediction
    savingToHistory.value = true;
    try {
      await predictionStore.savePrediction({
        homeTeam: selectedHomeTeam.value.name,
        awayTeam: selectedAwayTeam.value.name,
        homeLeague: selectedHomeTeam.value.league,
        awayLeague: selectedAwayTeam.value.league,
        predictedHomeScore: result.most_likely_home_score,
        predictedAwayScore: result.most_likely_away_score,
        expectedHomeGoals: result.expected_home_goals,
        expectedAwayGoals: result.expected_away_goals,
        homeWinProbability: result.probabilities.home_win,
        drawProbability: result.probabilities.draw,
        awayWinProbability: result.probabilities.away_win,
        predictedResult: getPredictionResult(),
      });
      savedToHistory.value = true;
      console.log('✅ Prediction automatically saved to history');
    } catch (saveError) {
      console.warn('⚠️ Failed to save prediction to history, but prediction was successful:', saveError);
      // Don't throw the error as the prediction itself was successful
    } finally {
      savingToHistory.value = false;
    }
    
  } catch (err: any) {
    console.error('Error making prediction:', err);
    predictionError.value = 'Failed to generate prediction. This could be due to teams not being in the training data or service unavailability.';
  } finally {
    predicting.value = false;
  }
};

// Lifecycle
onMounted(() => {
  loadData();
});
</script>

<style scoped>
.team-list {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 #f1f5f9;
}

.team-list::-webkit-scrollbar {
  width: 8px;
}

.team-list::-webkit-scrollbar-track {
  background: #f8fafc;
  border-radius: 4px;
}

.team-list::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

.team-list::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

@media (max-width: 768px) {
  .team-list {
    max-height: 240px;
  }
}
</style> 