<template>
  <Navbar v-if="showLayout" />
  <router-view />
  <Footer v-if="showLayout" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';

const route = useRoute();

// Show layout only if the route path doesn't start with /auth
const showLayout = computed(() => !route.path.startsWith('/auth'));
</script>