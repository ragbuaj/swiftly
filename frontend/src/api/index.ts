import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: (import.meta.env.VITE_API_URL as string) || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // IMPORTANT: Send cookies with every request
})

// Add response interceptor for automated security token management
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    const authStore = useAuthStore()
    const isAuthRequest = originalRequest.url.includes('/auth/login') || originalRequest.url.includes('/auth/logout')

    // Handle 401 Unauthorized errors by attempting a transparent token refresh.
    // We skip this for login/logout requests to prevent infinite loops and improve UX.
    if (error.response?.status === 401 && !originalRequest._retry && !isAuthRequest) {
      originalRequest._retry = true

      try {
        // Attempt to obtain new tokens using the Refresh Token (sent automatically via HttpOnly cookie)
        await axios.post(
          `${api.defaults.baseURL}/auth/refresh`,
          {},
          { withCredentials: true }
        )
        
        // If refresh succeeded, retry the original failed request
        return api(originalRequest)
      } catch (refreshError: any) {
        // If the Refresh Token is also invalid/expired, the session is officially dead.
        const errorMessage = refreshError.response?.data?.message || ''
        const path = window.location.pathname
        
        // Don't show the expiration popup if we are on auth-related pages
        const isAuthPage = path.includes('/auth') || path.includes('/forgot-password') || path.includes('/reset-password') || path.includes('/verify-otp')

        // We ignore generic "Authentication required" messages which occur for guests
        const isRealSessionError = errorMessage !== 'Authentication required' && errorMessage !== ''
        
        if (isRealSessionError && !isAuthPage && !isAuthRequest) {
          authStore.setSessionExpired(true)
        }
        
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
