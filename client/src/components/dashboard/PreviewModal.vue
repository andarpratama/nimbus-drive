<script setup>
import { computed } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  file: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close'])

const isImage = computed(() => {
  if (!props.file) return false
  const imageTypes = ['image', 'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  return imageTypes.includes(props.file.type)
})

const imageUrl = computed(() => {
  if (!props.file || !isImage.value) return ''
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  return `${API_BASE_URL}/api/files/${props.file.fileId}/image`
})

const handleClose = () => {
  emit('close')
}

const handleBackdropClick = (event) => {
  if (event.target === event.currentTarget) {
    handleClose()
  }
}

const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    handleClose()
  }
}
</script>

<template>
  <div 
    v-if="visible" 
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75"
    @click="handleBackdropClick"
    @keydown="handleKeydown"
    tabindex="0"
  >
    <!-- Modal -->
    <div class="relative max-w-4xl max-h-[90vh] w-full mx-4 bg-white dark:bg-gray-800 rounded-lg shadow-2xl overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center space-x-3">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ file?.name || 'Preview' }}
          </h3>
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ file?.size || '' }}
          </span>
        </div>
        <button
          @click="handleClose"
          class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <span class="text-2xl">×</span>
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto">
        <!-- Image Preview -->
        <div v-if="isImage" class="flex justify-center items-center p-6">
          <img
            :src="imageUrl"
            :alt="file?.name"
            class="max-w-full max-h-[70vh] object-contain rounded-lg shadow-lg"
            @error="$event.target.style.display = 'none'"
          />
        </div>

        <!-- Unsupported File Type -->
        <div v-else class="flex justify-center items-center py-20">
          <div class="text-center">
            <div class="text-6xl mb-4">📄</div>
            <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
              Preview not available
            </h3>
            <p class="text-gray-500 dark:text-gray-400">
              Preview is not supported for this file type
            </p>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between p-4 border-t border-gray-200 dark:border-gray-700">
        <div class="text-sm text-gray-500 dark:text-gray-400">
          {{ file?.modified || '' }}
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="handleClose"
            class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Focus styles for keyboard navigation */
div:focus {
  outline: none;
}
</style> 