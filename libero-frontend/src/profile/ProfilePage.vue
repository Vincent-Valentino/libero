<template>
  <div class="prediction-history-page container mx-auto px-4 py-8 pt-32 max-w-6xl">
    <!-- Header Section -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-4xl font-bold text-gray-900 mb-2">Prediction History</h1>
          <p class="text-gray-600">Track your football match predictions and analyze your prediction patterns</p>
        </div>
        <div class="flex gap-3">
          <button 
            @click="refreshHistory"
            :disabled="predictionStore.isLoading"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refresh
          </button>
          <button 
            v-if="predictionStore.predictions.length > 0"
            @click="clearAllHistory"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H8a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Clear All
          </button>
        </div>
      </div>

      <!-- Statistics Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-600 text-sm font-medium">Total Predictions</p>
              <p class="text-3xl font-bold text-gray-900">{{ stats.total }}</p>
            </div>
            <div class="bg-blue-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-600 text-sm font-medium">Home Wins</p>
              <p class="text-3xl font-bold text-green-600">{{ stats.homeWins }}</p>
              <p class="text-sm text-green-600">{{ (stats.homeWinPercentage || 0).toFixed(1) }}%</p>
            </div>
            <div class="bg-green-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
              </svg>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-600 text-sm font-medium">Draws</p>
              <p class="text-3xl font-bold text-gray-600">{{ stats.draws }}</p>
              <p class="text-sm text-gray-600">{{ (stats.drawPercentage || 0).toFixed(1) }}%</p>
            </div>
            <div class="bg-gray-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
              </svg>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-600 text-sm font-medium">Away Wins</p>
              <p class="text-3xl font-bold text-red-600">{{ stats.awayWins }}</p>
              <p class="text-sm text-red-600">{{ (stats.awayWinPercentage || 0).toFixed(1) }}%</p>
            </div>
            <div class="bg-red-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="predictionStore.isLoading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Loading prediction history...</span>
    </div>

    <!-- Error State -->
    <div v-else-if="predictionStore.error" class="bg-red-50 border border-red-200 rounded-xl p-6 mb-8">
      <div class="flex items-center">
        <svg class="w-6 h-6 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <div class="flex-1">
          <h3 class="text-lg font-medium text-red-800">{{ predictionStore.error.includes('log in') ? 'Authentication Required' : 'Error Loading History' }}</h3>
          <p class="text-red-700">{{ predictionStore.error }}</p>
          <div v-if="predictionStore.error.includes('log in')" class="mt-4">
            <router-link 
              to="/auth" 
              class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
              </svg>
              Log In to Continue
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="predictionStore.predictions.length === 0" class="text-center py-16">
      <div class="bg-gray-50 rounded-2xl p-12 border-2 border-dashed border-gray-300">
        <svg class="w-16 h-16 text-gray-400 mx-auto mb-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <h3 class="text-2xl font-bold text-gray-900 mb-4">No Predictions Yet</h3>
        <p class="text-gray-600 mb-6 max-w-md mx-auto">
          You haven't made any match predictions yet. Head to the home page to start predicting match outcomes!
        </p>
        <router-link 
          to="/" 
          class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          Make Your First Prediction
        </router-link>
      </div>
    </div>

    <!-- Predictions List -->
    <div v-else class="space-y-4">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900">Recent Predictions</h2>
        <div class="text-sm text-gray-600">
          {{ predictionStore.predictions.length }} prediction{{ predictionStore.predictions.length !== 1 ? 's' : '' }}
        </div>
      </div>

      <div class="grid gap-4">
        <div 
          v-for="prediction in predictionStore.predictions" 
          :key="prediction.id"
          class="bg-white rounded-xl shadow-lg border border-gray-100 p-6 hover:shadow-xl transition-shadow"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-4">
              <div class="text-center">
                <div class="text-lg font-bold text-gray-900">{{ prediction.homeTeam }}</div>
                <div class="text-sm text-gray-600">{{ getLeagueDisplayName(prediction.homeLeague) }}</div>
              </div>
              
              <div class="text-center mx-6">
                <div class="text-3xl font-bold text-indigo-600 mb-1">
                  {{ prediction.predictedHomeScore ?? 0 }} - {{ prediction.predictedAwayScore ?? 0 }}
                </div>
                <div class="text-sm font-semibold text-indigo-700">{{ prediction.predictedResult ?? 'Unknown' }}</div>
              </div>
              
              <div class="text-center">
                <div class="text-lg font-bold text-gray-900">{{ prediction.awayTeam }}</div>
                <div class="text-sm text-gray-600">{{ getLeagueDisplayName(prediction.awayLeague) }}</div>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <div class="text-right">
                <div class="text-sm text-gray-600">Predicted on</div>
                <div class="text-sm font-medium text-gray-900">{{ formatDate(prediction.createdAt) }}</div>
              </div>
              <button 
                @click="deletePrediction(prediction.id)"
                class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H8a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Detailed Predictions -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 pt-4 border-t border-gray-100">
            <div>
              <h4 class="text-sm font-semibold text-gray-700 mb-3">Expected Goals</h4>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">{{ prediction.homeTeam }}</span>
                <span class="font-medium">{{ (prediction.expectedHomeGoals ?? 0).toFixed(1) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600">{{ prediction.awayTeam }}</span>
                <span class="font-medium">{{ (prediction.expectedAwayGoals ?? 0).toFixed(1) }}</span>
              </div>
            </div>

            <div>
              <h4 class="text-sm font-semibold text-gray-700 mb-3">Win Probabilities</h4>
              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-600">Home Win</span>
                  <span class="font-medium text-green-600">{{ ((prediction.homeWinProbability ?? 0) * 100).toFixed(1) }}%</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-600">Draw</span>
                  <span class="font-medium text-gray-600">{{ ((prediction.drawProbability ?? 0) * 100).toFixed(1) }}%</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-600">Away Win</span>
                  <span class="font-medium text-red-600">{{ ((prediction.awayWinProbability ?? 0) * 100).toFixed(1) }}%</span>
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
import { onMounted, computed, ref } from 'vue';
import { usePredictionStore, type PredictionStatistics } from '@/stores/prediction';

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

// Local stats state for API statistics
const apiStats = ref<PredictionStatistics | null>(null);
const statsLoading = ref(false);

// Computed properties
const stats = computed(() => {
  const currentStats = apiStats.value || predictionStore.getStats();
  
  // Ensure all required properties exist with defaults
  return {
    total: currentStats?.total ?? 0,
    homeWins: currentStats?.homeWins ?? 0,
    draws: currentStats?.draws ?? 0,
    awayWins: currentStats?.awayWins ?? 0,
    homeWinPercentage: currentStats?.homeWinPercentage ?? 0,
    drawPercentage: currentStats?.drawPercentage ?? 0,
    awayWinPercentage: currentStats?.awayWinPercentage ?? 0,
  };
});

// Methods
const getLeagueDisplayName = (leagueCode: string): string => {
  return leagueDisplayNames[leagueCode] || leagueCode;
};

const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric',
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const refreshHistory = async () => {
  await Promise.all([
    predictionStore.fetchPredictions(),
    fetchApiStats()
  ]);
};

const fetchApiStats = async () => {
  statsLoading.value = true;
  try {
    apiStats.value = await predictionStore.fetchStatistics();
  } catch (error) {
    console.warn('Failed to fetch API statistics, using local calculation');
  } finally {
    statsLoading.value = false;
  }
};

const deletePrediction = async (id: number) => {
  if (confirm('Are you sure you want to delete this prediction?')) {
    await predictionStore.deletePrediction(id);
    // Refresh stats after deletion
    await fetchApiStats();
  }
};

const clearAllHistory = async () => {
  if (confirm('Are you sure you want to clear all prediction history? This action cannot be undone.')) {
    await predictionStore.clearAllPredictions();
    apiStats.value = null; // Reset stats
  }
};

// Lifecycle
onMounted(async () => {
  await Promise.all([
    predictionStore.fetchPredictions(),
    fetchApiStats()
  ]);
});
</script>

<!-- Style block removed as Tailwind classes are used directly -->