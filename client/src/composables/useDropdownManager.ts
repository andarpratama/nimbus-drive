import { ref } from 'vue'

// Global dropdown state management
const activeDropdown = ref<string | null>(null)

export const useDropdownManager = () => {
  const openDropdown = (dropdownId: string) => {
    // If clicking the same dropdown, close it
    if (activeDropdown.value === dropdownId) {
      activeDropdown.value = null
      return false
    }
    
    // Close any other open dropdown and open the new one
    activeDropdown.value = dropdownId
    return true
  }

  const closeDropdown = (dropdownId: string) => {
    if (activeDropdown.value === dropdownId) {
      activeDropdown.value = null
    }
  }

  const closeAllDropdowns = () => {
    activeDropdown.value = null
  }

  const isDropdownOpen = (dropdownId: string) => {
    return activeDropdown.value === dropdownId
  }

  const getActiveDropdown = () => {
    return activeDropdown.value
  }

  return {
    openDropdown,
    closeDropdown,
    closeAllDropdowns,
    isDropdownOpen,
    getActiveDropdown
  }
} 