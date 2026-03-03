<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

const route = useRoute();
const authStore = useAuthStore();

const sidebarLinks = [
  { name: 'Profile Settings', path: '/profile/settings', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
  { name: 'My Orders', path: '/profile/orders', icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z' },
  { name: 'Addresses', path: '/profile/addresses', icon: 'M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z M15 11a3 3 0 11-6 0 3 3 0 016 0z' },
];

</script>

<template>
  <div class="min-h-screen bg-[#f8fafc] font-sans text-gray-900">
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="flex flex-col lg:flex-row gap-8">
        
        <!-- Sidebar Navigation -->
        <aside class="w-full lg:w-64 flex-shrink-0">
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <!-- User Summary (Sidebar Header) -->
            <div class="p-6 border-b border-gray-100 flex items-center gap-4">
              <img 
                :src="authStore.user?.avatar_url || 'https://api.dicebear.com/7.x/initials/svg?seed=' + (authStore.user?.full_name || 'User')" 
                alt="User Avatar" 
                class="w-12 h-12 rounded-full object-cover border border-gray-200 bg-gray-50"
              />
              <div class="overflow-hidden">
                <h3 class="font-bold text-gray-900 truncate">{{ authStore.user?.full_name || 'Loading...' }}</h3>
                <p class="text-xs text-gray-500 truncate">{{ authStore.user?.email || 'Loading...' }}</p>
              </div>
            </div>

            <!-- Navigation Links -->
            <nav class="p-4 flex flex-col gap-1">
              <router-link 
                v-for="link in sidebarLinks" 
                :key="link.path"
                :to="link.path"
                class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-colors no-underline"
                :class="[
                  route.path === link.path || route.path.startsWith(link.path) 
                    ? 'bg-[#7b00ff]/10 text-[#7b00ff]' 
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                ]"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 opacity-75" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" :d="link.icon" />
                </svg>
                {{ link.name }}
              </router-link>
            </nav>
          </div>
        </aside>

        <!-- Main Content Area -->
        <div class="flex-1 min-w-0">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>

      </div>
    </main>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
