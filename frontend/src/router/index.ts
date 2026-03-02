import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AuthView from '../views/AuthView.vue'

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
    meta: { hideNavbar: true }
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('../views/ForgotPasswordView.vue'),
    meta: { hideNavbar: true }
  },
  {
    path: '/reset-password',
    name: 'reset-password',
    component: () => import('../views/ResetPasswordView.vue'),
    meta: { hideNavbar: true }
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
