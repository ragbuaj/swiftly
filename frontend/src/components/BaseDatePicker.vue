<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { format, parseISO, isValid } from 'date-fns'
import { Calendar as CalendarIcon } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Label } from '@/components/ui/label'
import { 
  getLocalTimeZone, 
  parseDate, 
  today,
  type DateValue 
} from '@internationalized/date'
import { toDate } from 'reka-ui/date'

const props = defineProps<{
  modelValue?: string | number;
  label?: string;
  placeholder?: string;
  id?: string;
  error?: string;
  disabled?: boolean;
  readonly?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

// internal reka-ui date state
const date = ref<DateValue | undefined>(
  props.modelValue ? parseDate(String(props.modelValue).split('T')[0]) : undefined
)

// Sync from parent prop
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    try {
      date.value = parseDate(String(newVal).split('T')[0]);
    } catch (e) {
      date.value = undefined;
    }
  } else {
    date.value = undefined;
  }
})

// Update parent when selection changes
const handleDateUpdate = (newDate: DateValue | undefined) => {
  date.value = newDate;
  if (newDate) {
    // Gunakan toString() untuk mendapatkan format YYYY-MM-DD murni tanpa pergeseran timezone
    emit('update:modelValue', newDate.toString());
  } else {
    emit('update:modelValue', '');
  }
};

const formattedDate = computed(() => {
  if (!date.value) return props.placeholder || 'Pick a date';
  const jsDate = toDate(date.value);
  return isValid(jsDate) ? format(jsDate, 'PPP') : props.placeholder || 'Pick a date'
})
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

    <Popover>
      <PopoverTrigger as-child>
        <Button
          variant="outline"
          :id="id"
          :disabled="disabled || readonly"
          :class="cn(
            'w-full !h-11 justify-start text-left font-normal rounded-xl px-4 text-sm transition-all',
            !date && 'text-muted-foreground',
            error && 'border-destructive focus:ring-destructive/20',
            readonly && 'bg-transparent border-gray-100 shadow-none ring-0 opacity-100 cursor-default pointer-events-none'
          )"
        >
          <CalendarIcon v-if="!readonly" class="mr-2 h-4 w-4 opacity-50" />
          {{ formattedDate }}
        </Button>
      </PopoverTrigger>
      <PopoverContent v-if="!readonly" class="w-auto p-0 rounded-xl" align="start">
        <Calendar 
          :model-value="date" 
          @update:model-value="handleDateUpdate"
          layout="month-and-year" 
          initial-focus 
        />
      </PopoverContent>
    </Popover>

    <p v-if="error" class="mt-1.5 text-[11px] font-medium text-destructive animate-in fade-in slide-in-from-top-1">
      {{ error }}
    </p>
  </div>
</template>
