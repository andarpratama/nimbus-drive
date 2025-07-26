<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from './dashboard/Sidebar.vue'
import Toolbar from './dashboard/Toolbar.vue'
import ContentArea from './dashboard/ContentArea.vue'
import { useFileManager } from '../composables/useFileManager.js'

// Props to receive user data from parent
const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['logout'])

// Use the file manager composable
const {
  files,
  folders,
  currentFolderId,
  currentFolder,
  breadcrumbs,
  loading,
  error,
  selectedItems,
  allItems,
  fetchFolderContents,
  navigateToFolder,
  navigateToBreadcrumb,
  navigateToRoot,
  selectItem,
  isSelected,
  clearSelection,
  toggleStar
} = useFileManager()

// State management
const currentView = ref('my-drive') // my-drive, shared, recent, starred, trash
const viewMode = ref('grid') // grid, list
const searchQuery = ref('')

// Filtered items based on search
const filteredItems = computed(() => {
  if (!searchQuery.value) return allItems.value
  return allItems.value.filter(item => 
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Event handlers
const handleViewChange = (viewId) => {
  currentView.value = viewId
  if (viewId === 'my-drive') {
    navigateToRoot()
  }
}

const handleSearchChange = (query) => {
  searchQuery.value = query
}

const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

const handleItemSelect = (itemId) => {
  selectItem(itemId)
}

const handleItemDoubleClick = (item) => {
  if (item.type === 'folder') {
    navigateToFolder(item.folderId)
  } else {
    // Handle file click (download, preview, etc.)
    console.log('File clicked:', item)
  }
}

const handleItemStarToggle = (itemId) => {
  toggleStar(itemId)
}

const handleRetry = () => {
  fetchFolderContents(currentFolderId.value)
}

const handleLogout = () => {
  emit('logout')
}

// Initialize on mount
onMounted(() => {
  fetchFolderContents()
})
</script>

<template>
  <div class="flex h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Left Sidebar -->
    <Sidebar 
      :current-view="currentView"
      @view-change="handleViewChange"
      @navigate-root="navigateToRoot"
    />

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Top Toolbar -->
      <Toolbar 
        :breadcrumbs="breadcrumbs"
        :search-query="searchQuery"
        :view-mode="viewMode"
        :user="user"
        @search-change="handleSearchChange"
        @view-mode-change="handleViewModeChange"
        @navigate-breadcrumb="navigateToBreadcrumb"
        @logout="handleLogout"
      />

      <!-- Content Area -->
      <ContentArea 
        :items="filteredItems"
        :selected-items="selectedItems"
        :view-mode="viewMode"
        :loading="loading"
        :error="error"
        :search-query="searchQuery"
        @item-select="handleItemSelect"
        @item-double-click="handleItemDoubleClick"
        @item-star-toggle="handleItemStarToggle"
        @retry="handleRetry"
      />
    </div>
  </div>
</template>

<style scoped>
/* Close dropdowns when clicking outside */
</style> 