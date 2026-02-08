import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

api.interceptors.request.use(config => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

api.interceptors.response.use(
  response => {
    const data = response.data
    if (data.code !== 0) {
      return Promise.reject(new Error(data.message))
    }
    return data.data
  },
  error => {
    // extract error message from backend response
    const backendMsg = error.response?.data?.message
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/admin/login'
    }
    return Promise.reject(new Error(backendMsg || error.message || '请求失败'))
  }
)

export default api
