<script setup>
import HelloWorld from '../components/HelloWorld.vue'
</script>

<template>
  <main class="main-content">
    <div class="hero-section">
      <a href="https://vite.dev" target="_blank">
        <img src="/vite.svg" class="logo" alt="Vite logo" />
      </a>
      <a href="https://vuejs.org/" target="_blank">
        <img src="../assets/vue.svg" class="logo vue" alt="Vue logo" />
      </a>
    </div>
    <HelloWorld msg="Swiftly + Vite + Vue" />
  </main>
</template>

<style scoped>
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.hero-section {
  margin-bottom: 2rem;
}

.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #7b00ffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
