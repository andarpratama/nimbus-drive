<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import Notification from './Notification.vue'
import { apiRequest } from '../../composables/useApi'

const props = defineProps({
  user: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['logout'])

// Initialize activeTab from URL hash or default to 'profile'
const activeTab = ref(getTabFromHash() || 'profile')

// Function to get tab from URL hash
function getTabFromHash() {
  const hash = window.location.hash
  if (hash === '#password') return 'password'
  if (hash === '#profile') return 'profile'
  return null
}

// Function to validate tab name
function isValidTab(tab) {
  return tab === 'profile' || tab === 'password'
}

// Function to update URL hash when tab changes
function updateHash(tab) {
  window.location.hash = tab
}
const loading = ref(false)
const error = ref('')

// Profile form data
const profileForm = ref({
  name: '',
  username: '',
  email: '',
  bio: ''
})

// Function to initialize profile form with user data
function initializeProfileForm() {
  if (props.user) {
    profileForm.value = {
      name: props.user.name || '',
      username: props.user.username || '',
      email: props.user.email || '',
      bio: props.user.bio || ''
    }
  }
}

// Computed property to ensure form data is always in sync with user data
const userData = computed(() => {
  return {
    name: props.user?.name || '',
    username: props.user?.username || '',
    email: props.user?.email || '',
    bio: props.user?.bio || ''
  }
})

// Password form data
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const notification = ref({
  visible: false,
  type: 'success',
  title: '',
  message: ''
})

const showNotification = (type, title, message) => {
  notification.value = {
    visible: true,
    type,
    title,
    message
  }
  setTimeout(() => {
    notification.value.visible = false
  }, 3000)
}

const handleLogout = () => {
  emit('logout')
}

const updateProfile = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await apiRequest('/api/users/profile', {
      method: 'PUT',
      body: JSON.stringify({
        name: profileForm.value.name,
        username: profileForm.value.username,
        email: profileForm.value.email,
        bio: profileForm.value.bio
      })
    })
    
    showNotification('success', 'Profile Updated', 'Your profile has been updated successfully.')
  } catch (err) {
    console.error('Error updating profile:', err)
    error.value = 'Failed to update profile'
    showNotification('error', 'Update Failed', 'Failed to update profile. Please try again.')
  } finally {
    loading.value = false
  }
}

const resetProfileForm = () => {
  initializeProfileForm()
  showNotification('info', 'Form Reset', 'Form has been reset to current values.')
}

const updatePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    showNotification('error', 'Password Mismatch', 'New password and confirm password do not match.')
    return
  }
  
  if (passwordForm.value.newPassword.length < 6) {
    showNotification('error', 'Password Too Short', 'New password must be at least 6 characters long.')
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await apiRequest('/api/users/password', {
      method: 'PUT',
      body: JSON.stringify({
        current_password: passwordForm.value.currentPassword,
        new_password: passwordForm.value.newPassword
      })
    })
    
    showNotification('success', 'Password Updated', 'Your password has been updated successfully.')
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (err) {
    console.error('Error updating password:', err)
    error.value = 'Failed to update password'
    showNotification('error', 'Update Failed', 'Failed to update password. Please check your current password.')
  } finally {
    loading.value = false
  }
}

const switchTab = (tab) => {
  if (isValidTab(tab)) {
    activeTab.value = tab
    updateHash(tab)
  }
}

// Watch for hash changes (when user navigates with browser back/forward)
watch(() => window.location.hash, (newHash) => {
  const tab = getTabFromHash()
  if (tab && isValidTab(tab) && tab !== activeTab.value) {
    activeTab.value = tab
  }
})

// Watch for user prop changes to update form data
watch(() => props.user, (newUser) => {
  if (newUser) {
    initializeProfileForm()
  }
}, { immediate: true })

onMounted(() => {
  // Initialize form with user data
  initializeProfileForm()
  
  // Set initial hash if none exists
  if (!window.location.hash) {
    updateHash(activeTab.value)
  }
})
</script>

<template>
  <div class="h-full flex flex-col">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <svg class="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Profile Settings</h1>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 p-6 overflow-auto">
      <div class="max-w-2xl mx-auto">
        <!-- Tab Navigation -->
        <div class="border-b border-gray-200 dark:border-gray-700 mb-6">
          <nav class="flex space-x-8">
            <button
              @click="switchTab('profile')"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'profile'
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              Profile Information
            </button>
            <button
              @click="switchTab('password')"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'password'
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              Change Password
            </button>
          </nav>
        </div>

        <!-- Profile Tab -->
        <div v-if="activeTab === 'profile'" class="space-y-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Personal Information</h2>
            
            <form @submit.prevent="updateProfile" class="space-y-4">
              <div>
                <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Full Name
                </label>
                <input
                  id="name"
                  v-model="profileForm.name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div>
                <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Username
                </label>
                <input
                  id="username"
                  v-model="profileForm.username"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div>
                <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Email Address
                </label>
                <input
                  id="email"
                  v-model="profileForm.email"
                  type="email"
                  required
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div>
                <label for="bio" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Bio
                </label>
                <textarea
                  id="bio"
                  v-model="profileForm.bio"
                  rows="4"
                  placeholder="Tell us about yourself..."
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                ></textarea>
              </div>

              <div class="flex justify-end space-x-3">
                <button
                  type="button"
                  @click="resetProfileForm"
                  class="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
                >
                  Reset
                </button>
                <button
                  type="submit"
                  :disabled="loading"
                  class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ loading ? 'Updating...' : 'Update Profile' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Password Tab -->
        <div v-if="activeTab === 'password'" class="space-y-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Change Password</h2>
            
            <form @submit.prevent="updatePassword" class="space-y-4">
              <div>
                <label for="currentPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Current Password
                </label>
                <input
                  id="currentPassword"
                  v-model="passwordForm.currentPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div>
                <label for="newPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  New Password
                </label>
                <input
                  id="newPassword"
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  minlength="6"
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div>
                <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Confirm New Password
                </label>
                <input
                  id="confirmPassword"
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                />
              </div>

              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="loading"
                  class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ loading ? 'Updating...' : 'Update Password' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Notification -->
    <Notification
      v-if="notification.visible"
      :type="notification.type"
      :title="notification.title"
      :message="notification.message"
    />
  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
</style> 