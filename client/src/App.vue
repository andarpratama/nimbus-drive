<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDropdownManager } from './composables/useDropdownManager'

const router = useRouter()
const user = ref(null)
const loading = ref(true)

// Use the correct API URL for Docker environment
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const isAuthenticated = computed(() => {
  return localStorage.getItem('token') !== null
})

const handleLoginSuccess = (userData) => {
  user.value = userData
  router.push('/dashboard')
}

const logout = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    // No token to logout, just clear local state
    user.value = null
    router.push('/login')
    return
  }

  try {
    // Call logout API
    const response = await fetch(`${API_BASE_URL}/api/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (response.ok) {
      console.log('Logout successful')
    } else {
      console.warn('Logout API call failed, but clearing local state')
    }
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    // Always clear local state regardless of API response
    localStorage.removeItem('token')
    user.value = null
    router.push('/login')
  }
}

// Check authentication on mount
const checkAuth = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    loading.value = false
    return
  }
  
  try {
    // Fetch user data
    const response = await fetch(`${API_BASE_URL}/api/user`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (response.ok) {
      user.value = await response.json()
      if (router.currentRoute.value.path === '/login' || router.currentRoute.value.path === '/register') {
        router.push('/dashboard')
      }
    } else {
      // Token invalid, redirect to login
      localStorage.removeItem('token')
      router.push('/login')
    }
  } catch (error) {
    console.error('Error fetching user data:', error)
    localStorage.removeItem('token')
    router.push('/login')
  } finally {
    loading.value = false
  }
}

// Initialize dropdown manager
const { closeAllDropdowns } = useDropdownManager()

onMounted(() => {
  checkAuth()
  
  // Global click handler to close dropdowns when clicking outside
  document.addEventListener('click', (event) => {
    if (!event.target.closest('.dropdown-container')) {
      closeAllDropdowns()
    }
  })
})

// Listen for storage changes (login/logout)
window.addEventListener('storage', checkAuth)
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
    <!-- Page Content -->
    <main>
      <router-view :user="user" @login-success="handleLoginSuccess" @logout="logout" />
    </main>
  </div>
</template>

<style scoped>
.logo {
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
</style>
