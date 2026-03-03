<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface Option {
  value: string | number;
  label: string;
}

defineProps<{
  label?: string;
  modelValue: string | number;
  options: Option[];
  placeholder?: string;
  error?: string;
  id?: string;
  disabled?: boolean;
  readonly?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void;
}>();

const handleUpdate = (value: string) => {
  emit('update:modelValue', value);
};
</script>

<template>
  <div class="w-full mb-4">
    <Label 
      v-if="label" 
      :for="id" 
      :class="cn('block mb-1.5', error && 'text-destructive', (disabled && !readonly) && 'opacity-70')"
    >
      {{ label }}
    </Label>
    
    <Select :model-value="String(modelValue)" @update:model-value="handleUpdate" :disabled="disabled || readonly">
      <SelectTrigger 
        :id="id"
        :class="cn(
          'w-full !h-11 rounded-xl px-4 text-sm transition-all',
          error && 'border-destructive focus:ring-destructive/20',
          !modelValue && 'text-muted-foreground',
          readonly && 'bg-transparent border-gray-100 shadow-none ring-0 opacity-100 cursor-default pointer-events-none'
        )"
      >
        <SelectValue :placeholder="placeholder" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectItem 
            v-for="option in options" 
            :key="option.value" 
            :value="String(option.value)"
          >
            {{ option.label }}
          </SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
    
    <p 
      v-if="error" 
      class="mt-1.5 text-[11px] font-medium text-destructive animate-in fade-in slide-in-from-top-1"
    >
      {{ error }}
    </p>
  </div>
</template>
