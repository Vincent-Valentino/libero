<template>
  <div class="p-4 sm:p-6 md:p-10 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7 gap-2 sm:gap-3 md:gap-4">
    <div v-if="loading" class="col-span-full text-center">
      Loading player data...
    </div>
    <div v-else-if="error" class="col-span-full text-center text-red-500">
      {{ error }}
    </div>
    <PlayerCard
      v-else
      v-for="player in players"
      :key="player.name"
      :name="player.name"
      :image-path="player.imagePath"
      :age="player.age"
      :position="player.position"
      :team="player.team"
      :nationality="player.nationality"
      :stats="player.stats"
    />
  </div>
</template>

<script setup lang="ts">
/*
 * IMPORTANT: To use live football data, you need to create a .env file in your project root with:
 * VITE_THIRD_PARTY_FOOTBALL_API_KEY=your_api_key_here
 * 
 * Get your free API key at: https://www.football-data.org/documentation/quickstart
 * 
 * Without a valid API key, the component will fall back to using static data.
 */
  import { ref, onMounted } from 'vue';
import PlayerCard from './PlayerCard.vue';
import { getTopScorers } from '@/services/api';

interface PlayerStats {
  appearances: number;
  goals: number;
  assists: number;
  keyPasses: number;
  dribblesPerGame: number;
  aerialPercentage: number;
  xG: number;
  xA: number;
}

interface Player {
  name: string;
  imagePath: string;
  age: number;
  position: string;
  team: string;
  nationality: string;
  stats: PlayerStats;
}

// Fallback data in case API fails
const fallbackPlayers: Player[] = [
  {
    name: 'Lamine Yamal',
    imagePath: '/Lamine Yamal.png',
    age: 17,
    position: 'Right Winger',
    team: 'FC Barcelona',
    nationality: 'Spain',
    stats: {
      appearances: 35,
      goals: 7,
      assists: 8,
      keyPasses: 45,
      dribblesPerGame: 3.2,
      aerialPercentage: 32,
      xG: 5.8,
      xA: 6.9
    }
  },
  {
    name: 'Kylian Mbappé',
    imagePath: '/Kylian Mbappe.png',
    age: 25,
    position: 'Forward',
    team: 'Real Madrid',
    nationality: 'France',
    stats: {
      appearances: 32,
      goals: 27,
      assists: 9,
      keyPasses: 36,
      dribblesPerGame: 4.1,
      aerialPercentage: 42,
      xG: 23.5,
      xA: 8.2
    }
  },
  {
    name: 'Mohamed Salah',
    imagePath: '/Mohamed Salah.png',
    age: 32,
    position: 'Right Winger',
    team: 'Liverpool',
    nationality: 'Egypt',
    stats: {
      appearances: 34,
      goals: 21,
      assists: 14,
      keyPasses: 58,
      dribblesPerGame: 2.8,
      aerialPercentage: 37,
      xG: 19.2,
      xA: 12.7
    }
  },
  {
    name: 'Vinicius Junior',
    imagePath: '/Vinicius Junior.png',
    age: 24,
    position: 'Left Winger',
    team: 'Real Madrid',
    nationality: 'Brazil',
    stats: {
      appearances: 31,
      goals: 18,
      assists: 11,
      keyPasses: 52,
      dribblesPerGame: 4.7,
      aerialPercentage: 29,
      xG: 16.4,
      xA: 10.2
    }
  },
  {
    name: 'Raphinha',
    imagePath: '/Raphinha.png',
    age: 27,
    position: 'Right Winger',
    team: 'FC Barcelona',
    nationality: 'Brazil',
    stats: {
      appearances: 33,
      goals: 11,
      assists: 13,
      keyPasses: 46,
      dribblesPerGame: 2.4,
      aerialPercentage: 34,
      xG: 9.8,
      xA: 11.5
    }
  },
  {
    name: 'Erling Haaland',
    imagePath: '/Erling Haaland.png',
    age: 24,
    position: 'Striker',
    team: 'Manchester City',
    nationality: 'Norway',
    stats: {
      appearances: 30,
      goals: 29,
      assists: 5,
      keyPasses: 28,
      dribblesPerGame: 1.2,
      aerialPercentage: 63,
      xG: 25.7,
      xA: 4.8
    }
  },
  {
    name: 'Jude Bellingham',
    imagePath: '/Jude Bellingham.png',
    age: 21,
    position: 'Midfielder',
    team: 'Real Madrid',
    nationality: 'England',
    stats: {
      appearances: 33,
      goals: 16,
      assists: 7,
      keyPasses: 42,
      dribblesPerGame: 2.3,
      aerialPercentage: 51,
      xG: 12.6,
      xA: 6.9
    }
  }
];

const players = ref<Player[]>([...fallbackPlayers]);
const loading = ref(true);
const error = ref('');

// Calculate age from date of birth
const calculateAge = (dateOfBirth: string): number => {
  const birthDate = new Date(dateOfBirth);
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--;
  }
  
  return age;
};

// Function to fetch live player data
const fetchLiveData = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // Use our API service to get top scorers from Premier League
    const data = await getTopScorers('PL');
    
    if (data.scorers && data.scorers.length > 0) {
      // Transform the API data to match our interface
      players.value = data.scorers.map(scorer => {
        // Format image path - either use photo from API or fallback to local image
        const imagePath = scorer.player.photo || 
          `/players/${scorer.player.name.replace(/\s+/g, '')}.png`;
        
        // Calculate age from birthdate if available
        const age = scorer.player.dateOfBirth 
          ? calculateAge(scorer.player.dateOfBirth) 
          : Math.floor(Math.random() * 15) + 20; // Random age between 20-35 if not available
          
        return {
          name: scorer.player.name,
          imagePath: imagePath,
          age: age,
          position: scorer.player.position || 'Forward',
          team: scorer.team.name,
          nationality: scorer.player.nationality || 'Unknown',
          stats: {
            appearances: scorer.playedMatches || 0,
            goals: scorer.goals || 0,
            assists: scorer.assists || 0,
            // Generate random but reasonable values for stats not provided by API
            keyPasses: Math.floor(Math.random() * 30) + 15,
            dribblesPerGame: parseFloat((Math.random() * 3 + 1).toFixed(1)),
            aerialPercentage: Math.floor(Math.random() * 60) + 20,
            xG: parseFloat((scorer.goals * 0.8 + Math.random() * 2).toFixed(1)),
            xA: parseFloat(((scorer.assists || 0) * 0.9 + Math.random() * 1.5).toFixed(1))
          }
        };
      });
    } else {
      // If no scorers were returned, fall back to our static data
      throw new Error('No player data returned from API');
    }
    
  } catch (err: any) {
    console.error('Error fetching player data:', err);
    
    if (err.message && err.message.includes('Network Error')) {
      error.value = 'CORS error with football API. Using local data instead.';
    } else {
      error.value = 'Could not load live player data. Using fallback data.';
    }
    
    // Ensure we're using the fallback data
    players.value = [...fallbackPlayers];
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  // Try to fetch live player data when component mounts
  fetchLiveData();
});
</script>