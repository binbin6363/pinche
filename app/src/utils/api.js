import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import router from '@/router'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code !== 0) {
      const appStore = useAppStore()
      appStore.showToast(data.message, 'error')
      return Promise.reject(new Error(data.message))
    }
    return data.data
  },
  (error) => {
    const appStore = useAppStore()
    // extract error message from backend response
    const backendMsg = error.response?.data?.message
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      router.push('/login')
      appStore.showToast(backendMsg || '登录已过期，请重新登录', 'error')
    } else {
      appStore.showToast(backendMsg || '网络错误', 'error')
    }
    return Promise.reject(new Error(backendMsg || error.message))
  }
)

// upload image file
// bizType: 'images' for chat images, 'avatar' for user avatar
// for images: returns object key (e.g., "images/xxx.jpg")
// for avatar: returns public URL directly
export async function uploadImage(file, bizType = 'images') {
  const formData = new FormData()
  formData.append('file', file)

  const userStore = useUserStore()
  try {
    const response = await axios.post(`/api/upload/image?biz_type=${bizType}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${userStore.token}`
      },
      timeout: 30000
    })

    if (response.data.code !== 0) {
      throw new Error(response.data.message)
    }
    // for avatar and trip: return url; for images: return key
    if (bizType === 'avatar' || bizType === 'trip') {
      return response.data.data.url
    }
    return response.data.data.key
  } catch (error) {
    const msg = error.response?.data?.message || error.message || '上传失败'
    throw new Error(msg)
  }
}

// get signed URL for resource
// key: object key returned from upload (e.g., "images/xxx.jpg")
export async function getResourceUrl(key) {
  const userStore = useUserStore()
  try {
    const response = await axios.get('/api/resource/url', {
      params: { key },
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })

    if (response.data.code !== 0) {
      throw new Error(response.data.message)
    }
    return response.data.data.url
  } catch (error) {
    const msg = error.response?.data?.message || error.message || '获取资源失败'
    throw new Error(msg)
  }
}

export default api
