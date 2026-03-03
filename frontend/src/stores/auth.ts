import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'
import type { 
  User, 
  AuthResponse, 
  UserProfileResponse, 
  LoginRequest, 
  RegisterRequest,
  GoogleLoginRequest 
} from '../types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref<boolean>(false)
  const error = ref<string | null>(null)

  // We check if user object exists since we can't read the HttpOnly cookie via JS
  const isAuthenticated = computed(() => !!user.value)

  async function login(credentials: LoginRequest) {
    isLoading.value = true
    error.value = null
    try {
      await api.post<AuthResponse>('/auth/login', credentials)
      await fetchUserProfile()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function register(userData: RegisterRequest) {
    isLoading.value = true
    error.value = null
    try {
      await api.post<AuthResponse>('/auth/register', userData)
      await fetchUserProfile()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function loginWithGoogle(payload: GoogleLoginRequest) {
    isLoading.value = true
    error.value = null
    try {
      await api.post<AuthResponse>('/auth/google/token', payload)
      await fetchUserProfile()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Google login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function fetchUserProfile() {
    try {
      const response = await api.get<UserProfileResponse>('/users/profile')
      user.value = response.data.data
    } catch (err) {
      user.value = null
    }
  }

  async function updateProfile(profileData: Partial<UpdateProfileRequest>) {
    isLoading.value = true
    error.value = null
    try {
      await api.put('/users/profile', profileData)
      await fetchUserProfile()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to update profile'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function uploadAvatar(file: File) {
    isLoading.value = true
    error.value = null
    try {
      const formData = new FormData()
      formData.append('avatar', file)
      
      await api.post('/users/profile/avatar', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      await fetchUserProfile()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to upload avatar'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } finally {
      user.value = null
    }
  }

  return { 
    user, 
    isLoading, 
    error,
    isAuthenticated, 
    login,
    register,
    loginWithGoogle,
    fetchUserProfile,
    updateProfile,
    uploadAvatar,
    logout 
  }
})
