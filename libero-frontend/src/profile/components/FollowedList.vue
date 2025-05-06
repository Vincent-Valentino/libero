<template>
  <div class="followed-list">
    <ul v-if="items && items.length > 0" class="list-disc pl-5 space-y-2">
      <li v-for="item in items" :key="item.id" class="flex justify-between items-center">
        <span>{{ item.name }}</span>
        <button
          @click="handleRemoveClick(item.id)"
          class="ml-4 px-2 py-1 text-xs font-medium text-red-700 bg-red-100 rounded hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          aria-label="`Remove ${item.name}`"
        >
          Remove
        </button>
      </li>
    </ul>
    <!-- Optional: Message when the list is empty is handled in the parent (PreferencesManager) -->
    <!-- <p v-else class="text-gray-500 italic">
      No {{ itemType }}s followed yet.
    </p> -->
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

// Interface for items (can represent Team or Player)
interface ListItem {
  id: number;
  name: string;
}

// Define props
const props = defineProps<{
  items: ListItem[];
  itemType: 'team' | 'player' | 'competition'; // Added 'competition'
}>();

// Define emits
const emit = defineEmits<{
  (e: 'remove-item', id: number): void;
}>();

// Method to handle remove button click
const handleRemoveClick = (itemId: number) => {
  console.log(`Emitting remove-${props.itemType} for ID:`, itemId); // Log for debugging
  emit('remove-item', itemId);
};

// Log props to satisfy TS 'declared but never read' -- props are used in template
console.log(`FollowedList (${props.itemType}) received items:`, props.items);

</script>

<!-- Style block removed as Tailwind classes are used directly -->