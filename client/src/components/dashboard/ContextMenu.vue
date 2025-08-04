<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
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
const menuPosition = ref({ x: 100, y: 100 }) 
const isImageFile = (item) => {
  if (!item || item.type === 'folder') return false
  const imageTypes = ['image', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  return imageTypes.includes(item.type)
}
const calculatePosition = async () => {
  if (!menuRef.value) return
  await nextTick()
  console.log('Calculating position for context menu, props:', { x: props.x, y: props.y })
  const menu = menuRef.value
  const menuRect = menu.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  const estimatedWidth = 200 
  const estimatedHeight = 300 
  const menuWidth = menuRect.width || estimatedWidth
  const menuHeight = menuRect.height || estimatedHeight
  let newX = (typeof props.x === 'number' && !isNaN(props.x)) ? props.x : viewportWidth / 2
  let newY = (typeof props.y === 'number' && !isNaN(props.y)) ? props.y : viewportHeight / 2
  console.log('Initial position:', { newX, newY, menuWidth, menuHeight, viewportWidth, viewportHeight })
  if (newX + menuWidth > viewportWidth) {
    newX = viewportWidth - menuWidth - 10 
  }
  if (newX < 10) {
    newX = 10
  }
  if (newY + menuHeight > viewportHeight) {
    newY = viewportHeight - menuHeight - 10 
  }
  if (newY < 10) {
    newY = 10
  }
  menuPosition.value = { 
    x: Math.max(0, Math.min(newX, viewportWidth - 50)), 
    y: Math.max(0, Math.min(newY, viewportHeight - 50)) 
  }
  console.log('Final position:', menuPosition.value)
  if (menuPosition.value.x < 0 || menuPosition.value.y < 0) {
    const fallbackX = (typeof props.x === 'number' && !isNaN(props.x)) ? Math.max(10, props.x) : Math.max(10, (viewportWidth - menuWidth) / 2)
    const fallbackY = (typeof props.y === 'number' && !isNaN(props.y)) ? Math.max(10, props.y) : Math.max(10, (viewportHeight - menuHeight) / 2)
    menuPosition.value = {
      x: fallbackX,
      y: fallbackY
    }
    console.log('Using fallback position:', menuPosition.value)
  }
}
const updatePosition = async () => {
  if (props.visible) {
    await new Promise(resolve => setTimeout(resolve, 10))
    await calculatePosition()
  }
}
const handleClickOutside = (event) => {
  if (menuRef.value && !menuRef.value.contains(event.target)) {
    emit('close')
  }
}
const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    emit('close')
  }
}
const handleAction = (action) => {
  emit('action', { action, item: props.item })
  emit('close')
}
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
  window.addEventListener('resize', updatePosition)
  window.addEventListener('scroll', updatePosition, true)
  if (props.visible) {
    updatePosition()
  }
})
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', updatePosition)
  window.removeEventListener('scroll', updatePosition, true)
})
watch(() => [props.x, props.y, props.visible], async () => {
  console.log('ContextMenu watch triggered:', { x: props.x, y: props.y, visible: props.visible })
  if (props.visible) {
    console.log('Menu becoming visible with props:', { x: props.x, y: props.y })
    await updatePosition()
  } else {
    menuPosition.value = { x: 0, y: 0 }
  }
}, { immediate: true })
const getMenuItems = () => {
  if (!props.item) {
    return [
      { id: 'new-folder', label: 'New folder', icon: '📁' },
      { id: 'upload', label: 'Upload files', icon: '📤' },
      { id: 'separator' },
      { id: 'select-all', label: 'Select all', icon: '☑️' }
    ]
  }
  const isFolder = props.item.type === 'folder'
  const isInTrash = props.item.deletedAt 
  if (isInTrash) {
    return [
      { id: 'restore', label: 'Restore', icon: '🔄' },
      { id: 'preview', label: 'Preview', icon: '👁️' },
      { id: 'separator' },
      { id: 'delete-permanent', label: 'Delete permanently', icon: '🗑️', danger: true }
    ]
  }
  const commonItems = [
    { id: 'rename', label: 'Rename', icon: '✏️' },
    { id: 'star', label: props.item.starred ? 'Remove star' : 'Add star', icon: props.item.starred ? '⭐' : '☆' },
    { id: 'share', label: 'Share', icon: '👥' },
    { id: 'move', label: 'Move to...', icon: '📁' }
  ]
  const folderItems = [
    { id: 'new-folder', label: 'New folder', icon: '📁' },
    { id: 'upload', label: 'Upload files', icon: '📤' },
    { id: 'select-all', label: 'Select all', icon: '☑️' }
  ]
  const fileItems = [
    { id: 'download', label: 'Download', icon: '⬇️' }
  ]
  if (isImageFile(props.item)) {
    fileItems.push({ id: 'preview', label: 'Preview', icon: '👁️' })
  }
  fileItems.push(
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
    :style="{ left: menuPosition.x + 'px', top: menuPosition.y + 'px' }"
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