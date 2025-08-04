import { ref, computed } from 'vue'
import type { File, Folder, FileItem, FoldersResponse, FilesResponse } from './types'
import { apiRequest } from './useApi'
import { getFileType, formatFileSize, formatDate } from './useFileUtils'
import { useStarred } from './useStarred'
export function useFileData() {
  const files = ref<File[]>([])
  const folders = ref<Folder[]>([])
  const currentFolderId = ref<string | null>(null)
  const currentFolder = ref<Folder | null>(null)
  const loading = ref(false)
  const error = ref('')
  const starred = useStarred()
  const fetchFolderContents = async (folderId: string | null = null, fetchStarred: boolean = true, skipFileFolderFetch: boolean = false) => {
    loading.value = true
    error.value = ''
    try {
      console.log('Fetching folder contents for folderId:', folderId, 'fetchStarred:', fetchStarred, 'skipFileFolderFetch:', skipFileFolderFetch)
      currentFolderId.value = folderId
      
      // Only fetch starred items if requested
      const starredPromise = fetchStarred ? starred.fetchStarredItems() : Promise.resolve()
      
      // Skip file/folder fetching if requested (for profile view)
      if (!skipFileFolderFetch) {
        console.log('Making API calls for folders and files')
        const foldersData = await apiRequest<FoldersResponse>('/api/folders')
        const allFolders = foldersData.folders || []
        folders.value = allFolders.filter((folder: Folder) => {
          const parentId = folder.ParentID || folder.parent_id
          if (folderId === null) {
            return parentId === null || parentId === undefined
          }
          return parentId === folderId
        })
        if (folderId) {
          currentFolder.value = allFolders.find((f: Folder) => f.ID === folderId) || null
        } else {
          currentFolder.value = null
        }
        const filesEndpoint = folderId ? `/api/files?folder_id=${folderId}` : '/api/files'
        const filesData = await apiRequest<FilesResponse>(filesEndpoint)
        files.value = filesData.files || []
      } else {
        // For profile view, just set empty arrays - NO API calls
        console.log('Skipping API calls for profile view - setting empty arrays')
        folders.value = []
        files.value = []
        currentFolder.value = null
      }
      
      await starredPromise
    } catch (err) {
      console.error('Error fetching data:', err)
      error.value = 'Failed to load files and folders'
    } finally {
      loading.value = false
    }
  }
  const allItems = computed<FileItem[]>(() => {
    const items: FileItem[] = []
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
    files,
    folders,
    currentFolderId,
    currentFolder,
    loading,
    error,
    allItems,
    fetchFolderContents
  }
} 