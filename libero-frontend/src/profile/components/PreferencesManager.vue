<template>
  <div class="preferences-manager border rounded-lg p-4 shadow-sm bg-white">
    <h2 class="text-xl font-semibold mb-4">Manage Preferences</h2>

    <!-- Section to Add New Preferences -->
    <div class="add-preference-section mb-6">
      <!-- Added Tailwind classes for border and padding -->
      <h3 class="text-lg font-medium mb-2 border-b border-gray-200 pb-2">Follow New Teams/Players/Competitions</h3>
      <AddPreference @add-team="emitAddTeam" @add-player="emitAddPlayer" @add-competition="emitAddCompetition" />
    </div>

    <!-- Section to Display Followed Teams -->
      <!-- Added Tailwind classes for border and padding -->
    <div class="followed-teams-section mb-6">
      <h3 class="text-lg font-medium mb-2 border-b border-gray-200 pb-2">Followed Teams</h3>
      <FollowedList
        :items="followedTeams"
        item-type="team"
        @remove-item="emitRemoveTeam"
      />
      <p v-if="!followedTeams || followedTeams.length === 0" class="text-gray-500 italic">
        You are not following any teams yet.
      </p>
    </div>

    <!-- Section to Display Followed Players -->
      <!-- Added Tailwind classes for border and padding -->
    <div class="followed-players-section">
      <h3 class="text-lg font-medium mb-2 border-b border-gray-200 pb-2">Followed Players</h3>
      <FollowedList
        :items="followedPlayers"
        item-type="player"
        @remove-item="emitRemovePlayer"
      />
       <p v-if="!followedPlayers || followedPlayers.length === 0" class="text-gray-500 italic">
        You are not following any players yet.
      </p>
    </div>

    <!-- Section to Display Followed Competitions -->
    <div class="followed-competitions-section mt-6"> <!-- Added margin-top -->
      <h3 class="text-lg font-medium mb-2 border-b border-gray-200 pb-2">Followed Competitions</h3>
      <FollowedList
        :items="followedCompetitions"
        item-type="competition"
        @remove-item="emitRemoveCompetition"
      />
      <p v-if="!followedCompetitions || followedCompetitions.length === 0" class="text-gray-500 italic">
        You are not following any competitions yet.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import FollowedList from './FollowedList.vue'; // Placeholder import
import AddPreference from './AddPreference.vue'; // Placeholder import

// Define interfaces for props (could be imported if shared)
interface Team {
  id: number;
  name: string;
}

interface Player {
  id: number;
  name: string;
}

// Add Competition interface
interface Competition {
  id: number;
  name: string;
}

// Define props
const props = defineProps<{
  followedTeams: Team[];
  followedPlayers: Player[];
  followedCompetitions: Competition[]; // Added competitions prop
}>();

// Log props to satisfy TS 'declared but never read' -- props are used in template
console.log('PreferencesManager received props:', props.followedTeams, props.followedPlayers, props.followedCompetitions);

// Define emits for actions triggered within this component or its children
const emit = defineEmits<{
  (e: 'remove-team', id: number): void;
  (e: 'remove-player', id: number): void;
  (e: 'add-team', id: number): void;
  (e: 'add-player', id: number): void;
  (e: 'add-competition', id: number): void; // Added competition emit
  (e: 'remove-competition', id: number): void; // Added competition emit
}>();

// --- Emitter Functions ---
// These functions are called by child components (FollowedList, AddPreference)
// and simply relay the events up to the parent (ProfilePage)

const emitRemoveTeam = (teamId: number) => {
  emit('remove-team', teamId);
};

const emitRemovePlayer = (playerId: number) => {
  emit('remove-player', playerId);
};

const emitAddTeam = (teamId: number) => { // Assuming AddPreference emits the ID
  emit('add-team', teamId);
};

const emitAddPlayer = (playerId: number) => { // Assuming AddPreference emits the ID
  emit('add-player', playerId);
};

const emitAddCompetition = (competitionId: number) => {
  emit('add-competition', competitionId);
};

const emitRemoveCompetition = (competitionId: number) => {
  emit('remove-competition', competitionId);
};

</script>

<!-- Style block removed as Tailwind classes are used directly -->