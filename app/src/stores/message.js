import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useMessageStore = defineStore('message', () => {
  const conversations = ref([])
  const currentMessages = ref([])
  const unreadCount = ref(0)

  // fetch all conversations
  async function fetchConversations() {
    const data = await api.get('/conversations')
    conversations.value = data.list || []
    return data
  }

  // fetch messages with a specific user (peerOpenId is the open_id string)
  async function fetchMessages(peerOpenId, page = 1, pageSize = 20) {
    const data = await api.get('/messages', {
      params: { peer_id: peerOpenId, page, page_size: pageSize }
    })
    return data
  }

  // send a text message (receiverOpenId is the open_id string)
  async function sendTextMessage(receiverOpenId, content) {
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: content,
      msg_type: 1
    })
    return msg
  }

  // send an image message (receiverOpenId is the open_id string)
  async function sendImageMessage(receiverOpenId, imageUrl) {
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: imageUrl,
      msg_type: 2
    })
    return msg
  }

  // send a voice message (receiverOpenId is the open_id string)
  async function sendVoiceMessage(receiverOpenId, voiceKey, duration) {
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: voiceKey,
      msg_type: 3,
      duration: duration
    })
    return msg
  }

  // send an emoji message (receiverOpenId is the open_id string)
  async function sendEmojiMessage(receiverOpenId, emojiCode) {
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: emojiCode,
      msg_type: 4
    })
    return msg
  }

  // send a call record message (receiverOpenId is the open_id string)
  // callType: 'audio' | 'video', duration: seconds, status: 'completed' | 'missed' | 'rejected' | 'cancelled'
  async function sendCallRecordMessage(receiverOpenId, callType, duration, status) {
    const content = JSON.stringify({
      call_type: callType,
      duration: duration,
      status: status
    })
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: content,
      msg_type: 5
    })
    return msg
  }

  // send a video message (receiverOpenId is the open_id string)
  // videoInfo: { key, thumbnail, duration, width, height }
  async function sendVideoMessage(receiverOpenId, videoInfo) {
    const content = JSON.stringify({
      key: videoInfo.key,
      thumbnail: videoInfo.thumbnail,
      duration: videoInfo.duration,
      width: videoInfo.width || 0,
      height: videoInfo.height || 0
    })
    const msg = await api.post('/messages', {
      receiver_id: receiverOpenId,
      content: content,
      msg_type: 6,
      duration: videoInfo.duration
    })
    return msg
  }

  // mark messages from a peer as read (peerOpenId is the open_id string)
  async function markAsRead(peerOpenId) {
    await api.put('/messages/read', null, {
      params: { peer_id: peerOpenId }
    })
  }

  // get total unread message count
  async function fetchUnreadCount() {
    const data = await api.get('/messages/unread-count')
    unreadCount.value = data.count || 0
    return data.count
  }

  // add a new message to current conversation (called from websocket)
  function addMessage(message) {
    currentMessages.value.unshift(message)
  }

  // clear current messages when leaving chat page
  function clearCurrentMessages() {
    currentMessages.value = []
  }

  // increment unread count by 1
  function incrementUnreadCount() {
    unreadCount.value++
  }

  return {
    conversations,
    currentMessages,
    unreadCount,
    fetchConversations,
    fetchMessages,
    sendTextMessage,
    sendImageMessage,
    sendVoiceMessage,
    sendEmojiMessage,
    sendCallRecordMessage,
    sendVideoMessage,
    markAsRead,
    fetchUnreadCount,
    addMessage,
    clearCurrentMessages,
    incrementUnreadCount
  }
})
