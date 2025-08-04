<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const props = defineProps({
  currentView: {
    type: String,
    required: true
  }
})
const emit = defineEmits(['view-change', 'navigate-root'])
const navigationItems = [
  { id: 'my-drive', label: 'My Drive', icon: 'folder', route: '/dashboard' },
  { id: 'shared', label: 'Shared with me', icon: 'users', route: '/shared' },
  { id: 'recent', label: 'Recent', icon: 'clock', route: '/recent' },
  { id: 'starred', label: 'Starred', icon: 'star', route: '/starred' },
  { id: 'trash', label: 'Trash', icon: 'trash', route: '/trash' }
]
const handleViewChange = (viewId) => {
  const item = navigationItems.find(nav => nav.id === viewId)
  if (item) {
    if (item.route) {
      router.push(item.route)
    } else if (viewId === 'my-drive') {
      emit('navigate-root')
    }
    emit('view-change', viewId)
  }
}
</script>
<template>
  <div class="w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col">
    <!-- Logo/Brand -->
    <div class="p-4 border-b border-gray-200 dark:border-gray-700">
      <h1 class="text-xl font-semibold text-gray-900 dark:text-white">Nimbus Drive</h1>
    </div>
    <!-- Navigation -->
    <nav class="flex-1 p-4 space-y-2">
      <button 
        v-for="item in navigationItems"
        :key="item.id"
        @click="handleViewChange(item.id)"
        :class="[
          'w-full flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors',
          currentView === item.id 
            ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300' 
            : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
        ]"
      >
        <!-- SVG Icons -->
        <svg v-if="item.icon === 'folder'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-5l-2-2H5a2 2 0 00-2 2z"></path>
        </svg>
        <svg v-else-if="item.icon === 'users'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
        </svg>
        <svg v-else-if="item.icon === 'clock'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <svg v-else-if="item.icon === 'star'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
        </svg>
        <svg v-else-if="item.icon === 'trash'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
        </svg>
        <svg v-else-if="item.icon === 'user'" class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
        </svg>
        {{ item.label }}
      </button>
    </nav>
    <!-- Storage Info -->
    <div class="p-4 border-t border-gray-200 dark:border-gray-700">
      <div class="text-sm text-gray-600 dark:text-gray-400 mb-2">Storage</div>
      <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 mb-2">
        <div class="bg-blue-600 h-2 rounded-full" style="width: 25%"></div>
      </div>
      <div class="text-xs text-gray-500 dark:text-gray-400">2.5 GB of 10 GB used</div>
    </div>
  </div>
</template> 