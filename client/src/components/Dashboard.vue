<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from './dashboard/Sidebar.vue'
import Toolbar from './dashboard/Toolbar.vue'
import ContentArea from './dashboard/ContentArea.vue'
import TrashView from './dashboard/TrashView.vue'
import StarredView from './dashboard/StarredView.vue'
import ContextMenu from './dashboard/ContextMenu.vue'
import ConfirmModal from './dashboard/ConfirmModal.vue'
import Notification from './dashboard/Notification.vue'
import UploadModal from './dashboard/UploadModal.vue'
import Breadcrumb from './dashboard/Breadcrumb.vue'
import PreviewModal from './dashboard/PreviewModal.vue'
import NewFolderModal from './dashboard/NewFolderModal.vue'
import RenameModal from './dashboard/RenameModal.vue'
import MoveModal from './dashboard/MoveModal.vue'
import { useFileManager } from '../composables/useFileManager'

// Props to receive user data from parent
const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['logout'])

// Handle logout
const handleLogout = () => {
  emit('logout')
}

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

// Ref for TrashView component
const trashViewRef = ref(null)

// Ref for UploadModal component
const uploadModalRef = ref(null)

// Initialize view from URL on mount
const initializeViewFromURL = () => {
  const urlParams = new URLSearchParams(window.location.search)
  const view = urlParams.get('view') || 'my-drive'
  currentView.value = view
  
  // If we're in trash view, navigate to root first
  if (view === 'trash') {
    navigateToRoot()
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

// Upload modal state
const uploadModal = ref({
  visible: false
})

// New folder modal state
const newFolderModal = ref({
  visible: false
})

// Preview modal state
const previewModal = ref({
  visible: false,
  file: null
})

// Rename modal state
const renameModal = ref({
  visible: false,
  item: null
})

// Move modal state
const moveModal = ref({
  visible: false,
  item: null
})

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
  
  // Update URL to reflect current view
  const url = new URL(window.location)
  url.searchParams.set('view', viewId)
  window.history.pushState({}, '', url)
  
  if (viewId === 'my-drive') {
    navigateToRoot()
  }
  // Note: starred view is handled by StarredView component
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
    // Handle file click (preview for images, download for others)
    if (isImageFile(item)) {
      openPreview(item)
    } else {
      console.log('File clicked:', item)
      // TODO: Implement download functionality
    }
  }
}

const handleItemStarToggle = async (item) => {
  try {
    const success = await toggleStar(item)
    if (success) {
      // Refresh the current view data
      if (currentView.value === 'my-drive') {
        await fetchFolderContents(currentFolderId.value)
      }
      showNotification('success', 'Star updated', `${item.name} has been ${item.starred ? 'unstarred' : 'starred'}.`)
    } else {
      showNotification('error', 'Star update failed', 'Failed to update star status.')
    }
  } catch (error) {
    console.error('Star toggle error:', error)
    showNotification('error', 'Star update failed', 'An error occurred while updating star status.')
  }
}

const handleRetry = () => {
  fetchFolderContents(currentFolderId.value)
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
    case 'rename':
      openRenameModal(item)
      break
    case 'delete':
      showDeleteConfirmation(item)
      break
    case 'restore':
      restoreItem(item)
      break
    case 'delete-permanent':
      showPermanentDeleteConfirmation(item)
      break
    case 'download':
      console.log('Download item:', item.name)
      break
    case 'star':
      toggleStar(item.id)
      break
    case 'share':
      console.log('Share item:', item.name)
      break
    case 'new-folder':
      openNewFolderModal()
      break
    case 'upload':
      console.log('Upload files')
      break
    case 'select-all':
      console.log('Select all items')
      break
    case 'preview':
      openPreview(item)
      break
    case 'move':
      openMoveModal(item)
      break
    case 'copy':
      console.log('Copy item:', item.name)
      break
    default:
      console.log('Unknown action:', action)
  }
}

const showDeleteConfirmation = (item) => {
  const isFolder = item.type === 'folder'
  confirmModal.value = {
    visible: true,
    title: `Delete ${isFolder ? 'Folder' : 'File'}`,
    message: `Are you sure you want to delete "${item.name}"? This action cannot be undone.`,
    item: item,
    action: 'delete'
  }
}

const showPermanentDeleteConfirmation = (item) => {
  const isFolder = item.type === 'folder'
  confirmModal.value = {
    visible: true,
    title: `Permanently Delete ${isFolder ? 'Folder' : 'File'}`,
    message: `Are you sure you want to permanently delete "${item.name}"? This action cannot be undone and the item will be lost forever.`,
    item: item,
    action: 'delete-permanent'
  }
}

const handleConfirmAction = async () => {
  const { action, item } = confirmModal.value
  
  if (action === 'delete') {
    await deleteItem(item)
  } else if (action === 'delete-permanent') {
    await permanentlyDeleteItem(item)
  }
  
  confirmModal.value.visible = false
}

const deleteItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    const url = item.type === 'folder' 
      ? `${API_BASE_URL}/api/folders/${item.folderId}`
      : `${API_BASE_URL}/api/files/${item.fileId}`
    
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (response.ok) {
      // Show success notification
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been deleted.`
      }
      
      // Refresh the current folder contents
      await fetchFolderContents(currentFolderId.value)
      
      // Clear selection
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to delete item')
    }
  } catch (error) {
    console.error('Delete error:', error)
    
    // Show error notification
    notification.value = {
      visible: true,
      type: 'error',
      title: 'Error',
      message: error.message || 'Failed to delete item. Please try again.'
    }
  }
}

const handleConfirmCancel = () => {
  confirmModal.value.visible = false
}

const handleConfirmClose = () => {
  confirmModal.value.visible = false
}

const restoreItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    const url = item.type === 'folder' 
      ? `${API_BASE_URL}/api/folders/${item.folderId}/restore`
      : `${API_BASE_URL}/api/files/${item.fileId}/restore`
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (response.ok) {
      // Show success notification
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been restored.`
      }
      
      // Refresh the current view
      if (currentView.value === 'trash') {
        // Refresh trash view
        if (trashViewRef.value) {
          trashViewRef.value.refresh()
        }
      } else {
        // Refresh the current folder contents
        await fetchFolderContents(currentFolderId.value)
      }
      
      // Clear selection
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to restore item')
    }
  } catch (error) {
    console.error('Restore error:', error)
    
    // Show error notification
    notification.value = {
      visible: true,
      type: 'error',
      title: 'Error',
      message: error.message || 'Failed to restore item. Please try again.'
    }
  }
}

const permanentlyDeleteItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    const url = item.type === 'folder' 
      ? `${API_BASE_URL}/api/folders/${item.folderId}/permanent`
      : `${API_BASE_URL}/api/files/${item.fileId}/permanent`
    
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (response.ok) {
      // Show success notification
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been permanently deleted.`
      }
      
      // Refresh the current view
      if (currentView.value === 'trash') {
        // Refresh trash view
        if (trashViewRef.value) {
          trashViewRef.value.refresh()
        }
      }
      
      // Clear selection
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to permanently delete item')
    }
  } catch (error) {
    console.error('Permanent delete error:', error)
    
    // Show error notification
    notification.value = {
      visible: true,
      type: 'error',
      title: 'Error',
      message: error.message || 'Failed to permanently delete item. Please try again.'
    }
  }
}

const handleNotificationClose = () => {
  notification.value.visible = false
}

const handleUploadFiles = () => {
  uploadModal.value.visible = true
}

const handleUploadClose = () => {
  uploadModal.value.visible = false
  // Reset the upload modal form
  if (uploadModalRef.value) {
    uploadModalRef.value.resetForm()
  }
}

const handleUploadComplete = async () => {
  // Show success notification
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Upload Complete',
    message: 'Files have been uploaded successfully.'
  }
  
  // Refresh the current folder contents
  await fetchFolderContents(currentFolderId.value)
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

const openRenameModal = (item) => {
  renameModal.value = {
    visible: true,
    item: item
  }
}

const handleRenameSuccess = (updatedItem) => {
  console.log('Rename success handler called with:', updatedItem)
  console.log('Current folders:', folders.value)
  console.log('Current files:', files.value)
  
  // Update the files/folders arrays based on type
  if (updatedItem.type === 'folder') {
    const folderIndex = folders.value.findIndex(f => f.ID === updatedItem.folderId)
    console.log('Folder index found:', folderIndex)
    if (folderIndex !== -1) {
      folders.value[folderIndex] = { ...folders.value[folderIndex], Name: updatedItem.name }
      console.log('Folder updated:', folders.value[folderIndex])
    }
  } else {
    const fileIndex = files.value.findIndex(f => f.ID === updatedItem.fileId)
    console.log('File index found:', fileIndex)
    if (fileIndex !== -1) {
      files.value[fileIndex] = { ...files.value[fileIndex], Name: updatedItem.name }
      console.log('File updated:', files.value[fileIndex])
    }
  }
  
  // Show success notification
  showNotification('success', 'Item Renamed', `"${updatedItem.name}" has been renamed successfully.`)
}

const handleRenameClose = () => {
  renameModal.value.visible = false
  renameModal.value.item = null
}

const showNotification = (type, title, message) => {
  notification.value = {
    visible: true,
    type: type,
    title: title,
    message: message
  }
}

const openNewFolderModal = () => {
  newFolderModal.value.visible = true
}

const handleNewFolderModalClose = () => {
  newFolderModal.value.visible = false
}

const handleFolderCreated = async (folder) => {
  // Show success notification
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Success',
    message: `Folder "${folder.Name}" has been created.`
  }
  
  // Refresh the current folder contents
  await fetchFolderContents(currentFolderId.value)
}

const closePreview = () => {
  previewModal.value.visible = false
  previewModal.value.file = null
}

const openMoveModal = (item) => {
  moveModal.value = {
    visible: true,
    item: item
  }
}

const handleMoveSuccess = async (result) => {
  // Show success notification
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Success',
    message: `${moveModal.value.item.type === 'folder' ? 'Folder' : 'File'} "${moveModal.value.item.name}" has been moved successfully.`
  }
  
  // Refresh the current folder contents
  await fetchFolderContents(currentFolderId.value)
  
  // Clear selection
  clearSelection()
}

const handleMoveModalClose = () => {
  moveModal.value.visible = false
  moveModal.value.item = null
}

// Handle browser back/forward buttons
const handlePopState = () => {
  initializeViewFromURL()
}

// Initialize on mount
onMounted(async () => {
  initializeViewFromURL()
  await navigateToRoot()
  
  // Listen for browser back/forward
  window.addEventListener('popstate', handlePopState)
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
        :search-query="searchQuery"
        :view-mode="viewMode"
        :user="user"
        @search-change="handleSearchChange"
        @view-mode-change="handleViewModeChange"
        @logout="handleLogout"
        @upload-files="handleUploadFiles"
        @new-folder="openNewFolderModal"
      />

      <!-- Breadcrumb -->
      <Breadcrumb 
        v-if="currentView === 'my-drive'"
        :breadcrumbs="breadcrumbs"
        @navigate-breadcrumb="navigateToBreadcrumb"
      />

      <!-- Content Area -->
      <div class="flex-1 flex flex-col">
        <!-- Content Area -->
        <ContentArea 
          v-if="currentView === 'my-drive'"
          :items="filteredItems"
          :selected-items="selectedItems"
          :view-mode="viewMode"
          :loading="loading"
          :error="error"
          :search-query="searchQuery"
          @item-select="handleItemSelect"
          @item-double-click="handleItemDoubleClick"
          @item-star-toggle="handleItemStarToggle"
          @context-menu="handleContextMenu"
          @retry="handleRetry"
        />
        
        <!-- Starred View -->
        <StarredView
          v-else-if="currentView === 'starred'"
          :user="user"
          @logout="handleLogout"
        />
        
        <!-- Trash View -->
        <TrashView
          v-else-if="currentView === 'trash'"
          ref="trashViewRef"
          :view-mode="viewMode"
          @item-select="handleItemSelect"
          @item-double-click="handleItemDoubleClick"
          @item-star-toggle="handleItemStarToggle"
          @context-menu="handleContextMenu"
        />
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
    
    <!-- Upload Modal -->
    <UploadModal
      ref="uploadModalRef"
      :visible="uploadModal.visible"
      :current-folder-id="currentFolderId"
      @close="handleUploadClose"
      @upload-complete="handleUploadComplete"
    />
    
    <!-- Preview Modal -->
    <PreviewModal
      :visible="previewModal.visible"
      :file="previewModal.file"
      @close="closePreview"
    />
    
    <!-- New Folder Modal -->
    <NewFolderModal
      :visible="newFolderModal.visible"
      :current-folder="currentFolder"
      @close="handleNewFolderModalClose"
      @folder-created="handleFolderCreated"
    />
    
    <!-- Rename Modal -->
    <RenameModal
      :visible="renameModal.visible"
      :item="renameModal.item"
      @close="handleRenameClose"
      @rename-success="handleRenameSuccess"
    />
    
    <!-- Move Modal -->
    <MoveModal
      :visible="moveModal.visible"
      :item="moveModal.item"
      @close="handleMoveModalClose"
      @move="handleMoveSuccess"
    />
  </div>
</template>

<style scoped>
/* Close dropdowns when clicking outside */
</style> 