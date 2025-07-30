# Starred API Documentation

This document describes the API endpoints for the starred functionality in Nimbus Drive.

## Authentication

All starred endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### 1. Star a File
**POST** `/api/starred/files/:id/star`

Stars a file for the authenticated user.

**Parameters:**
- `id` (path): File UUID

**Response:**
```json
{
  "message": "file starred successfully"
}
```

**Error Responses:**
- `400`: Invalid file ID or user ID
- `404`: File not found or doesn't belong to user
- `409`: File already starred
- `500`: Internal server error

### 2. Unstar a File
**DELETE** `/api/starred/files/:id/star`

Removes star from a file.

**Parameters:**
- `id` (path): File UUID

**Response:**
```json
{
  "message": "file unstarred successfully"
}
```

**Error Responses:**
- `400`: Invalid file ID or user ID
- `404`: File not starred
- `500`: Internal server error

### 3. Star a Folder
**POST** `/api/starred/folders/:id/star`

Stars a folder for the authenticated user.

**Parameters:**
- `id` (path): Folder UUID

**Response:**
```json
{
  "message": "folder starred successfully"
}
```

**Error Responses:**
- `400`: Invalid folder ID or user ID
- `404`: Folder not found or doesn't belong to user
- `409`: Folder already starred
- `500`: Internal server error

### 4. Unstar a Folder
**DELETE** `/api/starred/folders/:id/star`

Removes star from a folder.

**Parameters:**
- `id` (path): Folder UUID

**Response:**
```json
{
  "message": "folder unstarred successfully"
}
```

**Error Responses:**
- `400`: Invalid folder ID or user ID
- `404`: Folder not starred
- `500`: Internal server error

### 5. List Starred Items
**GET** `/api/starred`

Returns all starred items (files and folders) for the authenticated user.

**Response:**
```json
{
  "starred_items": [
    {
      "id": "starred-item-uuid",
      "type": "file",
      "file_id": "file-uuid",
      "name": "document.pdf",
      "size": 1024,
      "mime_type": "application/pdf",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "starred-item-uuid",
      "type": "folder",
      "folder_id": "folder-uuid",
      "name": "My Documents",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400`: Invalid user ID
- `500`: Internal server error

### 6. Check Starred Status
**GET** `/api/starred/:id/status?type=file|folder`

Checks if a file or folder is starred by the user.

**Parameters:**
- `id` (path): Item UUID
- `type` (query): Either "file" or "folder"

**Response:**
```json
{
  "is_starred": true
}
```

**Error Responses:**
- `400`: Invalid item ID, user ID, or item type
- `500`: Internal server error

## Database Schema

The starred functionality uses a `starred` table with the following structure:

```sql
CREATE TABLE starred (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    user_id CHAR(36) NOT NULL,
    file_id CHAR(36) NULL,
    folder_id CHAR(36) NULL,
    
    -- Constraints ensure only one of file_id or folder_id is set
    CONSTRAINT chk_starred_item CHECK (
        (file_id IS NOT NULL AND folder_id IS NULL) OR 
        (file_id IS NULL AND folder_id IS NOT NULL)
    ),
    
    -- Foreign key constraints
    CONSTRAINT fk_starred_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_starred_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    CONSTRAINT fk_starred_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE,
    
    -- Unique constraints prevent duplicate stars
    CONSTRAINT uk_starred_user_file UNIQUE (user_id, file_id),
    CONSTRAINT uk_starred_user_folder UNIQUE (user_id, folder_id)
);
```

## Usage Examples

### Star a file
```bash
curl -X POST \
  http://localhost:8080/api/starred/files/123e4567-e89b-12d3-a456-426614174000/star \
  -H "Authorization: Bearer your-jwt-token"
```

### List all starred items
```bash
curl -X GET \
  http://localhost:8080/api/starred \
  -H "Authorization: Bearer your-jwt-token"
```

### Check if a file is starred
```bash
curl -X GET \
  "http://localhost:8080/api/starred/123e4567-e89b-12d3-a456-426614174000/status?type=file" \
  -H "Authorization: Bearer your-jwt-token"
``` 