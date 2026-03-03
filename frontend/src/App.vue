<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import Navbar from './components/Navbar.vue'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const goToLogin = () => {
  authStore.setSessionExpired(false)
  router.push('/auth')
}
</script>

<template>
  <div class="app-container">
    <!-- Only show Navbar if not on Auth page AND session is still valid -->
    <Navbar v-if="!route.meta.hideNavbar && !authStore.isSessionExpired" />
    
    <!-- 
      Render Current Route Content ONLY if session is valid.
      This globally prevents crashes in child components when user becomes null.
    -->
    <router-view v-if="!authStore.isSessionExpired" />

    <!-- Global Session Expired Popup -->
    <div v-if="authStore.isSessionExpired" class="session-modal-overlay">
      <div class="session-modal">
        <h3>Session Expired</h3>
        <p>Your session has ended. Please log in again to continue accessing your account.</p>
        <button @click="goToLogin" class="login-btn">Login Now</button>
      </div>
    </div>
  </div>
</template>

<style>
/* Global reset or layout */
body {
  margin: 0;
  padding: 0;
  font-family: 'Inter', sans-serif;
  background-color: #f8fafc;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Modal Styling */
.session-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.session-modal {
  background: var(--color-background);
  padding: 2rem;
  border-radius: var(--radius-xl);
  max-width: 400px;
  width: 90%;
  text-align: center;
  box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
  border: 1px solid var(--color-border);
}

.session-modal h3 {
  margin-top: 0;
  color: var(--color-foreground);
  font-size: 1.25rem;
  font-weight: 700;
}

.session-modal p {
  color: var(--color-muted-foreground);
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.login-btn {
  background: var(--color-primary);
  color: var(--color-primary-foreground);
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
  width: 100%;
}

.login-btn:hover {
  opacity: 0.9;
}
</style>
