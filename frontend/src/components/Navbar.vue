<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const searchQuery = ref('');
const cartItemCount = ref(1);
const isUserMenuOpen = ref(false);

const navLinks = [
  { name: 'Home', path: '/' },
  { name: 'Catalog', path: '/catalog' },
  { name: 'Deals', path: '#' },
  { name: 'Support', path: '#' }
];

const toggleUserMenu = () => {
  isUserMenuOpen.value = !isUserMenuOpen.value;
};

const closeMenu = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.user-menu-container')) {
    isUserMenuOpen.value = false;
  }
};

const goToAuth = () => {
  router.push('/auth');
  isUserMenuOpen.value = false;
};

const handleLogout = async () => {
  await authStore.logout();
  isUserMenuOpen.value = false;
  router.push('/auth');
};

onMounted(() => {
  window.addEventListener('click', closeMenu);
});

onUnmounted(() => {
  window.removeEventListener('click', closeMenu);
});
</script>

<template>
  <nav class="sticky top-0 z-50 w-full h-[72px] bg-white border-b border-gray-100 px-4 md:px-8 flex items-center">
    <div class="w-full flex justify-between items-center gap-4">
      
      <!-- Logo Section -->
      <div class="flex items-center gap-2 flex-shrink-0">
        <div class="text-[#7b00ff]">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
          </svg>
        </div>
        <router-link to="/" class="no-underline cursor-pointer group">
          <span class="text-xl font-bold text-gray-900 tracking-tight group-hover:text-[#7b00ff] transition-colors">Swiftly</span>
        </router-link>
      </div>

      <!-- Search Bar Section -->
      <div class="flex-1 max-w-[400px] relative hidden md:block">
        <div class="absolute inset-y-0 left-4 flex items-center pointer-events-none">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-gray-400">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </div>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search for products..." 
          class="w-full h-10 py-2 pl-10 pr-4 bg-[#f3f4f6] border border-transparent rounded-full text-sm focus:outline-none focus:bg-white focus:border-[#7b00ff]/30 focus:ring-4 focus:ring-[#7b00ff]/5 transition-all placeholder:text-gray-400"
        />
      </div>

      <!-- Navigation & Actions Section -->
      <div class="flex items-center gap-4 lg:gap-8">
        <!-- Navigation Links -->
        <ul class="hidden lg:flex items-center gap-6 list-none p-0 m-0">
          <li v-for="link in navLinks" :key="link.name">
            <router-link 
              :to="link.path" 
              class="text-[14px] font-medium transition-colors no-underline whitespace-nowrap cursor-pointer px-2 py-1 rounded-md hover:bg-gray-50"
              :class="route.path === link.path ? 'text-[#7b00ff]' : 'text-gray-600 hover:text-gray-900'"
            >
              {{ link.name }}
            </router-link>
          </li>
        </ul>

        <!-- Action Icons -->
        <div class="flex items-center gap-2 md:gap-4 border-gray-200 lg:border-l lg:pl-8">
          <!-- Wishlist -->
          <button class="p-2 text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] rounded-full transition-all cursor-pointer" title="Wishlist">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.29 1.5 4.04 3 5.5l7 7Z" />
            </svg>
          </button>

          <!-- Shopping Bag (Cart) -->
          <button class="relative p-2 text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] rounded-full transition-all cursor-pointer" title="Shopping Bag">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4Z"></path>
              <path d="M3 6h18"></path>
              <path d="M16 10a4 4 0 0 1-8 0"></path>
            </svg>
            <span v-if="cartItemCount > 0" class="absolute top-1 right-1 bg-[#7b00ff] text-white text-[9px] font-bold min-w-[15px] h-[15px] rounded-full flex items-center justify-center px-0.5 border border-white">
              {{ cartItemCount }}
            </span>
          </button>

          <!-- User Profile Dropdown -->
          <div class="relative user-menu-container">
            <button 
              @click="toggleUserMenu"
              class="p-2 text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] rounded-full transition-all cursor-pointer flex items-center gap-1" 
              :class="{'bg-gray-50 text-[#7b00ff]': isUserMenuOpen}"
              title="Account"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              <svg v-if="authStore.isAuthenticated" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" class="mt-0.5"><path d="m6 9 6 6 6-6"/></svg>
            </button>

            <!-- Dropdown Menu -->
            <div 
              v-if="isUserMenuOpen"
              class="absolute right-0 mt-2 w-52 bg-white border border-gray-100 rounded-xl shadow-xl py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-200"
            >
              <!-- State: Authenticated -->
              <template v-if="authStore.isAuthenticated">
                <div class="px-4 py-2 border-b border-gray-50 mb-1">
                  <p class="text-[10px] text-gray-400 uppercase font-bold tracking-wider">Signed in as</p>
                  <p class="text-sm font-semibold text-gray-900 truncate">{{ authStore.user?.full_name || 'User' }}</p>
                </div>
                <a href="#" class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] transition-colors no-underline">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                  My Profile
                </a>
                <a href="#" class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] transition-colors no-underline">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l8.84-8.84 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                  Wishlist
                </a>
                <a href="#" class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-[#7b00ff] transition-colors no-underline">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
                  Settings
                </a>
                <div class="h-px bg-gray-100 my-1 mx-2"></div>
                <button 
                  @click="handleLogout"
                  class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors cursor-pointer text-left border-none bg-transparent font-medium"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
                  Sign Out
                </button>
              </template>

              <!-- State: Guest -->
              <template v-else>
                <div class="px-4 py-2 mb-1">
                  <p class="text-xs text-gray-500 leading-relaxed">Sign in to manage your orders and profile.</p>
                </div>
                <div class="px-2">
                  <button 
                    @click="goToAuth"
                    class="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-[#7b00ff] text-white text-sm font-bold rounded-lg hover:bg-[#6a00e0] transition-colors cursor-pointer border-none"
                  >
                    Sign In / Register
                  </button>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
/* Tooltip/Dropdown Animation */
.animate-in {
  animation: fadeInSlide 0.2s ease-out;
}

@keyframes fadeInSlide {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
