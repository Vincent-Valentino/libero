<!-- src/leagues/components/LeagueTable.vue -->
<template>
  <div class="league-table-section p-4">
    <h2 class="text-xl font-semibold mb-3" :style="{ color: themeColor }">League Table</h2>
    <div class="overflow-x-auto shadow rounded">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-8">Pos</th>
            <th scope="col" class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Team</th>
            <th scope="col" class="px-2 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">P</th>
            <th scope="col" class="px-2 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">W</th>
            <th scope="col" class="px-2 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">D</th>
            <th scope="col" class="px-2 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">L</th>
            <th scope="col" class="px-2 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">GD</th>
            <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Pts</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-if="!tableData || tableData.length === 0">
            <td colspan="8" class="px-4 py-3 text-center text-sm text-gray-500 italic">No table data available.</td>
          </tr>
          <tr v-for="row in tableData" :key="row.team.id" class="odd:bg-white even:bg-gray-50 hover:bg-gray-100">
            <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900 text-center">{{ row.position }}</td>
            <td class="px-4 py-2 whitespace-nowrap text-sm font-medium text-gray-900 flex items-center space-x-2">
               <img :src="row.team.logo" :alt="row.team.name" class="h-5 w-5 object-contain">
               <span>{{ row.team.name }}</span>
            </td>
            <td class="px-2 py-2 whitespace-nowrap text-sm text-gray-500 text-center">{{ row.played }}</td>
            <td class="px-2 py-2 whitespace-nowrap text-sm text-gray-500 text-center">{{ row.won }}</td>
            <td class="px-2 py-2 whitespace-nowrap text-sm text-gray-500 text-center">{{ row.drawn }}</td>
            <td class="px-2 py-2 whitespace-nowrap text-sm text-gray-500 text-center">{{ row.lost }}</td>
            <td class="px-2 py-2 whitespace-nowrap text-sm text-gray-500 text-center">{{ row.goalDifference }}</td>
            <td class="px-3 py-2 whitespace-nowrap text-sm font-bold text-gray-900 text-center">{{ row.points }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import type { LeagueTableRow } from "../mockData";

defineProps<{
  tableData: LeagueTableRow[];
  themeColor: string;
}>();
</script>

<style scoped>
/* Add minor styling adjustments if needed */
th, td {
  vertical-align: middle;
}
</style>