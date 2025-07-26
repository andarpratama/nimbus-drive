<script setup>
import { ref, onMounted, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'success' // success, error, warning, info
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  duration: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['close'])

const isVisible = ref(false)

// Auto-hide after duration
onMounted(() => {
  if (props.visible && props.duration > 0) {
    setTimeout(() => {
      emit('close')
    }, props.duration)
  }
})

// Watch for visibility changes
watch(() => props.visible, (visible) => {
  isVisible.value = visible
  
  // Auto-hide after duration when notification becomes visible
  if (visible && props.duration > 0) {
    setTimeout(() => {
      emit('close')
    }, props.duration)
  }
})

const handleClose = () => {
  emit('close')
}

const getIcon = () => {
  switch (props.type) {
    case 'success':
      return '✅'
    case 'error':
      return '❌'
    case 'warning':
      return '⚠️'
    case 'info':
      return 'ℹ️'
    default:
      return '✅'
  }
}

const getClasses = () => {
  const baseClasses = 'flex items-center p-4 rounded-lg shadow-lg border-l-4'
  
  switch (props.type) {
    case 'success':
      return `${baseClasses} bg-green-50 dark:bg-green-900/20 border-green-400 text-green-800 dark:text-green-200`
    case 'error':
      return `${baseClasses} bg-red-50 dark:bg-red-900/20 border-red-400 text-red-800 dark:text-red-200`
    case 'warning':
      return `${baseClasses} bg-yellow-50 dark:bg-yellow-900/20 border-yellow-400 text-yellow-800 dark:text-yellow-200`
    case 'info':
      return `${baseClasses} bg-blue-50 dark:bg-blue-900/20 border-blue-400 text-blue-800 dark:text-blue-200`
    default:
      return `${baseClasses} bg-green-50 dark:bg-green-900/20 border-green-400 text-green-800 dark:text-green-200`
  }
}
</script>

<template>
  <Transition
    enter-active-class="transition ease-out duration-300"
    enter-from-class="transform translate-y-2 opacity-0"
    enter-to-class="transform translate-y-0 opacity-100"
    leave-active-class="transition ease-in duration-200"
    leave-from-class="transform translate-y-0 opacity-100"
    leave-to-class="transform translate-y-2 opacity-0"
  >
    <div
      v-if="visible"
      class="fixed bottom-4 right-4 z-50 max-w-sm w-full"
    >
      <div :class="getClasses()">
        <div class="flex-shrink-0 mr-3">
          <span class="text-lg">{{ getIcon() }}</span>
        </div>
        
        <div class="flex-1">
          <h4 v-if="title" class="text-sm font-medium mb-1">
            {{ title }}
          </h4>
          <p v-if="message" class="text-sm">
            {{ message }}
          </p>
        </div>
        
        <button
          @click="handleClose"
          class="flex-shrink-0 ml-3 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          ✕
        </button>
      </div>
    </div>
  </Transition>
</template> 