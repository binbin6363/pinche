<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="callStore.hasIncomingCall" class="incoming-call-overlay">
        <div class="incoming-call-card">
          <!-- 头像区域 -->
          <div class="avatar-section">
            <div class="avatar-wrapper">
              <div class="avatar">
                <img v-if="callStore.peerInfo.avatar" :src="callStore.peerInfo.avatar" />
                <span v-else class="avatar-text">{{ peerInitial }}</span>
              </div>
              <div class="pulse-ring"></div>
              <div class="pulse-ring delay-1"></div>
            </div>
          </div>

          <!-- 信息区域 -->
          <div class="info-section">
            <h2 class="caller-name">{{ callStore.peerInfo.nickname }}</h2>
            <p class="call-type">
              <svg v-if="callStore.callType === 'video'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
              {{ callStore.callType === 'video' ? '视频通话' : '语音通话' }}
            </p>
          </div>

          <!-- 操作按钮 -->
          <div class="action-section">
            <button class="action-btn reject-btn" @click="rejectCall">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M5 3a2 2 0 00-2 2v1c0 8.284 6.716 15 15 15h1a2 2 0 002-2v-3.28a1 1 0 00-.684-.948l-4.493-1.498a1 1 0 00-1.21.502l-1.13 2.257a11.042 11.042 0 01-5.516-5.517l2.257-1.128a1 1 0 00.502-1.21L9.228 3.683A1 1 0 008.279 3H5z" />
              </svg>
              <span>拒绝</span>
            </button>

            <button class="action-btn accept-btn" @click="acceptCall">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
              <span>接听</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCallStore } from '@/stores/call'
import { sendCallAnswer, sendCallEnd } from '@/utils/websocket'

const router = useRouter()
const callStore = useCallStore()

const peerInitial = computed(() => {
  return callStore.peerInfo.nickname?.charAt(0) || '?'
})

// 接听通话
function acceptCall() {
  // 发送接听信令
  sendCallAnswer(callStore.peerInfo.userId, callStore.callId, true)

  // 更新状态
  callStore.acceptCall()

  // 跳转到通话页面
  router.push('/call')
}

// 拒绝通话
function rejectCall() {
  // 发送拒绝信令
  sendCallAnswer(callStore.peerInfo.userId, callStore.callId, false)

  // 结束通话
  callStore.endCall('rejected')
}
</script>

<style scoped>
.incoming-call-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, .85);
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.incoming-call-card {
  width: 90%;
  max-width: 320px;
  padding: 40px 24px;
  background: linear-gradient(135deg, rgba(255, 255, 255, .1) 0%, rgba(255, 255, 255, .05) 100%);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, .1);
  text-align: center;
}

/* 头像区域 */
.avatar-section {
  margin-bottom: 24px;
}

.avatar-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  font-size: 40px;
  font-weight: 600;
  color: #fff;
}

.pulse-ring {
  position: absolute;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 2px solid rgba(102, 126, 234, .6);
  animation: pulse 2s ease-out infinite;
}

.pulse-ring.delay-1 {
  animation-delay: 1s;
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(1.8);
    opacity: 0;
  }
}

/* 信息区域 */
.info-section {
  margin-bottom: 32px;
}

.caller-name {
  font-size: 24px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px;
}

.call-type {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: rgba(255, 255, 255, .7);
  margin: 0;
}

/* 操作按钮 */
.action-section {
  display: flex;
  justify-content: center;
  gap: 48px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  color: #fff;
}

.action-btn > svg {
  width: 64px;
  height: 64px;
  padding: 16px;
  border-radius: 50%;
}

.action-btn > span {
  font-size: 14px;
  color: rgba(255, 255, 255, .8);
}

.reject-btn > svg {
  background: #ef4444;
}

.reject-btn:active > svg {
  background: #dc2626;
  transform: scale(.95);
}

.accept-btn > svg {
  background: #10b981;
}

.accept-btn:active > svg {
  background: #059669;
  transform: scale(.95);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity .3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .incoming-call-card {
  animation: slideUp .3s ease;
}

.fade-leave-active .incoming-call-card {
  animation: slideDown .3s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideDown {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(20px);
  }
}
</style>
