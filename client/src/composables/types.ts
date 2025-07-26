import { Ref } from 'vue'

// Core data types
export interface File {
  ID: number
  Name: string
  Size: number
  UpdatedAt: string
  updated_at?: string  // GORM might return this as snake_case
  UserID: number
  user_id?: number     // GORM might return this as snake_case
  FolderID?: number
  folder_id?: number   // GORM might return this as snake_case
  Path: string
  deleted_at?: string
}

export interface Folder {
  ID: number
  Name: string
  ParentID?: number
  parent_id?: number  // GORM might return this as snake_case
  UserID: number
  user_id?: number    // GORM might return this as snake_case
  UpdatedAt: string
  updated_at?: string  // GORM might return this as snake_case
  file_count?: number
  subfolder_count?: number
  total_items?: number
  deleted_at?: string
}

export interface Breadcrumb {
  id: number | null
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
  folderId?: number
  fileId?: number
  itemCount?: number
  fileCount?: number
  subfolderCount?: number
  rawSize?: number
  deletedAt?: string
}

// API response types
export interface FoldersResponse {
  folders: Folder[]
}

export interface FilesResponse {
  files: File[]
}

export interface FolderResponse {
  folder: Folder
}

// State types
export interface FileManagerState {
  files: Ref<File[]>
  folders: Ref<Folder[]>
  currentFolderId: Ref<number | null>
  currentFolder: Ref<Folder | null>
  breadcrumbs: Ref<Breadcrumb[]>
  loading: Ref<boolean>
  error: Ref<string>
  selectedItems: Ref<string[]>
  allItems: Ref<FileItem[]>
} 