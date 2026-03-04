<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../../components/BaseButton.vue';
import BaseModal from '../../components/BaseModal.vue';
import SessionCard from '../../components/SessionCard.vue';

const authStore = useAuthStore();
const isRevoking = ref<string | null>(null);

// Modal state
const modalConfig = ref({
  open: false,
  title: '',
  description: '',
  onConfirm: () => {},
  isLoading: false,
});

onMounted(() => {
  authStore.fetchSessions();
});

const handleRevokeSession = (sessionID: string) => {
  modalConfig.value = {
    open: true,
    title: 'Logout device?',
    description: 'Are you sure you want to log out this device? You will need to sign in again to access your account from this device.',
    isLoading: false,
    onConfirm: async () => {
      modalConfig.value.isLoading = true;
      isRevoking.value = sessionID;
      try {
        await authStore.revokeSession(sessionID);
        modalConfig.value.open = false;
      } finally {
        isRevoking.value = null;
        modalConfig.value.isLoading = false;
      }
    },
  };
};

const handleRevokeAllOthers = () => {
  modalConfig.value = {
    open: true,
    title: 'Logout all other devices?',
    description: 'This will log you out from all your other active sessions. This action cannot be undone.',
    isLoading: false,
    onConfirm: async () => {
      modalConfig.value.isLoading = true;
      try {
        const others = authStore.activeSessions.filter(s => !s.is_current);
        for (const s of others) {
          await authStore.revokeSession(s.id);
        }
        modalConfig.value.open = false;
      } finally {
        modalConfig.value.isLoading = false;
      }
    },
  };
};
</script>

<template>
  <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Logout Confirmation Modal -->
    <BaseModal
      v-model:open="modalConfig.open"
      :title="modalConfig.title"
      :description="modalConfig.description"
      :is-loading="modalConfig.isLoading"
      confirm-text="Logout Device"
      variant="danger"
      @confirm="modalConfig.onConfirm"
    />
    <!-- Header Section -->
    <div class="p-6 md:p-8 border-b border-gray-100 bg-gradient-to-r from-white to-gray-50/30">
      <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
        <div class="flex-1 min-w-0">
          <h1 class="text-xl md:text-2xl font-bold text-gray-900 tracking-tight">Devices & Security</h1>
          <p class="text-sm text-gray-500 mt-1 max-w-lg leading-relaxed">
            Manage your active sessions and logout from other devices if you notice any suspicious activity.
          </p>
        </div>
        <BaseButton 
          v-if="authStore.activeSessions.length > 1"
          variant="outline" 
          size="sm"
          class="inline-flex items-center self-start text-red-600 border-red-100 hover:bg-red-50 font-semibold h-8 px-3 rounded-lg text-[12px] shrink-0"
          @click="handleRevokeAllOthers"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="mr-1"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          Logout all other sessions
        </BaseButton>
      </div>
    </div>

    <!-- Sessions List -->
    <div class="p-6 md:p-8">
      <!-- Loading State -->
      <div v-if="authStore.isLoading && authStore.activeSessions.length === 0" class="space-y-4">
        <div v-for="i in 2" :key="i" class="h-24 bg-gray-50 animate-pulse rounded-2xl border border-gray-100"></div>
      </div>

      <!-- Content -->
      <div v-else class="grid grid-cols-1 gap-4">
        <SessionCard 
          v-for="session in authStore.activeSessions" 
          :key="session.id"
          :session="session"
          :is-revoking="isRevoking === session.id"
          @revoke="handleRevokeSession"
        />
      </div>

      <!-- Empty State -->
      <div v-if="!authStore.isLoading && authStore.activeSessions.length === 0" class="text-center py-20 bg-gray-50/30 rounded-3xl border border-dashed border-gray-200">
        <div class="w-20 h-20 bg-white rounded-2xl flex items-center justify-center mx-auto mb-6 border border-gray-100 shadow-sm text-gray-300 transform -rotate-12">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
        </div>
        <h3 class="text-lg font-bold text-gray-900">No active devices</h3>
        <p class="text-gray-500 mt-1">We couldn't find any active sessions for your account.</p>
      </div>
    </div>

    <!-- Security Info Section -->
    <div class="mx-6 md:mx-8 mb-8 p-5 bg-slate-900 rounded-2xl text-white flex flex-col sm:flex-row gap-4 items-center overflow-hidden relative">
      <div class="absolute top-0 right-0 p-4 opacity-10 transform translate-x-4 -translate-y-4 hidden md:block">
        <svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
      </div>
      <div class="w-10 h-10 rounded-full bg-[#7b00ff] flex items-center justify-center flex-shrink-0 shadow-lg shadow-[#7b00ff]/20">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
      </div>
      <div class="flex-1 text-center sm:text-left relative z-10">
        <h4 class="text-base font-bold">Keep your account safe</h4>
        <p class="text-[11px] text-slate-400 mt-0.5 leading-relaxed">
          Logout from devices you no longer use or suspect are compromised.
        </p>
      </div>
      <BaseButton variant="outline" size="sm" class="bg-slate-800 border-slate-700 text-white hover:bg-slate-700 rounded-lg px-4 h-8 text-[10px] font-bold uppercase tracking-wider relative z-10 shrink-0">
        Change Password
      </BaseButton>
    </div>
  </div>
</template>
