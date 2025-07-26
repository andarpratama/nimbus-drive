<script setup>
import { computed } from 'vue'

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

const emit = defineEmits(['select', 'double-click', 'star-toggle'])

const getFileIcon = (type) => {
  const icons = {
    document: '📄',
    pdf: '📕',
    folder: '📁',
    image: '🖼️',
    spreadsheet: '📊',
    presentation: '📽️',
    video: '🎥',
    audio: '🎵',
    archive: '📦'
  }
  return icons[type] || '📄'
}

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
  emit('star-toggle', props.item.id)
}
</script>

<template>
  <div
    @click="handleClick"
    @dblclick="handleDoubleClick"
    :class="[
      'flex items-center p-3 rounded-lg cursor-pointer transition-colors',
      isSelected
        ? 'bg-blue-100 dark:bg-blue-900/20'
        : 'hover:bg-gray-100 dark:hover:bg-gray-800'
    ]"
  >
    <!-- Selection checkbox -->
    <div class="mr-3">
      <div :class="[
        'w-4 h-4 rounded border-2 flex items-center justify-center',
        isSelected
          ? 'bg-blue-600 border-blue-600'
          : 'border-gray-300 dark:border-gray-600'
      ]">
        <span v-if="isSelected" class="text-white text-xs">✓</span>
      </div>
    </div>
    
    <!-- File icon or image preview -->
    <div class="mr-4">
      <!-- Image preview -->
      <img 
        v-if="isImage && getImageUrl" 
        :src="getImageUrl" 
        :alt="item.name"
        class="w-8 h-8 object-cover rounded border border-gray-200 dark:border-gray-600"
        @error="$event.target.style.display = 'none'"
      />
      <!-- File icon for non-images -->
      <div v-else class="text-2xl">{{ getFileIcon(item.type) }}</div>
    </div>
    
    <!-- File info -->
    <div class="flex-1 min-w-0">
      <div class="flex items-center">
        <div class="text-sm font-medium text-gray-900 dark:text-white truncate">
          {{ item.name }}
        </div>
        <button
          @click="handleStarToggle"
          :class="[
            'ml-2 p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700',
            item.starred ? 'text-yellow-500' : 'text-gray-400'
          ]"
        >
          {{ item.starred ? '⭐' : '☆' }}
        </button>
      </div>
      <div class="text-xs text-gray-500 dark:text-gray-400">
        {{ item.size }} • {{ item.modified }}
      </div>
    </div>
    
    <!-- Shared indicator -->
    <div v-if="item.shared" class="mr-3 text-blue-600">
      👥
    </div>
    
    <!-- More options -->
    <button class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-400">
      ⋮
    </button>
  </div>
</template> 