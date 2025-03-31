<template>
  <div class="bg-neutral-50 w-full h-full flex flex-col">
    <Hero />
  </div>
  <div class="bg-white w-full h-full pb-20">
    <UpcomingMatches />
  </div>
  <div class="bg-neutral-50 w-full h-full">
    <div class="flex flex-col p-10 md:p-20 w-full h-full">
      <h1 class="font-title text-5xl mb-4">The MVPs</h1>
      <PlayerList />
    </div>
  </div>
</template>

<script setup lang="ts">
import Hero from './components/Hero.vue';
import UpcomingMatches from './components/UpcomingMatches.vue';
import PlayerList from './components/PlayerList.vue';
</script>