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

    <!-- 来电提醒组件 -->
    <IncomingCall />
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import { useFriendStore } from '@/stores/friend'
import { useCallStore } from '@/stores/call'
import { setOnMessageCallback, addCallListener, removeCallListener } from '@/utils/websocket'
import IncomingCall from '@/components/IncomingCall.vue'

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()
const messageStore = useMessageStore()
const friendStore = useFriendStore()
const callStore = useCallStore()
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
    case 'trip_grabbed':
      appStore.showToast(message.data?.grabber_name ? `${message.data.grabber_name} 想搭您的车` : '有人想搭您的车', 'info')
      if (message.data?.notification) {
        appStore.addNotification(message.data.notification)
      }
      break
    case 'new_message':
      // increment unread count if not currently viewing the chat
      if (!route.path.startsWith('/chat/')) {
        messageStore.incrementUnreadCount()
      }
      break
    case 'friend_request':
      // new friend request received
      appStore.showToast('收到新的好友申请', 'info')
      friendStore.fetchFriendCount()
      break
  }
}

// 处理通话信令消息
function handleCallSignaling(message) {
  const { type, data } = message

  // 只处理来电邀请，其他信令由 VideoCall 组件处理
  if (type === 'call_invite') {
    const { call_id, call_type, caller_info, from_open_id } = data
    
    // 触发来电 (使用 from_open_id 作为用户标识)
    const accepted = callStore.receiveCall(
      call_id,
      from_open_id,
      caller_info?.nickname || '未知用户',
      caller_info?.avatar || '',
      call_type || 'audio'
    )

    if (!accepted) {
      console.log('App: Incoming call rejected (already in call)')
    }
  }
}

onMounted(async () => {
  // set websocket message callback
  setOnMessageCallback(handleWebSocketMessage)
  
  // 添加通话信令监听
  addCallListener(handleCallSignaling)
  
  // initialize WebSocket after app is mounted
  userStore.initWebSocket()
  // fetch latest user profile on app start
  if (userStore.isLoggedIn) {
    await userStore.fetchProfile()
    // fetch friend counts for badge display
    await friendStore.fetchFriendCount()
  }
})

onUnmounted(() => {
  removeCallListener(handleCallSignaling)
})
</script>
