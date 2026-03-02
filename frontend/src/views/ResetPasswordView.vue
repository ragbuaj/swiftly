<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import api from '../api';
import type { ResetPasswordRequest } from '../types';
import BaseInput from '../components/BaseInput.vue';
import BaseButton from '../components/BaseButton.vue';

const router = useRouter();
const route = useRoute();
const isLoading = ref(false);
const isVerifying = ref(true); // New state for initial verification
const isSuccess = ref(false);
const apiError = ref('');

const token = computed(() => route.query.token as string);

const form = reactive({
  password: '',
  confirmPassword: ''
});

const errors = reactive({
  password: '',
  confirmPassword: ''
});

const passwordRequirements = computed(() => {
  return [
    { label: 'At least 8 characters', met: form.password.length >= 8 },
    { label: 'One uppercase letter', met: /[A-Z]/.test(form.password) },
    { label: 'One number', met: /[0-9]/.test(form.password) },
    { label: 'One special character', met: /[^A-Za-z0-9]/.test(form.password) }
  ];
});

const passwordsMatch = computed(() => form.password.length > 0 && form.password === form.confirmPassword);
const isPasswordValid = computed(() => passwordRequirements.value.every(req => req.met) && passwordsMatch.value);

const validate = (): boolean => {
  errors.password = '';
  errors.confirmPassword = '';
  apiError.value = '';

  if (!isPasswordValid.value) {
    errors.password = 'Please meet all password requirements';
    return false;
  }

  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match';
    return false;
  }

  return true;
};

const verifyToken = async () => {
  if (!token.value) {
    apiError.value = 'Invalid or missing reset token.';
    isVerifying.value = false;
    return;
  }

  try {
    await api.get(`/auth/validate-reset-token?token=${token.value}`);
    isVerifying.value = false;
  } catch (err: any) {
    apiError.value = err.response?.data?.message || 'Invalid or expired reset link.';
    isVerifying.value = false;
  }
};

const handleSubmit = async () => {
  if (!validate()) return;
  
  isLoading.value = true;
  try {
    await api.post('/auth/reset-password', {
      token: token.value,
      new_password: form.password
    } as ResetPasswordRequest);
    
    isSuccess.value = true;
    setTimeout(() => {
      router.push('/auth');
    }, 3000);
  } catch (err: any) {
    apiError.value = err.response?.data?.message || 'Failed to reset password.';
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  verifyToken();
});
</script>

<template>
  <div class="min-h-screen w-full bg-[#fafafa] flex flex-col items-center justify-start md:justify-center p-6 sm:p-10">
    <!-- Branding -->
    <div class="mb-8 text-center flex-shrink-0 animate-in fade-in zoom-in duration-700">
      <router-link to="/" class="inline-flex items-center gap-2.5 text-black no-underline">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="#7b00ff">
          <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
        </svg>
        <span class="text-2xl font-black tracking-tighter">Swiftly</span>
      </router-link>
    </div>

    <!-- Card -->
    <div class="w-full max-w-[480px] bg-white border border-gray-100 rounded-[32px] shadow-2xl shadow-gray-200/50 p-8 md:p-12 transition-all duration-500 mb-10 text-left">
      
      <!-- Verifying State -->
      <div v-if="isVerifying" class="text-center py-10 animate-pulse">
        <div class="w-12 h-12 border-4 border-gray-100 border-t-[#7b00ff] rounded-full animate-spin mx-auto mb-4"></div>
        <p class="text-gray-400 font-medium">Verifying reset link...</p>
      </div>

      <!-- Error State (Show if token is invalid) -->
      <div v-else-if="apiError && !isSuccess && form.password === ''" class="text-center animate-in fade-in zoom-in duration-500">
        <div class="w-20 h-20 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
        </div>
        <h2 class="text-2xl font-extrabold text-gray-900 mb-3">Invalid Link</h2>
        <p class="text-gray-400 text-sm leading-relaxed mb-8">{{ apiError }}</p>
        <BaseButton @click="router.push('/forgot-password')">
          Request New Link
        </BaseButton>
      </div>

      <!-- Form State -->
      <div v-else-if="!isSuccess">
        <div class="mb-8">
          <h1 class="text-2xl font-extrabold text-gray-900 mb-2 tracking-tight">Create New Password</h1>
          <p class="text-gray-400 text-sm leading-relaxed">Please enter a new, strong password for your Swiftly account.</p>
        </div>

        <div v-if="apiError" class="mb-6 p-3 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold rounded-xl flex items-center gap-2 animate-in fade-in slide-in-from-top-1">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
          {{ apiError }}
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="mb-2">
            <BaseInput 
              label="New Password" 
              v-model="form.password" 
              type="password"
              placeholder="••••••••" 
              :error="errors.password"
              id="password"
            />
            
            <!-- Checklist -->
            <div v-if="form.password.length > 0" class="px-1 mb-6 animate-in fade-in slide-in-from-top-2">
              <p class="text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-2">Security Checklist</p>
              <div class="grid grid-cols-2 gap-y-1.5 gap-x-4">
                <div v-for="req in passwordRequirements" :key="req.label" class="flex items-center gap-2">
                  <div class="w-3.5 h-3.5 rounded-full flex items-center justify-center transition-colors" :class="req.met ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-300'">
                    <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  </div>
                  <span class="text-[11px] font-medium transition-colors" :class="req.met ? 'text-green-600' : 'text-gray-400'">{{ req.label }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <div class="w-3.5 h-3.5 rounded-full flex items-center justify-center transition-colors shrink-0" :class="passwordsMatch ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-300'">
                    <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  </div>
                  <span class="text-[11px] font-medium transition-colors whitespace-nowrap" :class="passwordsMatch ? 'text-green-600' : 'text-gray-400'">Passwords match</span>
                </div>
              </div>
            </div>
          </div>

          <BaseInput 
            label="Confirm New Password" 
            v-model="form.confirmPassword" 
            type="password"
            placeholder="••••••••" 
            :error="errors.confirmPassword"
            id="confirmPassword"
          />

          <BaseButton type="submit" :loading="isLoading" :disabled="!isPasswordValid || !token" class="shadow-lg shadow-[#7b00ff]/20 mt-4">
            Reset Password
          </BaseButton>
        </form>
      </div>

      <!-- Success State -->
      <div v-else class="text-center animate-in fade-in zoom-in duration-500">
        <div class="w-20 h-20 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        </div>
        <h2 class="text-2xl font-extrabold text-gray-900 mb-3">Password Updated!</h2>
        <p class="text-gray-400 text-sm leading-relaxed mb-8">
          Your password has been successfully reset. You will be redirected to the login page in a few seconds.
        </p>
        <BaseButton @click="router.push('/auth')">
          Go to Login Now
        </BaseButton>
      </div>

    </div>
  </div>
</template>
