<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { sanitize, normalizeEmail } from '../utils/sanitizer';
import api from '../api';
import type { CheckEmailResponse, RegisterRequest, LoginRequest } from '../types';
import BaseInput from '../components/BaseInput.vue';
import BaseButton from '../components/BaseButton.vue';

// Local interfaces
interface AuthForm {
  email: string;
  username: string;
  phone: string;
  password: string;
  confirmPassword: string;
  fullName: string;
}

const router = useRouter();
const authStore = useAuthStore();

const isLogin = ref(true);
const apiError = ref('');
const googleBtn = ref<HTMLElement | null>(null);
const turnstileToken = ref('');

const isCheckingEmail = ref(false);
const emailAvailable = ref<boolean | null>(null);

const form = reactive<AuthForm>({
  email: '',
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
  fullName: ''
});

const errors = reactive<Record<keyof AuthForm, string>>({
  email: '',
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
  fullName: ''
});

const resetForm = (): void => {
  form.email = '';
  form.username = '';
  form.phone = '';
  form.password = '';
  form.confirmPassword = '';
  form.fullName = '';
  (Object.keys(errors) as Array<keyof AuthForm>).forEach(key => errors[key] = '');
  apiError.value = '';
  emailAvailable.value = null;
  isCheckingEmail.value = false;
  turnstileToken.value = '';
  if (window.turnstile) window.turnstile.reset();
};

const toggleMode = (loginMode: boolean): void => {
  if (isLogin.value === loginMode) return;
  isLogin.value = loginMode;
  resetForm();
};

const passwordRequirements = computed(() => {
  return [
    { label: 'At least 8 characters', met: form.password.length >= 8 },
    { label: 'One uppercase letter', met: /[A-Z]/.test(form.password) },
    { label: 'One number', met: /[0-9]/.test(form.password) },
    { label: 'One special character', met: /[^A-Za-z0-9]/.test(form.password) }
  ];
});

const passwordsMatch = computed(() => form.password.length > 0 && form.password === form.confirmPassword);

const isFormValid = computed(() => {
  // Turnstile token is required for both Login and Register
  if (!turnstileToken.value) return false;

  if (isLogin.value) {
    return form.email && form.password;
  }
  return passwordRequirements.value.every(req => req.met) && passwordsMatch.value && form.fullName && form.username;
});

// Watch email for availability check
let debounceTimer: any = null;
watch(() => form.email, (newEmail) => {
  if (isLogin.value) return;
  emailAvailable.value = null;
  errors.email = '';
  if (debounceTimer) clearTimeout(debounceTimer);
  if (newEmail && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(newEmail)) {
    debounceTimer = setTimeout(checkEmailAvailability, 600);
  }
});

const checkEmailAvailability = async () => {
  isCheckingEmail.value = true;
  try {
    const cleanEmail = normalizeEmail(form.email);
    const response = await api.get<CheckEmailResponse>(`/users/check-email?email=${cleanEmail}`);
    emailAvailable.value = response.data.data.available;
    if (!emailAvailable.value) errors.email = 'This email is already registered';
  } catch (err) {
    console.error('Email check failed', err);
  } finally {
    isCheckingEmail.value = false;
  }
};

const handleSubmit = async () => {
  apiError.value = '';
  
  const cleanEmail = normalizeEmail(form.email);

  try {
    if (isLogin.value) {
      await authStore.login({ 
        email: cleanEmail, 
        password: form.password,
        captcha_token: turnstileToken.value 
      } as LoginRequest);
      router.push('/');
    } else {
      const payload = {
        email: cleanEmail,
        password: form.password,
        full_name: sanitize(form.fullName),
        username: form.username.toLowerCase().trim(),
        phone_number: form.phone.replace(/[^0-9+]/g, ''),
        captcha_token: turnstileToken.value
      };
      await authStore.register(payload as any);
      router.push(`/verify-otp?email=${cleanEmail}`);
    }
  } catch (err: any) {
    apiError.value = authStore.error || 'Authentication failed.';
    if (window.turnstile) window.turnstile.reset();
    turnstileToken.value = ''; // Reset token on failure
  }
};

const handleGoogleCallback = async (response: any) => {
  try {
    apiError.value = '';
    await authStore.loginWithGoogle({ id_token: response.credential });
    router.push('/');
  } catch (err) {
    apiError.value = 'Google sign-in failed.';
  }
};

onMounted(() => {
  // Init Google
  if (window.google) {
    window.google.accounts.id.initialize({
      client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID,
      callback: handleGoogleCallback,
      auto_select: false,
      cancel_on_tap_outside: true,
    });
    window.google.accounts.id.renderButton(googleBtn.value!, { theme: "outline", size: "large", width: 360, shape: "pill", text: "signin_with" });
  }

  // Init Turnstile
  if (window.turnstile) {
    window.turnstile.render('#turnstile-container', {
      sitekey: import.meta.env.VITE_TURNSTILE_SITE_KEY,
      appearance: 'interaction-only', // Only show if necessary
      theme: 'light',
      callback: (token: string) => {
        turnstileToken.value = token;
      },
    });
  }
});
</script>

<template>
  <div class="min-h-screen w-full bg-[#fafafa] flex flex-col items-center justify-start md:justify-center p-6 sm:p-10">
    <div class="mb-8 text-center flex-shrink-0 mt-4 md:mt-0 animate-in fade-in zoom-in duration-700">
      <router-link to="/" class="inline-flex items-center gap-2.5 text-black no-underline">
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="#7b00ff">
          <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
        </svg>
        <span class="text-2xl font-black tracking-tighter">Swiftly</span>
      </router-link>
    </div>

    <div class="w-full max-w-[520px] bg-white border border-gray-100 rounded-[32px] shadow-2xl shadow-gray-200/50 p-8 md:p-12 transition-all duration-500 mb-10 text-left">
      <div class="mb-8">
        <h1 class="text-2xl font-extrabold text-gray-900 mb-2 tracking-tight">{{ isLogin ? 'Welcome Back' : 'Get Started' }}</h1>
        <p class="text-gray-400 text-sm leading-relaxed">{{ isLogin ? 'Enter your details to access your account' : 'Create an account to start shopping' }}</p>
      </div>

      <div v-if="apiError" class="mb-6 p-3 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold rounded-xl flex items-center gap-2 animate-in fade-in slide-in-from-top-1">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        {{ apiError }}
      </div>

      <div class="flex p-1 bg-gray-100/80 rounded-2xl mb-8">
        <button @click="toggleMode(true)" class="flex-1 py-2.5 text-sm font-bold rounded-xl transition-all duration-300 cursor-pointer" :class="isLogin ? 'bg-white shadow-sm text-gray-900' : 'text-gray-400 hover:text-gray-600'">Sign In</button>
        <button @click="toggleMode(false)" class="flex-1 py-2.5 text-sm font-bold rounded-xl transition-all duration-300 cursor-pointer" :class="!isLogin ? 'bg-white shadow-sm text-gray-900' : 'text-gray-400 hover:text-gray-600'">Register</button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-1">
        <BaseInput v-if="!isLogin" label="Full Name" v-model="form.fullName" placeholder="John Doe" :error="errors.fullName" id="name" />
        
        <div v-if="!isLogin" class="grid grid-cols-1 sm:grid-cols-2 gap-x-4">
          <BaseInput label="Username" v-model="form.username" placeholder="johndoe" :error="errors.username" id="username" />
          <BaseInput label="Phone" v-model="form.phone" placeholder="+62..." :error="errors.phone" id="phone" />
        </div>

        <BaseInput label="Email Address" v-model="form.email" type="email" placeholder="name@example.com" :error="errors.email" id="email">
          <template #right>
            <div v-if="!isLogin && (isCheckingEmail || emailAvailable !== null)">
              <div v-if="isCheckingEmail" class="w-4 h-4 border-2 border-gray-200 border-t-[#7b00ff] rounded-full animate-spin"></div>
              <svg v-else-if="emailAvailable === true" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" class="text-green-500 bg-green-50 rounded-full p-0.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
            </div>
          </template>
        </BaseInput>

        <div class="mb-2">
          <BaseInput label="Password" v-model="form.password" type="password" placeholder="••••••••" :error="errors.password" id="password" />
          <BaseInput v-if="!isLogin" label="Confirm Password" v-model="form.confirmPassword" type="password" placeholder="••••••••" :error="errors.confirmPassword" id="confirmPassword" />

          <div v-if="!isLogin && form.password.length > 0" class="px-1 mb-6 animate-in fade-in slide-in-from-top-2">
            <p class="text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-2 text-left">Security Checklist</p>
            <div class="grid grid-cols-2 gap-y-1.5 gap-x-4">
              <div v-for="req in passwordRequirements" :key="req.label" class="flex items-center gap-2">
                <div class="w-3.5 h-3.5 rounded-full flex items-center justify-center transition-colors shrink-0" :class="req.met ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-300'">
                  <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </div>
                <span class="text-[11px] font-medium transition-colors whitespace-nowrap" :class="req.met ? 'text-green-600' : 'text-gray-400'">{{ req.label }}</span>
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

        <div v-if="isLogin" class="flex justify-end pb-6">
          <router-link to="/forgot-password" class="text-xs font-bold text-[#7b00ff] hover:text-[#6a00e0] no-underline">Forgot password?</router-link>
        </div>

        <!-- Cloudflare Turnstile Container (Visible for both Login and Register) -->
        <div class="py-4 flex justify-center">
          <div id="turnstile-container"></div>
        </div>

        <BaseButton type="submit" :loading="authStore.isLoading" :disabled="!isFormValid" class="w-full h-11 shadow-lg shadow-[#7b00ff]/20">{{ isLogin ? 'Sign In' : 'Create Account' }}</BaseButton>
      </form>

      <div class="relative my-8">
        <div class="absolute inset-0 flex items-center"><div class="w-full border-t border-gray-100"></div></div>
        <div class="relative flex justify-center text-[11px]"><span class="bg-white px-4 text-gray-300 font-bold tracking-widest uppercase tracking-tighter">OR</span></div>
      </div>

      <div ref="googleBtn" class="w-full flex justify-center pb-2">
        <div id="google-signin-btn"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
form { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
