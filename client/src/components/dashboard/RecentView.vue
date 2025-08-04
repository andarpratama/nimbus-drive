<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import ContentArea from './ContentArea.vue'
import ContextMenu from './ContextMenu.vue'
import ConfirmModal from './ConfirmModal.vue'
import Notification from './Notification.vue'
import Breadcrumb from './Breadcrumb.vue'
import PreviewModal from './PreviewModal.vue'
import { useFileManager } from '../../composables/useFileManager'
import { useRecent } from '../../composables/useRecent'

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

const {
  recentFiles,
  loading: recentLoading,
  error: recentError,
  fetchRecentFiles
} = useRecent()

const currentView = ref('recent')
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

const loading = computed(() => recentLoading.value)
const error = computed(() => recentError.value)

// Debug computed property for preview modal
const previewModalDebug = computed(() => {
  console.log('Preview modal state changed:', previewModal.value)
  return previewModal.value
})

const filteredItems = computed(() => {
  let items = recentFiles.value
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

const handleNavigateToFolder = (folderId) => {
  emit('navigate-to-folder', folderId)
}

const handleSearchChange = (query) => {
  emit('search-change', query)
}

const handleViewModeChange = (mode) => {
  emit('view-mode-change', mode)
}

const handleContextMenu = (data) => {
  console.log('handleContextMenu called with data:', data)
  const { event, item } = data
  
  if (event) {
    event.preventDefault()
  }
  
  let x = event?.clientX || 0
  let y = event?.clientY || 0
  
  if (event?.type === 'click') {
    x = (event.clientX || 0) - 200
    y = (event.clientY || 0) + 10
  }
  
  console.log('Setting context menu position:', { x, y, item })
  contextMenu.value = {
    visible: true,
    x: x,
    y: y,
    item: item
  }
  console.log('Context menu state:', contextMenu.value)
}

const handleContextMenuClose = () => {
  console.log('handleContextMenuClose called')
  contextMenu.value.visible = false
}

const handleContextMenuAction = (data) => {
  console.log('handleContextMenuAction called with data:', data)
  const { action, item } = data
  console.log('Context menu action:', action, 'on item:', item)
  
  switch (action) {
    case 'preview':
      console.log('Opening preview for item:', item)
      openPreview(item)
      break
    case 'download':
      console.log('Download item:', item.name)
      break
    case 'rename':
      console.log('Rename item:', item.name)
      break
    case 'delete':
      console.log('Delete item:', item.name)
      break
    case 'star':
      console.log('Star item:', item.name)
      break
    case 'share':
      console.log('Share item:', item.name)
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

const handleItemDoubleClick = (item) => {
  console.log('handleItemDoubleClick called with item:', item)
  if (item.type === 'folder') {
    navigateToFolder(item.id)
  } else {
    // Handle file double click - could open preview or download
    if (isImageFile(item)) {
      console.log('Item is image file, opening preview')
      openPreview(item)
    } else {
      // For non-image files, you might want to download or show a different preview
      console.log('File clicked:', item)
    }
  }
}

const handleItemSelect = (itemId) => {
  selectItem(itemId)
}

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

const isImageFile = (item) => {
  console.log('isImageFile called with item:', item)
  if (!item || item.type === 'folder') {
    console.log('Item is null or folder, returning false')
    return false
  }
  const imageTypes = ['image', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  const isImage = imageTypes.includes(item.type)
  console.log('Item type:', item.type, 'isImage:', isImage)
  return isImage
}

const openPreview = (item) => {
  console.log('openPreview called with item:', item)
  if (item.type === 'folder') return
  console.log('Setting preview modal visible with file:', item)
  previewModal.value = {
    visible: true,
    file: item
  }
  console.log('Preview modal state:', previewModal.value)
}

const closePreviewModal = () => {
  previewModal.value.visible = false
  previewModal.value.file = null
}

const handleConfirmAction = () => {
  // Handle confirmed actions
  confirmModal.value.visible = false
}

const closeConfirmModal = () => {
  confirmModal.value.visible = false
}

// Watch for prop changes and refetch data
watch(() => props.resetKey, () => {
  if (currentView.value === 'recent') {
    fetchRecentFiles()
  }
})

onMounted(() => {
  fetchRecentFiles()
})
</script>

<template>
  <div class="h-full flex flex-col">
    <!-- Content Area -->
    <ContentArea
      :items="filteredItems"
      :loading="loading"
      :error="error"
      :view-mode="viewMode"
      :selected-items="selectedItems"
      @item-double-click="handleItemDoubleClick"
      @item-select="handleItemSelect"
      @context-menu="handleContextMenu"
      @retry="fetchRecentFiles"
    />

    <!-- Context Menu -->
    <ContextMenu
      :visible="contextMenu.visible"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :item="contextMenu.item"
      @close="handleContextMenuClose"
      @action="handleContextMenuAction"
    />

    <!-- Confirm Modal -->
    <ConfirmModal
      v-if="confirmModal.visible"
      :title="confirmModal.title"
      :message="confirmModal.message"
      @confirm="handleConfirmAction"
      @close="closeConfirmModal"
    />

    <!-- Notification -->
    <Notification
      v-if="notification.visible"
      :type="notification.type"
      :title="notification.title"
      :message="notification.message"
    />

    <!-- Preview Modal -->
    <PreviewModal
      :visible="previewModal.visible"
      :file="previewModal.file"
      @close="closePreviewModal"
    />
  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
</style> 