<script setup>
import FileGrid from './FileGrid.vue'
import FileList from './FileList.vue'

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  selectedItems: {
    type: Array,
    required: true
  },
  viewMode: {
    type: String,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  searchQuery: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'item-select', 
  'item-double-click', 
  'item-star-toggle',
  'retry'
])

const handleItemSelect = (itemId) => {
  emit('item-select', itemId)
}

const handleItemDoubleClick = (item) => {
  emit('item-double-click', item)
}

const handleItemStarToggle = (itemId) => {
  emit('item-star-toggle', itemId)
}

const handleRetry = () => {
  emit('retry')
}
</script>

<template>
  <div class="flex-1 p-6 overflow-auto">
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>
    
    <!-- Error state -->
    <div v-else-if="error" class="text-center py-12">
      <div class="text-red-600 dark:text-red-400 mb-4">{{ error }}</div>
      <button 
        @click="handleRetry"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        Try Again
      </button>
    </div>
    
    <!-- Content -->
    <div v-else>
      <!-- Grid View -->
      <FileGrid
        v-if="viewMode === 'grid'"
        :items="items"
        :selected-items="selectedItems"
        @item-select="handleItemSelect"
        @item-double-click="handleItemDoubleClick"
        @item-star-toggle="handleItemStarToggle"
      />

      <!-- List View -->
      <FileList
        v-else
        :items="items"
        :selected-items="selectedItems"
        @item-select="handleItemSelect"
        @item-double-click="handleItemDoubleClick"
        @item-star-toggle="handleItemStarToggle"
      />

      <!-- Empty state -->
      <div v-if="items.length === 0 && !loading" class="text-center py-12">
        <div class="text-6xl mb-4">📁</div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
          {{ searchQuery ? 'No files found' : 'No files yet' }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400">
          {{ searchQuery ? 'Try adjusting your search terms' : 'Upload files or create new documents to get started' }}
        </p>
      </div>
    </div>
  </div>
</template> 