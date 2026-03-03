import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'
import type { 
  User, 
  AuthResponse, 
  UserProfileResponse, 
  LoginRequest, 
  RegisterRequest,
  GoogleLoginRequest,
  Session,
  SessionListResponse
} from '../types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref<boolean>(false)
  const error = ref<string | null>(null)
  
  // Security & Session State
  const isSessionExpired = ref<boolean>(false)
  const activeSessions = ref<Session[]>([])

  // Authentication status derived from user state
  const isAuthenticated = computed(() => !!user.value)

  /**
   * Updates the global session expiration status.
   * When expired, local user data is wiped for security.
   */
  function setSessionExpired(value: boolean) {
    isSessionExpired.value = value
    if (value) {
      user.value = null
    }
  }

  async function login(credentials: LoginRequest) {
    isLoading.value = true
    error.value = null
    isSessionExpired.value = false
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

  async function fetchSessions() {
    isLoading.value = true
    try {
      const response = await api.get<SessionListResponse>('/auth/sessions')
      activeSessions.value = response.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to fetch sessions'
    } finally {
      isLoading.value = false
    }
  }

  async function revokeSession(sessionID: string) {
    try {
      await api.delete(`/auth/sessions/${sessionID}`)
      await fetchSessions()
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to revoke session'
      return false
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
      isSessionExpired.value = false
    }
  }

  return { 
    user, 
    isLoading, 
    error,
    isSessionExpired,
    activeSessions,
    isAuthenticated, 
    setSessionExpired,
    login,
    register,
    loginWithGoogle,
    fetchUserProfile,
    fetchSessions,
    revokeSession,
    updateProfile,
    uploadAvatar,
    logout 
  }
})
