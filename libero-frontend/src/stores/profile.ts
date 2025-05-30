import { defineStore } from 'pinia';
import { getUserProfile, updateUserPreferences, type UserProfile, type UpdatePreferencesPayload } from '@/services/api';

// Define interfaces based on api.ts (could be imported if shared)
interface Team {
  id: number;
  name: string;
}

interface Player {
  id: number;
  name: string;
}

// Add Competition interface (mirroring api.ts)
interface Competition {
  id: number;
  name: string;
}

// Define the state structure
interface ProfileState {
  user: { id: number; name: string; email: string } | null;
  followedTeams: Team[];
  followedPlayers: Player[];
  followedCompetitions: Competition[]; // Added competitions
  isLoading: boolean;
  error: string | null;
}

export const useProfileStore = defineStore('profile', {
  state: (): ProfileState => ({
    user: null,
    followedTeams: [],
    followedPlayers: [],
    followedCompetitions: [], // Initialize competitions
    isLoading: false,
    error: null,
  }),

  actions: {
    async fetchProfile() {
      this.isLoading = true;
      this.error = null;
      try {
        const profileData: UserProfile = await getUserProfile();
        // Extract basic user info and preferences
        this.user = {
          id: profileData.id,
          name: profileData.name || profileData.username || 'N/A', // Handle potential missing name/username
          email: profileData.email,
        };
        this.followedTeams = profileData.preferences.followed_teams || [];
        this.followedPlayers = profileData.preferences.followed_players || [];
        this.followedCompetitions = profileData.preferences.followed_competitions || []; // Populate competitions
      } catch (err: any) {
        console.error('Failed to fetch profile:', err);
        this.error = err.response?.data?.message || err.message || 'Failed to load profile data.';
        // Reset state on error? Optional, depends on desired UX
        // this.user = null;
        // this.followedTeams = [];
        // this.followedPlayers = [];
      } finally {
        this.isLoading = false;
      }
    },

    async updatePreferences(payload: UpdatePreferencesPayload) {
      this.isLoading = true; // Indicate loading for preference update
      this.error = null;
      try {
        // Call the API to update preferences
        await updateUserPreferences(payload);

        // Option 1: Re-fetch the entire profile to ensure data consistency
        await this.fetchProfile();

        // Option 2: Optimistic update or update based on response (if backend returns updated profile)
        // if (response) {
        //   // Assuming response is the updated UserProfile
        //   this.user = { id: response.id, name: response.name || response.username || 'N/A', email: response.email };
        //   this.followedTeams = response.preferences.followed_teams || [];
        //   this.followedPlayers = response.preferences.followed_players || [];
        // } else {
        //    // If no response body, re-fetch is safer
        //    await this.fetchProfile();
        // }

      } catch (err: any) {
        console.error('Failed to update preferences:', err);
        this.error = err.response?.data?.message || err.message || 'Failed to update preferences.';
        // Optionally re-fetch profile even on error to get the last known good state
        // await this.fetchProfile();
      } finally {
        // Set isLoading to false *after* potential re-fetch in Option 1
        // If using Option 2 without re-fetch, set it here.
        // For Option 1 (current implementation), isLoading is handled by fetchProfile.
        // We might want a separate loading state for preference updates if fetchProfile takes long.
        // For simplicity now, we rely on fetchProfile's loading state.
      }
    },

    // --- Convenience Actions ---

    async addTeam(teamId: number) {
      await this.updatePreferences({ add_teams: [teamId] });
    },

    async removeTeam(teamId: number) {
      await this.updatePreferences({ remove_teams: [teamId] });
    },

    async addPlayer(playerId: number) {
      await this.updatePreferences({ add_players: [playerId] });
    },

    async removePlayer(playerId: number) {
      await this.updatePreferences({ remove_players: [playerId] });
    },

    // --- Competition Actions ---
    async addCompetition(competitionId: number) {
      await this.updatePreferences({ add_competitions: [competitionId] });
    },

    async removeCompetition(competitionId: number) {
      await this.updatePreferences({ remove_competitions: [competitionId] });
    },
  },

  getters: {
    // Example getter: Check if profile data is loaded
    isProfileLoaded: (state): boolean => state.user !== null && !state.isLoading,
    // Add other getters as needed
  },
});