<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  item: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'rename-success'])

const newName = ref('')
const loading = ref(false)
const error = ref('')

// Reset form when modal opens/closes
watch(() => props.visible, (newVal) => {
  if (newVal && props.item) {
    // For files, show only the filename without extension
    if (props.item.type !== 'folder') {
      newName.value = getFilenameWithoutExtension(props.item.name) || ''
    } else {
      // For folders, show the full name
      newName.value = props.item.name || ''
    }
    error.value = ''
  }
})

// Helper function to get filename without extension
const getFilenameWithoutExtension = (filename) => {
  if (!filename) return ''
  const lastDotIndex = filename.lastIndexOf('.')
  return lastDotIndex > 0 ? filename.substring(0, lastDotIndex) : filename
}

// Helper function to get file extension
const getFileExtension = (filename) => {
  if (!filename) return ''
  const lastDotIndex = filename.lastIndexOf('.')
  return lastDotIndex > 0 ? filename.substring(lastDotIndex) : ''
}

// Focus input when modal opens
watch(() => props.visible, async (newVal) => {
  if (newVal && props.item) {
    await nextTick()
    setTimeout(() => {
      const input = document.getElementById('rename-input')
      if (input) {
        input.focus()
        
        // For files, select only the filename part (without extension)
        if (props.item.type !== 'folder') {
          const filenameWithoutExt = getFilenameWithoutExtension(props.item.name)
          const startPos = 0
          const endPos = filenameWithoutExt.length
          
          input.setSelectionRange(startPos, endPos)
        } else {
          // For folders, select all text
          input.select()
        }
      }
    }, 50)
  }
}, { immediate: true })

const handleSubmit = async () => {
  if (!newName.value.trim()) {
    error.value = 'Name cannot be empty'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    
    // For files, reconstruct the full filename with extension
    let finalName = newName.value.trim()
    if (props.item.type !== 'folder') {
      const originalExtension = getFileExtension(props.item.name)
      // Only add extension if user didn't type it
      if (!finalName.includes('.') && originalExtension) {
        finalName = finalName + originalExtension
      }
    }
    
    const endpoint = props.item.type === 'folder' 
      ? `${API_BASE_URL}/api/folders/${props.item.folderId}/rename`
      : `${API_BASE_URL}/api/files/${props.item.fileId}/rename`
    
    const response = await fetch(endpoint, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: finalName
      })
    })

    const data = await response.json()

    if (response.ok) {
      emit('rename-success', {
        ...props.item,
        name: finalName
      })
      emit('close')
    } else {
      error.value = data.error || 'Failed to rename item'
    }
  } catch (err) {
    console.error('Rename error:', err)
    error.value = 'Network error. Please try again.'
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  if (!loading.value) {
    emit('close')
  }
}

const handleKeydown = (event) => {
  if (event.key === 'Escape' && !loading.value) {
    handleClose()
  } else if (event.key === 'Enter' && !loading.value) {
    handleSubmit()
  }
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
      class="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4"
      @click.stop
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
          Rename {{ item?.type === 'folder' ? 'Folder' : 'File' }}
        </h3>
        <button
          v-if="!loading"
          @click="handleClose"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          ✕
        </button>
      </div>
      
      <!-- Content -->
      <div class="p-6">
        <div class="mb-4">
          <label for="rename-input" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            New name
          </label>
          <input
            id="rename-input"
            v-model="newName"
            @keydown="handleKeydown"
            type="text"
            :disabled="loading"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 disabled:opacity-50"
            :placeholder="item?.type === 'folder' ? 'Enter new folder name' : 'Enter new filename (extension will be preserved)'"
          />
        </div>
        
        <!-- Help text for files -->
        <div v-if="item?.type !== 'folder'" class="mb-4">
          <p class="text-xs text-gray-500 dark:text-gray-400">
            💡 Only the filename is shown. 
            <span v-if="getFileExtension(item?.name || '')">
              The file extension ({{ getFileExtension(item?.name || '') }}) will be preserved automatically.
            </span>
            <span v-else>
              This file has no extension.
            </span>
          </p>
        </div>
        
        <!-- Error Message -->
        <div v-if="error" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>
        </div>
      </div>
      
      <!-- Actions -->
      <div class="flex justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          v-if="!loading"
          @click="handleClose"
          class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
        >
          Cancel
        </button>
        <button
          @click="handleSubmit"
          :disabled="loading || !newName.trim()"
          :class="[
            'px-4 py-2 text-sm font-medium rounded-lg transition-colors',
            loading || !newName.trim()
              ? 'bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-400 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700'
          ]"
        >
          <span v-if="loading">Renaming...</span>
          <span v-else>Rename</span>
        </button>
      </div>
    </div>
  </div>
</template> 