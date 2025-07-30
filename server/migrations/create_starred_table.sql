-- Create starred table migration
-- This table stores starred files and folders for users

CREATE TABLE IF NOT EXISTS starred (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    user_id CHAR(36) NOT NULL,
    file_id CHAR(36) NULL,
    folder_id CHAR(36) NULL,
    
    -- Ensure at least one of file_id or folder_id is set, but not both
    CONSTRAINT chk_starred_item CHECK (
        (file_id IS NOT NULL AND folder_id IS NULL) OR 
        (file_id IS NULL AND folder_id IS NOT NULL)
    ),
    
    -- Foreign key constraints
    CONSTRAINT fk_starred_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_starred_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    CONSTRAINT fk_starred_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE,
    
    -- Unique constraint to prevent duplicate stars
    CONSTRAINT uk_starred_user_file UNIQUE (user_id, file_id),
    CONSTRAINT uk_starred_user_folder UNIQUE (user_id, folder_id),
    
    -- Indexes for better performance
    INDEX idx_starred_user_id (user_id),
    INDEX idx_starred_file_id (file_id),
    INDEX idx_starred_folder_id (folder_id),
    INDEX idx_starred_deleted_at (deleted_at)
); 