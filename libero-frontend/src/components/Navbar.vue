<template>
  <div> <!-- Root element -->
    <!-- Main Navbar Container -->
    <nav class="flex py-8 md:fixed px-6 font-roboto-condensed items-center gap-10 text-sm w-full bg-white relative z-20 border-b border-gray-100">
      <h1 class="text-xl md:text-2xl font-title tracking-tight">Libero</h1>

      <!-- Desktop Navigation -->
      <div class="hidden md:flex gap-8 font-roboto-condensed">
        <router-link to="/" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Home</router-link>
        <router-link to="/premier-league" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Premier League</router-link>
        <router-link to="/la-liga" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">La Liga</router-link>
        <router-link to="/serie-a" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Serie A</router-link>
        <router-link to="/bundesliga" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Bundesliga</router-link>
        <router-link to="/ligue-1" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Ligue 1</router-link>
        <router-link to="/ucl" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">UCL</router-link>
        <router-link to="/uel" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">UEL</router-link>
        <router-link to="/player" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Player</router-link>
        <router-link to="/team" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Team</router-link>
        <router-link to="/nations" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Nations</router-link>
        <router-link to="/awards" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Awards</router-link>
      </div>

      <button class="hidden md:block ml-auto bg-black rounded-md text-white text-sm font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-gray-800 transition-colors">SIGN IN</button>

      <!-- Mobile Menu Button -->
      <button
        type="button"
        class="md:hidden ml-auto flex flex-col gap-1 p-2 z-30"
        @click="toggleMenu"
        aria-label="Toggle menu">
        <span
          class="block w-5 h-0.5 bg-black transition-transform duration-300"
          :class="{'rotate-45 translate-y-1.5': isMenuOpen}">
        </span>
        <span
          class="block w-5 h-0.5 bg-black transition-opacity duration-300"
          :class="{'opacity-0': isMenuOpen}">
        </span>
        <span
          class="block w-5 h-0.5 bg-black transition-transform duration-300"
          :class="{'-rotate-45 -translate-y-1.5': isMenuOpen}">
        </span>
      </button> <!-- Correctly closed button tag -->
    </nav> <!-- Closing nav -->

    <!-- Mobile Menu Overlay -->
    <Transition name="fade">
      <div
        v-if="isMenuOpen"
        class="fixed inset-0 bg-black bg-opacity-50 z-10"
        @click="closeMenu">
      </div>
    </Transition>

    <!-- Mobile Menu -->
    <Transition name="slide-right">
      <aside
        v-if="isMenuOpen"
        class="fixed top-0 right-0 h-full w-64 bg-white shadow-lg z-20" >
        <div class="p-6 flex flex-col gap-3">
          <router-link @click="closeMenu" to="/" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Home</router-link>
          <router-link @click="closeMenu" to="/premier-league" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Premier League</router-link>
          <router-link @click="closeMenu" to="/la-liga" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">La Liga</router-link>
          <router-link @click="closeMenu" to="/serie-a" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Serie A</router-link>
          <router-link @click="closeMenu" to="/bundesliga" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Bundesliga</router-link>
          <router-link @click="closeMenu" to="/ligue-1" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Ligue 1</router-link>
          <router-link @click="closeMenu" to="/ucl" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">UCL</router-link>
          <router-link @click="closeMenu" to="/uel" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">UEL</router-link>
          <router-link @click="closeMenu" to="/player" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Player</router-link>
          <router-link @click="closeMenu" to="/team" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Team</router-link>
          <router-link @click="closeMenu" to="/nations" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Nations</router-link>
          <router-link @click="closeMenu" to="/awards" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Awards</router-link>
          <button class="mt-4 w-full bg-black rounded-md text-white text-xs font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-gray-800 transition-colors">SIGN IN</button>
        </div>
      </aside> <!-- Closing aside -->
    </Transition>
  </div> <!-- Closing root div -->
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

const isMenuOpen = ref(false);
const route = useRoute();

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value;
  setBodyScroll(isMenuOpen.value);
};

const closeMenu = () => {
  isMenuOpen.value = false;
  setBodyScroll(false);
};

const setBodyScroll = (isOpen: boolean) => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = isOpen ? 'hidden' : '';
  }
};

const handleEscKey = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isMenuOpen.value) {
    closeMenu();
  }
};

if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handleEscKey);

  watch(route, () => {
    if (isMenuOpen.value) {
      closeMenu();
    }
  });

  onUnmounted(() => {
    window.removeEventListener('keydown', handleEscKey);
    setBodyScroll(false);
  });
}
</script>

<style scoped>
/* Style for the active router link */
.router-link-exact-active {
  text-decoration: underline;
  text-underline-offset: 4px;
  color: black;
}

/* Transitions for Mobile Menu Slide */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(100%);
}
.slide-right-enter-to,
.slide-right-leave-from {
  transform: translateX(0);
}

/* Transitions for Overlay Fade */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease-in-out;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.fade-enter-to,
.fade-leave-from {
  opacity: 1;
}
</style>