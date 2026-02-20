<template>
  <div class="video-call-page">
    <!-- 视频区域 -->
    <div class="video-container">
      <!-- 远端视频/头像 (大画面) -->
      <div class="remote-video-wrapper">
        <video
          v-if="callStore.callType === 'video' && remoteStream"
          ref="remoteVideoRef"
          class="remote-video"
          autoplay
          playsinline
        ></video>
        <div v-else class="remote-avatar-wrapper">
          <div class="remote-avatar">
            <img v-if="callStore.peerInfo.avatar" :src="callStore.peerInfo.avatar" />
            <span v-else class="avatar-text">{{ peerInitial }}</span>
          </div>
          <div class="pulse-ring"></div>
          <div class="pulse-ring delay-1"></div>
          <div class="pulse-ring delay-2"></div>
        </div>
      </div>

      <!-- 本地预览 (小画面) -->
      <div
        v-if="callStore.callType === 'video' && localStream"
        class="local-video-wrapper"
        :style="localVideoStyle"
        @touchstart="startDrag"
        @touchmove="onDrag"
        @touchend="endDrag"
      >
        <video ref="localVideoRef" class="local-video" autoplay playsinline muted></video>
        <div v-if="!isVideoEnabled" class="video-off-overlay">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- 顶部状态栏 -->
    <div class="top-bar safe-area-top">
      <div class="peer-info">
        <h2 class="peer-name">{{ callStore.peerInfo.nickname }}</h2>
        <p class="call-status">{{ statusText }}</p>
      </div>
      <button v-if="callStore.callType === 'video'" class="switch-camera-btn" @click="switchCamera">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>

    <!-- 底部控制栏 -->
    <div class="bottom-bar safe-area-bottom">
      <div class="control-buttons">
        <!-- 静音按钮 -->
        <button class="control-btn" :class="{ active: isMuted }" @click="toggleMute">
          <svg v-if="!isMuted" class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
          </svg>
          <svg v-else class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2" />
          </svg>
          <span class="btn-label">{{ isMuted ? '取消静音' : '静音' }}</span>
        </button>

        <!-- 挂断按钮 -->
        <button class="hangup-btn" @click="hangup">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M5 3a2 2 0 00-2 2v1c0 8.284 6.716 15 15 15h1a2 2 0 002-2v-3.28a1 1 0 00-.684-.948l-4.493-1.498a1 1 0 00-1.21.502l-1.13 2.257a11.042 11.042 0 01-5.516-5.517l2.257-1.128a1 1 0 00.502-1.21L9.228 3.683A1 1 0 008.279 3H5z" />
          </svg>
        </button>

        <!-- 视频开关 (仅视频通话) -->
        <button v-if="callStore.callType === 'video'" class="control-btn" :class="{ active: !isVideoEnabled }" @click="toggleVideo">
          <svg v-if="isVideoEnabled" class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
          </svg>
          <span class="btn-label">{{ isVideoEnabled ? '关闭视频' : '开启视频' }}</span>
        </button>

        <!-- 扬声器切换 (仅语音通话) -->
        <button v-else class="control-btn" :class="{ active: isSpeakerOn }" @click="toggleSpeaker">
          <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
          </svg>
          <span class="btn-label">扬声器</span>
        </button>
      </div>
    </div>

    <!-- 通话结束提示 -->
    <div v-if="callStore.status === 'ended'" class="end-overlay">
      <div class="end-content">
        <p class="end-text">{{ callStore.getEndReasonText() }}</p>
        <p v-if="callStore.duration > 0" class="end-duration">通话时长 {{ callStore.formattedDuration }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCallStore } from '@/stores/call'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import { WebRTCManager } from '@/utils/webrtc'
import {
  addCallListener,
  removeCallListener,
  sendWebRTCOffer,
  sendWebRTCAnswer,
  sendIceCandidate,
  sendCallEnd
} from '@/utils/websocket'

const router = useRouter()
const callStore = useCallStore()
const userStore = useUserStore()
const messageStore = useMessageStore()

// refs
const localVideoRef = ref(null)
const remoteVideoRef = ref(null)

// WebRTC
let webrtc = null
const localStream = ref(null)
const remoteStream = ref(null)

// UI state
const isMuted = ref(false)
const isVideoEnabled = ref(true)
const isSpeakerOn = ref(true)

// 本地视频位置
const localVideoPos = ref({ x: 16, y: 100 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

const localVideoStyle = computed(() => ({
  right: `${localVideoPos.value.x}px`,
  top: `${localVideoPos.value.y}px`
}))

// 对方昵称首字母
const peerInitial = computed(() => {
  return callStore.peerInfo.nickname?.charAt(0) || '?'
})

// 状态文本
const statusText = computed(() => {
  switch (callStore.status) {
    case 'calling':
      return '正在呼叫...'
    case 'incoming':
      return '来电中...'
    case 'connecting':
      return '连接中...'
    case 'connected':
      return callStore.formattedDuration
    case 'ended':
      return callStore.getEndReasonText()
    default:
      return ''
  }
})

// 初始化 WebRTC
async function initWebRTC() {
  webrtc = new WebRTCManager({
    onLocalStream: (stream) => {
      localStream.value = stream
      if (localVideoRef.value) {
        localVideoRef.value.srcObject = stream
      }
    },
    onRemoteStream: (stream) => {
      console.log('VideoCall: Received remote stream with', stream.getTracks().length, 'tracks')
      stream.getTracks().forEach(track => {
        console.log('VideoCall: Remote track:', track.kind, 'enabled:', track.enabled, 'readyState:', track.readyState)
      })
      remoteStream.value = stream
      if (remoteVideoRef.value) {
        console.log('VideoCall: Setting remote video srcObject in callback')
        remoteVideoRef.value.srcObject = stream
        // 移动端需要显式调用 play()
        remoteVideoRef.value.play().catch(err => {
          console.warn('VideoCall: Auto-play failed:', err)
        })
      }
    },
    onIceCandidate: (candidate) => {
      sendIceCandidate(callStore.peerInfo.userId, callStore.callId, candidate)
    },
    onConnectionStateChange: (state) => {
      console.log('VideoCall: Connection state:', state)
      if (state === 'connected') {
        callStore.setConnected()
      } else if (state === 'failed' || state === 'disconnected') {
        handleCallEnd('network_error')
      }
    },
    onError: (code, message) => {
      console.error('VideoCall: WebRTC error:', code, message)
    }
  })

  // 初始化 PeerConnection
  webrtc.initPeerConnection()

  // 获取本地媒体流
  try {
    console.log('VideoCall: Getting local stream...')
    await webrtc.getLocalStream(callStore.callType)
    console.log('VideoCall: Adding local stream to PC...')
    webrtc.addLocalStreamToPC()
    console.log('VideoCall: Local stream added, senders:', webrtc.pc.getSenders().length)
  } catch (err) {
    console.error('VideoCall: Failed to get media stream:', err)
    handleCallEnd('media_error')
    return
  }

  // 如果是发起方，创建 offer
  if (callStore.isCaller && callStore.status === 'connecting') {
    try {
      const offer = await webrtc.createOffer()
      sendWebRTCOffer(callStore.peerInfo.userId, callStore.callId, offer)
    } catch (err) {
      console.error('VideoCall: Failed to create offer:', err)
    }
  }
}

// 处理信令消息
function handleSignaling(message) {
  const { type, data } = message

  // 验证 callId
  if (data.call_id !== callStore.callId) {
    console.log('VideoCall: Ignoring signaling for different call')
    return
  }

  switch (type) {
    case 'call_answer':
      handleCallAnswer(data.accept)
      break
    case 'webrtc_offer':
      handleOffer(data.sdp)
      break
    case 'webrtc_answer':
      handleAnswer(data.sdp)
      break
    case 'ice_candidate':
      handleIceCandidate(data.candidate)
      break
    case 'call_end':
      handleCallEnd(data.reason || 'peer_ended')
      break
  }
}

// 处理对方接听/拒绝 (发起方收到)
async function handleCallAnswer(accept) {
  if (!callStore.isCaller) return

  console.log('VideoCall: handleCallAnswer', { accept, webrtcReady: !!webrtc })

  if (accept) {
    // 对方接受，更新状态
    callStore.peerAccepted()
    
    // 等待 WebRTC 初始化完成后再创建 offer
    // 使用延迟确保接收方也已初始化
    await waitForWebRTC()
    
    try {
      console.log('VideoCall: Creating offer after peer accepted')
      const offer = await webrtc.createOffer()
      sendWebRTCOffer(callStore.peerInfo.userId, callStore.callId, offer)
    } catch (err) {
      console.error('VideoCall: Failed to create offer:', err)
    }
  } else {
    // 对方拒绝
    handleCallEnd('rejected')
  }
}

// 等待 WebRTC 初始化完成（包括本地流已添加到 PC）
function waitForWebRTC() {
  return new Promise((resolve) => {
    const isReady = () => {
      if (!webrtc || !webrtc.pc) return false
      // 确保本地流已添加到 PeerConnection
      const senders = webrtc.pc.getSenders()
      return senders.length > 0
    }
    
    if (isReady()) {
      console.log('VideoCall: WebRTC already ready')
      resolve()
    } else {
      const check = setInterval(() => {
        if (isReady()) {
          console.log('VideoCall: WebRTC now ready')
          clearInterval(check)
          resolve()
        }
      }, 50)
      // 超时 5 秒
      setTimeout(() => {
        clearInterval(check)
        console.warn('VideoCall: WebRTC wait timeout, proceeding anyway')
        resolve()
      }, 5000)
    }
  })
}

// 处理 offer (接收方)
async function handleOffer(sdp) {
  console.log('VideoCall: handleOffer called', { 
    webrtcReady: !!webrtc, 
    pcReady: !!webrtc?.pc,
    sendersCount: webrtc?.pc?.getSenders()?.length || 0
  })
  
  // 等待 WebRTC 初始化完成（确保本地流已添加）
  await waitForWebRTC()
  
  if (!webrtc) {
    console.error('VideoCall: WebRTC not ready for offer')
    return
  }

  try {
    console.log('VideoCall: Setting remote description (offer)')
    await webrtc.setRemoteDescription(sdp)
    
    console.log('VideoCall: Creating answer, senders:', webrtc.pc.getSenders().length)
    const answer = await webrtc.createAnswer()
    
    console.log('VideoCall: Sending answer')
    sendWebRTCAnswer(callStore.peerInfo.userId, callStore.callId, answer)
  } catch (err) {
    console.error('VideoCall: Failed to handle offer:', err)
  }
}

// 处理 answer (发起方)
async function handleAnswer(sdp) {
  console.log('VideoCall: handleAnswer', { webrtcReady: !!webrtc })
  
  await waitForWebRTC()
  
  if (!webrtc) {
    console.error('VideoCall: WebRTC not ready for answer')
    return
  }

  try {
    await webrtc.setRemoteDescription(sdp)
  } catch (err) {
    console.error('VideoCall: Failed to handle answer:', err)
  }
}

// 处理 ICE candidate
async function handleIceCandidate(candidate) {
  await waitForWebRTC()
  
  if (!webrtc) {
    console.warn('VideoCall: WebRTC not ready for ICE candidate')
    return
  }

  try {
    await webrtc.addIceCandidate(candidate)
  } catch (err) {
    console.error('VideoCall: Failed to add ICE candidate:', err)
  }
}

// 处理通话结束
async function handleCallEnd(reason) {
  const wasConnected = callStore.status === 'connected'
  const callDuration = callStore.duration
  const peerId = callStore.peerInfo.userId
  const callType = callStore.callType
  const wasCaller = callStore.isCaller
  
  callStore.endCall(reason)
  cleanup()

  // 发送通话记录消息（仅发起方发送，避免重复）
  if (peerId && wasCaller) {
    try {
      let status = 'completed'
      if (reason === 'rejected') {
        status = 'rejected'
      } else if (reason === 'timeout' || reason === 'no_answer') {
        status = 'missed'
      } else if (reason === 'cancelled' || (reason === 'normal' && !wasConnected)) {
        status = 'cancelled'
      }
      
      await messageStore.sendCallRecordMessage(
        peerId,
        callType,
        wasConnected ? callDuration : 0,
        status
      )
      console.log('VideoCall: Call record sent', { callType, callDuration, status })
    } catch (err) {
      console.error('VideoCall: Failed to send call record:', err)
    }
  }

  // 延迟返回
  setTimeout(() => {
    router.back()
  }, 2000)
}

// 挂断
function hangup() {
  sendCallEnd(callStore.peerInfo.userId, callStore.callId, 'normal')
  handleCallEnd('normal')
}

// 切换静音
function toggleMute() {
  if (webrtc) {
    isMuted.value = webrtc.toggleMute()
  }
}

// 切换视频
function toggleVideo() {
  if (webrtc) {
    isVideoEnabled.value = webrtc.toggleVideo()
  }
}

// 切换摄像头
async function switchCamera() {
  if (webrtc) {
    await webrtc.switchCamera()
  }
}

// 切换扬声器
function toggleSpeaker() {
  isSpeakerOn.value = !isSpeakerOn.value
  // 实际切换扬声器需要使用 AudioContext 或原生 API
}

// 拖动本地视频
function startDrag(e) {
  isDragging.value = true
  const touch = e.touches[0]
  dragStart.value = {
    x: touch.clientX + localVideoPos.value.x,
    y: touch.clientY - localVideoPos.value.y
  }
}

function onDrag(e) {
  if (!isDragging.value) return
  const touch = e.touches[0]
  localVideoPos.value = {
    x: Math.max(0, dragStart.value.x - touch.clientX),
    y: Math.max(50, touch.clientY - dragStart.value.y + localVideoPos.value.y)
  }
}

function endDrag() {
  isDragging.value = false
}

// 清理资源
function cleanup() {
  if (webrtc) {
    webrtc.close()
    webrtc = null
  }
  localStream.value = null
  remoteStream.value = null
}

// 监听视频元素
watch(localVideoRef, (el) => {
  if (el && localStream.value) {
    el.srcObject = localStream.value
  }
})

watch(remoteVideoRef, (el) => {
  if (el && remoteStream.value) {
    console.log('VideoCall: Setting remote video srcObject (ref changed)')
    el.srcObject = remoteStream.value
    // 确保在移动端可以播放
    el.play().catch(err => {
      console.warn('VideoCall: Auto-play failed for remote video:', err)
    })
  }
})

// 监听 remoteStream 变化，确保 video 元素更新
watch(remoteStream, (stream) => {
  if (stream && remoteVideoRef.value) {
    console.log('VideoCall: Setting remote video srcObject (stream changed)')
    remoteVideoRef.value.srcObject = stream
    // 确保在移动端可以播放
    remoteVideoRef.value.play().catch(err => {
      console.warn('VideoCall: Auto-play failed for remote video:', err)
    })
  }
})

onMounted(() => {
  // 验证通话状态
  if (!callStore.isInCall) {
    console.error('VideoCall: Not in call, redirecting back')
    router.back()
    return
  }

  // 添加信令监听
  addCallListener(handleSignaling)

  // 初始化 WebRTC
  initWebRTC()
})

onUnmounted(() => {
  removeCallListener(handleSignaling)
  cleanup()
})
</script>

<style scoped>
.video-call-page {
  position: fixed;
  inset: 0;
  background: #000;
  z-index: 9999;
}

.video-container {
  position: absolute;
  inset: 0;
}

/* 远端视频/头像 */
.remote-video-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
}

.remote-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remote-avatar-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remote-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  z-index: 1;
}

.remote-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  font-size: 48px;
  font-weight: 600;
  color: #fff;
}

/* 脉冲动画 */
.pulse-ring {
  position: absolute;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 2px solid rgba(102, 126, 234, .5);
  animation: pulse 2s ease-out infinite;
}

.pulse-ring.delay-1 {
  animation-delay: .5s;
}

.pulse-ring.delay-2 {
  animation-delay: 1s;
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(2);
    opacity: 0;
  }
}

/* 本地视频预览 */
.local-video-wrapper {
  position: absolute;
  width: 100px;
  height: 140px;
  border-radius: 12px;
  overflow: hidden;
  background: #000;
  box-shadow: 0 4px 20px rgba(0, 0, 0, .5);
  border: 2px solid rgba(255, 255, 255, .2);
  z-index: 10;
}

.local-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: scaleX(-1);
}

.video-off-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, .8);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

/* 顶部状态栏 */
.top-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 16px 20px;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  background: linear-gradient(to bottom, rgba(0, 0, 0, .6), transparent);
  z-index: 20;
}

.peer-info {
  text-align: left;
}

.peer-name {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.call-status {
  font-size: 14px;
  color: rgba(255, 255, 255, .7);
  margin: 4px 0 0;
}

.switch-camera-btn {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: rgba(255, 255, 255, .2);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border: none;
}

.switch-camera-btn:active {
  background: rgba(255, 255, 255, .3);
}

/* 底部控制栏 */
.bottom-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24px 20px 32px;
  background: linear-gradient(to top, rgba(0, 0, 0, .8), transparent);
  z-index: 20;
}

.control-buttons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32px;
}

.control-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  color: #fff;
}

.control-btn > svg {
  width: 56px;
  height: 56px;
  padding: 14px;
  background: rgba(255, 255, 255, .2);
  backdrop-filter: blur(10px);
  border-radius: 50%;
}

.control-btn.active > svg {
  background: rgba(255, 255, 255, .9);
  color: #333;
}

.control-btn:active > svg {
  transform: scale(.95);
}

.btn-label {
  font-size: 12px;
  color: rgba(255, 255, 255, .8);
}

.hangup-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #ef4444;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border: none;
  box-shadow: 0 4px 20px rgba(239, 68, 68, .5);
}

.hangup-btn:active {
  transform: scale(.95);
  background: #dc2626;
}

/* 通话结束遮罩 */
.end-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, .9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 30;
}

.end-content {
  text-align: center;
}

.end-text {
  font-size: 20px;
  color: #fff;
  margin: 0;
}

.end-duration {
  font-size: 14px;
  color: rgba(255, 255, 255, .6);
  margin-top: 8px;
}

/* 安全区域 */
.safe-area-top {
  padding-top: max(16px, env(safe-area-inset-top));
}

.safe-area-bottom {
  padding-bottom: max(32px, env(safe-area-inset-bottom));
}
</style>
