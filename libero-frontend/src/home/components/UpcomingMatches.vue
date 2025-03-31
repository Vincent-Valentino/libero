<template>
  <div class="p-10 md:p-20 w-full flex flex-col">
    <div class="flex flex-col gap-2">
      <h1 class="text-4xl font-title mb-4">Upcoming Matches</h1>
      
      <LeagueTabs @tab-change="handleTabChange" />
      
      <!-- All Leagues View - Horizontal League Kanban -->
      <div v-if="selectedTab === 'all'" class="flex justify-start overflow-x-auto">
        <div class="flex gap-2 max-w-full">
          <!-- League Columns - Each with independent height -->
          <div v-for="league in leagueColumns" :key="league.id" class="w-56 flex-shrink-0 self-start">
            <LeagueSection 
              :leagueName="league.name" 
              :leagueLogoPath="league.logo">
              
              <div class="space-y-2">
                <div class="flex flex-col">
                  <MatchCard v-for="(match, index) in league.matches" 
                    :key="index"
                    :class="{ 'mt-2': index > 0 }"
                    :homeTeam="match.homeTeam" 
                    :awayTeam="match.awayTeam" 
                    :matchStarted="match.matchStarted" 
                    :matchDate="match.matchDate" 
                    :matchStatus="match.matchStatus" 
                    :stadium="match.stadium"
                  />
                  
                  <!-- Empty state if no matches -->
                  <div v-if="league.matches.length === 0" class="text-center py-3 text-gray-500 text-xs">
                    No matches scheduled
                  </div>
                </div>
              </div>
            </LeagueSection>
          </div>
        </div>
      </div>
      
      <!-- Specific League View - Date-based Kanban -->
      <div v-else class="flex justify-start overflow-x-auto">
        <div class="flex gap-2 max-w-full">
          <!-- Date Columns - Each with independent height -->
          <div v-for="column in dateColumns" :key="column.id" class="w-56 flex-shrink-0 self-start">
            <div class="bg-neutral-100 rounded-2xl p-2">
              <div class="flex items-center mb-1.5">
                <h2 class="font-bold text-xs">{{ column.name }}</h2>
              </div>
              
              <div class="flex flex-col">
                <!-- Display matches for this date column and selected league -->
                <template v-for="(match, index) in getMatchesForColumn(column.id)" :key="index">
                  <MatchCard 
                    :class="{ 'mt-2': index > 0 }"
                    :homeTeam="match.homeTeam" 
                    :awayTeam="match.awayTeam" 
                    :matchStarted="match.matchStarted" 
                    :matchDate="match.matchDate" 
                    :matchStatus="match.matchStatus" 
                    :stadium="match.stadium"
                  />
                </template>
                
                <!-- Empty state if no matches for this date column -->
                <div v-if="getMatchesForColumn(column.id).length === 0" 
                     class="text-center py-3 text-gray-500 text-xs">
                  No matches scheduled
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import LeagueTabs from './LeagueTabs.vue';
import LeagueSection from './LeagueSection.vue';
import MatchCard from './MatchCard.vue';

// Define types for league IDs and date columns
type LeagueId = 'premier-league' | 'laliga' | 'serie-a' | 'bundesliga' | 'ligue-1' | 'all';
type DateColumnId = 'today' | 'tomorrow' | 'upcoming';

// Track the selected tab
const selectedTab = ref<LeagueId>('all');

// Handle tab change from LeagueTabs component
const handleTabChange = (tabId: string) => {
  selectedTab.value = tabId as LeagueId;
};

// Date columns for specific league view
const dateColumns = [
  { id: 'today', name: 'Today' },
  { id: 'tomorrow', name: 'Tomorrow' },
  { id: 'upcoming', name: 'Upcoming' }
];

// League columns configuration
const leagueColumns = [
  { 
    id: 'premier-league' as LeagueId, 
    name: 'Premier League', 
    logo: '/Premier League.svg',
    matches: [
      {
        homeTeam: { 
          name: 'Arsenal', 
          acronym: 'ARS', 
          colorClass: 'bg-gray-100',
          score: 2,
          logo: 'https://resources.premierleague.com/premierleague/badges/t3.svg'
        },
        awayTeam: { 
          name: 'Chelsea', 
          acronym: 'CHE', 
          colorClass: 'bg-blue-100',
          score: 1,
          logo: 'https://resources.premierleague.com/premierleague/badges/t8.svg'
        },
        matchStarted: true,
        matchStatus: "75'",
        stadium: "Emirates Stadium"
      },
      {
        homeTeam: { 
          name: 'Man United', 
          acronym: 'MUN', 
          colorClass: 'bg-red-100',
          logo: 'https://resources.premierleague.com/premierleague/badges/t1.svg'
        },
        awayTeam: { 
          name: 'Liverpool', 
          acronym: 'LIV', 
          colorClass: 'bg-red-700',
          textColorClass: 'text-white',
          logo: 'https://resources.premierleague.com/premierleague/badges/t14.svg'
        },
        matchStarted: false,
        matchDate: "Today",
        matchStatus: "20:45",
        stadium: "Old Trafford"
      }
    ]
  },
  { 
    id: 'laliga' as LeagueId, 
    name: 'La Liga', 
    logo: '/LaLiga.svg',
    matches: [
      {
        homeTeam: { 
          name: 'Barcelona', 
          acronym: 'BAR', 
          colorClass: 'bg-blue-100',
          logo: 'https://assets.laliga.com/squad/2023/t178/p82464/32/512x512/p82464_t178_2023_1_003_000.png'
        },
        awayTeam: { 
          name: 'Real Madrid', 
          acronym: 'RMA', 
          colorClass: 'bg-gray-100',
          logo: 'https://assets.laliga.com/squad/2023/t186/p48772/32/512x512/p48772_t186_2023_1_003_000.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "18:00",
        stadium: "Camp Nou"
      },
      {
        homeTeam: { 
          name: 'Sevilla', 
          acronym: 'SEV', 
          colorClass: 'bg-red-200',
          score: 0,
          logo: 'https://assets.laliga.com/squad/2023/t179/p49888/32/512x512/p49888_t179_2023_1_003_000.png'
        },
        awayTeam: { 
          name: 'Atletico Madrid', 
          acronym: 'ATM', 
          colorClass: 'bg-red-500',
          textColorClass: 'text-white',
          score: 0,
          logo: 'https://assets.laliga.com/squad/2023/t175/p77906/32/512x512/p77906_t175_2023_1_003_000.png'
        },
        matchStarted: true,
        matchStatus: "12'",
        stadium: "Ramón Sánchez Pizjuán"
      }
    ]
  },
  { 
    id: 'serie-a' as LeagueId, 
    name: 'Serie A', 
    logo: '/Lega Serie A.svg',
    matches: [
      {
        homeTeam: { 
          name: 'Juventus', 
          acronym: 'JUV', 
          colorClass: 'bg-black',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/496.png'
        },
        awayTeam: { 
          name: 'Inter', 
          acronym: 'INT', 
          colorClass: 'bg-blue-800',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/505.png'
        },
        matchStarted: false,
        matchDate: "Today",
        matchStatus: "21:00",
        stadium: "Allianz Stadium"
      },
      {
        homeTeam: { 
          name: 'AC Milan', 
          acronym: 'MIL', 
          colorClass: 'bg-red-600',
          textColorClass: 'text-white',
          score: 3,
          logo: 'https://media.api-sports.io/football/teams/489.png'
        },
        awayTeam: { 
          name: 'Napoli', 
          acronym: 'NAP', 
          colorClass: 'bg-blue-500',
          textColorClass: 'text-white',
          score: 1,
          logo: 'https://media.api-sports.io/football/teams/492.png'
        },
        matchStarted: true,
        matchStatus: "90'+2",
        stadium: "San Siro"
      },
      {
        homeTeam: { 
          name: 'AC Milan', 
          acronym: 'MIL', 
          colorClass: 'bg-red-600',
          textColorClass: 'text-white',
          score: 3,
          logo: 'https://media.api-sports.io/football/teams/489.png'
        },
        awayTeam: { 
          name: 'Napoli', 
          acronym: 'NAP', 
          colorClass: 'bg-blue-500',
          textColorClass: 'text-white',
          score: 1,
          logo: 'https://media.api-sports.io/football/teams/492.png'
        },
        matchStarted: true,
        matchStatus: "90'+2",
        stadium: "San Siro"
      },
      {
        homeTeam: { 
          name: 'AC Milan', 
          acronym: 'MIL', 
          colorClass: 'bg-red-600',
          textColorClass: 'text-white',
          score: 3,
          logo: 'https://media.api-sports.io/football/teams/489.png'
        },
        awayTeam: { 
          name: 'Napoli', 
          acronym: 'NAP', 
          colorClass: 'bg-blue-500',
          textColorClass: 'text-white',
          score: 1,
          logo: 'https://media.api-sports.io/football/teams/492.png'
        },
        matchStarted: true,
        matchStatus: "90'+2",
        stadium: "San Siro"
      }
    ]
  },
  { 
    id: 'bundesliga' as LeagueId, 
    name: 'Bundesliga', 
    logo: '/Bundesliga.svg',
    matches: [
      {
        homeTeam: { 
          name: 'Bayern Munich', 
          acronym: 'BAY', 
          colorClass: 'bg-red-500',
          textColorClass: 'text-white',
          score: 2,
          logo: 'https://media.api-sports.io/football/teams/157.png'
        },
        awayTeam: { 
          name: 'Dortmund', 
          acronym: 'BVB', 
          colorClass: 'bg-yellow-400',
          score: 0,
          logo: 'https://media.api-sports.io/football/teams/165.png'
        },
        matchStarted: true,
        matchStatus: "45'",
        stadium: "Allianz Arena"
      },
      {
        homeTeam: { 
          name: 'RB Leipzig', 
          acronym: 'RBL', 
          colorClass: 'bg-blue-200',
          logo: 'https://media.api-sports.io/football/teams/173.png'
        },
        awayTeam: { 
          name: 'Leverkusen', 
          acronym: 'LEV', 
          colorClass: 'bg-red-600',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/168.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "16:30",
        stadium: "Red Bull Arena"
      }
    ]
  },
  { 
    id: 'ligue-1' as LeagueId, 
    name: 'Ligue 1', 
    logo: '/Ligue 1 Uber Eats.svg',
    matches: [
      {
        homeTeam: { 
          name: 'PSG', 
          acronym: 'PSG', 
          colorClass: 'bg-blue-800',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/85.png'
        },
        awayTeam: { 
          name: 'Marseille', 
          acronym: 'MAR', 
          colorClass: 'bg-blue-500',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/81.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "20:00",
        stadium: "Parc des Princes"
      },
      {
        homeTeam: { 
          name: 'Lyon', 
          acronym: 'LYO', 
          colorClass: 'bg-red-800',
          textColorClass: 'text-white',
          score: 1,
          logo: 'https://media.api-sports.io/football/teams/80.png'
        },
        awayTeam: { 
          name: 'Monaco', 
          acronym: 'MON', 
          colorClass: 'bg-red-500',
          textColorClass: 'text-white',
          score: 1,
          logo: 'https://media.api-sports.io/football/teams/91.png'
        },
        matchStarted: true,
        matchStatus: "60'",
        stadium: "Groupama Stadium"
      }
    ]
  }
];

// All matches data organized by league and date
const allMatches: Record<Exclude<LeagueId, 'all'>, Record<DateColumnId, any[]>> = {
  'premier-league': {
    'today': [
      {
        homeTeam: { 
          name: 'Man United', 
          acronym: 'MUN', 
          colorClass: 'bg-red-100',
          logo: 'https://resources.premierleague.com/premierleague/badges/t1.svg'
        },
        awayTeam: { 
          name: 'Liverpool', 
          acronym: 'LIV', 
          colorClass: 'bg-red-700',
          textColorClass: 'text-white',
          logo: 'https://resources.premierleague.com/premierleague/badges/t14.svg'
        },
        matchStarted: false,
        matchDate: "Today",
        matchStatus: "20:45",
        stadium: "Old Trafford"
      }
    ],
    'tomorrow': [],
    'upcoming': [
      {
        homeTeam: { 
          name: 'Tottenham', 
          acronym: 'TOT', 
          colorClass: 'bg-blue-900',
          textColorClass: 'text-white',
          logo: 'https://resources.premierleague.com/premierleague/badges/t6.svg'
        },
        awayTeam: { 
          name: 'Man City', 
          acronym: 'MCI', 
          colorClass: 'bg-blue-400',
          logo: 'https://resources.premierleague.com/premierleague/badges/t43.svg'
        },
        matchStarted: false,
        matchDate: "Oct 24",
        matchStatus: "16:30",
        stadium: "Tottenham Hotspur Stadium"
      }
    ]
  },
  'laliga': {
    'today': [],
    'tomorrow': [
      {
        homeTeam: { 
          name: 'Barcelona', 
          acronym: 'BAR', 
          colorClass: 'bg-blue-100',
          logo: 'https://assets.laliga.com/squad/2023/t178/p82464/32/512x512/p82464_t178_2023_1_003_000.png'
        },
        awayTeam: { 
          name: 'Real Madrid', 
          acronym: 'RMA', 
          colorClass: 'bg-gray-100',
          logo: 'https://assets.laliga.com/squad/2023/t186/p48772/32/512x512/p48772_t186_2023_1_003_000.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "18:00",
        stadium: "Camp Nou"
      }
    ],
    'upcoming': []
  },
  'serie-a': {
    'today': [
      {
        homeTeam: { 
          name: 'Juventus', 
          acronym: 'JUV', 
          colorClass: 'bg-black',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/496.png'
        },
        awayTeam: { 
          name: 'Inter', 
          acronym: 'INT', 
          colorClass: 'bg-blue-800',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/505.png'
        },
        matchStarted: false,
        matchDate: "Today",
        matchStatus: "21:00",
        stadium: "Allianz Stadium"
      }
    ],
    'tomorrow': [],
    'upcoming': []
  },
  'bundesliga': {
    'today': [],
    'tomorrow': [
      {
        homeTeam: { 
          name: 'RB Leipzig', 
          acronym: 'RBL', 
          colorClass: 'bg-blue-200',
          logo: 'https://media.api-sports.io/football/teams/173.png'
        },
        awayTeam: { 
          name: 'Leverkusen', 
          acronym: 'LEV', 
          colorClass: 'bg-red-600',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/168.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "16:30",
        stadium: "Red Bull Arena"
      }
    ],
    'upcoming': []
  },
  'ligue-1': {
    'today': [],
    'tomorrow': [
      {
        homeTeam: { 
          name: 'PSG', 
          acronym: 'PSG', 
          colorClass: 'bg-blue-800',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/85.png'
        },
        awayTeam: { 
          name: 'Marseille', 
          acronym: 'MAR', 
          colorClass: 'bg-blue-500',
          textColorClass: 'text-white',
          logo: 'https://media.api-sports.io/football/teams/81.png'
        },
        matchStarted: false,
        matchDate: "Tomorrow",
        matchStatus: "20:00",
        stadium: "Parc des Princes"
      }
    ],
    'upcoming': []
  }
};

// Get matches for a specific column and the selected league
const getMatchesForColumn = (columnId: string) => {
  if (!selectedTab.value || selectedTab.value === 'all') {
    return [];
  }
  
  const league = selectedTab.value as Exclude<LeagueId, 'all'>;
  return allMatches[league][columnId as DateColumnId] || [];
};
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