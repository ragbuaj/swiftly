<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api';
import type { ForgotPasswordRequest, ForgotPasswordResponse } from '../types';
import BaseInput from '../components/BaseInput.vue';
import BaseButton from '../components/BaseButton.vue';

const router = useRouter();
const isLoading = ref(false);
const isSubmitted = ref(false);
const apiError = ref('');
const emailHint = ref('');

const form = reactive({
  identifier: ''
});

const errors = reactive({
  identifier: ''
});

const validate = (): boolean => {
  errors.identifier = '';
  apiError.value = '';
  if (!form.identifier) {
    errors.identifier = 'Please enter your email or verified phone number';
    return false;
  }
  return true;
};

const handleSubmit = async () => {
  if (!validate()) return;
  
  isLoading.value = true;
  try {
    const response = await api.post<ForgotPasswordResponse>('/auth/forgot-password', {
      identifier: form.identifier.trim()
    } as ForgotPasswordRequest);
    
    isSubmitted.value = true;
    if (response.data.data) {
      emailHint.value = response.data.data.email_hint;
    }
  } catch (err: any) {
    apiError.value = err.response?.data?.message || 'Something went wrong. Please try again.';
  } finally {
    isLoading.value = false;
  }
};
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
      
      <div v-if="!isSubmitted">
        <div class="mb-8">
          <h1 class="text-2xl font-extrabold text-gray-900 mb-2 tracking-tight">Find Your Account</h1>
          <p class="text-gray-400 text-sm leading-relaxed">Enter your registered email address or verified phone number to receive reset instructions.</p>
        </div>

        <div v-if="apiError" class="mb-6 p-3 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold rounded-xl flex items-center gap-2 animate-in fade-in slide-in-from-top-1">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
          {{ apiError }}
        </div>

        <form @submit.prevent="handleSubmit">
          <BaseInput 
            label="Email or Phone Number" 
            v-model="form.identifier" 
            placeholder="e.g. name@example.com or +62..." 
            :error="errors.identifier"
            id="identifier"
          />

          <BaseButton type="submit" :loading="isLoading" class="shadow-lg shadow-[#7b00ff]/20 mt-4">
            Send Reset Link
          </BaseButton>
        </form>

        <div class="mt-8 text-center border-t border-gray-50 pt-6">
          <router-link to="/auth" class="text-sm font-bold text-gray-400 hover:text-[#7b00ff] transition-colors no-underline">
            Remember your account? <span class="text-[#7b00ff]">Sign In</span>
          </router-link>
        </div>
      </div>

      <!-- Success State -->
      <div v-else class="text-center animate-in fade-in zoom-in duration-500">
        <div class="w-20 h-20 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m22 2-7 20-4-9-9-4Z"/><path d="M22 2 11 13"/></svg>
        </div>
        <h2 class="text-2xl font-extrabold text-gray-900 mb-3">Link Sent!</h2>
        <p class="text-gray-400 text-sm leading-relaxed mb-8">
          If <span class="text-gray-900 font-bold">"{{ form.identifier }}"</span> is verified, we've sent reset instructions to <span class="text-[#7b00ff] font-bold">{{ emailHint }}</span>.
        </p>
        <BaseButton variant="outline" @click="isSubmitted = false">
          Try another account
        </BaseButton>
        <div class="mt-8">
          <router-link to="/auth" class="text-sm font-bold text-[#7b00ff] no-underline">
            Return to Login
          </router-link>
        </div>
      </div>

    </div>
  </div>
</template>
