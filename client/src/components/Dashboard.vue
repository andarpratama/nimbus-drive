<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from './dashboard/Sidebar.vue'
import Toolbar from './dashboard/Toolbar.vue'
import ContentArea from './dashboard/ContentArea.vue'
import TrashView from './dashboard/TrashView.vue'
import StarredView from './dashboard/StarredView.vue'
import RecentView from './dashboard/RecentView.vue'
import ProfileView from './dashboard/ProfileView.vue'
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

const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['logout'])

const route = useRoute()
const router = useRouter()

const handleLogout = () => {
  emit('logout')
}
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

const currentView = ref('my-drive')
const viewMode = ref('grid')
const searchQuery = ref('')
const starredViewKey = ref(0)

const trashViewRef = ref(null)

const uploadModalRef = ref(null)
const initializeViewFromURL = async () => {
  console.log('Initializing view from URL')
  console.log('Current route path:', route.path)
  console.log('Route query:', route.query)
  
  const view = route.query.view || 'my-drive'
  const folderId = route.query.folder || null
  
  console.log('Route query - view:', view, 'folder:', folderId)
  
  currentView.value = view
  
  // Only fetch data for views that actually need it
  if (view === 'my-drive') {
    if (folderId) {
      console.log('Navigating to folder:', folderId)
      await navigateToFolder(folderId, true)
    } else {
      console.log('Navigating to root')
      await navigateToRoot(true)
    }
  } else if (view === 'trash') {
    // Trash view needs to fetch its own data
    console.log('Initializing trash view - TrashView component will handle its own data fetching')
  } else if (view === 'recent') {
    // Recent view needs to fetch its own data
    console.log('Initializing recent view')
  } else if (view === 'profile') {
    // Profile view doesn't need any file/folder data
    console.log('Initializing profile view - no data fetching needed')
    // No need to call any file manager functions for profile view
  } else if (view === 'starred') {
    // Starred view needs to fetch its own data
    console.log('Initializing starred view')
  }
}

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  item: null
})

const confirmModal = ref({
  visible: false,
  title: '',
  message: '',
  item: null,
  action: ''
})

const notification = ref({
  visible: false,
  type: 'success',
  title: '',
  message: ''
})

const uploadModal = ref({
  visible: false
})

const newFolderModal = ref({
  visible: false
})

const previewModal = ref({
  visible: false,
  file: null
})

const renameModal = ref({
  visible: false,
  item: null
})

const moveModal = ref({
  visible: false,
  item: null
})
const filteredItems = computed(() => {
  if (!searchQuery.value) return allItems.value
  return allItems.value.filter(item => 
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const handleViewChange = (viewId) => {
  currentView.value = viewId
  
  if (viewId === 'my-drive') {
    router.push({
      path: '/dashboard',
      query: { view: 'my-drive' }
    })
    handleNavigateToRoot()
  } else if (viewId === 'starred') {
    router.push({
      path: '/dashboard',
      query: { view: 'starred' }
    })
    starredViewKey.value++
  } else if (viewId === 'recent') {
    router.push({
      path: '/dashboard',
      query: { view: 'recent' }
    })
  } else if (viewId === 'profile') {
    router.push({
      path: '/dashboard',
      query: { view: 'profile' }
    })
  } else if (viewId === 'trash') {
    router.push({
      path: '/dashboard',
      query: { view: 'trash' }
    })
  }
}

const handleSearchChange = (query) => {
  searchQuery.value = query
}

const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

const handleProfile = () => {
  currentView.value = 'profile'
  router.push({
    path: '/dashboard',
    query: { view: 'profile' }
  })
}

const handleItemSelect = (itemId) => {
  selectItem(itemId)
}

const handleItemDoubleClick = (item) => {
  if (item.type === 'folder') {
    handleNavigateToFolder(item.folderId)
  } else {
    if (isImageFile(item)) {
      openPreview(item)
    } else {
      console.log('File clicked:', item)
    }
  }
}

const handleItemStarToggle = async (item) => {
      try {
      const success = await toggleStar(item)
      if (success) {
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
  
  let x = event.clientX
  let y = event.clientY
  
  if (event.type === 'click') {
    x = event.clientX - 200
    y = event.clientY + 10
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
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been deleted.`
      }
      
      await fetchFolderContents(currentFolderId.value)
      
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to delete item')
    }
  } catch (error) {
    console.error('Delete error:', error)
    
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
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been restored.`
      }
      
      if (currentView.value === 'trash') {
        if (trashViewRef.value) {
          trashViewRef.value.refresh()
        }
      } else {
        await fetchFolderContents(currentFolderId.value)
      }
      
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to restore item')
    }
  } catch (error) {
    console.error('Restore error:', error)
    
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
      notification.value = {
        visible: true,
        type: 'success',
        title: 'Success',
        message: `${item.type === 'folder' ? 'Folder' : 'File'} "${item.name}" has been permanently deleted.`
      }
      
      if (currentView.value === 'trash') {
        if (trashViewRef.value) {
          trashViewRef.value.refresh()
        }
      }
      
      clearSelection()
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to permanently delete item')
    }
  } catch (error) {
    console.error('Permanent delete error:', error)
    
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
  if (uploadModalRef.value) {
    uploadModalRef.value.resetForm()
  }
}

const handleUploadComplete = async () => {
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Upload Complete',
    message: 'Files have been uploaded successfully.'
  }
  
  await fetchFolderContents(currentFolderId.value)
}

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
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Success',
    message: `Folder "${folder.Name}" has been created.`
  }
  
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
  notification.value = {
    visible: true,
    type: 'success',
    title: 'Success',
    message: `${moveModal.value.item.type === 'folder' ? 'Folder' : 'File'} "${moveModal.value.item.name}" has been moved successfully.`
  }
  
  await fetchFolderContents(currentFolderId.value)
  
  clearSelection()
}

const handleMoveModalClose = () => {
  moveModal.value.visible = false
  moveModal.value.item = null
}

const handleNavigateToFolder = async (folderId) => {
  console.log('handleNavigateToFolder called with folderId:', folderId)
  
  currentView.value = 'my-drive'
  await navigateToFolder(folderId)
  
  router.push({
    path: '/dashboard',
    query: { 
      view: 'my-drive',
      folder: folderId 
    }
  })
  
  console.log('URL updated for folder navigation')
}

const handleNavigateToRoot = async () => {
  await navigateToRoot()
  router.push({
    path: '/dashboard',
    query: { view: 'my-drive' }
  })
}

const handleNavigateToBreadcrumb = async (folderId) => {
  await navigateToBreadcrumb(folderId)
  
  if (folderId === null) {
    router.push({
      path: '/dashboard',
      query: { view: 'my-drive' }
    })
  } else {
    router.push({
      path: '/dashboard',
      query: { 
        view: 'my-drive',
        folder: folderId 
      }
    })
  }
}

const handlePopState = () => {
  console.log('Pop state detected')
}

watch(() => route.query, async (newQuery, oldQuery) => {
  console.log('Route query changed:', oldQuery, '->', newQuery)
  
  const view = newQuery.view || 'my-drive'
  const folderId = newQuery.folder || null
  
  console.log('New view from query:', view, 'folder:', folderId)
  
  if (view !== currentView.value) {
    currentView.value = view
  }
  
  if (currentView.value === 'my-drive') {
    if (folderId) {
      await navigateToFolder(folderId)
    } else {
      await navigateToRoot()
    }
  }
}, { immediate: false })

onMounted(async () => {
  await initializeViewFromURL()
  
  window.addEventListener('popstate', handlePopState)
})

onUnmounted(() => {
  window.removeEventListener('popstate', handlePopState)
})
</script>

<template>
  <div class="flex h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Left Sidebar -->
    <Sidebar 
      :current-view="currentView"
      @view-change="handleViewChange"
      @navigate-root="handleNavigateToRoot"
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
        @profile="handleProfile"
      />

      <!-- Breadcrumb -->
      <Breadcrumb 
        v-if="currentView === 'my-drive'"
        :breadcrumbs="breadcrumbs"
        @navigate-breadcrumb="handleNavigateToBreadcrumb"
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
          :key="starredViewKey"
          :user="user"
          :reset-key="starredViewKey"
          :search-query="searchQuery"
          :view-mode="viewMode"
          :is-active="currentView === 'starred'"
          @logout="handleLogout"
          @navigate-to-folder="handleNavigateToFolder"
          @search-change="handleSearchChange"
          @view-mode-change="handleViewModeChange"
        />
        
        <!-- Recent View -->
        <RecentView
          v-else-if="currentView === 'recent'"
          :user="user"
          :search-query="searchQuery"
          :view-mode="viewMode"
          :is-active="currentView === 'recent'"
          @logout="handleLogout"
          @navigate-to-folder="handleNavigateToFolder"
          @search-change="handleSearchChange"
          @view-mode-change="handleViewModeChange"
        />
        
        <!-- Profile View -->
        <ProfileView
          v-else-if="currentView === 'profile'"
          :user="user"
          :is-active="currentView === 'profile'"
          @logout="handleLogout"
        />
        
        <!-- Trash View -->
        <TrashView
          v-else-if="currentView === 'trash'"
          ref="trashViewRef"
          :user="user"
          :search-query="searchQuery"
          :view-mode="viewMode"
          :is-active="currentView === 'trash'"
          @logout="handleLogout"
          @item-select="handleItemSelect"
          @item-double-click="handleItemDoubleClick"
          @item-star-toggle="handleItemStarToggle"
          @context-menu="handleContextMenu"
          @search-change="handleSearchChange"
          @view-mode-change="handleViewModeChange"
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