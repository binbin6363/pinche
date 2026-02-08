<template>
  <div class="login-page min-h-screen flex flex-col bg-gradient-to-b from-primary-500 to-primary-700">
    <!-- Logo区域 - 紧凑化设计 -->
    <div class="flex flex-col items-center px-6 pt-12 pb-6">
      <div class="w-16 h-16 bg-white/95 rounded-2xl flex items-center justify-center shadow-lg mb-3 backdrop-blur-sm">
        <svg class="w-9 h-9 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-white mb-1">春节拼车</h1>
      <p class="text-primary-100 text-sm">回家的路，一起走</p>
    </div>

    <!-- 表单区域 - 占据更多空间 -->
    <div class="flex-1 bg-white rounded-t-3xl px-6 pt-8 pb-8 safe-area-bottom">
      <h2 class="text-xl font-semibold text-gray-800 mb-6">登录账号</h2>
      
      <form @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-gray-600 mb-1">手机号</label>
            <input
              v-model="form.phone"
              type="tel"
              maxlength="11"
              placeholder="请输入手机号"
              class="input"
            />
          </div>
          
          <div>
            <label class="block text-sm text-gray-600 mb-1">密码</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              class="input"
            />
          </div>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="btn btn-primary w-full mt-6"
        >
          <span v-if="loading" class="loading-spinner mr-2"></span>
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <p class="text-center text-gray-500 mt-6">
        还没有账号？
        <router-link to="/register" class="text-primary-500 font-medium">立即注册</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { hashPassword } from '@/utils/crypto'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const loading = ref(false)
const form = reactive({
  phone: '',
  password: ''
})

async function handleSubmit() {
  if (!form.phone || form.phone.length !== 11) {
    appStore.showToast('请输入正确的手机号', 'error')
    return
  }
  if (!form.password || form.password.length < 6) {
    appStore.showToast('密码至少6位', 'error')
    return
  }

  loading.value = true
  try {
    await userStore.login({
      phone: form.phone,
      password: hashPassword(form.password)
    })
    appStore.showToast('登录成功', 'success')
    const redirect = route.query.redirect || '/'
    router.replace(redirect)
  } catch (e) {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>
