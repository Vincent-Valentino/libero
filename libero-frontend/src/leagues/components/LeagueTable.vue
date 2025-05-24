<!-- src/leagues/components/LeagueTable.vue -->
<template>
  <div class="league-table-section p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-bold" :style="{ color: themeColor }">League Table</h2>
      <div v-if="isLoading" class="text-sm text-gray-500">Refreshing...</div>
    </div>

    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-4">
      <p class="text-red-700">{{ error }}</p>
    </div>

    <div class="overflow-x-auto">
      <table class="min-w-full bg-white">
        <thead>
          <tr class="text-sm border-b">
            <th class="text-left p-3">#</th>
            <th class="text-left p-3">Team</th>
            <th class="text-center p-3">MP</th>
            <th class="text-center p-3">W</th>
            <th class="text-center p-3">D</th>
            <th class="text-center p-3">L</th>
            <th class="text-center p-3">GD</th>
            <th class="text-center p-3">PTS</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="!isLoading && !error && tableData.length > 0">
            <tr v-for="row in tableData" :key="row.position" class="text-sm border-b hover:bg-gray-50">
              <td class="p-3">{{ row.position }}</td>
              <td class="p-3 flex items-center">
                <img
                  :src="imgSrc(row.team.logo)"
                  :alt="row.team.name"
                  class="w-6 h-6 mr-2"
                  @error="onImgError($event, row.team)"
                >
                {{ row.team.name }}
              </td>
              <td class="text-center p-3">{{ row.played }}</td>
              <td class="text-center p-3">{{ row.won }}</td>
              <td class="text-center p-3">{{ row.drawn }}</td>
              <td class="text-center p-3">{{ row.lost }}</td>
              <td class="text-center p-3">{{ row.goalDifference }}</td>
              <td class="text-center p-3 font-bold">{{ row.points }}</td>
            </tr>
          </template>
          <tr v-else-if="isLoading">
            <td colspan="8" class="text-center p-4">
              <div class="flex justify-center space-x-2">
                <div class="animate-spin rounded-full h-5 w-5 border-2" :style="{ borderColor: `${themeColor} transparent` }"></div>
                <span class="text-gray-500">Loading table data...</span>
              </div>
            </td>
          </tr>
          <tr v-else-if="error">
            <td colspan="8" class="text-center p-4 text-gray-500">
              Unable to load table data
            </td>
          </tr>
          <tr v-else>
            <td colspan="8" class="text-center p-4 text-gray-500">
              No table data available
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, ref } from "vue";
import type { LeagueTableRow } from "../mockData";

defineProps<{
  tableData: LeagueTableRow[];
  themeColor: string;
  isLoading?: boolean;
  error?: string;
}>();

const fallbackTeamLogo = '/fallback-team.png';
const erroredTeams = ref<Record<string, boolean>>({});
function imgSrc(logo: string) {
  return erroredTeams.value[logo] ? fallbackTeamLogo : logo;
}
function onImgError(event: Event, team: any) {
  if (!erroredTeams.value[team.logo]) {
    erroredTeams.value[team.logo] = true;
    const target = event.target as HTMLImageElement | null;
    if (target) target.src = fallbackTeamLogo;
  }
}
</script>

<style scoped>
/* Add minor styling adjustments if needed */
th, td {
  vertical-align: middle;
}
</style>