import { defineStore } from 'pinia'
import { ref, reactive, computed, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  const toast = reactive({
    show: false,
    message: '',
    type: 'info'
  })

  const notifications = ref([])
  const unreadCount = ref(0)

  // theme: light | dark | spring
  const theme = ref(localStorage.getItem('app_theme') || 'spring')

  // theme config
  const themeConfig = computed(() => {
    switch (theme.value) {
      case 'dark':
        return {
          name: '深色',
          primary: '#60A5FA',
          primaryGradient: 'from-gray-800 via-gray-700 to-gray-600',
          bgPrimary: '#1F2937',
          bgSecondary: '#111827',
          textPrimary: '#F9FAFB',
          textSecondary: '#9CA3AF',
          cardBg: '#374151',
          tabBarBg: '#1F2937',
          headerGradient: 'from-gray-800 via-gray-700 to-gray-600'
        }
      case 'spring':
        return {
          name: '新春',
          primary: '#EF4444',
          primaryGradient: 'from-red-600 via-red-500 to-orange-500',
          bgPrimary: '#FFFFFF',
          bgSecondary: '#F9FAFB',
          textPrimary: '#111827',
          textSecondary: '#6B7280',
          cardBg: '#FFFFFF',
          tabBarBg: '#FFFFFF',
          headerGradient: 'from-red-600 via-red-500 to-orange-500'
        }
      default: // light
        return {
          name: '浅色',
          primary: '#3B82F6',
          primaryGradient: 'from-blue-600 via-blue-500 to-cyan-500',
          bgPrimary: '#FFFFFF',
          bgSecondary: '#F3F4F6',
          textPrimary: '#111827',
          textSecondary: '#6B7280',
          cardBg: '#FFFFFF',
          tabBarBg: '#FFFFFF',
          headerGradient: 'from-blue-600 via-blue-500 to-cyan-500'
        }
    }
  })

  // apply theme to document
  function applyTheme() {
    const root = document.documentElement
    const config = themeConfig.value
    
    root.setAttribute('data-theme', theme.value)
    
    // update CSS variables
    root.style.setProperty('--theme-primary', config.primary)
    root.style.setProperty('--theme-bg-primary', config.bgPrimary)
    root.style.setProperty('--theme-bg-secondary', config.bgSecondary)
    root.style.setProperty('--theme-text-primary', config.textPrimary)
    root.style.setProperty('--theme-text-secondary', config.textSecondary)
    root.style.setProperty('--theme-card-bg', config.cardBg)
    
    // dark mode class
    if (theme.value === 'dark') {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem('app_theme', newTheme)
    applyTheme()
  }

  // watch theme changes
  watch(theme, applyTheme, { immediate: true })

  let toastTimer = null

  function showToast(message, type = 'info', duration = 3000) {
    if (toastTimer) {
      clearTimeout(toastTimer)
    }
    toast.show = true
    toast.message = message
    toast.type = type
    toastTimer = setTimeout(() => {
      toast.show = false
    }, duration)
  }

  function addNotification(notification) {
    notifications.value.unshift(notification)
    unreadCount.value++
  }

  function setUnreadCount(count) {
    unreadCount.value = count
  }

  return {
    toast,
    notifications,
    unreadCount,
    theme,
    themeConfig,
    showToast,
    addNotification,
    setUnreadCount,
    setTheme,
    applyTheme
  }
})
