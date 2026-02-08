<template>
  <div class="register-page min-h-screen flex flex-col bg-gradient-to-b from-primary-500 to-primary-700">
    <!-- Logo区域 - 紧凑化设计 -->
    <div class="flex flex-col items-center px-6 pt-10 pb-4">
      <div class="w-14 h-14 bg-white/95 rounded-xl flex items-center justify-center shadow-lg mb-2 backdrop-blur-sm">
        <svg class="w-8 h-8 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
        </svg>
      </div>
      <h1 class="text-xl font-bold text-white">加入春节拼车</h1>
    </div>

    <!-- 表单区域 - 占据更多空间 -->
    <div class="flex-1 bg-white rounded-t-3xl px-6 pt-6 pb-8 safe-area-bottom overflow-y-auto">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">注册账号</h2>
      
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
            <label class="block text-sm text-gray-600 mb-1">昵称</label>
            <input
              v-model="form.nickname"
              type="text"
              maxlength="20"
              placeholder="请输入昵称"
              class="input"
            />
          </div>
          
          <div>
            <label class="block text-sm text-gray-600 mb-1">密码</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              class="input"
            />
          </div>

          <div>
            <label class="block text-sm text-gray-600 mb-1">确认密码</label>
            <input
              v-model="form.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              class="input"
            />
          </div>

          <!-- 身份选择 -->
          <div>
            <label class="block text-sm text-gray-600 mb-2">您的出行目的</label>
            <div class="grid grid-cols-2 gap-3">
              <div
                @click="form.identity = 2"
                class="flex flex-col items-center p-4 rounded-xl border-2 cursor-pointer transition-all"
                :class="form.identity === 2 ? 'border-primary-500 bg-primary-50' : 'border-gray-200 bg-gray-50'"
              >
                <span class="text-2xl mb-1">🙋</span>
                <span class="text-sm font-medium" :class="form.identity === 2 ? 'text-primary-600' : 'text-gray-600'">人找车</span>
                <span class="text-xs text-gray-400 mt-1">我是乘客</span>
              </div>
              <div
                @click="form.identity = 1"
                class="flex flex-col items-center p-4 rounded-xl border-2 cursor-pointer transition-all"
                :class="form.identity === 1 ? 'border-primary-500 bg-primary-50' : 'border-gray-200 bg-gray-50'"
              >
                <span class="text-2xl mb-1">🚗</span>
                <span class="text-sm font-medium" :class="form.identity === 1 ? 'text-primary-600' : 'text-gray-600'">车找人</span>
                <span class="text-xs text-gray-400 mt-1">我是车主</span>
              </div>
            </div>
          </div>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="btn btn-primary w-full mt-6"
        >
          <span v-if="loading" class="loading-spinner mr-2"></span>
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <p class="text-center text-gray-500 mt-6">
        已有账号？
        <router-link to="/login" class="text-primary-500 font-medium">立即登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { hashPassword } from '@/utils/crypto'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const loading = ref(false)
const form = reactive({
  phone: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  identity: 2 // 默认人找车（乘客）
})

async function handleSubmit() {
  if (!form.phone || form.phone.length !== 11) {
    appStore.showToast('请输入正确的手机号', 'error')
    return
  }
  if (!form.nickname || form.nickname.length < 2) {
    appStore.showToast('昵称至少2个字符', 'error')
    return
  }
  if (!form.password || form.password.length < 6) {
    appStore.showToast('密码至少6位', 'error')
    return
  }
  if (form.password !== form.confirmPassword) {
    appStore.showToast('两次密码输入不一致', 'error')
    return
  }

  loading.value = true
  try {
    await userStore.register({
      phone: form.phone,
      nickname: form.nickname,
      password: hashPassword(form.password),
      identity: form.identity
    })
    // 注册成功后设置默认身份
    userStore.setIdentity(form.identity)
    appStore.showToast('注册成功，请登录', 'success')
    router.replace('/login')
  } catch (e) {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>
