<template>
  <div class="user-info-card border rounded-lg p-6 shadow-sm bg-white">
    <div class="flex items-start gap-6">
      <!-- User Avatar Section -->
      <div class="flex-shrink-0">
        <div class="w-20 h-20 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-full flex items-center justify-center">
          <svg class="w-10 h-10 text-white" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>

      <!-- User Information Section -->
      <div class="flex-1">
        <div class="flex items-center gap-4 mb-4">
          <h2 class="text-2xl font-bold text-gray-900">Football Fan Profile</h2>
          <span class="px-3 py-1 text-xs font-medium text-green-800 bg-green-100 rounded-full">
            Active Follower
          </span>
        </div>

        <div v-if="user" class="space-y-3">
          <!-- Name -->
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 text-gray-400">
              <svg fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-600">Name</p>
              <p class="text-lg font-medium text-gray-900">{{ user.name || 'Football Fan' }}</p>
            </div>
          </div>

          <!-- Email -->
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 text-gray-400">
              <svg fill="currentColor" viewBox="0 0 20 20">
                <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
                <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-600">Email</p>
              <p class="text-lg font-medium text-gray-900">{{ user.email }}</p>
            </div>
          </div>

          <!-- User ID -->
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 text-gray-400">
              <svg fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2H4zm3 6V7h6v3H7z" clip-rule="evenodd" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-600">User ID</p>
              <p class="text-lg font-medium text-gray-900">#{{ user.id }}</p>
            </div>
          </div>

          <!-- Member Since -->
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 text-gray-400">
              <svg fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-600">Member Since</p>
              <p class="text-lg font-medium text-gray-900">{{ formatMemberSince() }}</p>
            </div>
          </div>
        </div>

        <div v-else class="text-gray-500">
          <p>User data not available.</p>
        </div>
      </div>
    </div>

    <!-- Team Following Stats Section -->
    <div v-if="user" class="mt-6 pt-6 border-t border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Your Football Following</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Teams Followed -->
        <div class="bg-gradient-to-r from-indigo-50 to-blue-50 rounded-lg p-4 border border-indigo-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-3xl font-bold text-indigo-600">{{ followedTeamsCount || 0 }}</p>
              <p class="text-sm text-gray-700 font-medium">Teams Following</p>
            </div>
            <div class="text-indigo-500">
              <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Upcoming Matches -->
        <div class="bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg p-4 border border-green-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-3xl font-bold text-green-600">{{ getUpcomingMatchesCount() }}</p>
              <p class="text-sm text-gray-700 font-medium">Upcoming Matches</p>
            </div>
            <div class="text-green-500">
              <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Predictions Available -->
        <div class="bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg p-4 border border-purple-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-3xl font-bold text-purple-600">{{ getPredictionsCount() }}</p>
              <p class="text-sm text-gray-700 font-medium">Predictions Available</p>
            </div>
            <div class="text-purple-500">
              <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clip-rule="evenodd" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Action Hint -->
      <div class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
        <div class="flex items-center">
          <svg class="w-5 h-5 text-blue-600 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
          </svg>
          <p class="text-sm text-blue-800">
            <span class="font-medium">Pro Tip:</span> Follow teams from the top 5 leagues (Premier League, La Liga, Serie A, Bundesliga, Ligue 1) to get AI-powered match predictions!
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';

// Define the expected structure for the user prop
interface UserInfo {
  id: number;
  name: string;
  email: string;
}

// Define props using defineProps - only team following now
const props = defineProps<{
  user: UserInfo | null;
  followedTeamsCount?: number;
}>();

// Helper function to format member since date
const formatMemberSince = (): string => {
  // In a real app, this would come from user creation date
  // For now, we'll show a placeholder
  return new Date().getFullYear().toString();
};

// Helper function to estimate upcoming matches
const getUpcomingMatchesCount = (): number => {
  // Estimate 2-3 upcoming matches per team on average
  const teamsCount = props.followedTeamsCount || 0;
  return Math.min(teamsCount * 2, 10); // Cap at 10 for display
};

// Helper function to estimate predictions available
const getPredictionsCount = (): number => {
  // Assume 80% of matches from followed teams will have predictions available
  const upcomingMatches = getUpcomingMatchesCount();
  return Math.floor(upcomingMatches * 0.8);
};

// Log props to satisfy TS 'declared but never read' -- props are used in template
console.log('UserInfoCard received props:', props.user);
</script>

<!-- Style block removed as Tailwind classes are used directly -->