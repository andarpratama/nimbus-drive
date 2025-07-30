<script setup>
import { computed } from 'vue'
import FileIcon from '../FileIcon.vue'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  isSelected: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['select', 'double-click', 'star-toggle', 'context-menu'])

// Check if file is an image
const isImage = computed(() => {
  return props.item.type === 'image'
})

// Get image URL for preview
const getImageUrl = computed(() => {
  if (!isImage.value || !props.item.fileId) return null
  
  // Construct the image URL for the image
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${API_BASE_URL}/api/files/${props.item.fileId}/image`
})

const handleClick = () => {
  emit('select', props.item.id)
}

const handleDoubleClick = () => {
  emit('double-click', props.item)
}

const handleStarToggle = (event) => {
  event.stopPropagation()
  emit('star-toggle', props.item)
}

const handleContextMenu = (event) => {
  event.preventDefault()
  event.stopPropagation()
  emit('context-menu', { event, item: props.item })
}
</script>

<template>
  <div
    @click="handleClick"
    @dblclick="handleDoubleClick"
    @contextmenu="handleContextMenu"
    :class="[
      'relative p-4 rounded-lg border-2 cursor-pointer transition-all hover:shadow-md group',
      isSelected
        ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
        : 'border-transparent hover:border-gray-300 dark:hover:border-gray-600 bg-white dark:bg-gray-800'
    ]"
  >
    <!-- Selection indicator -->
    <div v-if="isSelected" class="absolute top-2 left-2 w-4 h-4 bg-blue-600 rounded-full flex items-center justify-center">
      <span class="text-white text-xs">✓</span>
    </div>
    
    <!-- File icon or image preview -->
    <div class="text-center mb-3">
      <!-- Image preview -->
      <div v-if="isImage && getImageUrl" class="mb-2">
        <img 
          :src="getImageUrl" 
          :alt="item.name"
          class="w-16 h-16 object-cover rounded mx-auto border border-gray-200 dark:border-gray-600"
          @error="$event.target.style.display = 'none'"
        />
      </div>
      <!-- File icon for non-images -->
      <div v-else class="mb-2 flex justify-center">
        <FileIcon v-if="item.type !== 'folder'" :filename="item.name" size="xl" />
        <div v-else class="text-4xl">📁</div>
      </div>
      <div class="text-xs text-gray-500 dark:text-gray-400">
        <span v-if="item.type === 'folder'">{{ item.size }} {{ item.size === 1 ? 'Item' : 'Items' }}</span>
        <span v-else>{{ item.size }}</span>
      </div>
    </div>
    
    <!-- File name -->
    <div class="text-sm font-medium text-gray-900 dark:text-white text-center truncate" :title="item.name">
      {{ item.name }}
    </div>
    
    <!-- File metadata -->
    <div class="text-xs text-gray-500 dark:text-gray-400 text-center mt-1">
      {{ item.modified }}
    </div>
    
    <!-- Action buttons (show on hover) -->
    <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
      <button
        @click="handleStarToggle"
        :class="[
          'p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700',
          item.starred ? 'text-yellow-500' : 'text-gray-400'
        ]"
      >
        {{ item.starred ? '⭐' : '☆' }}
      </button>
    </div>
  </div>
</template> 