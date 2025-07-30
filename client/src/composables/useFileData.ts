import { ref, computed } from 'vue'
import type { File, Folder, FileItem, FoldersResponse, FilesResponse } from './types'
import { apiRequest } from './useApi'
import { getFileType, formatFileSize, formatDate } from './useFileUtils'
import { useStarred } from './useStarred'

export function useFileData() {
  // State
  const files = ref<File[]>([])
  const folders = ref<Folder[]>([])
  const currentFolderId = ref<string | null>(null)
  const currentFolder = ref<Folder | null>(null)
  const loading = ref(false)
  const error = ref('')
  
  // Use starred composable
  const starred = useStarred()

  // Fetch folder contents
  const fetchFolderContents = async (folderId: string | null = null) => {
    loading.value = true
    error.value = ''
    
    try {
      console.log('Fetching folder contents for folderId:', folderId)
      
      // Update current folder ID
      currentFolderId.value = folderId
      
      // Fetch starred items in parallel
      const starredPromise = starred.fetchStarredItems()
      
      // Fetch folders
      const foldersData = await apiRequest<FoldersResponse>('/api/folders')
      // Handle null response by defaulting to empty array
      const allFolders = foldersData.folders || []
      folders.value = allFolders.filter((folder: Folder) => {
        // Handle both PascalCase and snake_case field names
        const parentId = folder.ParentID || folder.parent_id
        if (folderId === null) {
          return parentId === null || parentId === undefined
        }
        return parentId === folderId
      })
      
      // Set current folder if we're in a specific folder
      if (folderId) {
        currentFolder.value = allFolders.find((f: Folder) => f.ID === folderId) || null
      } else {
        currentFolder.value = null
      }
      
      // Fetch files with folder filtering
      const filesEndpoint = folderId ? `/api/files?folder_id=${folderId}` : '/api/files'
      const filesData = await apiRequest<FilesResponse>(filesEndpoint)
      
      // Handle null response by defaulting to empty array
      files.value = filesData.files || []
      
      // Wait for starred items to be fetched
      await starredPromise
      
    } catch (err) {
      console.error('Error fetching data:', err)
      error.value = 'Failed to load files and folders'
    } finally {
      loading.value = false
    }
  }

  // Combined files and folders for display
  const allItems = computed<FileItem[]>(() => {
    const items: FileItem[] = []
    
    // Add folders first
    folders.value.forEach((folder: Folder) => {
      if (!folder || !folder.Name) return
      
      const totalItems = folder.total_items || 0
      const fileCount = folder.file_count || 0
      const subfolderCount = folder.subfolder_count || 0
      
      items.push({
        id: `folder-${folder.ID}`,
        name: folder.Name,
        type: 'folder',
        size: totalItems > 0 ? `${totalItems} item${totalItems !== 1 ? 's' : ''}` : 'Empty',
        modified: formatDate(folder.UpdatedAt),
        starred: starred.starredItems.value.some(s => s.folder_id === folder.ID),
        shared: false,
        folderId: folder.ID,
        itemCount: totalItems,
        fileCount: fileCount,
        subfolderCount: subfolderCount
      })
    })
    
    // Add files
    files.value.forEach((file: File) => {
      if (!file || !file.Name) return
      
      items.push({
        id: `file-${file.ID}`,
        name: file.Name,
        type: getFileType(file.Name),
        size: formatFileSize(file.Size),
        modified: formatDate(file.UpdatedAt),
        starred: starred.starredItems.value.some(s => s.file_id === file.ID),
        shared: false,
        fileId: file.ID,
        rawSize: file.Size,
        deletedAt: file.deleted_at
      })
    })
    
    return items
  })

  return {
    // State
    files,
    folders,
    currentFolderId,
    currentFolder,
    loading,
    error,
    allItems,
    
    // Methods
    fetchFolderContents
  }
} 