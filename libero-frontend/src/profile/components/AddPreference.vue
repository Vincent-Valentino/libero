<template>
  <div class="add-preference">
    <!-- Team Search and Selection -->
    <div class="space-y-4">
      <div class="flex flex-col md:flex-row gap-4">
        <!-- League Selection -->
        <div class="flex-1">
          <label for="league-select" class="block text-sm font-medium text-gray-700 mb-2">Select League</label>
          <select
            id="league-select"
            v-model="selectedLeague"
            @change="filterTeamsByLeague"
            class="block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          >
            <option value="">All Leagues</option>
            <option v-for="league in availableLeagues" :key="league" :value="league">
              {{ getLeagueDisplayName(league) }}
            </option>
          </select>
        </div>

        <!-- Team Search -->
        <div class="flex-1">
          <label for="team-search" class="block text-sm font-medium text-gray-700 mb-2">Search Teams</label>
          <input
            type="text"
            id="team-search"
            v-model="teamSearchQuery"
            placeholder="Type to search teams..."
            class="block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          />
        </div>
      </div>

      <!-- Teams Display -->
      <div v-if="loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600 mx-auto"></div>
        <p class="mt-2 text-gray-600 text-sm">Loading teams...</p>
      </div>

      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
        <p class="text-red-800 text-sm">{{ error }}</p>
        <button @click="loadTeams" class="mt-2 px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700 transition-colors">
          Retry
        </button>
      </div>

      <div v-else-if="filteredTeams.length > 0" class="border border-gray-200 rounded-lg">
        <div class="max-h-64 overflow-y-auto">
          <button
            v-for="team in filteredTeams"
            :key="team.name"
            @click="handleAddTeam(team)"
            :disabled="isTeamAlreadyFollowed(team)"
            class="w-full text-left px-4 py-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-gray-50"
          >
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-medium text-gray-900">{{ team.name }}</h4>
                <p class="text-sm text-gray-600">{{ getLeagueDisplayName(team.league) }}</p>
              </div>
              <div class="flex items-center">
                <span v-if="isTeamAlreadyFollowed(team)" class="text-xs text-green-600 font-medium bg-green-50 px-2 py-1 rounded">
                  Following
                </span>
                <svg v-else class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
              </div>
            </div>
          </button>
        </div>
      </div>

      <div v-else-if="teamSearchQuery || selectedLeague" class="text-center py-8 text-gray-600">
        <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
        </svg>
        <p class="text-sm">No teams found matching your search.</p>
      </div>

      <div v-else class="text-center py-8 text-gray-600">
        <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
        </svg>
        <p class="text-sm">Select a league or search to see available teams.</p>
      </div>

      <!-- Feedback message -->
      <div v-if="feedbackMessage" :class="[
        'p-3 rounded-lg text-sm',
        feedbackType === 'success' ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'
      ]">
        {{ feedbackMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineEmits, defineProps, watch } from 'vue';
import { getAvailableTeams, getAvailableLeagues } from '@/services/api';

// Define interfaces
interface Team {
  name: string;
  league: string;
}

interface FollowedTeam {
  id: number;
  name: string;
}

// Props
const props = defineProps<{
  followedTeams?: FollowedTeam[];
}>();

// Define emits - only team-related
const emit = defineEmits<{
  (e: 'add-team', id: number): void;
}>();

// Reactive data
const loading = ref(false);
const error = ref('');
const allTeams = ref<Team[]>([]);
const availableLeagues = ref<string[]>([]);
const selectedLeague = ref('');
const teamSearchQuery = ref('');
const feedbackMessage = ref('');
const feedbackType = ref<'success' | 'error'>('success');

// League display names
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

// Team to league mapping (same as in MatchPredictor)
const teamLeagueMap: Record<string, string> = {
  // Premier League teams
  'Arsenal': 'E0', 'Aston Villa': 'E0', 'Brighton': 'E0', 'Burnley': 'E0',
  'Chelsea': 'E0', 'Crystal Palace': 'E0', 'Everton': 'E0', 'Fulham': 'E0',
  'Liverpool': 'E0', 'Luton Town': 'E0', 'Manchester City': 'E0', 'Manchester United': 'E0',
  'Newcastle': 'E0', 'Nottingham Forest': 'E0', 'Sheffield United': 'E0', 'Tottenham': 'E0',
  'West Ham': 'E0', 'Wolves': 'E0', 'Bournemouth': 'E0', 'Brentford': 'E0',
  'Leicester': 'E0', 'Leeds': 'E0', 'Watford': 'E0', 'Norwich': 'E0',
  
  // La Liga teams
  'Barcelona': 'SP1', 'Real Madrid': 'SP1', 'Atletico Madrid': 'SP1', 'Athletic Bilbao': 'SP1',
  'Real Sociedad': 'SP1', 'Villarreal': 'SP1', 'Valencia': 'SP1', 'Sevilla': 'SP1',
  'Real Betis': 'SP1', 'Osasuna': 'SP1', 'Getafe': 'SP1', 'Las Palmas': 'SP1',
  'Girona': 'SP1', 'Alaves': 'SP1', 'Celta Vigo': 'SP1', 'Mallorca': 'SP1',
  'Cadiz': 'SP1', 'Granada': 'SP1', 'Almeria': 'SP1', 'Rayo Vallecano': 'SP1',
  'Ath Madrid': 'SP1', 'Ath Bilbao': 'SP1', 'Espanol': 'SP1',
  
  // Serie A teams
  'Juventus': 'I1', 'AC Milan': 'I1', 'Inter Milan': 'I1', 'AS Roma': 'I1',
  'Napoli': 'I1', 'Lazio': 'I1', 'Atalanta': 'I1', 'Fiorentina': 'I1',
  'Bologna': 'I1', 'Torino': 'I1', 'Udinese': 'I1', 'Monza': 'I1',
  'Genoa': 'I1', 'Lecce': 'I1', 'Cagliari': 'I1', 'Hellas Verona': 'I1',
  'Frosinone': 'I1', 'Sassuolo': 'I1', 'Empoli': 'I1', 'Salernitana': 'I1',
  'Inter': 'I1', 'Milan': 'I1', 'Roma': 'I1', 'Verona': 'I1',
  
  // Bundesliga teams
  'Bayern Munich': 'D1', 'Borussia Dortmund': 'D1', 'RB Leipzig': 'D1', 'Union Berlin': 'D1',
  'SC Freiburg': 'D1', 'Bayer Leverkusen': 'D1', 'Eintracht Frankfurt': 'D1', 'Wolfsburg': 'D1',
  'Borussia Monchengladbach': 'D1', 'Mainz': 'D1', 'Augsburg': 'D1', 'Werder Bremen': 'D1',
  'Hoffenheim': 'D1', 'FC Koln': 'D1', 'VfL Bochum': 'D1', 'Heidenheim': 'D1',
  'VfB Stuttgart': 'D1', 'Darmstadt': 'D1', 'Bayern': 'D1', 'Dortmund': 'D1',
  'Leipzig': 'D1', 'Leverkusen': 'D1', 'Frankfurt': 'D1', 'Gladbach': 'D1',
  'Stuttgart': 'D1', 'Freiburg': 'D1', 'Koln': 'D1', 'Bochum': 'D1',
  
  // Ligue 1 teams
  'Paris Saint-Germain': 'F1', 'Olympique Marseille': 'F1', 'AS Monaco': 'F1', 'Lille': 'F1',
  'Lyon': 'F1', 'Nice': 'F1', 'Lens': 'F1', 'Rennes': 'F1',
  'Montpellier': 'F1', 'Strasbourg': 'F1', 'Nantes': 'F1', 'Reims': 'F1',
  'Toulouse': 'F1', 'Le Havre': 'F1', 'Brest': 'F1', 'Clermont': 'F1',
  'Lorient': 'F1', 'Metz': 'F1', 'Angers': 'F1', 'Monaco': 'F1',
  'Marseille': 'F1', 'PSG': 'F1', 'Paris SG': 'F1', 'St Etienne': 'F1'
};

// Computed properties
const filteredTeams = computed(() => {
  let teams = allTeams.value;
  
  // Filter by league if selected
  if (selectedLeague.value) {
    teams = teams.filter(team => team.league === selectedLeague.value);
  }
  
  // Filter by search query
  if (teamSearchQuery.value) {
    const query = teamSearchQuery.value.toLowerCase();
    teams = teams.filter(team => 
      team.name.toLowerCase().includes(query)
    );
  }
  
  return teams.slice(0, 20); // Limit to 20 results for performance
});

// Methods
const getLeagueDisplayName = (leagueCode: string): string => {
  return leagueDisplayNames[leagueCode] || leagueCode;
};

const isTeamAlreadyFollowed = (team: Team): boolean => {
  return props.followedTeams?.some(followedTeam => followedTeam.name === team.name) || false;
};

const buildTeamLeagueMapping = (teams: string[], leagues: string[]): Team[] => {
  return teams.map(teamName => {
    const league = teamLeagueMap[teamName];
    if (!league) {
      console.warn(`Unknown team: "${teamName}"`);
      return null;
    }
    return {
      name: teamName,
      league: league
    };
  }).filter((team): team is Team => team !== null);
};

const loadTeams = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const [teams, leagues] = await Promise.all([
      getAvailableTeams(),
      getAvailableLeagues()
    ]);
    
    allTeams.value = buildTeamLeagueMapping(teams, leagues);
    availableLeagues.value = leagues.sort();
    
  } catch (err: any) {
    console.error('Error loading teams:', err);
    error.value = 'Failed to load teams. Please try again.';
  } finally {
    loading.value = false;
  }
};

const filterTeamsByLeague = () => {
  teamSearchQuery.value = ''; // Clear search when changing league
};

const handleAddTeam = (team: Team) => {
  if (isTeamAlreadyFollowed(team)) {
    showFeedback(`You are already following ${team.name}`, 'error');
    return;
  }
  
  // Generate a placeholder ID (in real app, this would come from search results)
  const placeholderTeamId = Math.floor(Math.random() * 10000);
  console.log('Emitting add-team:', { id: placeholderTeamId, name: team.name });
  
  emit('add-team', placeholderTeamId);
  showFeedback(`Successfully started following ${team.name}!`, 'success');
};

const showFeedback = (message: string, type: 'success' | 'error') => {
  feedbackMessage.value = message;
  feedbackType.value = type;
  setTimeout(() => {
    feedbackMessage.value = '';
  }, 3000);
};

// Lifecycle
onMounted(() => {
  loadTeams();
});

// Clear search when teams change
watch(() => props.followedTeams, () => {
  // Could refresh the UI or clear search if needed
}, { deep: true });
</script>

<!-- No scoped styles needed as Tailwind is used -->