<template>
  <footer class="bg-black text-neutral-50 py-12 px-6 sm:px-10 lg:px-20">
    <div class="container mx-auto flex flex-col md:flex-row justify-between gap-10 mb-8">

      <div class="md:w-1/3">
        <h3 class="font-semibold text-xl mb-4 uppercase tracking-wider">Our Team</h3>
        <div class="flex flex-col gap-1 text-gray-300">
          <p class="text-base hover:text-white transition-colors duration-200">Egip Sinargo</p>
          <p class="text-base hover:text-white transition-colors duration-200">Vincent Valentino</p>
          <p class="text-base hover:text-white transition-colors duration-200">William Zonata</p>
          <p class="text-base hover:text-white transition-colors duration-200">Aleshi Agung Wicaksono</p>
          <p class="text-base hover:text-white transition-colors duration-200">Alvin Liandy</p>
        </div>
      </div>

      <div class="md:w-1/3">
        <h3 class="font-semibold text-xl mb-4 uppercase tracking-wider">Tech Stacks</h3>
        <div class="flex flex-col gap-3 text-gray-300">
          <div class="flex items-center gap-3 hover:text-white transition-colors duration-200">
            <img src="https://cdn.simpleicons.org/vuedotjs/white" alt="Vue.js Logo" class="h-5 w-5 flex-shrink-0 object-contain">
            <span class="text-base">Vue ( Frontend )</span>
          </div>
          <div class="flex items-center gap-3 hover:text-white transition-colors duration-200">
            <img src="https://cdn.simpleicons.org/typescript/white" alt="TypeScript Logo" class="h-5 w-5 flex-shrink-0 object-contain">
            <span class="text-base">TypeScript ( Frontend )</span>
          </div>
          <div class="flex items-center gap-3 hover:text-white transition-colors duration-200">
            <img src="https://cdn.simpleicons.org/go/white" alt="Go Logo" class="h-5 w-5 flex-shrink-0 object-contain">
            <span class="text-base">Golang ( Backend )</span>
          </div>
          <div class="flex items-center gap-3 hover:text-white transition-colors duration-200">
            <img src="https://cdn.simpleicons.org/python/white" alt="Python Logo" class="h-5 w-5 flex-shrink-0 object-contain">
            <span class="text-base">Python ( Machine Learning )</span>
          </div>
          <div class="flex items-center gap-3 hover:text-white transition-colors duration-200">
            <img src="https://cdn.simpleicons.org/postgresql/white" alt="PostgreSQL Logo" class="h-5 w-5 flex-shrink-0 object-contain">
            <span class="text-base">PostgreSQL ( Database )</span>
          </div>
        </div>
      </div>

      <div class="md:w-1/3 text-gray-400">
         <h3 class="font-semibold text-xl mb-4 uppercase tracking-wider">Libero</h3> <p class="text-base">
           Advanced football analytics and predictions.
           </p>
         </div>

    </div>

    <div class="mt-8 pt-8 border-t border-gray-700 flex flex-col sm:flex-row justify-between items-center text-sm text-gray-400">
      <p>&copy; {{ new Date().getFullYear() }} Libero Project. All rights reserved.</p>
    </div>
  </footer>
</template>