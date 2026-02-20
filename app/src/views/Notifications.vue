<template>
  <div class="notifications-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部 -->
    <div class="border-b safe-area-top" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-center h-12 px-4">
        <button @click="goBack" class="p-2 -ml-2">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : ''">消息通知</h1>
        <button 
          v-if="notifications.length > 0"
          @click="markAllRead"
          class="text-sm text-primary-500"
        >
          全部已读
        </button>
        <div v-else class="w-14"></div>
      </div>
    </div>

    <div class="px-4 py-4">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="notifications.length === 0" class="text-center py-12">
        <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        </div>
        <p class="text-gray-500">暂无消息</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          @click="handleNotificationClick(notification)"
          class="card p-4 cursor-pointer"
          :class="{ 'bg-blue-50': notification.is_read === 0 }"
        >
          <div class="flex items-start gap-3">
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
              :class="notification.is_read === 0 ? 'bg-primary-100' : 'bg-gray-100'"
            >
              <svg class="w-5 h-5" :class="notification.is_read === 0 ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-1">
                <span class="text-sm font-medium text-gray-800">{{ notification.title }}</span>
                <span class="text-xs text-gray-400">{{ formatTime(notification.created_at) }}</span>
              </div>
              <p class="text-sm text-gray-600 line-clamp-2">{{ notification.content }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const appStore = useAppStore()

const notifications = ref([])
const loading = ref(true)

async function fetchNotifications() {
  try {
    const result = await api.get('/notifications')
    notifications.value = result.list || []
    appStore.setUnreadCount(result.unread || 0)
  } catch (e) {
    // ignore
  }
}

onMounted(async () => {
  await fetchNotifications()
  loading.value = false
})

// watch appStore.notifications for new websocket notifications
// when a new notification is added via websocket, refresh the list
watch(
  () => appStore.notifications.length,
  (newLen, oldLen) => {
    if (newLen > oldLen) {
      fetchNotifications()
    }
  }
)

function goBack() {
  router.back()
}

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
}

async function handleNotificationClick(notification) {
  // 标记已读
  if (notification.is_read === 0) {
    try {
      await api.put(`/notifications/${notification.id}/read`)
      notification.is_read = 1
      appStore.setUnreadCount(Math.max(0, appStore.unreadCount - 1))
    } catch (e) {
      // ignore
    }
  }
  
  // 跳转：优先 match_id，其次 trip_id（抢单通知）
  if (notification.match_id && notification.match_id > 0) {
    router.push(`/match/${notification.match_id}`)
  } else if (notification.trip_id && notification.trip_id > 0) {
    router.push(`/my-trip/${notification.trip_id}`)
  }
}

async function markAllRead() {
  try {
    await api.put('/notifications/read-all')
    notifications.value.forEach(n => n.is_read = 1)
    appStore.setUnreadCount(0)
    appStore.showToast('已全部标记为已读', 'success')
  } catch (e) {
    // error handled in interceptor
  }
}
</script>
