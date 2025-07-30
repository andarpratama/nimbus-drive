import { ref, computed } from 'vue'
import { apiRequest, getAuthHeaders } from './useApi'
import type { FileItem } from './types'

export interface StarredItem {
  id: string
  type: 'file' | 'folder'
  file_id?: string
  folder_id?: string
  name: string
  size?: string // For files: file size, for folders: item count
  mime_type?: string
  created_at: string
}

export interface StarredResponse {
  starred_items: StarredItem[]
}

export function useStarred() {
  const starredItems = ref<StarredItem[]>([])
  const loading = ref(false)
  const error = ref('')

  // Fetch all starred items
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

  // Star a file
  const starFile = async (fileId: string): Promise<boolean> => {
    try {
      await apiRequest(`/api/starred/files/${fileId}/star`, {
        method: 'POST'
      })
      return true
    } catch (err) {
      console.error('Error starring file:', err)
      return false
    }
  }

  // Unstar a file
  const unstarFile = async (fileId: string): Promise<boolean> => {
    try {
      await apiRequest(`/api/starred/files/${fileId}/star`, {
        method: 'DELETE'
      })
      return true
    } catch (err) {
      console.error('Error unstarring file:', err)
      return false
    }
  }

  // Star a folder
  const starFolder = async (folderId: string): Promise<boolean> => {
    try {
      await apiRequest(`/api/starred/folders/${folderId}/star`, {
        method: 'POST'
      })
      return true
    } catch (err) {
      console.error('Error starring folder:', err)
      return false
    }
  }

  // Unstar a folder
  const unstarFolder = async (folderId: string): Promise<boolean> => {
    try {
      await apiRequest(`/api/starred/folders/${folderId}/star`, {
        method: 'DELETE'
      })
      return true
    } catch (err) {
      console.error('Error unstarring folder:', err)
      return false
    }
  }

  // Check if an item is starred
  const checkStarredStatus = async (itemId: string, type: 'file' | 'folder'): Promise<boolean> => {
    try {
      const response = await apiRequest<{ is_starred: boolean }>(`/api/starred/${itemId}/status?type=${type}`)
      return response.is_starred
    } catch (err) {
      console.error('Error checking starred status:', err)
      return false
    }
  }

  // Toggle star for an item
  const toggleStar = async (item: FileItem): Promise<boolean> => {
    if (item.type === 'folder' && item.folderId) {
      const isStarred = await checkStarredStatus(item.folderId, 'folder')
      if (isStarred) {
        return await unstarFolder(item.folderId)
      } else {
        return await starFolder(item.folderId)
      }
    } else if (item.fileId) {
      const isStarred = await checkStarredStatus(item.fileId, 'file')
      if (isStarred) {
        return await unstarFile(item.fileId)
      } else {
        return await starFile(item.fileId)
      }
    }
    return false
  }

  // Convert starred items to FileItem format for display
  const starredFileItems = computed<FileItem[]>(() => {
    return starredItems.value.map((item: StarredItem) => {
      // Handle size for files vs folders
      let displaySize = 'Unknown'
      let rawSize = 0
      
      if (item.type === 'file' && item.size) {
        // For files, use the actual file size
        displaySize = item.size
        rawSize = parseInt(item.size) || 0
      } else if (item.type === 'folder' && item.size) {
        // For folders, use the item count from backend
        displaySize = item.size
        rawSize = 0 // Folders don't have raw size
      }
      
      return {
        id: item.type === 'file' ? `file-${item.file_id}` : `folder-${item.folder_id}`,
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
    // State
    starredItems,
    starredFileItems,
    loading,
    error,
    
    // Methods
    fetchStarredItems,
    starFile,
    unstarFile,
    starFolder,
    unstarFolder,
    checkStarredStatus,
    toggleStar
  }
} 