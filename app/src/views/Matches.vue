<template>
  <div class="messages-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-100'">
    <!-- 微信风格搜索框 -->
    <div class="px-3 py-2" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-[#EDEDED]'">
      <div class="relative">
        <div class="flex items-center rounded-md px-3 py-2" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-white'">
          <svg class="w-4 h-4 flex-shrink-0" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索"
            class="flex-1 ml-2 text-sm bg-transparent outline-none"
            :class="appStore.theme === 'dark' ? 'placeholder-gray-500 text-white' : 'placeholder-gray-400'"
          />
          <button
            v-if="searchKeyword"
            @click="searchKeyword = ''"
            class="p-0.5 rounded-full bg-gray-300 hover:bg-gray-400"
          >
            <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 微信风格的消息列表 -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="loading-spinner"></div>
    </div>

    <div v-else-if="filteredConversations.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
        <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <p class="text-gray-500 mb-2">{{ searchKeyword ? '未找到相关聊天' : '暂无消息' }}</p>
      <p v-if="!searchKeyword" class="text-gray-400 text-sm">在行程页面点击"私聊"开始沟通</p>
    </div>

    <!-- 微信风格会话列表 -->
    <div v-else :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
      <div
        v-for="(conv, index) in filteredConversations"
        :key="conv.peer_id"
        class="flex items-center px-4 py-3 transition-colors cursor-pointer"
        :class="[
          appStore.theme === 'dark' ? 'active:bg-gray-700' : 'active:bg-gray-100',
          index < filteredConversations.length - 1 ? (appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100') : ''
        ]"
      >
        <!-- 头像 -->
        <div class="relative flex-shrink-0 mr-3" @click="goUserProfile(conv)">
          <!-- 系统通知特殊图标 -->
          <div
            v-if="conv.is_system"
            class="w-12 h-12 rounded-lg bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center"
          >
            <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
          </div>
          <!-- 普通用户头像 -->
          <div
            v-else
            class="w-12 h-12 rounded-lg flex items-center justify-center overflow-hidden"
            :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"
          >
            <img v-if="conv.peer_avatar" :src="conv.peer_avatar" class="w-full h-full object-cover" />
            <span v-else class="text-lg font-medium text-gray-500">{{ getPeerInitial(conv) }}</span>
          </div>
          <!-- 未读红点 -->
          <div
            v-if="conv.unread_count > 0"
            class="absolute -top-1 -right-1 min-w-[18px] h-[18px] px-1 bg-red-500 rounded-full flex items-center justify-center"
          >
            <span class="text-[11px] text-white font-medium">{{ conv.unread_count > 99 ? '99+' : conv.unread_count }}</span>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="flex-1 min-w-0" @click="goChat(conv)">
          <div class="flex items-center justify-between">
            <span class="font-medium truncate text-[15px]" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
              {{ conv.is_system ? '系统通知' : (conv.peer_nickname || '用户') }}
            </span>
            <span class="text-xs text-gray-400 flex-shrink-0 ml-2">{{ formatTime(conv.last_message_time) }}</span>
          </div>
          <div class="flex items-center mt-1">
            <!-- 消息前缀 -->
            <span v-if="conv.last_message_type === 2" class="text-[13px] text-gray-400">[图片]</span>
            <span v-else-if="conv.last_message_type === 3" class="text-[13px] text-gray-400">[语音]</span>
            <span v-else-if="conv.last_message_type === 4" class="text-[13px] text-gray-400">[表情]</span>
            <span v-else-if="conv.last_message_type === 5" class="text-[13px] text-gray-400">[通话]</span>
            <span v-else-if="conv.last_message_type === 6" class="text-[13px] text-gray-400">[视频]</span>
            <span v-else class="text-[13px] text-gray-500 truncate">{{ conv.last_message_content || '暂无消息' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessageStore } from '@/stores/message'
import { useAppStore } from '@/stores/app'
import { addMessageListener, removeMessageListener } from '@/utils/websocket'

const router = useRouter()
const messageStore = useMessageStore()
const appStore = useAppStore()

const conversations = ref([])
const systemNotifications = ref([])
const loading = ref(true)
const searchKeyword = ref('')

// system user open_id - a fixed string
const SYSTEM_USER_OPEN_ID = 'system'

// merge system notification and normal conversations
const allConversations = computed(() => {
  const list = []
  
  // add system notification entry if has notifications
  if (systemNotifications.value.length > 0) {
    const latest = systemNotifications.value[0]
    list.push({
      peer_id: SYSTEM_USER_OPEN_ID,
      peer_nickname: '系统通知',
      peer_avatar: '',
      is_system: true,
      last_message_content: latest.title || latest.content,
      last_message_time: latest.created_at,
      last_message_type: 1,
      unread_count: systemNotifications.value.filter(n => !n.is_read).length
    })
  }
  
  // add normal conversations (exclude system user)
  // map conversation data to include peer info from peer object
  const normalConvs = conversations.value
    .filter(c => c.peer_id !== SYSTEM_USER_OPEN_ID)
    .map(c => ({
      peer_id: c.peer_id,
      peer_nickname: c.peer?.nickname || '用户',
      peer_avatar: c.peer?.avatar || '',
      is_system: false,
      last_message_content: c.last_message?.content || '',
      last_message_time: c.last_message_at,
      last_message_type: c.last_message?.msg_type || 1,
      unread_count: c.unread_count || 0
    }))
  list.push(...normalConvs)
  
  // sort by last message time
  list.sort((a, b) => {
    const timeA = new Date(a.last_message_time || 0).getTime()
    const timeB = new Date(b.last_message_time || 0).getTime()
    return timeB - timeA
  })
  
  return list
})

// filter conversations by search keyword
const filteredConversations = computed(() => {
  if (!searchKeyword.value.trim()) {
    return allConversations.value
  }
  const keyword = searchKeyword.value.toLowerCase()
  return allConversations.value.filter(conv => {
    const nickname = (conv.peer_nickname || '').toLowerCase()
    const content = (conv.last_message_content || '').toLowerCase()
    return nickname.includes(keyword) || content.includes(keyword)
  })
})

onMounted(async () => {
  try {
    const [convResult, notifyResult] = await Promise.all([
      messageStore.fetchConversations(),
      fetchSystemNotifications()
    ])
    conversations.value = convResult.list || []
    systemNotifications.value = notifyResult || []
  } finally {
    loading.value = false
  }
  
  // listen for new messages to update conversation list
  addMessageListener(handleNewMessage)
})

onUnmounted(() => {
  removeMessageListener(handleNewMessage)
})

// handle new message from websocket
function handleNewMessage(msg) {
  // find the conversation with this peer
  const peerId = msg.sender_id
  const convIndex = conversations.value.findIndex(c => c.peer_id === peerId)
  
  if (convIndex >= 0) {
    // update existing conversation
    const conv = conversations.value[convIndex]
    conv.last_message = {
      content: msg.content,
      msg_type: msg.msg_type
    }
    conv.last_message_at = msg.created_at
    conv.unread_count = (conv.unread_count || 0) + 1
    
    // move to top by removing and re-adding
    conversations.value.splice(convIndex, 1)
    conversations.value.unshift(conv)
  } else {
    // new conversation - refetch the list
    messageStore.fetchConversations().then(result => {
      conversations.value = result.list || []
    })
  }
}

async function fetchSystemNotifications() {
  try {
    const response = await fetch('/api/notifications?page=1&page_size=10', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    const data = await response.json()
    return data.data?.list || []
  } catch {
    return []
  }
}

function goChat(conv) {
  if (conv.is_system) {
    router.push('/notifications')
  } else {
    // conv.peer_id is open_id
    router.push({
      path: `/chat/${conv.peer_id}`,
      query: { 
        nickname: conv.peer_nickname,
        avatar: conv.peer_avatar
      }
    })
  }
}

function goUserProfile(conv) {
  if (!conv.is_system && conv.peer_id) {
    router.push(`/user/${conv.peer_id}`)
  } else if (conv.is_system) {
    router.push('/notifications')
  }
}

function getPeerInitial(conv) {
  if (conv.peer_nickname) {
    return conv.peer_nickname.charAt(0)
  }
  return 'U'
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diffMs = now - date
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    const hour = date.getHours().toString().padStart(2, '0')
    const minute = date.getMinutes().toString().padStart(2, '0')
    return `${hour}:${minute}`
  } else if (diffDays === 1) {
    return '昨天'
  } else if (diffDays < 7) {
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  } else {
    const month = date.getMonth() + 1
    const day = date.getDate()
    return `${month}/${day}`
  }
}
</script>
