import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

// Initialize theme on app startup
const initializeTheme = () => {
  const savedTheme = localStorage.getItem('theme')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  
  if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// Initialize theme before mounting the app
initializeTheme()

const app = createApp(App)
app.use(router)
app.mount('#app')
