<script setup lang="ts">
defineProps<{
  label?: string;
  modelValue: string;
  type?: string;
  placeholder?: string;
  error?: string;
  id?: string;
}>();

defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();
</script>

<template>
  <div class="w-full mb-4">
    <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1.5">{{ label }}</label>
    <div class="relative flex items-center p-0.5 -m-0.5">
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        :placeholder="placeholder"
        class="w-full h-11 px-4 bg-white border rounded-xl text-sm transition-all outline-none pr-10"
        :class="error ? 'border-red-500 focus:ring-4 focus:ring-red-500/10' : 'border-gray-200 focus:border-[#7b00ff] focus:ring-4 focus:ring-[#7b00ff]/10'"
      />
      <!-- Slot for right-side icons/indicators -->
      <div class="absolute right-3 flex items-center justify-center">
        <slot name="right"></slot>
      </div>
    </div>
    <p v-if="error" class="mt-1.5 text-[11px] font-medium text-red-500 animate-in fade-in slide-in-from-top-1">{{ error }}</p>
  </div>
</template>
