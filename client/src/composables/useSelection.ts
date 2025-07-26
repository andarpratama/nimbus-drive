import { ref } from 'vue'

export function useSelection() {
  const selectedItems = ref<string[]>([])

  // Selection management
  const selectItem = (itemId: string): void => {
    const index = selectedItems.value.indexOf(itemId)
    if (index > -1) {
      selectedItems.value.splice(index, 1)
    } else {
      selectedItems.value.push(itemId)
    }
  }

  const isSelected = (itemId: string): boolean => selectedItems.value.includes(itemId)

  const clearSelection = (): void => {
    selectedItems.value = []
  }

  return {
    selectedItems,
    selectItem,
    isSelected,
    clearSelection
  }
} 