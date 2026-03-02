<script setup>
import { useRoute } from 'vue-router'
import Navbar from './components/Navbar.vue'

const route = useRoute()
</script>

<template>
  <div class="app-container">
    <!-- Only show Navbar if not on Auth page -->
    <Navbar v-if="!route.meta.hideNavbar" />
    
    <!-- Render Current Route Content -->
    <router-view />
  </div>
</template>

<style>
/* Global reset or layout */
body {
  margin: 0;
  padding: 0;
  font-family: 'Inter', sans-serif;
  background-color: #fafafa;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>
