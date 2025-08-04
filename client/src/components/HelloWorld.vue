<script setup>
import { ref } from 'vue'
defineProps({
  msg: String,
})
const count = ref(0)
</script>
<template>
  <div class="text-center">
    <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-8">{{ msg }}</h1>
    <div class="card max-w-md mx-auto mb-8">
      <button 
        type="button" 
        @click="count++"
        class="btn-primary mb-4"
      >
        Count is {{ count }}
      </button>
      <p class="text-gray-600 dark:text-gray-300">
        Edit
        <code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-sm text-gray-800 dark:text-gray-200">components/HelloWorld.vue</code> to test HMR
      </p>
    </div>
    <div class="space-y-4 text-gray-600 dark:text-gray-300">
      <p>
        Check out
        <a 
          href="https:
          target="_blank"
          class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 underline"
        >
          create-vue
        </a>, the official Vue + Vite starter
      </p>
      <p>
        Learn more about IDE Support for Vue in the
        <a
          href="https:
          target="_blank"
          class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 underline"
        >
          Vue Docs Scaling up Guide
        </a>.
      </p>
      <p class="text-gray-500 dark:text-gray-400 text-sm">
        Click on the Vite and Vue logos to learn more
      </p>
    </div>
  </div>
</template>
<style scoped>
.read-the-docs {
  color: #888;
}
</style>
