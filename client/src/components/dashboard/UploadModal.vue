<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import FileIcon from '../FileIcon.vue'
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  currentFolderId: {
    type: [String, null],
    default: null
  }
})
const emit = defineEmits(['close', 'upload-complete'])
const dragActive = ref(false)
const files = ref([])
const uploadProgress = ref({})
const uploading = ref(false)
const error = ref('')
const fileInputRef = ref(null)
const hasFiles = computed(() => files.value.length > 0)
const totalProgress = computed(() => {
  if (files.value.length === 0) return 0
  const total = Object.values(uploadProgress.value).reduce((sum, progress) => sum + progress, 0)
  return Math.round(total / files.value.length)
})
const handleFileSelect = (event) => {
  const selectedFiles = Array.from(event.target.files)
  addFiles(selectedFiles)
}
const addFiles = (newFiles) => {
  const validFiles = newFiles.filter(file => {
    if (file.size > 100 * 1024 * 1024) {
      error.value = `File "${file.name}" is too large. Maximum size is 100MB.`
      return false
    }
    return true
  })
  files.value.push(...validFiles)
  validFiles.forEach(file => {
    uploadProgress.value[file.name] = 0
  })
  error.value = ''
}
const handleDrag = (e) => {
  e.preventDefault()
  e.stopPropagation()
}
const handleDragIn = (e) => {
  e.preventDefault()
  e.stopPropagation()
  if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
    dragActive.value = true
  }
}
const handleDragOut = (e) => {
  e.preventDefault()
  e.stopPropagation()
  dragActive.value = false
}
const handleDrop = (e) => {
  e.preventDefault()
  e.stopPropagation()
  dragActive.value = false
  if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
    addFiles(Array.from(e.dataTransfer.files))
  }
}
const removeFile = (index) => {
  const file = files.value[index]
  delete uploadProgress.value[file.name]
  files.value.splice(index, 1)
}
const clearFiles = () => {
  files.value = []
  uploadProgress.value = {}
  error.value = ''
}
const uploadFiles = async () => {
  if (files.value.length === 0) return
  uploading.value = true
  error.value = ''
  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    for (let i = 0; i < files.value.length; i++) {
      const file = files.value[i]
      const formData = new FormData()
      formData.append('file', file)
      if (props.currentFolderId) {
        formData.append('folder_id', props.currentFolderId.toString())
      }
      const response = await fetch(`${API_BASE_URL}/api/files/upload`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      })
      if (response.ok) {
        uploadProgress.value[file.name] = 100
      } else {
        const errorData = await response.json().catch(() => ({}))
        error.value = `Failed to upload "${file.name}": ${errorData.error || response.statusText}`
        throw new Error(`Upload failed for ${file.name}`)
      }
    }
    emit('upload-complete')
    clearFiles() 
    emit('close')
  } catch (err) {
    console.error('Upload error:', err)
    if (!error.value) {
      error.value = 'Upload failed. Please try again.'
    }
  } finally {
    uploading.value = false
  }
}
const handleClose = () => {
  if (!uploading.value) {
    clearFiles()
    emit('close')
  }
}
const resetForm = () => {
  clearFiles()
  uploading.value = false
  error.value = ''
}
defineExpose({
  resetForm
})
const handleKeydown = (event) => {
  if (event.key === 'Escape' && !uploading.value) {
    handleClose()
  }
}
onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>
<template>
  <div
    v-if="visible"
    class="fixed inset-0 z-50 flex items-center justify-center"
    @click="handleClose"
  >
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black bg-opacity-50"></div>
    <!-- Modal -->
    <div
      class="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-hidden"
      @click.stop
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
          Upload Files
        </h3>
        <button
          v-if="!uploading"
          @click="handleClose"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          ✕
        </button>
      </div>
      <!-- Content -->
      <div class="p-6">
        <!-- Drag & Drop Area -->
        <div
          @dragenter="handleDragIn"
          @dragleave="handleDragOut"
          @dragover="handleDrag"
          @drop="handleDrop"
          :class="[
            'border-2 border-dashed rounded-lg p-4 text-center transition-colors',
            dragActive
              ? 'border-blue-400 bg-blue-50 dark:bg-blue-900/20'
              : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
          ]"
        >
          <div class="text-3xl mb-2">📁</div>
          <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-1">
            Drop files here or click to select
          </h4>
          <p class="text-xs text-gray-600 dark:text-gray-400 mb-3">
            Maximum file size: 100MB
          </p>
          <button
            @click="fileInputRef?.click()"
            class="px-3 py-1.5 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
          >
            Choose Files
          </button>
          <input
            ref="fileInputRef"
            type="file"
            multiple
            @change="handleFileSelect"
            class="hidden"
          />
        </div>
        <!-- File List -->
        <div v-if="hasFiles" class="mt-6">
          <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-3">
            Files to Upload ({{ files.length }})
          </h4>
          <div class="space-y-2 max-h-64 overflow-y-auto">
            <div
              v-for="(file, index) in files"
              :key="file.name"
              class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg"
            >
              <div class="flex items-center flex-1 min-w-0">
                <div class="mr-3">
                  <FileIcon :filename="file.name" size="md" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900 dark:text-white truncate">
                    {{ file.name }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    {{ formatFileSize(file.size) }}
                  </p>
                </div>
              </div>
              <!-- Progress Bar -->
              <div v-if="uploading" class="flex items-center ml-4">
                <div class="w-16 h-2 bg-gray-200 dark:bg-gray-600 rounded-full mr-2">
                  <div
                    class="h-2 bg-blue-600 rounded-full transition-all duration-300"
                    :style="{ width: uploadProgress[file.name] + '%' }"
                  ></div>
                </div>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ uploadProgress[file.name] }}%
                </span>
              </div>
              <!-- Remove Button -->
              <button
                v-if="!uploading"
                @click="removeFile(index)"
                class="ml-2 text-gray-400 hover:text-red-600 dark:hover:text-red-400"
              >
                ✕
              </button>
            </div>
          </div>
        </div>
        <!-- Error Message -->
        <div v-if="error" class="mt-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>
        </div>
      </div>
      <!-- Actions -->
      <div class="flex justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          v-if="!uploading"
          @click="handleClose"
          class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
        >
          Cancel
        </button>
        <button
          @click="uploadFiles"
          :disabled="!hasFiles || uploading"
          :class="[
            'px-4 py-2 text-sm font-medium rounded-lg transition-colors',
            hasFiles && !uploading
              ? 'bg-blue-600 text-white hover:bg-blue-700'
              : 'bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-400 cursor-not-allowed'
          ]"
        >
          <span v-if="uploading">Uploading... ({{ totalProgress }}%)</span>
          <span v-else>Upload {{ files.length }} File{{ files.length !== 1 ? 's' : '' }}</span>
        </button>
      </div>
    </div>
  </div>
</template> 