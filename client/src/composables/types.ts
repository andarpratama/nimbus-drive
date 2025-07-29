import { Ref } from 'vue'

// Core data types
export interface File {
  ID: string
  Name: string
  Size: number
  UpdatedAt: string
  updated_at?: string  // GORM might return this as snake_case
  UserID: string
  user_id?: string     // GORM might return this as snake_case
  FolderID?: string
  folder_id?: string   // GORM might return this as snake_case
  Path: string
  deleted_at?: string
}

export interface Folder {
  ID: string
  Name: string
  ParentID?: string
  parent_id?: string  // GORM might return this as snake_case
  UserID: string
  user_id?: string    // GORM might return this as snake_case
  UpdatedAt: string
  updated_at?: string  // GORM might return this as snake_case
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

// API response types - allow null values to match actual API response
export interface FoldersResponse {
  folders: Folder[] | null
}

export interface FilesResponse {
  files: File[] | null
}

export interface FolderResponse {
  folder: Folder
}

// State types
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