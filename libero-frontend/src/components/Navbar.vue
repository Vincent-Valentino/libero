<template>
  <div> <!-- Root element -->
    <!-- Main Navbar Container -->
    <nav class="flex py-8 md:fixed px-6 font-roboto-condensed items-center gap-10 text-sm w-full bg-white relative z-20 border-b border-gray-100">
      <h1 class="text-xl md:text-2xl font-title tracking-tight">Libero</h1>

      <!-- Desktop Navigation -->
      <div class="hidden md:flex gap-8 font-roboto-condensed">
        <router-link to="/" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Home</router-link>
        <router-link to="/leagues/PL" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Premier League</router-link>
        <router-link to="/leagues/PD" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">La Liga</router-link>
        <router-link to="/leagues/SA" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Serie A</router-link>
        <router-link to="/leagues/BL1" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Bundesliga</router-link>
        <router-link to="/leagues/FL1" class="text-xs uppercase font-medium hover:text-gray-500 transition-colors">Ligue 1</router-link>
      </div>

      <!-- Conditional Sign In / Logout Button (Desktop) -->
      <div v-if="!authStore.isAuthenticated" class="hidden md:block ml-auto">
        <router-link to="/auth" class="bg-black rounded-md text-white text-sm font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-gray-800 transition-colors">SIGN IN</router-link>
      </div>
      <div v-else class="hidden md:flex ml-auto gap-2">
        <router-link to="/profile" class="bg-blue-600 rounded-md text-white text-sm font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-blue-700 transition-colors">PROFILE</router-link>
        <button @click="handleLogout" class="bg-red-600 rounded-md text-white text-sm font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-red-700 transition-colors">LOG OUT</button>
      </div>

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
          <router-link @click="closeMenu" to="/leagues/PL" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Premier League</router-link>
          <router-link @click="closeMenu" to="/leagues/PD" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">La Liga</router-link>
          <router-link @click="closeMenu" to="/leagues/SA" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Serie A</router-link>
          <router-link @click="closeMenu" to="/leagues/BL1" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Bundesliga</router-link>
          <router-link @click="closeMenu" to="/leagues/FL1" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Ligue 1</router-link>
          <router-link @click="closeMenu" to="/leagues/CL" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">UCL</router-link>
          <router-link @click="closeMenu" to="/leagues/EL" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">UEL</router-link>
          <router-link @click="closeMenu" to="/player" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Player</router-link>
          <router-link @click="closeMenu" to="/team" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Team</router-link>
          <router-link @click="closeMenu" to="/nations" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Nations</router-link>
          <router-link @click="closeMenu" to="/awards" class="py-1.5 text-xs uppercase font-medium hover:text-gray-500 transition-colors">Awards</router-link>
          <!-- Conditional Sign In / Logout Button (Mobile) -->
          <div v-if="!authStore.isAuthenticated" class="mt-4">
            <router-link @click="closeMenu" to="/auth" class="w-full text-center bg-black rounded-md text-white text-xs font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-gray-800 transition-colors block">SIGN IN</router-link>
          </div>
          <div v-else class="mt-4 flex flex-col gap-2">
            <router-link @click="closeMenu" to="/profile" class="w-full text-center bg-blue-600 rounded-md text-white text-xs font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-blue-700 transition-colors block">PROFILE</router-link>
            <button @click="handleLogout" class="w-full bg-red-600 rounded-md text-white text-xs font-medium tracking-wide px-4 py-1.5 cursor-pointer hover:bg-red-700 transition-colors">LOG OUT</button>
          </div>
        </div>
      </aside> <!-- Closing aside -->
    </Transition>
  </div> <!-- Closing root div -->
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router'; // Import useRouter
import { useAuthStore } from '@/stores/auth'; // Import auth store

const isMenuOpen = ref(false);
const route = useRoute();
const router = useRouter(); // Get router instance
const authStore = useAuthStore(); // Get auth store instance

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

// Logout Handler
const handleLogout = () => {
  authStore.logout(); // Call store action to clear auth state/token
  closeMenu(); // Close mobile menu if open
  router.push({ name: 'Home' }); // Redirect to home page
};
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