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
  { id: 'my-drive', label: 'My Drive', icon: '📁', route: '/dashboard' },
  { id: 'shared', label: 'Shared with me', icon: '👥', route: '/shared' },
  { id: 'recent', label: 'Recent', icon: '🕒', route: '/recent' },
  { id: 'starred', label: 'Starred', icon: '⭐', route: '/starred' },
  { id: 'trash', label: 'Trash', icon: '🗑️', route: '/trash' }
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
        <span class="mr-3">{{ item.icon }}</span>
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