<script setup>
import { computed } from 'vue'

const props = defineProps({
  filename: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'document'
  },
  size: {
    type: String,
    default: 'md' // sm, md, lg, xl
  }
})

const getFileExtension = (filename) => {
  if (!filename) return ''
  return filename.split('.').pop()?.toLowerCase() || ''
}

const iconClass = computed(() => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12'
  }
  return sizeClasses[props.size] || sizeClasses.md
})

const getIcon = computed(() => {
  const ext = getFileExtension(props.filename)
  
  // Microsoft Office icons
  if (ext === 'xlsx' || ext === 'xls') {
    return 'excel'
  }
  if (ext === 'pptx' || ext === 'ppt') {
    return 'powerpoint'
  }
  if (ext === 'docx' || ext === 'doc') {
    return 'word'
  }
  if (ext === 'pdf') {
    return 'pdf'
  }
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg'].includes(ext)) {
    return 'image'
  }
  if (['mp4', 'avi', 'mov', 'wmv', 'flv', 'webm'].includes(ext)) {
    return 'video'
  }
  if (['mp3', 'wav', 'flac', 'aac'].includes(ext)) {
    return 'audio'
  }
  if (['zip', 'rar', '7z', 'tar', 'gz'].includes(ext)) {
    return 'archive'
  }
  if (['js', 'ts', 'jsx', 'tsx', 'html', 'css', 'scss', 'sass', 'py', 'java', 'cpp', 'c', 'php', 'rb', 'go', 'rs', 'swift', 'kt', 'sql'].includes(ext)) {
    return 'code'
  }
  
  return 'document'
})
</script>

<template>
  <div :class="iconClass">
    <!-- Microsoft Excel -->
    <svg v-if="getIcon === 'excel'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#217346"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#217346"/>
      <path d="M10 10h2v2h-2v-2zm2 2h2v2h-2v-2z" fill="white"/>
    </svg>
    
    <!-- Microsoft PowerPoint -->
    <svg v-else-if="getIcon === 'powerpoint'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#D24726"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#D24726"/>
      <path d="M10 10h2v2h-2v-2zm2 2h2v2h-2v-2z" fill="white"/>
    </svg>
    
    <!-- Microsoft Word -->
    <svg v-else-if="getIcon === 'word'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#2B579A"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#2B579A"/>
      <path d="M10 10h2v2h-2v-2zm2 2h2v2h-2v-2z" fill="white"/>
    </svg>
    
    <!-- PDF -->
    <svg v-else-if="getIcon === 'pdf'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#DC3545"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#DC3545"/>
      <text x="12" y="16" text-anchor="middle" fill="white" font-size="8" font-weight="bold">PDF</text>
    </svg>
    
    <!-- Image -->
    <svg v-else-if="getIcon === 'image'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#28A745"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#28A745"/>
      <circle cx="11" cy="11" r="1" fill="white"/>
      <path d="M9 15l2-2 2 2 2-2" stroke="white" stroke-width="1" fill="none"/>
    </svg>
    
    <!-- Video -->
    <svg v-else-if="getIcon === 'video'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#6F42C1"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#6F42C1"/>
      <polygon points="12,10 16,12 12,14" fill="white"/>
    </svg>
    
    <!-- Audio -->
    <svg v-else-if="getIcon === 'audio'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#FD7E14"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#FD7E14"/>
      <path d="M11 11h2v2h-2v-2zm2 2h2v2h-2v-2z" fill="white"/>
      <circle cx="12" cy="12" r="1" fill="#FD7E14"/>
    </svg>
    
    <!-- Archive -->
    <svg v-else-if="getIcon === 'archive'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#6C757D"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#6C757D"/>
      <path d="M10 10h4v1h-4v-1zm0 2h4v1h-4v-1z" fill="white"/>
    </svg>
    
    <!-- Code -->
    <svg v-else-if="getIcon === 'code'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#17A2B8"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#17A2B8"/>
      <path d="M10 10l1 1-1 1m2-2l1 1-1 1" stroke="white" stroke-width="1" fill="none"/>
    </svg>
    
    <!-- Default Document -->
    <svg v-else viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="2" fill="#6C757D"/>
      <path d="M7 7h10v10H7V7z" fill="white"/>
      <path d="M9 9v6h6V9H9zm1 1h4v4h-4v-4z" fill="#6C757D"/>
      <path d="M10 10h4v1h-4v-1zm0 2h4v1h-4v-1zm0 2h2v1h-2v-1z" fill="white"/>
    </svg>
  </div>
</template> 