<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  searchQuery: {
    type: String,
    required: true
  },
  viewMode: {
    type: String,
    required: true
  },
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits([
  'search-change', 
  'view-mode-change', 
  'logout',
  'create-menu-toggle',
  'upload-menu-toggle',
  'upload-files',
  'new-folder'
])

const showCreateMenu = ref(false)
const showUploadMenu = ref(false)
const showUserMenu = ref(false)

const handleSearchChange = (event) => {
  emit('search-change', event.target.value)
}

const handleViewModeChange = (mode) => {
  emit('view-mode-change', mode)
}



const handleLogout = () => {
  emit('logout')
}

const toggleCreateMenu = () => {
  showCreateMenu.value = !showCreateMenu.value
  emit('create-menu-toggle', showCreateMenu.value)
}

const handleNewFolder = () => {
  showCreateMenu.value = false
  emit('new-folder')
}

const toggleUploadMenu = () => {
  showUploadMenu.value = !showUploadMenu.value
  emit('upload-menu-toggle', showUploadMenu.value)
}

const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
}

// Close dropdowns when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.dropdown-container')) {
    showCreateMenu.value = false
    showUploadMenu.value = false
    showUserMenu.value = false
  }
}

// Add click outside listener
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-4">
    <div class="flex items-center justify-between">
      <!-- Left side: Search -->
      <div class="flex items-center flex-1">
        <!-- Search -->
        <div class="relative max-w-md">
          <input
            :value="searchQuery"
            @input="handleSearchChange"
            type="text"
            placeholder="Search in Drive"
            class="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
          <span class="absolute left-3 top-2.5 text-gray-400">🔍</span>
        </div>
      </div>
      
      <!-- Right side: Actions -->
      <div class="flex items-center space-x-2">
        <!-- View toggle -->
        <div class="flex border border-gray-300 dark:border-gray-600 rounded-lg">
          <button
            @click="handleViewModeChange('grid')"
            :class="[
              'px-3 py-2 text-sm transition-colors',
              viewMode === 'grid' 
                ? 'bg-blue-600 text-white' 
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
            ]"
          >
            ⊞
          </button>
          <button
            @click="handleViewModeChange('list')"
            :class="[
              'px-3 py-2 text-sm transition-colors',
              viewMode === 'list' 
                ? 'bg-blue-600 text-white' 
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
            ]"
          >
            ☰
          </button>
        </div>
        
        <!-- Create button -->
        <div class="relative dropdown-container">
          <button
            @click="toggleCreateMenu"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors"
          >
            + New
          </button>
          
          <!-- Create dropdown -->
          <div v-if="showCreateMenu" class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-10">
            <div class="py-1">
              <button 
                @click="handleNewFolder"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                📁 New folder
              </button>
              <button class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                📄 Google Docs
              </button>
              <button class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                📊 Google Sheets
              </button>
              <button class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                📽️ Google Slides
              </button>
            </div>
          </div>
        </div>
        
        <!-- Upload button -->
        <div class="relative dropdown-container">
          <button
            @click="toggleUploadMenu"
            class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-lg transition-colors"
          >
            ↑ Upload
          </button>
          
          <!-- Upload dropdown -->
          <div v-if="showUploadMenu" class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-10">
            <div class="py-1">
              <button 
                @click="() => { emit('upload-files'); showUploadMenu = false; }"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                📁 Upload files
              </button>
              <button class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                📁 Upload folder
              </button>
            </div>
          </div>
        </div>
        
        <!-- User menu -->
        <div v-if="user" class="relative dropdown-container">
          <button
            @click="toggleUserMenu"
            class="flex items-center space-x-2 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-white text-sm font-medium">
              {{ user.name ? user.name.charAt(0).toUpperCase() : 'U' }}
            </div>
            <span class="hidden sm:block">{{ user.name }}</span>
            <span>▼</span>
          </button>
          
          <!-- User dropdown -->
          <div v-if="showUserMenu" class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-10">
            <div class="py-1">
              <div class="px-4 py-2 text-sm text-gray-500 dark:text-gray-400 border-b border-gray-200 dark:border-gray-700">
                {{ user.email }}
              </div>
              <button 
                @click="handleLogout"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                🚪 Sign out
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template> 