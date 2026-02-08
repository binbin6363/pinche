import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const admin = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    // simple admin login - in production should have dedicated admin auth
    const data = await api.post('/admin/login', { username, password })
    token.value = data.token
    admin.value = data.admin
    localStorage.setItem('admin_token', data.token)
    return data
  }

  function logout() {
    token.value = ''
    admin.value = null
    localStorage.removeItem('admin_token')
  }

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }

  return {
    token,
    admin,
    isLoggedIn,
    login,
    logout,
    setToken
  }
})
