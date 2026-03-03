<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import api from '../api';
import BaseButton from '../components/BaseButton.vue';

const router = useRouter();
const route = useRoute();
const isLoading = ref(false);
const apiError = ref('');
const email = ref((route.query.email as string) || '');

const otp = reactive(['', '', '', '', '', '']);
const inputs = ref<HTMLInputElement[]>([]);

const handleInput = (index: number, e: Event) => {
  const target = e.target as HTMLInputElement;
  const val = target.value;
  
  // Allow only numbers
  if (!/^\d*$/.test(val)) {
    otp[index] = '';
    return;
  }

  // Auto focus next
  if (val && index < 5) {
    inputs.value[index + 1].focus();
  }
};

const handleKeyDown = (index: number, e: KeyboardEvent) => {
  if (e.key === 'Backspace' && !otp[index] && index > 0) {
    inputs.value[index - 1].focus();
  }
};

const handleSubmit = async () => {
  const otpCode = otp.join('');
  if (otpCode.length < 6) {
    apiError.value = 'Please enter all 6 digits';
    return;
  }

  isLoading.value = true;
  apiError.value = '';
  try {
    await api.post('/auth/verify-otp', {
      email: email.value,
      otp: otpCode
    });
    // Success!
    router.push('/?verified=true');
  } catch (err: any) {
    apiError.value = err.response?.data?.message || 'Invalid OTP code. Please try again.';
  } finally {
    isLoading.value = false;
  }
};

const resendOtp = async () => {
  try {
    await api.post('/auth/resend-otp', { email: email.value });
    alert('A new OTP has been sent to your email (check console in dev)');
  } catch (err) {
    apiError.value = 'Failed to resend OTP.';
  }
};

onMounted(() => {
  if (!email.value) {
    router.push('/auth');
  }
  nextTick(() => {
    inputs.value[0]?.focus();
  });
});
</script>

<template>
  <div class="min-h-screen w-full bg-[#fafafa] flex flex-col items-center justify-center p-6">
    <div class="w-full max-w-[440px] bg-white border border-gray-100 rounded-[32px] shadow-2xl p-8 md:p-10 text-center">
      <div class="w-16 h-16 bg-[#7b00ff]/10 rounded-full flex items-center justify-center mx-auto mb-6 text-[#7b00ff]">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="11" x="4" y="9" rx="2" ry="2"/><path d="M7 9V5a5 5 0 0 1 10 0v4"/></svg>
      </div>
      
      <h1 class="text-2xl font-extrabold text-gray-900 mb-2">Verify your Phone</h1>
      <p class="text-gray-400 text-sm mb-8">We've sent a 6-digit code to your phone and email <span class="text-gray-900 font-medium">{{ email }}</span></p>

      <div v-if="apiError" class="mb-6 p-3 bg-red-50 text-red-600 text-xs font-semibold rounded-xl">{{ apiError }}</div>

      <form @submit.prevent="handleSubmit">
        <div class="flex justify-between gap-2 mb-8">
          <input
            v-for="(digit, index) in otp"
            :key="index"
            ref="inputs"
            v-model="otp[index]"
            type="text"
            maxlength="1"
            class="w-12 h-14 text-center text-xl font-bold bg-gray-50 border border-gray-100 rounded-xl focus:bg-white focus:border-[#7b00ff] focus:ring-4 focus:ring-[#7b00ff]/10 outline-none transition-all"
            @input="handleInput(index, $event)"
            @keydown="handleKeyDown(index, $event)"
          />
        </div>

        <BaseButton type="submit" :loading="isLoading" class="w-full h-11 shadow-lg shadow-[#7b00ff]/20">
          Verify OTP
        </BaseButton>

      </form>

      <div class="mt-8">
        <p class="text-sm text-gray-400">
          Didn't receive the code? 
          <button @click="resendOtp" class="text-[#7b00ff] font-bold hover:underline cursor-pointer bg-transparent border-none p-0">Resend Code</button>
        </p>
      </div>
    </div>
  </div>
</template>
