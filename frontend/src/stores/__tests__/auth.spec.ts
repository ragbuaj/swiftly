import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'
import api from '../../api'

// Mock the API
vi.mock('../../api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
  }
}))

describe('Auth Store - Session Management', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with default session state', () => {
    const authStore = useAuthStore()
    expect(authStore.isSessionExpired).toBe(false)
    expect(authStore.activeSessions).toEqual([])
  })

  it('should set session expired state correctly', () => {
    const authStore = useAuthStore()
    authStore.user = { id: '1', email: 'test@example.com', full_name: 'Test' } as any
    
    authStore.setSessionExpired(true)
    
    expect(authStore.isSessionExpired).toBe(true)
    expect(authStore.user).toBeNull() // Should clear user on expiration
  })

  it('should fetch active sessions successfully', async () => {
    const authStore = useAuthStore()
    const mockSessions = [
      { id: 'sess_1', device_type: 'Desktop', is_current: true },
      { id: 'sess_2', device_type: 'Mobile', is_current: false }
    ]
    
    ;(api.get as any).mockResolvedValueOnce({
      data: { data: mockSessions }
    })

    await authStore.fetchSessions()

    expect(api.get).toHaveBeenCalledWith('/auth/sessions')
    expect(authStore.activeSessions).toHaveLength(2)
    expect(authStore.activeSessions[0].id).toBe('sess_1')
  })

  it('should handle session revocation', async () => {
    const authStore = useAuthStore()
    const sessionID = 'sess_to_delete'
    
    ;(api.delete as any).mockResolvedValueOnce({})
    // Mock subsequent fetchSessions
    ;(api.get as any).mockResolvedValueOnce({ data: { data: [] } })

    const result = await authStore.revokeSession(sessionID)

    expect(api.delete).toHaveBeenCalledWith(`/auth/sessions/${sessionID}`)
    expect(result).toBe(true)
  })

  it('should reset session expired flag on logout', async () => {
    const authStore = useAuthStore()
    authStore.isSessionExpired = true
    
    ;(api.post as any).mockResolvedValueOnce({})

    await authStore.logout()

    expect(authStore.isSessionExpired).toBe(false)
    expect(authStore.user).toBeNull()
  })
})
