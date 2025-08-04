export const getFileType = (filename: string): string => {
  if (!filename || typeof filename !== 'string') {
    return 'document'
  }
  const ext = filename.split('.').pop()?.toLowerCase()
  if (!ext) return 'document'
  const typeMap: Record<string, string> = {
    'pdf': 'pdf',
    'doc': 'document',
    'docx': 'document',
    'txt': 'text',
    'rtf': 'document',
    'odt': 'document',
    'xls': 'spreadsheet',
    'xlsx': 'spreadsheet',
    'csv': 'spreadsheet',
    'ods': 'spreadsheet',
    'ppt': 'presentation',
    'pptx': 'presentation',
    'odp': 'presentation',
    'jpg': 'image',
    'jpeg': 'image',
    'png': 'image',
    'gif': 'image',
    'bmp': 'image',
    'webp': 'image',
    'svg': 'image',
    'ico': 'image',
    'mp4': 'video',
    'avi': 'video',
    'mov': 'video',
    'wmv': 'video',
    'flv': 'video',
    'webm': 'video',
    'mp3': 'audio',
    'wav': 'audio',
    'flac': 'audio',
    'aac': 'audio',
    'zip': 'archive',
    'rar': 'archive',
    '7z': 'archive',
    'tar': 'archive',
    'gz': 'archive',
    'js': 'code',
    'ts': 'code',
    'jsx': 'code',
    'tsx': 'code',
    'html': 'code',
    'css': 'code',
    'scss': 'code',
    'sass': 'code',
    'py': 'code',
    'java': 'code',
    'cpp': 'code',
    'c': 'code',
    'php': 'code',
    'rb': 'code',
    'go': 'code',
    'rs': 'code',
    'swift': 'code',
    'kt': 'code',
    'sql': 'code',
    'json': 'data',
    'xml': 'data',
    'yaml': 'data',
    'yml': 'data',
    'exe': 'executable',
    'msi': 'executable',
    'dmg': 'executable',
    'deb': 'executable',
    'rpm': 'executable'
  }
  return typeMap[ext] || 'document'
}
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
export const formatDate = (dateString: string): string => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffSeconds = Math.floor(diffTime / 1000)
  const diffMinutes = Math.floor(diffTime / (1000 * 60))
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60))
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  if (diffTime < 0) {
    return date.toLocaleDateString()
  }
  if (diffSeconds < 60) {
    return 'Just now'
  }
  if (diffMinutes < 60) {
    return diffMinutes === 1 ? '1 minute ago' : `${diffMinutes} minutes ago`
  }
  if (diffHours < 24) {
    return diffHours === 1 ? '1 hour ago' : `${diffHours} hours ago`
  }
  if (diffDays === 1) {
    return '1 day ago'
  }
  if (diffDays < 7) {
    return `${diffDays} days ago`
  }
  const diffWeeks = Math.floor(diffDays / 7)
  if (diffWeeks === 1) {
    return '1 week ago'
  }
  if (diffDays < 30) {
    return `${diffWeeks} weeks ago`
  }
  const diffMonths = Math.floor(diffDays / 30)
  if (diffMonths === 1) {
    return '1 month ago'
  }
  if (diffDays < 365) {
    return `${diffMonths} months ago`
  }
  const diffYears = Math.floor(diffDays / 365)
  if (diffYears === 1) {
    return '1 year ago'
  }
  return `${diffYears} years ago`
} 