import { Ref } from 'vue'
export interface File {
  ID: string
  Name: string
  Size: number
  UpdatedAt: string
  updated_at?: string  
  UserID: string
  user_id?: string     
  FolderID?: string
  folder_id?: string   
  Path: string
  deleted_at?: string
}
export interface Folder {
  ID: string
  Name: string
  ParentID?: string
  parent_id?: string  
  UserID: string
  user_id?: string    
  UpdatedAt: string
  updated_at?: string  
  file_count?: number
  subfolder_count?: number
  total_items?: number
  deleted_at?: string
}
export interface Breadcrumb {
  id: string | null
  name: string
}
export interface FileItem {
  id: string
  name: string
  type: string
  size: string
  modified: string
  starred: boolean
  shared: boolean
  folderId?: string
  fileId?: string
  itemCount?: number
  fileCount?: number
  subfolderCount?: number
  rawSize?: number
  deletedAt?: string
}
export interface FoldersResponse {
  folders: Folder[] | null
}
export interface FilesResponse {
  files: File[] | null
}
export interface FolderResponse {
  folder: Folder
}
export interface FileManagerState {
  files: Ref<File[]>
  folders: Ref<Folder[]>
  currentFolderId: Ref<string | null>
  currentFolder: Ref<Folder | null>
  breadcrumbs: Ref<Breadcrumb[]>
  loading: Ref<boolean>
  error: Ref<string>
  selectedItems: Ref<string[]>
  allItems: Ref<FileItem[]>
} 