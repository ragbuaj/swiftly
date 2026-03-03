<script setup lang="ts">
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

defineProps<{
  modelValue: boolean;
  label: string;
  id?: string;
  disabled?: boolean;
}>();

defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>();
</script>

<template>
  <div class="flex items-center gap-3">
    <Switch 
      :id="id"
      :model-value="modelValue"
      :disabled="disabled"
      @update:model-value="$emit('update:modelValue', $event)"
    />
    <Label v-if="label" :for="id" :class="cn('text-sm font-medium text-gray-700 cursor-pointer', disabled && 'opacity-70 cursor-not-allowed')">
      {{ label }}
    </Label>
  </div>
</template>
