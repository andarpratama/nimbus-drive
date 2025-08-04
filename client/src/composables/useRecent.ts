import { ref, computed } from 'vue'
import { apiRequest, getAuthHeaders } from './useApi'
import { formatFileSize } from './useFileUtils'
import type { FileItem } from './types'

export interface RecentFile {
  id: string
  name: string
  size: number
  mime_type: string
  created_at: string
  updated_at: string
  folder_id?: string
}

export interface RecentResponse {
  files: RecentFile[]
  total: number
  limit: number
}

export function useRecent() {
  const recentFiles = ref<FileItem[]>([])
  const loading = ref(false)
  const error = ref('')

  const fetchRecentFiles = async (limit: number = 20) => {
    loading.value = true
    error.value = ''
    
    try {
      const response = await apiRequest<RecentResponse>(`/api/files/recent?limit=${limit}`)
      
      // Transform the response to match FileItem interface
      recentFiles.value = response.files.map(file => {
        // Determine file type from MIME type
        let fileType = 'file'
        if (file.mime_type.startsWith('image/')) {
          fileType = 'image' // Set to 'image' for FileCard thumbnail support
        } else if (file.mime_type.startsWith('video/')) {
          fileType = 'video'
        } else if (file.mime_type.startsWith('audio/')) {
          fileType = 'audio'
        } else if (file.mime_type.includes('pdf')) {
          fileType = 'pdf'
        } else if (file.mime_type.includes('text/')) {
          fileType = 'text'
        }
        
        return {
          id: file.id,
          fileId: file.id,
          folderId: file.folder_id,
          name: file.name,
          size: formatFileSize(file.size),
          type: fileType,
          modified: file.updated_at,
          starred: false,
          shared: false,
          rawSize: file.size,
          mimeType: file.mime_type
        }
      })
    } catch (err) {
      console.error('Error fetching recent files:', err)
      error.value = 'Failed to load recent files'
      recentFiles.value = []
    } finally {
      loading.value = false
    }
  }

  const refreshRecentFiles = async () => {
    await fetchRecentFiles()
  }

  return {
    recentFiles: computed(() => recentFiles.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    fetchRecentFiles,
    refreshRecentFiles
  }
} 