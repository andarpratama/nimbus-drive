<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import ContentArea from './ContentArea.vue'
import ContextMenu from './ContextMenu.vue'
import ConfirmModal from './ConfirmModal.vue'
import Notification from './Notification.vue'
import Breadcrumb from './Breadcrumb.vue'
import PreviewModal from './PreviewModal.vue'
import { useFileManager } from '../../composables/useFileManager'
import { useStarred } from '../../composables/useStarred'
const props = defineProps({
  user: {
    type: Object,
    default: null
  },
  resetKey: {
    type: Number,
    default: 0
  },
  searchQuery: {
    type: String,
    default: ''
  },
  viewMode: {
    type: String,
    default: 'grid'
  }
})
const emit = defineEmits(['logout', 'navigate-to-folder', 'search-change', 'view-mode-change'])
const {
  files,
  folders,
  currentFolderId,
  currentFolder,
  loading: fileLoading,
  error: fileError,
  selectedItems,
  allItems,
  fetchFolderContents,
  navigateToFolder,
  navigateToBreadcrumb,
  navigateToRoot,
  selectItem,
  isSelected,
  clearSelection
} = useFileManager()
const starredBreadcrumbs = computed(() => {
  if (!isInFolder.value) {
    return [{ id: null, name: 'Starred' }]
  }
  return [
    { id: null, name: 'Starred' },
    { id: currentFolderId.value, name: currentFolder.value?.Name || 'Folder' }
  ]
})
const {
  starredFileItems,
  loading: starredLoading,
  error: starredError,
  fetchStarredItems,
  toggleStar
} = useStarred()
const currentView = ref('starred')
const isInFolder = ref(false) 
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
const previewModal = ref({
  visible: false,
  file: null
})
const loading = computed(() => starredLoading.value)
const error = computed(() => starredError.value)
const filteredItems = computed(() => {
  let items = isInFolder.value ? allItems.value : starredFileItems.value
  if (props.searchQuery) {
    items = items.filter(item => 
      item.name.toLowerCase().includes(props.searchQuery.toLowerCase())
    )
  }
  return items
})
const handleLogout = () => {
  emit('logout')
}
const handleViewChange = (viewId) => {
  currentView.value = viewId
  if (viewId === 'starred') {
    isInFolder.value = false
    navigateToRoot()
    fetchStarredItems()
  }
}
const handleSearchChange = (query) => {
  emit('search-change', query)
}
const handleViewModeChange = (mode) => {
  emit('view-mode-change', mode)
}
const handleItemSelect = (itemId) => {
  selectItem(itemId)
}
const handleItemDoubleClick = async (item) => {
  if (item.type === 'folder') {
    isInFolder.value = true
    await navigateToFolder(item.folderId)
    emit('navigate-to-folder', item.folderId)
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
      await fetchStarredItems()
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
  if (isInFolder.value) {
    fetchFolderContents(currentFolderId.value)
  } else {
    fetchStarredItems()
  }
}
const handleBreadcrumbNavigation = async (folderId) => {
  if (folderId === null) {
    isInFolder.value = false
  } else {
    await navigateToFolder(folderId)
  }
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
      console.log('Rename item:', item.name)
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
      console.log('Create new folder')
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
      console.log('Move item:', item.name)
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
    message: `Are you sure you want to permanently delete "${item.name}"? This action cannot be undone.`,
    item: item,
    action: 'delete-permanent'
  }
}
const handleConfirmAction = () => {
  const { action, item } = confirmModal.value
  switch (action) {
    case 'delete':
      deleteItem(item)
      break
    case 'delete-permanent':
      permanentlyDeleteItem(item)
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
const deleteItem = async (item) => {
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    const endpoint = item.type === 'folder' 
      ? `/api/folders/${item.folderId}` 
      : `/api/files/${item.fileId}`
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    if (response.ok) {
      showNotification('success', 'Item moved to trash', `${item.name} has been moved to trash.`)
      if (isInFolder.value) {
        await fetchFolderContents(currentFolderId.value)
      } else {
        await fetchStarredItems()
      }
    } else {
      showNotification('error', 'Delete failed', 'Failed to delete the item.')
    }
  } catch (error) {
    console.error('Delete error:', error)
    showNotification('error', 'Delete failed', 'An error occurred while deleting the item.')
  }
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
      if (isInFolder.value) {
        await fetchFolderContents(currentFolderId.value)
      } else {
        await fetchStarredItems()
      }
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
      if (isInFolder.value) {
        await fetchFolderContents(currentFolderId.value)
      } else {
        await fetchStarredItems()
      }
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
onMounted(async () => {
  await fetchStarredItems()
})
watch(() => props.resetKey, () => {
  isInFolder.value = false
  navigateToRoot()
  fetchStarredItems()
}, { immediate: true })
</script>
<template>
  <div class="flex-1 flex flex-col">
    <!-- Breadcrumb -->
    <Breadcrumb 
      v-if="isInFolder"
      :breadcrumbs="starredBreadcrumbs"
      @navigate-breadcrumb="handleBreadcrumbNavigation"
    />
    <!-- Content Area -->
    <div class="flex-1 flex flex-col">
      <!-- Content Area -->
      <ContentArea 
        :items="filteredItems"
        :selected-items="selectedItems"
        :view-mode="props.viewMode"
        :loading="loading"
        :error="error"
        :search-query="props.searchQuery"
        @item-select="handleItemSelect"
        @item-double-click="handleItemDoubleClick"
        @item-star-toggle="handleItemStarToggle"
        @context-menu="handleContextMenu"
        @retry="handleRetry"
      />
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
<style scoped>
/* Close dropdowns when clicking outside */
</style> 