<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { computed } from 'vue'

interface Props {
  variant?: 'primary' | 'outline' | 'social';
  size?: 'default' | 'sm' | 'lg' | 'icon';
  loading?: boolean;
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'default',
  loading: false,
  disabled: false,
  type: 'button'
});

const shadcnVariant = computed(() => {
  switch (props.variant) {
    case 'primary': return 'default';
    case 'outline': return 'outline';
    case 'social': return 'outline';
    default: return 'default';
  }
});
</script>

<template>
  <Button
    :type="type"
    :disabled="disabled || loading"
    :variant="shadcnVariant"
    :size="size"
    :class="cn(
      'font-semibold active:scale-[0.98] transition-all',
      variant === 'social' && 'hover:border-gray-300',
      props.class
    )"
  >
    <div v-if="loading" class="w-4 h-4 border-2 border-current/30 border-t-current rounded-full animate-spin"></div>
    <slot v-else></slot>
  </Button>
</template>
