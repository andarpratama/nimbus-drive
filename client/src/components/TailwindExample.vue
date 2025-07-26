<script setup>
import { ref } from 'vue'

const isOpen = ref(false)
const folders = ref([
  { id: 1, name: 'Documents', files: 5 },
  { id: 2, name: 'Pictures', files: 12 },
  { id: 3, name: 'Work', files: 8 }
])
</script>

<template>
  <div class="space-y-6">
    <!-- Header Section -->
    <div class="card">
      <h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-4">File Manager - Hot Reload Test</h2>
      <div class="flex space-x-4">
        <button class="btn-primary">
          <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          New Folder
        </button>
        <button class="btn-secondary">
          <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
          </svg>
          Upload File
        </button>
      </div>
    </div>

    <!-- Folders Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="folder in folders" 
        :key="folder.id"
        class="card hover:shadow-lg dark:hover:shadow-gray-900/30 transition-shadow duration-200 cursor-pointer"
      >
        <div class="flex items-center space-x-3">
          <div class="bg-blue-100 dark:bg-blue-900/30 p-3 rounded-lg">
            <svg class="w-6 h-6 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
            </svg>
          </div>
          <div class="flex-1">
            <h3 class="font-medium text-gray-900 dark:text-white">{{ folder.name }}</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ folder.files }} files</p>
          </div>
          <button class="text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Form Example -->
    <div class="card">
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Create New Folder</h3>
      <form class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Folder Name</label>
          <input 
            type="text" 
            class="input-field"
            placeholder="Enter folder name"
          >
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Parent Folder</label>
          <select class="input-field">
            <option value="">Root Directory</option>
            <option value="1">Documents</option>
            <option value="2">Pictures</option>
            <option value="3">Work</option>
          </select>
        </div>
        <div class="flex space-x-4">
          <button type="submit" class="btn-primary">Create Folder</button>
          <button type="button" class="btn-secondary">Cancel</button>
        </div>
      </form>
    </div>
  </div>
</template> 