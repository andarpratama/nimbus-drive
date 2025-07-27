-- Fixed UUID Migration Script
-- This script properly converts existing tables from integer IDs to UUIDs

-- Step 1: Drop existing foreign key constraints first
ALTER TABLE folders DROP FOREIGN KEY fk_folders_parent;
ALTER TABLE files DROP FOREIGN KEY fk_files_folder;
ALTER TABLE files DROP FOREIGN KEY fk_files_user;
ALTER TABLE folders DROP FOREIGN KEY fk_folders_user;
ALTER TABLE sessions DROP FOREIGN KEY fk_sessions_user;
ALTER TABLE shared_files DROP FOREIGN KEY fk_shared_files_file;

-- Step 2: Add UUID columns to existing tables
ALTER TABLE users ADD COLUMN uuid_id CHAR(36);
ALTER TABLE folders ADD COLUMN uuid_id CHAR(36);
ALTER TABLE files ADD COLUMN uuid_id CHAR(36);
ALTER TABLE sessions ADD COLUMN uuid_id CHAR(36);
ALTER TABLE shared_files ADD COLUMN uuid_id CHAR(36);

-- Step 3: Generate UUIDs for existing records
UPDATE users SET uuid_id = UUID();
UPDATE folders SET uuid_id = UUID();
UPDATE files SET uuid_id = UUID();
UPDATE sessions SET uuid_id = UUID();
UPDATE shared_files SET uuid_id = UUID();

-- Step 4: Add UUID columns for foreign keys
ALTER TABLE folders ADD COLUMN uuid_parent_id CHAR(36);
ALTER TABLE folders ADD COLUMN uuid_user_id CHAR(36);
ALTER TABLE files ADD COLUMN uuid_folder_id CHAR(36);
ALTER TABLE files ADD COLUMN uuid_user_id CHAR(36);
ALTER TABLE sessions ADD COLUMN uuid_user_id CHAR(36);
ALTER TABLE shared_files ADD COLUMN uuid_file_id CHAR(36);

-- Step 5: Update foreign key references to use UUIDs
-- Update folder references in files table
UPDATE files f 
JOIN folders fo ON f.folder_id = fo.id 
SET f.uuid_folder_id = fo.uuid_id;

-- Update parent folder references
UPDATE folders f1 
JOIN folders f2 ON f1.parent_id = f2.id 
SET f1.uuid_parent_id = f2.uuid_id;

-- Update user references
UPDATE folders SET uuid_user_id = (SELECT uuid_id FROM users WHERE id = folders.user_id);
UPDATE files SET uuid_user_id = (SELECT uuid_id FROM users WHERE id = files.user_id);
UPDATE sessions SET uuid_user_id = (SELECT uuid_id FROM users WHERE id = sessions.user_id);
UPDATE shared_files SET uuid_file_id = (SELECT uuid_id FROM files WHERE id = shared_files.file_id);

-- Step 6: Drop old columns and rename UUID columns
-- Users table
ALTER TABLE users DROP PRIMARY KEY;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users CHANGE uuid_id id CHAR(36) PRIMARY KEY;

-- Folders table
ALTER TABLE folders DROP PRIMARY KEY;
ALTER TABLE folders DROP COLUMN id;
ALTER TABLE folders CHANGE uuid_id id CHAR(36) PRIMARY KEY;
ALTER TABLE folders DROP COLUMN user_id;
ALTER TABLE folders CHANGE uuid_user_id user_id CHAR(36) NOT NULL;
ALTER TABLE folders DROP COLUMN parent_id;
ALTER TABLE folders CHANGE uuid_parent_id parent_id CHAR(36);

-- Files table
ALTER TABLE files DROP PRIMARY KEY;
ALTER TABLE files DROP COLUMN id;
ALTER TABLE files CHANGE uuid_id id CHAR(36) PRIMARY KEY;
ALTER TABLE files DROP COLUMN user_id;
ALTER TABLE files CHANGE uuid_user_id user_id CHAR(36) NOT NULL;
ALTER TABLE files DROP COLUMN folder_id;
ALTER TABLE files CHANGE uuid_folder_id folder_id CHAR(36);

-- Sessions table
ALTER TABLE sessions DROP PRIMARY KEY;
ALTER TABLE sessions DROP COLUMN id;
ALTER TABLE sessions CHANGE uuid_id id CHAR(36) PRIMARY KEY;
ALTER TABLE sessions DROP COLUMN user_id;
ALTER TABLE sessions CHANGE uuid_user_id user_id CHAR(36) NOT NULL;

-- Shared_files table
ALTER TABLE shared_files DROP PRIMARY KEY;
ALTER TABLE shared_files DROP COLUMN id;
ALTER TABLE shared_files CHANGE uuid_id id CHAR(36) PRIMARY KEY;
ALTER TABLE shared_files DROP COLUMN file_id;
ALTER TABLE shared_files CHANGE uuid_file_id file_id CHAR(36) NOT NULL;

-- Step 7: Add foreign key constraints back
ALTER TABLE folders ADD CONSTRAINT fk_folders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE folders ADD CONSTRAINT fk_folders_parent FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE SET NULL;
ALTER TABLE files ADD CONSTRAINT fk_files_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE files ADD CONSTRAINT fk_files_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL;
ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE shared_files ADD CONSTRAINT fk_shared_files_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE; 