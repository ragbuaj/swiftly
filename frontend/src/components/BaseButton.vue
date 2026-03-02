<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'outline' | 'social';
  loading?: boolean;
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
  loading: false,
  disabled: false,
  type: 'button'
});
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="w-full h-11 px-6 rounded-xl text-sm font-semibold transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
    :class="{
      'bg-[#7b00ff] text-white hover:bg-[#6a00e0] active:scale-[0.98]': variant === 'primary',
      'bg-white border border-gray-200 text-gray-700 hover:bg-gray-50 active:scale-[0.98]': variant === 'outline',
      'bg-white border border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50': variant === 'social'
    }"
  >
    <div v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
    <slot v-else></slot>
  </button>
</template>
