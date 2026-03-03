<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

defineProps<{
  label?: string;
  modelValue: string | number;
  type?: string;
  placeholder?: string;
  error?: string;
  id?: string;
  required?: boolean;
  disabled?: boolean;
  readonly?: boolean;
}>();

defineEmits<{
  (e: 'update:modelValue', value: string | number): void;
}>();
</script>

<template>
  <div class="w-full mb-4">
    <Label 
      v-if="label" 
      :for="id" 
      :class="cn('block mb-1.5', error && 'text-destructive', (disabled && !readonly) && 'opacity-70')"
    >
      {{ label }} <span v-if="required" class="text-destructive">*</span>
    </Label>
    
    <div class="relative flex items-center">
      <Input
        :id="id"
        :type="type"
        :model-value="modelValue"
        @update:model-value="$emit('update:modelValue', $event)"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :class="cn(
          'h-11 rounded-xl pr-10', 
          error && 'border-destructive focus-visible:ring-destructive/20',
          type === 'date' && 'date-input',
          readonly && 'bg-transparent border-gray-100 shadow-none focus-visible:ring-0 focus-visible:border-gray-100 cursor-default'
        )"
        :aria-invalid="!!error"
      />
      
      <!-- Slot for right-side icons/indicators -->
      <div class="absolute right-3 flex items-center justify-center">
        <slot name="right"></slot>
      </div>
    </div>
    
    <p 
      v-if="error" 
      class="mt-1.5 text-[11px] font-medium text-destructive animate-in fade-in slide-in-from-top-1"
    >
      {{ error }}
    </p>
  </div>
</template>

<style scoped>
/* Styling khusus untuk input tanggal agar selaras dengan tema */
.date-input::-webkit-calendar-picker-indicator {
  cursor: pointer;
  opacity: 0.5;
  transition: opacity 0.2s ease-in-out;
}

.date-input::-webkit-calendar-picker-indicator:hover {
  opacity: 1;
}

/* Warna placeholder untuk input date jika belum diisi */
.date-input:invalid {
  color: #9ca3af;
}
</style>