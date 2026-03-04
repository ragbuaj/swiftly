<script setup lang="ts">
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from 'reka-ui';
import BaseButton from './BaseButton.vue';

interface Props {
  open: boolean;
  title: string;
  description: string;
  confirmText?: string;
  cancelText?: string;
  isLoading?: boolean;
  variant?: 'danger' | 'primary';
}

withDefaults(defineProps<Props>(), {
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  isLoading: false,
  variant: 'primary',
});

defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'confirm'): void;
}>();
</script>

<template>
  <DialogRoot :open="open" @update:open="$emit('update:open', $event)">
    <DialogPortal>
      <DialogOverlay
        class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
      />
      <DialogContent
        class="fixed left-[50%] top-[50%] z-50 w-full max-w-[90vw] translate-x-[-50%] translate-y-[-50%] gap-4 border bg-white p-6 shadow-2xl duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:max-w-[440px] rounded-3xl"
      >
        <div class="flex flex-col gap-2 text-center sm:text-left">
          <DialogTitle class="text-xl font-bold text-gray-900 leading-tight">
            {{ title }}
          </DialogTitle>
          <DialogDescription class="text-sm text-gray-500 leading-relaxed">
            {{ description }}
          </DialogDescription>
        </div>

        <div class="flex flex-col-reverse sm:flex-row sm:justify-end gap-3 mt-8">
          <DialogClose as-child>
            <BaseButton
              variant="outline"
              size="lg"
              class="w-full sm:w-auto font-bold rounded-xl border-gray-100 text-gray-600 hover:bg-gray-50 transition-colors px-8"
              :disabled="isLoading"
            >
              {{ cancelText }}
            </BaseButton>
          </DialogClose>
          <BaseButton
            :variant="variant === 'danger' ? 'danger' : 'primary'"
            size="lg"
            class="w-full sm:w-auto font-bold rounded-xl shadow-lg shadow-red-500/10 px-8"
            :is-loading="isLoading"
            @click="$emit('confirm')"
          >
            {{ confirmText }}
          </BaseButton>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
