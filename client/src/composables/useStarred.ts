import { ref, computed } from 'vue'
import { apiRequest, getAuthHeaders } from './useApi'
import type { FileItem } from './types'

export interface StarredItem {
  id: string
  type: string
  file_id?: string
  folder_id?: string
  name: string
  size?: string
  mime_type?: string
  created_at: string
}

export interface StarredResponse {
  starred_items: StarredItem[]
}

export interface StarToggleResponse {
  message: string
  starred: boolean
}

export function useStarred() {
  const starredItems = ref<StarredItem[]>([])
  const loading = ref(false)
  const error = ref('')

  const fetchStarredItems = async () => {
    loading.value = true
    error.value = ''
    
    try {
      const response = await apiRequest<StarredResponse>('/api/starred')
      starredItems.value = response.starred_items || []
    } catch (err) {
      console.error('Error fetching starred items:', err)
      error.value = 'Failed to load starred items'
    } finally {
      loading.value = false
    }
  }

  const starItem = async (itemId: string, itemType: 'file' | 'folder'): Promise<boolean> => {
    try {
      const response = await apiRequest<StarToggleResponse>('/api/starred', {
        method: 'POST',
        body: JSON.stringify({
          item_id: itemId,
          item_type: itemType
        })
      })
      return response.starred
    } catch (err) {
      console.error('Error starring item:', err)
      return false
    }
  }

  const unstarItem = async (itemId: string, itemType: 'file' | 'folder'): Promise<boolean> => {
    try {
      const response = await apiRequest<StarToggleResponse>('/api/starred', {
        method: 'DELETE',
        body: JSON.stringify({
          item_id: itemId,
          item_type: itemType
        })
      })
      return !response.starred
    } catch (err) {
      console.error('Error unstarring item:', err)
      return false
    }
  }

  const toggleStar = async (item: FileItem): Promise<boolean> => {
    try {
      let itemId: string | undefined
      let itemType: 'file' | 'folder'
      
      if (item.folderId) {
        itemId = item.folderId
        itemType = 'folder'
      } else if (item.fileId) {
        itemId = item.fileId
        itemType = 'file'
      } else {
        console.error('No item ID found for toggle star')
        return false
      }

      const response = await apiRequest<StarToggleResponse>('/api/starred', {
        method: 'PATCH',
        body: JSON.stringify({
          item_id: itemId,
          item_type: itemType
        })
      })
      return true
    } catch (err) {
      console.error('Error toggling star:', err)
      return false
    }
  }

  const starFile = async (fileId: string): Promise<boolean> => {
    return await starItem(fileId, 'file')
  }

  const unstarFile = async (fileId: string): Promise<boolean> => {
    return await unstarItem(fileId, 'file')
  }

  const starFolder = async (folderId: string): Promise<boolean> => {
    return await starItem(folderId, 'folder')
  }

  const unstarFolder = async (folderId: string): Promise<boolean> => {
    return await unstarItem(folderId, 'folder')
  }

  const checkStarredStatus = async (itemId: string, type: 'file' | 'folder'): Promise<boolean> => {
    return starredItems.value.some(item => {
      if (type === 'file') {
        return item.file_id === itemId
      } else {
        return item.folder_id === itemId
      }
    })
  }

  const starredFileItems = computed<FileItem[]>(() => {
    return starredItems.value.map((item: StarredItem) => {
      let displaySize = 'Unknown'
      let rawSize = 0
      
      if (item.type === 'folder' && item.size) {
        displaySize = item.size
        rawSize = 0
      } else if (item.size) {
        displaySize = item.size
        rawSize = parseInt(item.size) || 0
      }
      
      return {
        id: item.type === 'folder' ? `folder-${item.folder_id}` : `file-${item.file_id}`,
        name: item.name,
        type: item.type,
        size: displaySize,
        modified: new Date(item.created_at).toLocaleDateString(),
        starred: true,
        shared: false,
        folderId: item.folder_id,
        fileId: item.file_id,
        rawSize: rawSize
      }
    })
  })

  return {
    starredItems,
    starredFileItems,
    loading,
    error,
    
    fetchStarredItems,
    starItem,
    unstarItem,
    toggleStar,
    starFile,
    unstarFile,
    starFolder,
    unstarFolder,
    checkStarredStatus
  }
} 