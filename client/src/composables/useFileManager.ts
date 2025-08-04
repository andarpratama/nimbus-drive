import { computed } from 'vue'
import { useFileData } from './useFileData'
import { useBreadcrumbs } from './useBreadcrumbs'
import { useSelection } from './useSelection'
import { useStarred } from './useStarred'
import { getFileType, formatFileSize, formatDate } from './useFileUtils'
export function useFileManager() {
  const fileData = useFileData()
  const breadcrumbs = useBreadcrumbs()
  const selection = useSelection()
  const starred = useStarred()
  const navigateToFolder = async (folderId: string) => {
    await fileData.fetchFolderContents(folderId)
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }
  const navigateToBreadcrumb = async (folderId: string | null) => {
    await fileData.fetchFolderContents(folderId)
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }
  const navigateToRoot = async () => {
    await fileData.fetchFolderContents()
    await breadcrumbs.updateBreadcrumbs(fileData.currentFolder.value)
  }
  const toggleStar = async (item: any): Promise<boolean> => {
    return await starred.toggleStar(item)
  }
  return {
    files: fileData.files,
    folders: fileData.folders,
    currentFolderId: fileData.currentFolderId,
    currentFolder: fileData.currentFolder,
    breadcrumbs: breadcrumbs.breadcrumbs,
    loading: fileData.loading,
    error: fileData.error,
    selectedItems: selection.selectedItems,
    allItems: fileData.allItems,
    fetchFolderContents: fileData.fetchFolderContents,
    updateBreadcrumbs: breadcrumbs.updateBreadcrumbs,
    navigateToFolder,
    navigateToBreadcrumb,
    navigateToRoot,
    selectItem: selection.selectItem,
    isSelected: selection.isSelected,
    clearSelection: selection.clearSelection,
    toggleStar,
     getFileType,
     formatFileSize,
     formatDate
  }
} 