import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'
import Dashboard from '../components/Dashboard.vue'
import TrashView from '../components/dashboard/TrashView.vue'
import StarredView from '../components/dashboard/StarredView.vue'
import RecentView from '../components/dashboard/RecentView.vue'
import ProfileView from '../components/dashboard/ProfileView.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/trash',
    name: 'Trash',
    component: TrashView,
    meta: { requiresAuth: true }
  },
  {
    path: '/starred',
    name: 'Starred',
    component: StarredView,
    meta: { requiresAuth: true }
  },
  {
    path: '/recent',
    name: 'Recent',
    component: RecentView,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfileView,
    meta: { requiresAuth: true }
  }
]
const router = createRouter({
  history: createWebHistory(),
  routes
})
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const isAuthenticated = !!token
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && isAuthenticated) {
    next('/dashboard')
  } else if (to.path === '/register' && isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})
export default router 