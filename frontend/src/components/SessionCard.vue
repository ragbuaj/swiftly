<script setup lang="ts">
import type { Session } from '../types';
import { getDeviceInfo, formatLastActive } from '../utils/device';
import BaseButton from './BaseButton.vue';

defineProps<{
  session: Session;
  isRevoking: boolean;
}>();

defineEmits<{
  (e: 'revoke', id: string): void;
}>();
</script>

<template>
  <div 
    class="group flex flex-col lg:flex-row lg:items-center justify-between p-5 rounded-2xl border transition-all duration-200"
    :class="session.is_current ? 'border-[#7b00ff]/20 bg-[#7b00ff]/[0.02]' : 'border-gray-100 bg-white hover:border-gray-200 hover:shadow-md hover:shadow-gray-100'"
  >
    <div class="flex items-start sm:items-center gap-4 sm:gap-5">
      <!-- Device Icon -->
      <div 
        class="w-12 h-12 sm:w-14 sm:h-14 rounded-xl flex items-center justify-center border shadow-sm flex-shrink-0 transition-transform duration-300 group-hover:scale-105"
        :class="session.is_current ? 'bg-white border-[#7b00ff]/20 text-[#7b00ff]' : 'bg-gray-50 border-gray-100 text-gray-400'"
      >
        <svg v-if="session.device_type === 'Mobile'" xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"></rect><line x1="12" y1="18" x2="12.01" y2="18"></line></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
      </div>

      <!-- Device Meta -->
      <div class="min-w-0">
        <div class="flex flex-wrap items-center gap-2">
          <span class="font-bold text-gray-900 text-base sm:text-lg truncate">
            {{ getDeviceInfo(session.user_agent).browser }} on {{ getDeviceInfo(session.user_agent).os }}
          </span>
          <span v-if="session.is_current" class="text-[10px] bg-green-500 text-white font-bold px-2 py-0.5 rounded-full uppercase tracking-tighter">
            Current
          </span>
        </div>
        
        <div class="flex flex-col sm:flex-row sm:items-center gap-x-4 gap-y-1 mt-1">
          <div class="text-xs sm:text-sm text-gray-600 font-medium flex items-center gap-1.5 truncate">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-gray-400"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
            {{ session.location }}
          </div>
          <div class="text-[10px] sm:text-xs text-gray-400 font-mono">
            IP: {{ session.ip_address.split(':')[0] }}
          </div>
        </div>

        <div class="text-[10px] sm:text-xs text-gray-400 mt-2 flex items-center gap-1.5 italic">
          <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          Last seen {{ formatLastActive(session.last_active_at) }}
        </div>
      </div>
    </div>
    
    <!-- Action Button -->
    <div class="mt-4 lg:mt-0">
      <BaseButton 
        v-if="!session.is_current"
        variant="outline"
        size="sm"
        @click="$emit('revoke', session.id)"
        class="w-full lg:w-auto text-red-500 border-red-100 hover:bg-red-50 font-bold px-5 h-9 rounded-xl text-xs"
        :is-loading="isRevoking"
        loading-text="Revoking..."
      >
        Logout Device
      </BaseButton>
      <div v-else class="hidden lg:flex items-center gap-1.5 text-xs font-bold text-green-600 px-4 py-2 bg-green-50 rounded-xl border border-green-100">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 11"></polyline></svg>
        Active Now
      </div>
    </div>
  </div>
</template>
