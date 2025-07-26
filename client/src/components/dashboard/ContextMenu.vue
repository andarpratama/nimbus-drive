<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  x: {
    type: Number,
    default: 0
  },
  y: {
    type: Number,
    default: 0
  },
  item: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'action'])

const menuRef = ref(null)

// Helper function to check if file is an image
const isImageFile = (item) => {
  if (!item || item.type === 'folder') return false
  const imageTypes = ['image', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  return imageTypes.includes(item.type)
}

// Handle click outside to close menu
const handleClickOutside = (event) => {
  if (menuRef.value && !menuRef.value.contains(event.target)) {
    emit('close')
  }
}

// Handle escape key to close menu
const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    emit('close')
  }
}

// Handle menu action
const handleAction = (action) => {
  emit('action', { action, item: props.item })
  emit('close')
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})

// Get menu items based on item type
const getMenuItems = () => {
  if (!props.item) {
    // Context menu for empty space
    return [
      { id: 'new-folder', label: 'New folder', icon: '📁' },
      { id: 'upload', label: 'Upload files', icon: '📤' },
      { id: 'separator' },
      { id: 'select-all', label: 'Select all', icon: '☑️' }
    ]
  }
  
  const isFolder = props.item.type === 'folder'
  const isInTrash = props.item.deletedAt // Check if item is in trash
  
  if (isInTrash) {
    // Trash-specific menu items
    return [
      { id: 'restore', label: 'Restore', icon: '🔄' },
      { id: 'preview', label: 'Preview', icon: '👁️' },
      { id: 'separator' },
      { id: 'delete-permanent', label: 'Delete permanently', icon: '🗑️', danger: true }
    ]
  }
  
  // Normal menu items
  const commonItems = [
    { id: 'rename', label: 'Rename', icon: '✏️' },
    { id: 'star', label: props.item.starred ? 'Remove star' : 'Add star', icon: props.item.starred ? '⭐' : '☆' },
    { id: 'share', label: 'Share', icon: '👥' }
  ]
  
  const folderItems = [
    { id: 'new-folder', label: 'New folder', icon: '📁' },
    { id: 'upload', label: 'Upload files', icon: '📤' },
    { id: 'select-all', label: 'Select all', icon: '☑️' }
  ]
  
  const fileItems = [
    { id: 'download', label: 'Download', icon: '⬇️' }
  ]
  
  // Add preview option only for image files
  if (isImageFile(props.item)) {
    fileItems.push({ id: 'preview', label: 'Preview', icon: '👁️' })
  }
  
  fileItems.push(
    { id: 'move', label: 'Move to...', icon: '📁' },
    { id: 'copy', label: 'Copy', icon: '📋' }
  )
  
  const deleteItem = { id: 'delete', label: 'Delete', icon: '🗑️', danger: true }
  
  return [
    ...commonItems,
    ...(isFolder ? folderItems : fileItems),
    { id: 'separator' },
    deleteItem
  ]
}
</script>

<template>
  <div
    v-if="visible"
    ref="menuRef"
    :style="{ left: x + 'px', top: y + 'px' }"
    class="fixed z-50 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg py-1 min-w-48"
  >
    <template v-for="item in getMenuItems()" :key="item.id">
      <!-- Menu item -->
      <div
        v-if="item.id !== 'separator'"
        @click="handleAction(item.id)"
        :class="[
          'flex items-center px-3 py-1 text-sm cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors',
          item.danger ? 'text-red-600 dark:text-red-400' : 'text-gray-700 dark:text-gray-300'
        ]"
      >
        <span class="mr-3 text-base">{{ item.icon }}</span>
        <span>{{ item.label }}</span>
      </div>
      
      <!-- Separator -->
      <div
        v-else
        class="border-t border-gray-200 dark:border-gray-700 my-0.5"
      ></div>
    </template>
  </div>
</template> 