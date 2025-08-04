<script setup>
import { ref, watch } from 'vue'
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    default: 'Are you sure you want to proceed?'
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  cancelText: {
    type: String,
    default: 'Cancel'
  },
  danger: {
    type: Boolean,
    default: false
  }
})
const emit = defineEmits(['confirm', 'cancel', 'close'])
const handleConfirm = () => {
  emit('confirm')
}
const handleCancel = () => {
  emit('cancel')
}
const handleClose = () => {
  emit('close')
}
const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    handleClose()
  }
}
watch(() => props.visible, (visible) => {
  if (visible) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})
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
          {{ title }}
        </h3>
        <button
          @click="handleClose"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          ✕
        </button>
      </div>
      <!-- Content -->
      <div class="p-6">
        <p class="text-gray-600 dark:text-gray-300">
          {{ message }}
        </p>
      </div>
      <!-- Actions -->
      <div class="flex justify-end gap-3 p-6 border-t border-gray-200 dark:border-gray-700">
        <button
          @click="handleCancel"
          class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
        >
          {{ cancelText }}
        </button>
        <button
          @click="handleConfirm"
          :class="[
            'px-4 py-2 text-sm font-medium rounded-lg transition-colors',
            danger
              ? 'bg-red-600 text-white hover:bg-red-700'
              : 'bg-blue-600 text-white hover:bg-blue-700'
          ]"
        >
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template> 