import { computed } from 'vue'
import { useFileData } from './useFileData'
import { useBreadcrumbs } from './useBreadcrumbs'
import { useSelection } from './useSelection'
import { getFileType, formatFileSize, formatDate } from './useFileUtils'

export function useFileManager() {
  // Use modular composables
  const fileData = useFileData()
  const breadcrumbs = useBreadcrumbs()
  const selection = useSelection()

  // Navigation methods
  const navigateToFolder = async (folderId: number) => {
    await fileData.fetchFolderContents(folderId)
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }

  const navigateToBreadcrumb = async (folderId: number | null) => {
    await fileData.fetchFolderContents(folderId)
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }

  const navigateToRoot = async () => {
    await fileData.fetchFolderContents()
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }

  // Star management
  const toggleStar = (itemId: string): void => {
    const item = fileData.allItems.value.find(i => i.id === itemId)
    if (item) item.starred = !item.starred
  }

  // Re-export everything for backward compatibility
  return {
    // State
    files: fileData.files,
    folders: fileData.folders,
    currentFolderId: fileData.currentFolderId,
    currentFolder: fileData.currentFolder,
    breadcrumbs: breadcrumbs.breadcrumbs,
    loading: fileData.loading,
    error: fileData.error,
    selectedItems: selection.selectedItems,
    allItems: fileData.allItems,
    
    // Methods
    fetchFolderContents: fileData.fetchFolderContents,
    updateBreadcrumbs: breadcrumbs.updateBreadcrumbs,
    navigateToFolder,
    navigateToBreadcrumb,
    navigateToRoot,
    selectItem: selection.selectItem,
    isSelected: selection.isSelected,
    clearSelection: selection.clearSelection,
    toggleStar,
    
         // Helper functions (re-exported from useFileUtils)
     getFileType,
     formatFileSize,
     formatDate
  }
} 