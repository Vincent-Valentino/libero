<template>
  <div class="bg-neutral-50 w-full h-full flex flex-col">
    <Hero />
  </div>
  <div class="bg-white w-full h-full pb-20">
    <div class="py-8">
      <MatchPredictor />
    </div>
  </div>
</template>

<script setup lang="ts">
import Hero from './components/Hero.vue';
import MatchPredictor from './components/MatchPredictor.vue';
</script>