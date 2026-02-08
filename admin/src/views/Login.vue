<template>
  <div class="min-h-screen bg-gray-100 flex flex-col justify-center px-6 safe-area-top safe-area-bottom">
    <div class="w-full max-w-sm mx-auto">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="w-20 h-20 mx-auto mb-4 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center shadow-lg">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </div>
        <h1 class="text-xl font-bold text-gray-900">运营后台</h1>
        <p class="text-gray-500 mt-1 text-sm">拼车平台管理系统</p>
      </div>

      <!-- 登录表单 -->
      <div class="card p-6">
        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
            <input
              v-model="username"
              type="text"
              class="input"
              placeholder="请输入管理员用户名"
              autocomplete="username"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
            <input
              v-model="password"
              type="password"
              class="input"
              placeholder="请输入密码"
              autocomplete="current-password"
              required
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full btn btn-primary"
          >
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>

        <p v-if="error" class="mt-4 text-center text-sm text-red-500">{{ error }}</p>
      </div>

      <div class="mt-6 text-center text-xs text-gray-400">
        <p>请在服务端配置 ADMIN_USERNAME 和 ADMIN_PASSWORD</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { hashPassword } from '@/utils/crypto'
import api from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  error.value = ''
  loading.value = true

  try {
    const data = await api.post('/admin/login', {
      username: username.value,
      password: hashPassword(password.value)
    })
    authStore.setToken(data.token)
    router.replace('/')
  } catch (e) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>
