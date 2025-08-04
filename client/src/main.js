import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
const initializeTheme = () => {
  const savedTheme = localStorage.getItem('theme')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}
initializeTheme()
const app = createApp(App)
app.use(router)
app.mount('#app')
