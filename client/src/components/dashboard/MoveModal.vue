<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
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
const emit = defineEmits(['close', 'move'])
const selectedFolderId = ref(null)
const folders = ref([])
const loading = ref(false)
const error = ref('')
const currentPath = ref([])
const currentFolderId = ref(null)
const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    emit('close')
  }
}
const fetchFolders = async () => {
  try {
    loading.value = true
    error.value = ''
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    const response = await fetch(`${API_BASE_URL}/api/folders`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    if (response.ok) {
      const data = await response.json()
      folders.value = data.folders.filter(folder => {
        if (!props.item) return true
        if (props.item.type === 'folder' && folder.ID === props.item.folderId) return false
        return true
      })
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to fetch folders')
    }
  } catch (err) {
    console.error('Error fetching folders:', err)
    error.value = err.message || 'Failed to fetch folders'
  } finally {
    loading.value = false
  }
}
const rootFolders = computed(() => {
  return folders.value.filter(folder => !folder.ParentID)
})
const getSubfolders = (parentId) => {
  return folders.value.filter(folder => folder.ParentID === parentId)
}
const navigateToFolder = (folderId) => {
  const folder = folders.value.find(f => f.ID === folderId)
  if (folder) {
    currentFolderId.value = folderId
    const newPath = []
    let currentFolder = folder
    while (currentFolder) {
      newPath.unshift(currentFolder)
      currentFolder = folders.value.find(f => f.ID === currentFolder.ParentID)
    }
    currentPath.value = newPath
  }
}
const navigateToRoot = () => {
  currentFolderId.value = null
  currentPath.value = []
}
const selectFolder = (folderId) => {
  selectedFolderId.value = folderId
}
const handleMove = async () => {
  if (!props.item) return
  try {
    loading.value = true
    error.value = ''
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const token = localStorage.getItem('token')
    let url, method, body
    if (props.item.type === 'folder') {
      url = `${API_BASE_URL}/api/folders/${props.item.folderId}/move`
      method = 'PATCH'
      body = JSON.stringify({
        new_parent_id: selectedFolderId.value
      })
    } else {
      url = `${API_BASE_URL}/api/files/${props.item.fileId}/move`
      method = 'PUT'
      body = JSON.stringify({
        folder_id: selectedFolderId.value
      })
    }
    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body
    })
    if (response.ok) {
      const data = await response.json()
      emit('move', { success: true, data })
      emit('close')
    } else {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to move item')
    }
  } catch (err) {
    console.error('Error moving item:', err)
    error.value = err.message || 'Failed to move item'
  } finally {
    loading.value = false
  }
}
const handleCancel = () => {
  selectedFolderId.value = null
  error.value = ''
  emit('close')
}
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    selectedFolderId.value = null
    error.value = ''
    currentFolderId.value = null
    currentPath.value = []
    fetchFolders()
  }
})
onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
const currentSubfolders = computed(() => {
  if (!currentFolderId.value) {
    return rootFolders.value
  }
  return getSubfolders(currentFolderId.value)
})
const selectedFolderName = computed(() => {
  if (!selectedFolderId.value) return 'Root Folder'
  const folder = folders.value.find(f => f.ID === selectedFolderId.value)
  return folder ? folder.Name : 'Unknown Folder'
})
const isFolderSelected = (folderId) => {
  return selectedFolderId.value === folderId
}
</script>
<template>
  <div
    v-if="visible"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="handleCancel"
  >
    <div
      class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] flex flex-col"
      @click.stop
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          Move {{ item?.type === 'folder' ? 'Folder' : 'File' }}
        </h3>
        <button
          @click="handleCancel"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      <!-- Content -->
      <div class="flex-1 overflow-hidden">
        <div class="p-6">
          <div class="mb-4">
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Move "{{ item?.name }}" to:
            </p>
            <!-- Breadcrumb -->
            <div class="flex items-center space-x-2 mb-4 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
              <button
                @click="navigateToRoot"
                class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm font-medium"
              >
                My Drive
              </button>
              <span v-for="(folder, index) in currentPath" :key="folder.ID" class="flex items-center">
                <svg class="w-4 h-4 text-gray-400 mx-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                </svg>
                <button
                  @click="navigateToFolder(folder.ID)"
                  class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm font-medium"
                >
                  {{ folder.Name }}
                </button>
              </span>
            </div>
            <!-- Folder List -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg max-h-96 overflow-y-auto">
              <!-- Root option -->
              <div
                @click="selectFolder(null)"
                :class="[
                  'flex items-center p-3 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors border-b border-gray-100 dark:border-gray-600',
                  isFolderSelected(null) ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800' : ''
                ]"
              >
                <div class="flex items-center flex-1">
                  <svg class="w-5 h-5 text-gray-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                  </svg>
                  <span class="text-sm text-gray-900 dark:text-white font-medium">Root Folder</span>
                </div>
                <div v-if="isFolderSelected(null)" class="text-blue-600">
                  <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                  </svg>
                </div>
              </div>
              <!-- Folders -->
              <div
                v-for="folder in currentSubfolders"
                :key="folder.ID"
                @click="selectFolder(folder.ID)"
                :class="[
                  'flex items-center p-3 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors border-b border-gray-100 dark:border-gray-600',
                  isFolderSelected(folder.ID) ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800' : ''
                ]"
              >
                <div class="flex items-center flex-1">
                  <svg class="w-5 h-5 text-gray-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                  </svg>
                  <span class="text-sm text-gray-900 dark:text-white">{{ folder.Name }}</span>
                  <span class="ml-2 text-xs text-gray-500 dark:text-gray-400">({{ folder.total_items }} items)</span>
                </div>
                <div class="flex items-center space-x-2">
                  <!-- Navigate button -->
                  <button
                    @click.stop="navigateToFolder(folder.ID)"
                    class="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                    title="Open folder"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                    </svg>
                  </button>
                  <!-- Selection indicator -->
                  <div v-if="isFolderSelected(folder.ID)" class="text-blue-600">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                </div>
              </div>
              <!-- Empty state -->
              <div v-if="currentSubfolders.length === 0" class="p-8 text-center text-gray-500 dark:text-gray-400">
                <svg class="w-12 h-12 mx-auto mb-3 text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                </svg>
                <p class="text-sm">No folders in this location</p>
              </div>
            </div>
            <!-- Selected destination -->
            <div v-if="selectedFolderId !== null || selectedFolderId === null" class="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
              <p class="text-sm text-blue-800 dark:text-blue-200">
                <span class="font-medium">Selected destination:</span> {{ selectedFolderName }}
              </p>
            </div>
          </div>
          <!-- Error message -->
          <div v-if="error" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <p class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>
          </div>
          <!-- Loading indicator -->
          <div v-if="loading" class="mb-4 flex items-center justify-center">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            <span class="ml-2 text-sm text-gray-600 dark:text-gray-400">Loading...</span>
          </div>
        </div>
      </div>
      <!-- Footer -->
      <div class="flex items-center justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          @click="handleCancel"
          class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
        >
          Cancel
        </button>
        <button
          @click="handleMove"
          :disabled="loading"
          class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <span v-if="loading">Moving...</span>
          <span v-else>Move</span>
        </button>
      </div>
    </div>
  </div>
</template> 