<template>
  <div class="p-4 sm:p-6 md:p-10 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7 gap-2 sm:gap-3 md:gap-4">
    <PlayerCard
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
import PlayerCard from './PlayerCard.vue';

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

// Sample data for MVP players
const players: Player[] = [
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
</script>