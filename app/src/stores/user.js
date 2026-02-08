import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'
import { connectWebSocket, disconnectWebSocket } from '@/utils/websocket'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  // identity: 1=driver, 2=passenger, default passenger
  const identity = ref(parseInt(localStorage.getItem('identity')) || 2)

  const isLoggedIn = computed(() => !!token.value)

  function setIdentity(newIdentity) {
    identity.value = newIdentity
    localStorage.setItem('identity', newIdentity.toString())
  }

  async function register(data) {
    const result = await api.post('/user/register', data)
    return result
  }

  async function login(data) {
    const result = await api.post('/user/login', data)
    token.value = result.token
    user.value = result.user
    localStorage.setItem('token', result.token)
    localStorage.setItem('user', JSON.stringify(result.user))
    connectWebSocket(result.token)
    return result
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    disconnectWebSocket()
  }

  async function fetchProfile() {
    if (!token.value) return
    const result = await api.get('/user/profile')
    user.value = result
    localStorage.setItem('user', JSON.stringify(result))
    return result
  }

  async function updateProfile(data) {
    const result = await api.put('/user/profile', data)
    user.value = result
    localStorage.setItem('user', JSON.stringify(result))
    return result
  }

  // initialize WebSocket connection (called from App.vue onMounted)
  function initWebSocket() {
    if (token.value) {
      connectWebSocket(token.value)
    }
  }

  // get user profile completeness
  const profileCompleteness = computed(() => {
    if (!user.value) return 0
    let score = 0
    if (user.value.nickname) score += 20
    if (user.value.avatar) score += 20
    if (user.value.phone) score += 20
    if (user.value.contact_wechat || user.value.contact_phone) score += 20
    if (user.value.car_number) score += 20
    return score
  })

  return {
    token,
    user,
    identity,
    isLoggedIn,
    profileCompleteness,
    register,
    login,
    logout,
    fetchProfile,
    updateProfile,
    setIdentity,
    initWebSocket
  }
})
