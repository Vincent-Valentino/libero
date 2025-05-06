<template>
  <div class="user-info-card border rounded-lg p-4 shadow-sm bg-white">
    <h2 class="text-xl font-semibold mb-3">User Information</h2>
    <div v-if="user">
      <p><strong>Name:</strong> {{ user.name }}</p>
      <p><strong>Email:</strong> {{ user.email }}</p>
      <!-- Add other user details here if needed -->
    </div>
    <div v-else>
      <p>User data not available.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';

// Define the expected structure for the user prop
interface UserInfo {
  id: number;
  name: string;
  email: string;
}

// Define props using defineProps
const props = defineProps<{
  user: UserInfo | null; // Accept UserInfo object or null
}>();

// Log props to satisfy TS 'declared but never read' -- props are used in template
console.log('UserInfoCard received props:', props.user);

</script>

<!-- Style block removed as Tailwind classes are used directly -->