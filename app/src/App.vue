<template>
  <div class="app-container h-full bg-gray-50">
    <router-view v-slot="{ Component }">
      <keep-alive :include="['Home', 'Trips', 'Matches', 'Profile']">
        <component :is="Component" />
      </keep-alive>
    </router-view>
    
    <!-- Toast通知 -->
    <div v-if="toast.show" class="fixed top-16 left-4 right-4 z-50">
      <div
        class="p-4 rounded-lg shadow-lg text-white text-center"
        :class="{
          'bg-green-500': toast.type === 'success',
          'bg-red-500': toast.type === 'error',
          'bg-primary-500': toast.type === 'info'
        }"
      >
        {{ toast.message }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { setOnMessageCallback } from '@/utils/websocket'

const appStore = useAppStore()
const userStore = useUserStore()
const toast = computed(() => appStore.toast)

// handle websocket messages
function handleWebSocketMessage(message) {
  switch (message.type) {
    case 'match_found':
      appStore.showToast('发现新的匹配！', 'info')
      if (message.data?.notification) {
        appStore.addNotification(message.data.notification)
      }
      break
    case 'match_success':
      appStore.showToast('拼车成功！查看联系方式', 'success')
      if (message.data?.notification) {
        appStore.addNotification(message.data.notification)
      }
      break
    case 'match_rejected':
      appStore.showToast('匹配未成功', 'info')
      if (message.data?.notification) {
        appStore.addNotification(message.data.notification)
      }
      break
    case 'new_message':
      appStore.showToast('收到新消息', 'info')
      break
  }
}

onMounted(() => {
  // set websocket message callback
  setOnMessageCallback(handleWebSocketMessage)
  // initialize WebSocket after app is mounted
  userStore.initWebSocket()
})
</script>
