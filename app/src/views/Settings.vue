<template>
  <div class="settings-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部导航 -->
    <div class="page-header safe-area-top">
      <div class="page-header-bg"></div>
      <div class="relative flex items-center h-12 px-4">
        <button @click="goBack" class="w-10 h-10 -ml-2 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-white">设置</h1>
        <div class="w-10"></div>
      </div>
    </div>

    <div class="px-4 py-4 space-y-4">
      <!-- 身份切换 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">当前身份</span>
          <p class="text-xs text-gray-400 mt-1">切换后行程页显示对应信息</p>
        </div>

        <div class="p-4">
          <div class="grid grid-cols-2 gap-3">
            <button
              @click="switchIdentity(1)"
              class="py-4 rounded-xl flex flex-col items-center gap-2 transition-all border-2"
              :class="userStore.identity === 1 
                ? 'bg-blue-50 border-blue-500 text-blue-600' 
                : appStore.theme === 'dark' 
                  ? 'bg-gray-700 border-gray-600 text-gray-300' 
                  : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'"
            >
              <span class="text-2xl">🚗</span>
              <span class="text-sm font-medium">我是司机</span>
              <span class="text-xs opacity-70">发布车找人</span>
            </button>

            <button
              @click="switchIdentity(2)"
              class="py-4 rounded-xl flex flex-col items-center gap-2 transition-all border-2"
              :class="userStore.identity === 2 
                ? 'bg-green-50 border-green-500 text-green-600' 
                : appStore.theme === 'dark' 
                  ? 'bg-gray-700 border-gray-600 text-gray-300' 
                  : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'"
            >
              <span class="text-2xl">🙋</span>
              <span class="text-sm font-medium">我是乘客</span>
              <span class="text-xs opacity-70">发布人找车</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 外观设置 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">外观设置</span>
        </div>

        <div class="p-4">
          <p class="text-sm text-gray-500 mb-3">选择应用主题</p>
          <div class="grid grid-cols-3 gap-3">
            <!-- 浅色主题 -->
            <button
              @click="setTheme('light')"
              class="relative p-3 rounded-xl border-2 transition-all"
              :class="appStore.theme === 'light' 
                ? 'border-blue-500 bg-blue-50' 
                : 'border-gray-200 bg-white hover:border-gray-300'"
            >
              <div class="w-full aspect-[3/4] rounded-lg overflow-hidden mb-2 shadow-sm">
                <div class="h-1/4 bg-gradient-to-r from-blue-500 to-cyan-500"></div>
                <div class="h-3/4 bg-gray-50 p-1.5">
                  <div class="w-full h-2 bg-gray-200 rounded mb-1"></div>
                  <div class="w-2/3 h-2 bg-gray-200 rounded"></div>
                </div>
              </div>
              <span class="text-xs font-medium">浅色</span>
              <div 
                v-if="appStore.theme === 'light'"
                class="absolute -top-1 -right-1 w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center"
              >
                <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </button>

            <!-- 深色主题 -->
            <button
              @click="setTheme('dark')"
              class="relative p-3 rounded-xl border-2 transition-all"
              :class="appStore.theme === 'dark' 
                ? 'border-blue-500 bg-gray-800' 
                : appStore.theme === 'dark' ? 'border-gray-600 bg-gray-700' : 'border-gray-200 bg-white hover:border-gray-300'"
            >
              <div class="w-full aspect-[3/4] rounded-lg overflow-hidden mb-2 shadow-sm">
                <div class="h-1/4 bg-gradient-to-r from-gray-700 to-gray-600"></div>
                <div class="h-3/4 bg-gray-800 p-1.5">
                  <div class="w-full h-2 bg-gray-600 rounded mb-1"></div>
                  <div class="w-2/3 h-2 bg-gray-600 rounded"></div>
                </div>
              </div>
              <span class="text-xs font-medium" :class="appStore.theme === 'dark' ? 'text-gray-300' : ''">深色</span>
              <div 
                v-if="appStore.theme === 'dark'"
                class="absolute -top-1 -right-1 w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center"
              >
                <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </button>

            <!-- 新春主题 -->
            <button
              @click="setTheme('spring')"
              class="relative p-3 rounded-xl border-2 transition-all"
              :class="appStore.theme === 'spring' 
                ? 'border-red-500 bg-red-50' 
                : appStore.theme === 'dark' ? 'border-gray-600 bg-gray-700' : 'border-gray-200 bg-white hover:border-gray-300'"
            >
              <div class="w-full aspect-[3/4] rounded-lg overflow-hidden mb-2 shadow-sm">
                <div class="h-1/4 bg-gradient-to-r from-red-600 to-orange-500 relative">
                  <span class="absolute top-0 right-0.5 text-[8px]">🧧</span>
                </div>
                <div class="h-3/4 bg-gray-50 p-1.5">
                  <div class="w-full h-2 bg-red-100 rounded mb-1"></div>
                  <div class="w-2/3 h-2 bg-red-100 rounded"></div>
                </div>
              </div>
              <span class="text-xs font-medium" :class="appStore.theme === 'dark' ? 'text-gray-300' : ''">新春</span>
              <div 
                v-if="appStore.theme === 'spring'"
                class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full flex items-center justify-center"
              >
                <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- 其他设置 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">其他</span>
        </div>

        <div 
          class="list-item border-b border-gray-100"
          @click="clearCache"
        >
          <span class="text-sm">清除缓存</span>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm">当前版本</span>
          <span class="text-sm text-gray-400">v1.0.0</span>
        </div>

        <div 
          class="list-item"
          @click="showAbout = true"
        >
          <span class="text-sm">关于我们</span>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>

      <!-- 退出登录 -->
      <button 
        v-if="userStore.isLoggedIn"
        @click="handleLogout" 
        class="btn btn-danger w-full"
      >
        退出登录
      </button>
    </div>

    <!-- 关于我们弹窗 -->
    <div v-if="showAbout" class="action-sheet" @click.self="showAbout = false">
      <div class="action-sheet-overlay" @click="showAbout = false"></div>
      <div class="action-sheet-content">
        <div class="p-6 text-center">
          <div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-red-500 to-orange-500 rounded-2xl flex items-center justify-center">
            <span class="text-3xl">🚗</span>
          </div>
          <h3 class="text-lg font-bold mb-2">春节拼车</h3>
          <p class="text-sm text-gray-500 mb-4">让回家的路不再孤单</p>
          
          <div class="text-xs text-gray-400 space-y-1">
            <p>本平台仅提供信息发布服务</p>
            <p>请在拼车前核实对方身份</p>
            <p>注意保护个人财产安全</p>
          </div>

          <button 
            @click="showAbout = false" 
            class="btn btn-primary w-full mt-6"
          >
            我知道了
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const showAbout = ref(false)

function goBack() {
  router.back()
}

function setTheme(theme) {
  appStore.setTheme(theme)
  appStore.showToast(`已切换到${appStore.themeConfig.name}主题`, 'success')
}

function switchIdentity(identity) {
  userStore.setIdentity(identity)
  appStore.showToast(identity === 1 ? '🚗 已切换为司机身份' : '🙋 已切换为乘客身份', 'success')
}

function clearCache() {
  localStorage.removeItem('safety_tip_closed')
  appStore.showToast('缓存已清除', 'success')
}

function handleLogout() {
  if (!confirm('确定要退出登录吗？')) return
  userStore.logout()
  router.replace('/login')
}
</script>
