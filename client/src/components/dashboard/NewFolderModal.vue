<template>
  <Transition name="modal-fade">
    <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="handleBackdropClick">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4" @click.stop>
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
          Create new folder
        </h3>
        <button
          @click="handleClose"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          ✕
        </button>
      </div>
      <!-- Content -->
      <div class="p-6">
        <div class="mb-4">
          <label for="folderName" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Folder name
          </label>
          <input
            id="folderName"
            v-model="folderName"
            @keyup.enter="handleCreate"
            @keyup.esc="handleClose"
            type="text"
            placeholder="Enter folder name"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': error }"
            autofocus
          >
          <p v-if="error" class="mt-1 text-sm text-red-600 dark:text-red-400">
            {{ error }}
          </p>
        </div>
      </div>
      <!-- Footer -->
      <div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          @click="handleClose"
          class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
        >
          Cancel
        </button>
        <button
          @click="handleCreate"
          :disabled="!folderName.trim() || loading"
          class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed rounded-lg transition-colors"
        >
          <span v-if="loading">Creating...</span>
          <span v-else>Create folder</span>
        </button>
      </div>
      </div>
    </div>
  </Transition>
</template>
<style scoped>
/* Modal fade animation */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
.modal-fade-enter-to,
.modal-fade-leave-from {
  opacity: 1;
}
</style>
<script setup>
import { ref, watch, nextTick } from 'vue'
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  currentFolder: {
    type: Object,
    default: null
  }
})
const emit = defineEmits(['close', 'folder-created'])
const folderName = ref('')
const error = ref('')
const loading = ref(false)
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    folderName.value = ''
    error.value = ''
    loading.value = false
    nextTick(() => {
      const input = document.getElementById('folderName')
      if (input) {
        input.focus()
      }
    })
  }
})
const validateFolderName = () => {
  const name = folderName.value.trim()
  if (!name) {
    error.value = 'Folder name is required'
    return false
  }
  if (name.length > 255) {
    error.value = 'Folder name must be less than 255 characters'
    return false
  }
  const invalidChars = /[<>:"/\\|?*]/
  if (invalidChars.test(name)) {
    error.value = 'Folder name contains invalid characters'
    return false
  }
  error.value = ''
  return true
}
const handleCreate = async () => {
  if (!validateFolderName()) {
    return
  }
  loading.value = true
  error.value = ''
  try {
          const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    const payload = {
      name: folderName.value.trim(),
      parent_id: props.currentFolder?.ID || null
    }
    const response = await fetch(`${API_BASE_URL}/api/folders`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })
    if (response.ok) {
      const data = await response.json()
      emit('folder-created', data.folder)
      handleClose()
    } else {
      const errorData = await response.json()
      error.value = errorData.error || 'Failed to create folder'
    }
  } catch (err) {
    console.error('Error creating folder:', err)
    error.value = 'Failed to create folder. Please try again.'
  } finally {
    loading.value = false
  }
}
const handleClose = () => {
  emit('close')
}
const handleBackdropClick = () => {
  emit('close')
}
</script> 