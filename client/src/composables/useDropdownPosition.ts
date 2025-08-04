import { ref, nextTick } from 'vue'
export function useDropdownPosition() {
  const position = ref({ x: 0, y: 0 })
  const calculatePosition = async (
    triggerElement: HTMLElement | null,
    dropdownElement: HTMLElement | null,
    preferredPosition: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left' = 'bottom-right',
    offset: { x: number; y: number } = { x: 0, y: 8 }
  ) => {
    if (!triggerElement || !dropdownElement) return
    await nextTick()
    const triggerRect = triggerElement.getBoundingClientRect()
    const dropdownRect = dropdownElement.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight
    const dropdownWidth = dropdownRect.width || 200
    const dropdownHeight = dropdownRect.height || 200
    let newX = 0
    let newY = 0
    switch (preferredPosition) {
      case 'bottom-right':
        newX = triggerRect.right + offset.x
        newY = triggerRect.bottom + offset.y
        break
      case 'bottom-left':
        newX = triggerRect.left - dropdownWidth + offset.x
        newY = triggerRect.bottom + offset.y
        break
      case 'top-right':
        newX = triggerRect.right + offset.x
        newY = triggerRect.top - dropdownHeight + offset.y
        break
      case 'top-left':
        newX = triggerRect.left - dropdownWidth + offset.x
        newY = triggerRect.top - dropdownHeight + offset.y
        break
    }
    if (newX + dropdownWidth > viewportWidth) {
      if (triggerRect.left - dropdownWidth >= 0) {
        newX = triggerRect.left - dropdownWidth - offset.x
      } else {
        newX = viewportWidth - dropdownWidth - 10
      }
    }
    if (newX < 10) {
      newX = 10
    }
    if (newY + dropdownHeight > viewportHeight) {
      if (triggerRect.top - dropdownHeight >= 0) {
        newY = triggerRect.top - dropdownHeight - offset.y
      } else {
        newY = viewportHeight - dropdownHeight - 10
      }
    }
    if (newY < 10) {
      newY = 10
    }
    position.value = {
      x: Math.max(0, Math.min(newX, viewportWidth - 50)),
      y: Math.max(0, Math.min(newY, viewportHeight - 50))
    }
  }
  const resetPosition = () => {
    position.value = { x: 0, y: 0 }
  }
  return {
    position,
    calculatePosition,
    resetPosition
  }
} 