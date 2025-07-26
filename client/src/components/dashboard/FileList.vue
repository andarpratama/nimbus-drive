<script setup>
import FileListItem from './FileListItem.vue'

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

const emit = defineEmits(['item-select', 'item-double-click', 'item-star-toggle'])

const handleItemSelect = (itemId) => {
  emit('item-select', itemId)
}

const handleItemDoubleClick = (item) => {
  emit('item-double-click', item)
}

const handleItemStarToggle = (itemId) => {
  emit('item-star-toggle', itemId)
}
</script>

<template>
  <div class="space-y-1">
    <FileListItem
      v-for="item in items"
      :key="item.id"
      :item="item"
      :is-selected="selectedItems.includes(item.id)"
      @select="handleItemSelect"
      @double-click="handleItemDoubleClick"
      @star-toggle="handleItemStarToggle"
    />
  </div>
</template> 