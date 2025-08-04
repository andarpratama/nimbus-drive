<script setup>
const props = defineProps({
  breadcrumbs: {
    type: Array,
    required: true
  }
})
const emit = defineEmits(['navigate-breadcrumb'])
const handleBreadcrumbClick = (folderId) => {
  emit('navigate-breadcrumb', folderId)
}
</script>
<template>
  <div class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 px-4 py-3 shadow-sm">
    <div class="flex items-center space-x-2 text-sm">
      <template v-for="(crumb, index) in breadcrumbs" :key="crumb.id">
        <button
          @click="handleBreadcrumbClick(crumb.id)"
          :class="[
            'hover:text-blue-600 dark:hover:text-blue-400 transition-colors px-2 py-1 rounded',
            index === breadcrumbs.length - 1 
              ? 'text-gray-900 dark:text-white font-medium bg-gray-100 dark:bg-gray-700' 
              : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700'
          ]"
        >
          {{ crumb.name }}
        </button>
        <span v-if="index < breadcrumbs.length - 1" class="text-gray-400">/</span>
      </template>
    </div>
  </div>
</template> 