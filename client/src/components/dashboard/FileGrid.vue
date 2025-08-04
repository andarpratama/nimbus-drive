<script setup>
import { computed } from 'vue'
import FileCard from './FileCard.vue'
const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  selectedItems: {
    type: Array,
    required: true
  }
})
const emit = defineEmits(['item-select', 'item-double-click', 'item-star-toggle', 'context-menu'])
const handleItemSelect = (itemId) => {
  emit('item-select', itemId)
}
const handleItemDoubleClick = (item) => {
  emit('item-double-click', item)
}
const handleItemStarToggle = (itemId) => {
  emit('item-star-toggle', itemId)
}
const handleContextMenu = (data) => {
  emit('context-menu', data)
}
</script>
<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
    <FileCard
      v-for="item in items"
      :key="item.id"
      :item="item"
      :is-selected="selectedItems.includes(item.id)"
      @select="handleItemSelect"
      @double-click="handleItemDoubleClick"
      @star-toggle="handleItemStarToggle"
      @context-menu="handleContextMenu"
    />
  </div>
</template> 