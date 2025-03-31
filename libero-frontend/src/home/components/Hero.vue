<template>
  <div class="w-11/12 mx-auto flex flex-col md:flex-row h-full pt-30">
    <div class="w-full md:w-1/2 p-5 h-full flex flex-col justify-center md:justify-start md:pt-50">
      <h1 class="text-3xl md:text-4xl font-title tracking-tight leading-tight">AI Football Match Predictor</h1>
      <p class="mt-3 text-sm text-gray-600 leading-relaxed">Available for Premier League, La Liga, Serie A and more leagues worldwide. Get accurate predictions for upcoming matches.</p>
      <div class="mt-6 flex flex-wrap gap-3">
        <button class="bg-black rounded-md text-xs text-white font-medium tracking-wider uppercase px-4 py-2 cursor-pointer hover:bg-gray-800 transition-colors">Premier League</button>
        <button class="bg-black rounded-md text-xs text-white font-medium tracking-wider uppercase px-4 py-2 cursor-pointer hover:bg-gray-800 transition-colors">La Liga</button>
        <button class="bg-black rounded-md text-xs text-white font-medium tracking-wider uppercase px-4 py-2 cursor-pointer hover:bg-gray-800 transition-colors">Serie A</button>
        <button class="bg-gray-100 rounded-md text-xs text-black font-medium tracking-wider uppercase px-4 py-2 cursor-pointer hover:bg-gray-200 transition-colors">More</button>
      </div>
    </div>
    <div class="w-full md:w-1/2 flex items-center justify-center pt-5 px-5">
      <div class="rounded-xl overflow-hidden w-full h-full">
        <img src="/uclball.png" alt="A Ball" class="w-full h-full object-contain">
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

</script>