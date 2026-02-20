<template>
  <div class="chat-page" :class="appStore.theme === 'dark' ? 'chat-page--dark' : ''">
    <!-- 顶部导航 - 固定 -->
    <div class="chat-header border-b safe-area-top" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-center h-12 px-4">
        <button @click="goBack" class="p-2 -ml-2">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold truncate" :class="appStore.theme === 'dark' ? 'text-white' : ''">{{ peerNickname }}</h1>
        <div class="w-10"></div>
      </div>
    </div>

    <!-- 消息列表 - 中间可滚动区域 -->
    <div ref="messageListRef" class="chat-messages" @scroll="handleScroll">
      <!-- 加载更多 -->
      <div v-if="hasMore" class="flex justify-center mb-4">
        <button
          @click="loadMoreMessages"
          :disabled="loadingMore"
          class="px-4 py-1 text-xs rounded-full"
          :class="appStore.theme === 'dark' ? 'text-gray-400 bg-gray-700' : 'text-gray-500 bg-gray-100'"
        >
          {{ loadingMore ? '加载中...' : '加载更多' }}
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="messages.length === 0" class="text-center py-8 text-gray-400 text-sm">
        暂无消息，开始聊天吧
      </div>

      <!-- 消息列表（倒序显示，最新在底部） -->
      <div v-else class="space-y-4">
        <div
          v-for="(msg, index) in reversedMessages"
          :key="msg.id"
        >
          <!-- 时间分隔 -->
          <div
            v-if="shouldShowTime(msg, index)"
            class="flex justify-center mb-3"
          >
            <span class="px-2 py-1 text-xs rounded" :class="appStore.theme === 'dark' ? 'text-gray-500 bg-gray-700' : 'text-gray-400 bg-gray-100'">
              {{ formatMessageTime(msg.created_at) }}
            </span>
          </div>

          <!-- 消息气泡 -->
          <div
            class="flex items-end gap-2"
            :class="msg.sender_id === currentUserOpenId ? 'flex-row-reverse' : ''"
          >
            <!-- 头像 -->
            <div class="w-8 h-8 rounded-full flex-shrink-0 flex items-center justify-center overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
              <img 
                v-if="msg.sender_id === currentUserOpenId && userStore.user?.avatar" 
                :src="userStore.user.avatar" 
                class="w-full h-full object-cover" 
              />
              <img 
                v-else-if="msg.sender_id !== currentUserOpenId && peerAvatar" 
                :src="peerAvatar" 
                class="w-full h-full object-cover" 
              />
              <span v-else class="text-xs font-medium text-gray-500">
                {{ msg.sender_id === currentUserOpenId ? myInitial : peerInitial }}
              </span>
            </div>

            <!-- 消息内容 -->
            <div
              class="max-w-[70%]"
              :class="(msg.msg_type === 1 || msg.msg_type === 3)
                ? (msg.sender_id === currentUserOpenId
                  ? 'px-3 py-2 rounded-2xl bg-primary-500 text-white rounded-br-sm'
                  : appStore.theme === 'dark' ? 'px-3 py-2 rounded-2xl bg-gray-700 text-gray-200 rounded-bl-sm' : 'px-3 py-2 rounded-2xl bg-white text-gray-800 rounded-bl-sm shadow-sm')
                : ''"
            >
              <!-- 文字消息 -->
              <p v-if="msg.msg_type === 1" class="text-sm whitespace-pre-wrap break-words">
                {{ msg.content }}
              </p>
              <!-- 图片消息 -->
              <img
                v-else-if="msg.msg_type === 2"
                :src="imageUrlCache.get(msg.content) || ''"
                @click="previewImage(msg.content)"
                @load="onImageLoad"
                class="max-w-full rounded-xl cursor-pointer"
                style="max-height: 200px"
                :data-key="msg.content"
              />
              <div
                v-if="msg.msg_type === 2 && !imageUrlCache.get(msg.content)"
                class="w-32 h-32 rounded-xl flex items-center justify-center"
                :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"
              >
                <span class="text-xs text-gray-400">加载中...</span>
              </div>
              <!-- 语音消息 -->
              <div
                v-else-if="msg.msg_type === 3"
                @click="playVoice(msg)"
                class="flex items-center gap-2 cursor-pointer min-w-20"
              >
                <svg 
                  class="w-5 h-5 flex-shrink-0" 
                  :class="[
                    playingVoiceId === msg.id ? 'animate-pulse' : '',
                    msg.sender_id === currentUserOpenId ? 'text-white' : ''
                  ]"
                  fill="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
                </svg>
                <span class="text-sm">{{ msg.duration || 0 }}"</span>
              </div>
              <!-- 表情消息 -->
              <div v-else-if="msg.msg_type === 4" class="text-4xl">
                {{ msg.content }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部输入区域 - 固定 -->
    <div class="chat-input border-t safe-area-bottom" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-end gap-2 px-4 py-3">
        <!-- 语音按钮 -->
        <button
          @click="toggleVoiceMode"
          class="p-2 transition-colors flex-shrink-0"
          :class="[
            voiceMode ? 'text-primary-500' : '',
            appStore.theme === 'dark' ? 'text-gray-400 hover:text-gray-300' : 'text-gray-500 hover:text-gray-700'
          ]"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
          </svg>
        </button>

        <!-- 图片选择 -->
        <button
          @click="selectImage"
          class="p-2 transition-colors flex-shrink-0"
          :class="appStore.theme === 'dark' ? 'text-gray-400 hover:text-gray-300' : 'text-gray-500 hover:text-gray-700'"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </button>
        <input
          ref="imageInputRef"
          type="file"
          accept="image/*"
          class="hidden"
          @change="handleImageSelected"
        />

        <!-- 文字输入 / 语音录制 -->
        <div class="flex-1 relative">
          <!-- 语音录制模式 -->
          <div
            v-if="voiceMode"
            @touchstart.prevent="startRecording"
            @touchend.prevent="stopRecording"
            @mousedown.prevent="startRecording"
            @mouseup.prevent="stopRecording"
            @mouseleave="cancelRecording"
            class="w-full px-4 py-2 text-sm border rounded-2xl text-center select-none transition-colors cursor-pointer"
            :class="[
              recording ? 'bg-primary-100 border-primary-400 text-primary-600' : '',
              appStore.theme === 'dark' ? 'bg-gray-700 border-gray-600 text-gray-300' : 'border-gray-200 text-gray-600'
            ]"
          >
            {{ recording ? `录音中... ${recordingDuration}s` : '按住说话' }}
          </div>
          <!-- 文字输入模式 -->
          <textarea
            v-else
            v-model="inputText"
            ref="textInputRef"
            rows="1"
            placeholder="输入消息..."
            class="w-full px-4 py-2 text-sm border rounded-2xl resize-none focus:outline-none focus:border-primary-400 transition-colors"
            :class="appStore.theme === 'dark' ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-500' : 'border-gray-200'"
            style="max-height: 100px"
            @keydown.enter.exact.prevent="sendMessage"
            @input="autoResize"
          ></textarea>
        </div>

        <!-- 表情按钮 -->
        <button
          @click="toggleEmojiPicker"
          class="p-2 transition-colors flex-shrink-0"
          :class="[
            showEmojiPicker ? 'text-primary-500' : '',
            appStore.theme === 'dark' ? 'text-gray-400 hover:text-gray-300' : 'text-gray-500 hover:text-gray-700'
          ]"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </button>

        <!-- 发送按钮 -->
        <button
          v-if="!voiceMode"
          @click="sendMessage"
          :disabled="!inputText.trim() || sending"
          class="p-2 bg-primary-500 text-white rounded-full disabled:opacity-50 disabled:cursor-not-allowed transition-opacity flex-shrink-0"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
        </button>
      </div>

      <!-- 表情选择器 -->
      <div v-if="showEmojiPicker" class="emoji-picker border-t" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
        <div class="grid grid-cols-8 gap-2 p-3 max-h-48 overflow-y-auto">
          <button
            v-for="emoji in emojiList"
            :key="emoji"
            @click="sendEmoji(emoji)"
            class="w-10 h-10 text-2xl flex items-center justify-center rounded-lg hover:bg-gray-100 active:bg-gray-200 transition-colors"
            :class="appStore.theme === 'dark' ? 'hover:bg-gray-700 active:bg-gray-600' : ''"
          >
            {{ emoji }}
          </button>
        </div>
      </div>
    </div>

    <!-- 图片预览遮罩 -->
    <div
      v-if="previewImageUrl"
      @click="previewImageUrl = ''"
      class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
    >
      <img :src="previewImageUrl" class="max-w-full max-h-full" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import { useAppStore } from '@/stores/app'
import { addMessageListener, removeMessageListener } from '@/utils/websocket'
import { uploadImage, getResourceUrl } from '@/utils/api'

// image URL cache: key -> signedUrl
const imageUrlCache = ref(new Map())

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const messageStore = useMessageStore()
const appStore = useAppStore()

const messageListRef = ref(null)
const textInputRef = ref(null)
const imageInputRef = ref(null)

const peerOpenId = ref('')
const peerNickname = ref('')
const peerAvatar = ref('')
const peerInitial = ref('')
const messages = ref([])
const loading = ref(true)
const loadingMore = ref(false)
const sending = ref(false)
const inputText = ref('')
const page = ref(1)
const total = ref(0)
const hasMore = ref(false)
const previewImageUrl = ref('')

// voice recording
const voiceMode = ref(false)
const recording = ref(false)
const recordingDuration = ref(0)
let mediaRecorder = null
let audioChunks = []
let recordingTimer = null
let recordingStartTime = 0

// voice playback
const playingVoiceId = ref(null)
let currentAudio = null

// emoji picker
const showEmojiPicker = ref(false)
const emojiList = [
  '😀', '😃', '😄', '😁', '😅', '😂', '🤣', '😊',
  '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘',
  '😗', '😙', '😚', '😋', '😛', '😜', '🤪', '😝',
  '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐',
  '😑', '😶', '😏', '😒', '🙄', '😬', '🤥', '😌',
  '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢',
  '🤮', '🤧', '🥵', '🥶', '🥴', '😵', '🤯', '🤠',
  '🥳', '😎', '🤓', '🧐', '😕', '😟', '🙁', '☹️',
  '😮', '😯', '😲', '😳', '🥺', '😦', '😧', '😨',
  '😰', '😥', '😢', '😭', '😱', '😖', '😣', '😞',
  '😓', '😩', '😫', '🥱', '😤', '😡', '😠', '🤬',
  '👍', '👎', '👏', '🙌', '👋', '🤝', '💪', '❤️'
]

// use open_id for comparison
const currentUserOpenId = computed(() => userStore.user?.open_id || '')
const myInitial = computed(() => userStore.user?.nickname?.charAt(0) || 'U')

// messages are fetched DESC, reverse for display (oldest first)
const reversedMessages = computed(() => [...messages.value].reverse())

onMounted(async () => {
  peerOpenId.value = route.params.peerId
  peerNickname.value = route.query.nickname || '用户'
  peerAvatar.value = route.query.avatar || ''
  peerInitial.value = peerNickname.value.charAt(0)

  await fetchMessages()
  scrollToBottom()

  // mark messages as read
  messageStore.markAsRead(peerOpenId.value)

  // listen for new messages from websocket
  addMessageListener(handleNewMessage)
})

onUnmounted(() => {
  removeMessageListener(handleNewMessage)
  messageStore.clearCurrentMessages()
})

function handleNewMessage(msg) {
  // only add if it's from current conversation
  // msg.sender_id and msg.receiver_id are open_ids from backend
  // case 1: peer sends to me -> sender_id = peer, receiver_id = me
  // case 2: I send to peer (won't receive via websocket, already added via API)
  const isFromPeer = msg.sender_id === peerOpenId.value && msg.receiver_id === currentUserOpenId.value
  const isToPeer = msg.sender_id === currentUserOpenId.value && msg.receiver_id === peerOpenId.value
  
  if (isFromPeer || isToPeer) {
    // avoid duplicate: check if message already exists
    const exists = messages.value.some(m => m.id === msg.id)
    if (exists) return
    
    messages.value.unshift(msg)
    // load image/voice URL if needed
    if (msg.msg_type === 2 || msg.msg_type === 3) {
      loadMediaUrls([msg])
    }
    nextTick(() => scrollToBottom())
    // mark as read immediately if from peer
    if (isFromPeer) {
      messageStore.markAsRead(peerOpenId.value)
    }
  }
}

async function fetchMessages() {
  loading.value = true
  try {
    const data = await messageStore.fetchMessages(peerOpenId.value, page.value, 20)
    messages.value = data.list || []
    total.value = data.total
    hasMore.value = messages.value.length < total.value
    // load media URLs for image/voice messages
    await loadMediaUrls(messages.value)
  } finally {
    loading.value = false
  }
}

async function loadMoreMessages() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value++
  try {
    const data = await messageStore.fetchMessages(peerOpenId.value, page.value, 20)
    const newMessages = data.list || []
    messages.value = [...messages.value, ...newMessages]
    hasMore.value = messages.value.length < data.total
    // load media URLs for new messages
    await loadMediaUrls(newMessages)
  } finally {
    loadingMore.value = false
  }
}

async function sendMessage() {
  const text = inputText.value.trim()
  if (!text || sending.value) return

  sending.value = true
  try {
    const msg = await messageStore.sendTextMessage(peerOpenId.value, text)
    messages.value.unshift(msg)
    inputText.value = ''
    nextTick(() => {
      scrollToBottom()
      autoResize()
    })
  } finally {
    sending.value = false
  }
}

function selectImage() {
  imageInputRef.value?.click()
}

async function handleImageSelected(e) {
  const file = e.target.files?.[0]
  if (!file) return

  // validate file type and size
  if (!file.type.startsWith('image/')) {
    appStore.showToast('请选择图片文件', 'error')
    return
  }
  if (file.size > 100 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过100MB', 'error')
    return
  }

  sending.value = true
  try {
    // upload image to server first, get object key
    const imageKey = await uploadImage(file, 'images')
    // send image message with the object key
    const msg = await messageStore.sendImageMessage(peerOpenId.value, imageKey)
    // get signed URL for display
    const signedUrl = await getResourceUrl(imageKey)
    imageUrlCache.value.set(imageKey, signedUrl)
    messages.value.unshift(msg)
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '图片发送失败', 'error')
  } finally {
    sending.value = false
  }

  // clear input
  e.target.value = ''
}

function scrollToBottom() {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

function handleScroll() {
  // could implement pull-to-load-more here
}

function autoResize() {
  const textarea = textInputRef.value
  if (textarea) {
    textarea.style.height = 'auto'
    textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px'
  }
}

function shouldShowTime(msg, index) {
  if (index === 0) return true
  const prev = reversedMessages.value[index - 1]
  if (!prev) return true
  // show time if more than 5 minutes apart
  const prevTime = new Date(prev.created_at).getTime()
  const currTime = new Date(msg.created_at).getTime()
  return currTime - prevTime > 5 * 60 * 1000
}

function formatMessageTime(time) {
  const date = new Date(time)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')

  if (isToday) {
    return `${hour}:${minute}`
  }
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日 ${hour}:${minute}`
}

function previewImage(key) {
  const url = imageUrlCache.value.get(key)
  if (url) {
    previewImageUrl.value = url
  }
}

// load signed URLs for image messages
async function loadMediaUrls(msgList) {
  const mediaMessages = msgList.filter(m => (m.msg_type === 2 || m.msg_type === 3) && !imageUrlCache.value.has(m.content))
  for (const msg of mediaMessages) {
    try {
      const url = await getResourceUrl(msg.content)
      imageUrlCache.value.set(msg.content, url)
    } catch (err) {
      console.error('Failed to load media URL:', err)
    }
  }
}

function onImageLoad() {
  // trigger reactivity update
  imageUrlCache.value = new Map(imageUrlCache.value)
}

function goBack() {
  router.back()
}

// Voice recording functions
function toggleVoiceMode() {
  voiceMode.value = !voiceMode.value
  showEmojiPicker.value = false
}

async function startRecording() {
  if (recording.value) return
  
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder = new MediaRecorder(stream)
    audioChunks = []
    
    mediaRecorder.ondataavailable = (e) => {
      audioChunks.push(e.data)
    }
    
    mediaRecorder.onstop = async () => {
      stream.getTracks().forEach(track => track.stop())
      
      const duration = Math.round((Date.now() - recordingStartTime) / 1000)
      if (duration < 1) {
        appStore.showToast('录音时间太短', 'error')
        return
      }
      
      const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })
      await sendVoiceMessage(audioBlob, duration)
    }
    
    mediaRecorder.start()
    recording.value = true
    recordingStartTime = Date.now()
    recordingDuration.value = 0
    
    recordingTimer = setInterval(() => {
      recordingDuration.value = Math.round((Date.now() - recordingStartTime) / 1000)
      if (recordingDuration.value >= 60) {
        stopRecording()
      }
    }, 100)
  } catch (err) {
    appStore.showToast('无法访问麦克风', 'error')
  }
}

function stopRecording() {
  if (!recording.value || !mediaRecorder) return
  
  clearInterval(recordingTimer)
  recording.value = false
  
  if (mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
  }
}

function cancelRecording() {
  if (!recording.value) return
  
  clearInterval(recordingTimer)
  recording.value = false
  
  if (mediaRecorder && mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
    audioChunks = []
  }
}

async function sendVoiceMessage(audioBlob, duration) {
  sending.value = true
  try {
    // upload voice file
    const file = new File([audioBlob], 'voice.webm', { type: 'audio/webm' })
    const voiceKey = await uploadImage(file, 'voices')
    // send voice message
    const msg = await messageStore.sendVoiceMessage(peerOpenId.value, voiceKey, duration)
    // get signed URL for playback
    const signedUrl = await getResourceUrl(voiceKey)
    imageUrlCache.value.set(voiceKey, signedUrl)
    messages.value.unshift(msg)
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '语音发送失败', 'error')
  } finally {
    sending.value = false
  }
}

// Voice playback
async function playVoice(msg) {
  const url = imageUrlCache.value.get(msg.content)
  if (!url) {
    try {
      const signedUrl = await getResourceUrl(msg.content)
      imageUrlCache.value.set(msg.content, signedUrl)
      playVoiceUrl(msg.id, signedUrl)
    } catch (err) {
      appStore.showToast('语音加载失败', 'error')
    }
    return
  }
  playVoiceUrl(msg.id, url)
}

function playVoiceUrl(msgId, url) {
  // stop current playing
  if (currentAudio) {
    currentAudio.pause()
    currentAudio = null
    if (playingVoiceId.value === msgId) {
      playingVoiceId.value = null
      return
    }
  }
  
  currentAudio = new Audio(url)
  playingVoiceId.value = msgId
  
  currentAudio.onended = () => {
    playingVoiceId.value = null
    currentAudio = null
  }
  
  currentAudio.onerror = () => {
    playingVoiceId.value = null
    currentAudio = null
    appStore.showToast('语音播放失败', 'error')
  }
  
  currentAudio.play()
}

// Emoji functions
function toggleEmojiPicker() {
  showEmojiPicker.value = !showEmojiPicker.value
  voiceMode.value = false
}

async function sendEmoji(emoji) {
  sending.value = true
  try {
    const msg = await messageStore.sendEmojiMessage(peerOpenId.value, emoji)
    messages.value.unshift(msg)
    showEmojiPicker.value = false
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '表情发送失败', 'error')
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.chat-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background-color: #f9fafb;
}

.chat-page--dark {
  background-color: #111827;
}

.chat-header {
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  -webkit-overflow-scrolling: touch;
}

.chat-input {
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.emoji-picker {
  animation: slideUp .2s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
