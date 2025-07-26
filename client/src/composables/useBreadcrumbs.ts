import { ref } from 'vue'
import type { Folder, Breadcrumb, FolderResponse } from './types'
import { apiRequest } from './useApi'

export function useBreadcrumbs() {
  const breadcrumbs = ref<Breadcrumb[]>([])

  // Update breadcrumbs based on current folder
  const updateBreadcrumbs = async (currentFolder: Folder | null) => {
    breadcrumbs.value = []
    
    if (!currentFolder) {
      breadcrumbs.value = [{ id: null, name: 'My Drive' }]
      return
    }
    
    // Build breadcrumb path
    const path: Breadcrumb[] = []
    let current: Folder | null = currentFolder
    
    while (current) {
      path.unshift({ id: current.ID, name: current.Name })
      if (current.ParentID) {
        try {
          const data = await apiRequest<FolderResponse>(`/api/folders/${current.ParentID}`)
          current = data.folder
        } catch (err) {
          break
        }
      } else {
        break
      }
    }
    
    breadcrumbs.value = [{ id: null, name: 'My Drive' }, ...path]
  }

  return {
    breadcrumbs,
    updateBreadcrumbs
  }
} 