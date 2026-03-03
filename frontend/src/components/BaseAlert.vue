<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
}>();

const alertClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-50 text-green-700 border-green-200';
    case 'error':
      return 'bg-red-50 text-red-700 border-red-200';
    case 'warning':
      return 'bg-yellow-50 text-yellow-700 border-yellow-200';
    case 'info':
    default:
      return 'bg-blue-50 text-blue-700 border-blue-200';
  }
});
</script>

<template>
  <div v-if="message" class="p-4 rounded-xl border flex items-center gap-3 transition-all" :class="alertClasses">
    <!-- Success Icon -->
    <svg v-if="type === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
    </svg>
    <!-- Error/Warning Icon -->
    <svg v-else-if="type === 'error' || type === 'warning'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
    </svg>
    <!-- Info Icon -->
    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
    </svg>
    
    <span class="text-sm font-medium">{{ message }}</span>
  </div>
</template>