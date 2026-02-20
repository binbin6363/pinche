<template>
  <div class="friend-requests-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-100'">
    <!-- 顶部导航 -->
    <div class="sticky top-0 z-10 px-4 py-3 flex items-center justify-between" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
      <div class="flex items-center gap-3">
        <button @click="goBack" class="w-8 h-8 flex items-center justify-center">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-700'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="text-lg font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">好友申请</h1>
      </div>
      <span class="text-sm text-gray-500">{{ friendStore.requestCount }} 条待处理</span>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <!-- 申请列表 -->
    <div v-else-if="friendStore.friendRequests.length > 0" class="p-4 space-y-3">
      <div 
        v-for="request in friendStore.friendRequests" 
        :key="request.id"
        class="rounded-xl p-4"
        :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'"
      >
        <!-- 用户信息 -->
        <div class="flex items-center gap-3 mb-3" @click="goToProfile(request.user?.open_id)">
          <div class="w-12 h-12 rounded-full flex items-center justify-center overflow-hidden flex-shrink-0 cursor-pointer" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
            <img v-if="request.user?.avatar" :src="request.user.avatar" class="w-full h-full object-cover" />
            <span v-else class="text-lg font-bold" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'">
              {{ request.user?.nickname?.charAt(0) || '?' }}
            </span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium truncate" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
              {{ request.user?.nickname }}
            </div>
            <div class="text-xs text-gray-500">
              {{ formatTime(request.created_at) }}
            </div>
          </div>
        </div>

        <!-- 申请留言 -->
        <div v-if="request.message" class="mb-3 p-3 rounded-lg text-sm" :class="appStore.theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-50 text-gray-600'">
          {{ request.message }}
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-3">
          <button 
            @click="reject(request.id)"
            :disabled="processing === request.id"
            class="flex-1 py-2 rounded-lg font-medium transition-colors"
            :class="appStore.theme === 'dark' ? 'bg-gray-700 text-white hover:bg-gray-600' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
          >
            拒绝
          </button>
          <button 
            @click="accept(request.id)"
            :disabled="processing === request.id"
            class="flex-1 py-2 rounded-lg font-medium text-white transition-colors disabled:opacity-50"
            style="background: var(--theme-primary);"
          >
            {{ processing === request.id ? '处理中...' : '同意' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-20">
      <svg class="w-20 h-20 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
      </svg>
      <p class="text-gray-500 mb-2">暂无好友申请</p>
      <p class="text-sm text-gray-400">当有人向你发送好友申请时会在这里显示</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useFriendStore } from '@/stores/friend'

const router = useRouter()
const appStore = useAppStore()
const friendStore = useFriendStore()

const loading = ref(true)
const processing = ref(null)

onMounted(async () => {
  await friendStore.fetchFriendRequests()
  loading.value = false
})

function goBack() {
  router.back()
}

function goToProfile(openId) {
  if (openId) {
    router.push(`/user/${openId}`)
  }
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

async function accept(requestId) {
  if (processing.value) return
  processing.value = requestId
  await friendStore.acceptFriendRequest(requestId)
  appStore.showToast('已添加好友', 'success')
  processing.value = null
}

async function reject(requestId) {
  if (processing.value) return
  processing.value = requestId
  await friendStore.rejectFriendRequest(requestId)
  appStore.showToast('已拒绝申请', 'info')
  processing.value = null
}
</script>
