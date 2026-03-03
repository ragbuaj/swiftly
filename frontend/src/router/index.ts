import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import { useAuthStore } from '../stores/auth'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue')
  },
  {
    path: '/auth',
    name: 'auth',
    component: AuthView,
    meta: { hideNavbar: true, requiresGuest: true }
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('../views/ForgotPasswordView.vue'),
    meta: { hideNavbar: true, requiresGuest: true }
  },
  {
    path: '/reset-password',
    name: 'reset-password',
    component: () => import('../views/ResetPasswordView.vue'),
    meta: { hideNavbar: true, requiresGuest: true }
  },
  {
    path: '/verify-otp',
    name: 'verify-otp',
    component: () => import('../views/VerifyOtpView.vue'),
    meta: { hideNavbar: true }
  },
  {
    path: '/catalog',
    name: 'catalog',
    component: () => import('../views/CatalogView.vue')
  },
  {
    path: '/profile',
    component: () => import('../views/profile/ProfileLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '', // Default route redirects to settings
        redirect: '/profile/settings'
      },
      {
        path: 'settings',
        name: 'profile-settings',
        component: () => import('../views/profile/ProfileSettings.vue')
      }
      // Future routes: /profile/orders, /profile/addresses can be added here
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Refresh user state on initial load if we don't have it but might have a cookie
  if (!authStore.user && !authStore.isLoading) {
    await authStore.fetchUserProfile()
  }

  const isAuth = authStore.isAuthenticated

  if (to.meta.requiresAuth && !isAuth) {
    // Route requires login, but user is not logged in.
    next('/auth')
  } else if (to.meta.requiresGuest && isAuth) {
    // Route is for guests only (like login/register), but user is logged in.
    next('/profile/settings')
  } else {
    next()
  }
})

export default router
