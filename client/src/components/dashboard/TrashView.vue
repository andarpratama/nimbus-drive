<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from './Sidebar.vue'
import Toolbar from './Toolbar.vue'
import ContentArea from './ContentArea.vue'
import ContextMenu from './ContextMenu.vue'
import ConfirmModal from './ConfirmModal.vue'
import Notification from './Notification.vue'
import PreviewModal from './PreviewModal.vue'
import { getFileType, formatFileSize, formatDate } from '../../composables/useFileUtils'

// Props to receive user data from parent
const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['logout'])

// State
const currentView = ref('trash')
const viewMode = ref('grid') // grid, list
const searchQuery = ref('')
const trashedItems = ref([])
const loading = ref(false)
const error = ref('')
const selectedItems = ref([])

// Fetch trashed items
const fetchTrashedItems = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    // Fetch trashed files
    const filesResponse = await fetch(`${API_BASE_URL}/api/files/trash`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!filesResponse.ok) {
      throw new Error('Failed to fetch trashed files')
    }
    
    const filesData = await filesResponse.json()
    console.log('Trashed files:', filesData)
    
    // Process trashed files
    const processedFiles = (filesData.files || filesData).map(file => ({
      id: `file-${file.ID}`,
      name: file.Name,
      type: getFileType(file.Name),
      size: formatFileSize(file.Size || 0),
      modified: formatDate(file.UpdatedAt),
      starred: false,
      shared: false,
      fileId: file.ID,
      rawSize: file.Size || 0,
      deletedAt: file.DeletedAt
    }))
    
    // TODO: Fetch trashed folders when backend supports it
    const processedFolders = []
    
    trashedItems.value = [...processedFolders, ...processedFiles]
    
  } catch (err) {
    console.error('Error fetching trashed items:', err)
    error.value = 'Failed to load trashed items'
  } finally {
    loading.value = false
  }
}

// Context menu state
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  item: null
})

// Confirmation modal state
const confirmModal = ref({
  visible: false,
  title: '',
  message: '',
  item: null,
  action: ''
})

// Notification state
const notification = ref({
  visible: false,
  type: 'success',
  title: '',
  message: ''
})

// Preview modal state
const previewModal = ref({
  visible: false,
  file: null
})

// Filtered items based on search
const filteredItems = computed(() => {
  if (!searchQuery.value) return trashedItems.value
  return trashedItems.value.filter(item => 
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Handle logout
const handleLogout = () => {
  emit('logout')
}

// Event handlers
const handleViewChange = (viewId) => {
  currentView.value = viewId
}

const handleSearchChange = (query) => {
  searchQuery.value = query
}

const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

const handleItemSelect = (itemId) => {
  const index = selectedItems.value.indexOf(itemId)
  if (index > -1) {
    selectedItems.value.splice(index, 1)
  } else {
    selectedItems.value.push(itemId)
  }
}

const handleItemDoubleClick = (item) => {
  if (item.type === 'folder') {
    // In trash view, folders can't be navigated
    console.log('Folder double-clicked in trash:', item)
  } else {
    // Handle file click (preview for images)
    if (isImageFile(item)) {
      openPreview(item)
    } else {
      console.log('File double-clicked in trash:', item)
    }
  }
}

const handleItemStarToggle = (itemId) => {
  // Toggle star functionality
  console.log('Toggle star for item:', itemId)
}

const handleRetry = () => {
  fetchTrashedItems()
}

const handleContextMenu = (data) => {
  const { event, item } = data
  
  // Adjust position for three dots button to prevent overlap
  let x = event.clientX
  let y = event.clientY
  
  // If it's a click event (three dots button), adjust position
  if (event.type === 'click') {
    // Position menu to the left of the button to avoid overlap
    x = event.clientX - 200 // Adjust based on menu width
    y = event.clientY + 10  // Small offset below the button
  }
  
  contextMenu.value = {
    visible: true,
    x: x,
    y: y,
    item: item
  }
}

const handleContextMenuClose = () => {
  contextMenu.value.visible = false
}

const handleContextMenuAction = (data) => {
  const { action, item } = data
  console.log('Context menu action:', action, 'on item:', item)
  
  // Handle different actions
  switch (action) {
    case 'restore':
      restoreItem(item)
      break
    case 'delete-permanent':
      showPermanentDeleteConfirmation(item)
      break
    case 'preview':
      openPreview(item)
      break
    default:
      console.log('Unknown action:', action)
  }
}

const showPermanentDeleteConfirmation = (item) => {
  const isFolder = item.type === 'folder'
  confirmModal.value = {
    visible: true,
    title: `Permanently Delete ${isFolder ? 'Folder' : 'File'}`,
    message: `Are you sure you want to permanently delete "${item.name}"? This action cannot be undone.`,
    item: item,
    action: 'delete-permanent'
  }
}

const handleConfirmAction = () => {
  const { action, item } = confirmModal.value
  
  switch (action) {
    case 'delete-permanent':
      permanentlyDeleteItem(item)
      break
    case 'bulk-delete-permanent':
      handleBulkPermanentDelete()
      break
  }
  
  confirmModal.value.visible = false
}

const handleConfirmCancel = () => {
  confirmModal.value.visible = false
}

const handleConfirmClose = () => {
  confirmModal.value.visible = false
}

const permanentlyDeleteItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    const endpoint = item.type === 'folder' 
      ? `/api/folders/${item.folderId}/permanent` 
      : `/api/files/${item.fileId}/permanent`
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (response.ok) {
      showNotification('success', 'Item permanently deleted', `${item.name} has been permanently deleted.`)
      // Refresh the data
      await fetchTrashedItems()
    } else {
      showNotification('error', 'Delete failed', 'Failed to permanently delete the item.')
    }
  } catch (error) {
    console.error('Permanent delete error:', error)
    showNotification('error', 'Delete failed', 'An error occurred while permanently deleting the item.')
  }
}

const restoreItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    const endpoint = item.type === 'folder' 
      ? `/api/folders/${item.folderId}/restore` 
      : `/api/files/${item.fileId}/restore`
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (response.ok) {
      showNotification('success', 'Item restored', `${item.name} has been restored.`)
      // Refresh the data
      await fetchTrashedItems()
    } else {
      showNotification('error', 'Restore failed', 'Failed to restore the item.')
    }
  } catch (error) {
    console.error('Restore error:', error)
    showNotification('error', 'Restore failed', 'An error occurred while restoring the item.')
  }
}

const showNotification = (type, title, message) => {
  notification.value = {
    visible: true,
    type,
    title,
    message
  }
}

const handleNotificationClose = () => {
  notification.value.visible = false
}

// Helper functions
const isImageFile = (item) => {
  if (!item || item.type === 'folder') return false
  const imageTypes = ['image', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  return imageTypes.includes(item.type)
}

const openPreview = (item) => {
  if (item.type === 'folder') return
  previewModal.value = {
    visible: true,
    file: item
  }
}

const closePreview = () => {
  previewModal.value.visible = false
  previewModal.value.file = null
}

// Bulk operations
const restoreSelectedItems = async () => {
  if (selectedItems.value.length === 0) return
  
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    // Restore each selected item
    const promises = selectedItems.value.map(async (itemId) => {
      const item = filteredItems.value.find(i => i.id === itemId)
      if (!item) return
      
      const endpoint = item.type === 'folder' 
        ? `/api/folders/${item.folderId}/restore` 
        : `/api/files/${item.fileId}/restore`
      
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
      
      return response.ok
    })
    
    const results = await Promise.all(promises)
    const successCount = results.filter(Boolean).length
    
    if (successCount > 0) {
      showNotification('success', 'Items restored', `${successCount} item${successCount !== 1 ? 's' : ''} have been restored.`)
      selectedItems.value = []
      await fetchTrashedItems()
    } else {
      showNotification('error', 'Restore failed', 'Failed to restore the selected items.')
    }
  } catch (error) {
    console.error('Bulk restore error:', error)
    showNotification('error', 'Restore failed', 'An error occurred while restoring the items.')
  }
}

const permanentlyDeleteSelectedItems = async () => {
  if (selectedItems.value.length === 0) return
  
  // Show confirmation modal
  confirmModal.value = {
    visible: true,
    title: 'Permanently Delete Items',
    message: `Are you sure you want to permanently delete ${selectedItems.value.length} item${selectedItems.value.length !== 1 ? 's' : ''}? This action cannot be undone.`,
    item: null,
    action: 'bulk-delete-permanent'
  }
}

const handleBulkPermanentDelete = async () => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    // Delete each selected item
    const promises = selectedItems.value.map(async (itemId) => {
      const item = filteredItems.value.find(i => i.id === itemId)
      if (!item) return
      
      const endpoint = item.type === 'folder' 
        ? `/api/folders/${item.folderId}/permanent` 
        : `/api/files/${item.fileId}/permanent`
      
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
      
      return response.ok
    })
    
    const results = await Promise.all(promises)
    const successCount = results.filter(Boolean).length
    
    if (successCount > 0) {
      showNotification('success', 'Items deleted', `${successCount} item${successCount !== 1 ? 's' : ''} have been permanently deleted.`)
      selectedItems.value = []
      await fetchTrashedItems()
    } else {
      showNotification('error', 'Delete failed', 'Failed to permanently delete the selected items.')
    }
  } catch (error) {
    console.error('Bulk delete error:', error)
    showNotification('error', 'Delete failed', 'An error occurred while deleting the items.')
  }
}



// Initialize on mount
onMounted(() => {
  fetchTrashedItems()
})

// Expose refresh function for parent component
defineExpose({
  refresh: fetchTrashedItems
})
</script>

<template>
  <div class="flex h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Left Sidebar -->
    <Sidebar 
      :current-view="currentView"
      @view-change="handleViewChange"
    />

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Top Toolbar -->
      <Toolbar 
        :search-query="searchQuery"
        :view-mode="viewMode"
        :user="user"
        @search-change="handleSearchChange"
        @view-mode-change="handleViewModeChange"
        @logout="handleLogout"
      />

      <!-- Content Area -->
      <div class="flex-1 flex flex-col">
        <!-- Bulk Actions Header -->
        <div v-if="selectedItems.length > 0" class="p-4 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
          <div class="flex items-center justify-end space-x-3">
            <button
              @click="restoreSelectedItems"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Restore Selected ({{ selectedItems.length }})
            </button>
            <button
              @click="permanentlyDeleteSelectedItems"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
            >
              Delete Forever ({{ selectedItems.length }})
            </button>
          </div>
        </div>
        
        <!-- Content Area -->
        <div class="flex-1 overflow-auto bg-gray-50 dark:bg-gray-900">
          <!-- Loading state -->
          <div v-if="loading" class="flex justify-center items-center py-20">
            <div class="text-center">
              <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p class="text-gray-600 dark:text-gray-400">Loading trashed items...</p>
            </div>
          </div>
          
          <!-- Error state -->
          <div v-else-if="error" class="flex justify-center items-center py-20">
            <div class="text-center">
              <div class="text-6xl mb-4">⚠️</div>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
                Failed to load trash
              </h3>
              <p class="text-gray-500 dark:text-gray-400 mb-4">{{ error }}</p>
              <button 
                @click="handleRetry"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                Try Again
              </button>
            </div>
          </div>
          
          <!-- Empty state -->
          <div v-else-if="filteredItems.length === 0 && !searchQuery" class="flex justify-center items-center py-20">
            <div class="text-center">
              <div class="text-8xl mb-6">🗑️</div>
              <h3 class="text-xl font-medium text-gray-900 dark:text-white mb-2">
                Trash is empty
              </h3>
              <p class="text-gray-500 dark:text-gray-400 mb-4">
                Items you delete will appear here
              </p>
              <p class="text-sm text-gray-400 dark:text-gray-500">
                Files and folders in trash are automatically deleted after 30 days
              </p>
            </div>
          </div>
          
          <!-- Search empty state -->
          <div v-else-if="filteredItems.length === 0 && searchQuery" class="flex justify-center items-center py-20">
            <div class="text-center">
              <div class="text-6xl mb-4">🔍</div>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
                No items found
              </h3>
              <p class="text-gray-500 dark:text-gray-400">
                No trashed items match your search "{{ searchQuery }}"
              </p>
            </div>
          </div>
          
          <!-- Content -->
          <div v-else class="p-6">
            <!-- Items count -->
            <div class="mb-4 text-sm text-gray-600 dark:text-gray-400">
              {{ filteredItems.length }} item{{ filteredItems.length !== 1 ? 's' : '' }} in trash
            </div>
            
            <!-- Content Area -->
            <ContentArea 
              :items="filteredItems"
              :selected-items="selectedItems"
              :view-mode="viewMode"
              :loading="false"
              :error="null"
              :search-query="searchQuery"
              @item-select="handleItemSelect"
              @item-double-click="handleItemDoubleClick"
              @item-star-toggle="handleItemStarToggle"
              @context-menu="handleContextMenu"
              @retry="handleRetry"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Context Menu -->
    <ContextMenu
      :visible="contextMenu.visible"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :item="contextMenu.item"
      @close="handleContextMenuClose"
      @action="handleContextMenuAction"
    />
    
    <!-- Confirmation Modal -->
    <ConfirmModal
      :visible="confirmModal.visible"
      :title="confirmModal.title"
      :message="confirmModal.message"
      confirm-text="Delete"
      cancel-text="Cancel"
      :danger="true"
      @confirm="handleConfirmAction"
      @cancel="handleConfirmCancel"
      @close="handleConfirmClose"
    />
    
    <!-- Notification -->
    <Notification
      :visible="notification.visible"
      :type="notification.type"
      :title="notification.title"
      :message="notification.message"
      @close="handleNotificationClose"
    />
    
    <!-- Preview Modal -->
    <PreviewModal
      :visible="previewModal.visible"
      :file="previewModal.file"
      @close="closePreview"
    />
  </div>
</template> 