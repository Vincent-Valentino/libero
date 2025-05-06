<template>
  <div class="flex overflow-x-auto gap-3 md:gap-4 mb-4 scrollbar-hide pb-5">
    <button 
      v-for="tab in tabs" 
      :key="tab.id"
      @click="setActiveTab(tab.id)"
      class="px-3 md:px-4 py-2 rounded-lg font-medium whitespace-nowrap flex-shrink-0 flex items-center gap-2 transition-colors"
      :class="activeTab === tab.id ? 'bg-black text-white' : 'bg-gray-200 hover:bg-gray-300'">
      <img v-if="tab.image" :src="tab.image" :alt="tab.name" class="w-5 h-5" />
      {{ tab.name }}
    </button>
  </div>
</template>

<script setup lang="ts">
// Define active tab state and tab switching logic
import { ref } from 'vue';

const emit = defineEmits(['tab-change']);

const activeTab = ref('all');

const tabs = [
  { id: 'all', name: 'All Leagues', image: null },
  { id: 'pl',   name: 'Premier League',      image: '/Premier League.svg' },
  { id: 'pd',   name: 'La Liga',             image: '/LaLiga.svg' },
  { id: 'sa',   name: 'Serie A',             image: '/Lega Serie A.svg' },
  { id: 'bl1',  name: 'Bundesliga',          image: '/Bundesliga.svg' },
  { id: 'fl1',  name: 'Ligue 1',             image: '/Ligue 1 Uber Eats.svg' },
  { id: 'cl',   name: 'Champions League',    image: '/UCL.svg' },
  { id: 'el',   name: 'Europa League',       image: '/UEL.svg' }
];

const setActiveTab = (tabId: string) => {
  activeTab.value = tabId;
  emit('tab-change', tabId);
};

defineExpose({
  activeTab
});
</script>

<style scoped>
/* Hide scrollbar for clean look */
.scrollbar-hide {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;  /* Chrome, Safari and Opera */
}
</style>