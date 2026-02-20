import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

/**
 * 通话状态管理
 * 
 * 通话状态流转:
 * idle -> calling (发起方) / incoming (接收方)
 * calling/incoming -> connecting (接受通话)
 * connecting -> connected (WebRTC 连接成功)
 * connected -> ended (通话结束)
 * calling/incoming -> ended (拒绝/取消)
 */
export const useCallStore = defineStore('call', () => {
  // 通话状态: idle | calling | incoming | connecting | connected | ended
  const status = ref('idle')
  
  // 通话ID
  const callId = ref(null)
  
  // 通话类型: audio | video
  const callType = ref('audio')
  
  // 是否为发起方
  const isCaller = ref(false)
  
  // 对方用户信息
  const peerInfo = ref({
    userId: null,
    nickname: '',
    avatar: ''
  })
  
  // 通话开始时间 (用于计算通话时长)
  const callStartTime = ref(null)
  
  // 通话时长 (秒)
  const duration = ref(0)
  
  // 通话时长定时器
  let durationTimer = null
  
  // 结束原因
  const endReason = ref('')

  // 计算属性：是否在通话中
  const isInCall = computed(() => {
    return ['calling', 'incoming', 'connecting', 'connected'].includes(status.value)
  })

  // 计算属性：是否有来电
  const hasIncomingCall = computed(() => {
    return status.value === 'incoming'
  })

  // 计算属性：格式化通话时长
  const formattedDuration = computed(() => {
    const mins = Math.floor(duration.value / 60)
    const secs = duration.value % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  })

  /**
   * 生成唯一通话ID
   */
  function generateCallId() {
    return `call_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  /**
   * 发起通话
   */
  function initiateCall(targetUserId, targetNickname, targetAvatar, type = 'audio') {
    callId.value = generateCallId()
    callType.value = type
    isCaller.value = true
    status.value = 'calling'
    peerInfo.value = {
      userId: targetUserId,
      nickname: targetNickname,
      avatar: targetAvatar
    }
    endReason.value = ''
    
    console.log('CallStore: Initiating call', { callId: callId.value, type, targetUserId })
    
    return callId.value
  }

  /**
   * 收到来电
   */
  function receiveCall(incomingCallId, callerUserId, callerNickname, callerAvatar, type = 'audio') {
    // 如果已在通话中，忽略来电
    if (isInCall.value) {
      console.log('CallStore: Already in call, ignoring incoming call')
      return false
    }

    callId.value = incomingCallId
    callType.value = type
    isCaller.value = false
    status.value = 'incoming'
    peerInfo.value = {
      userId: callerUserId,
      nickname: callerNickname,
      avatar: callerAvatar
    }
    endReason.value = ''
    
    console.log('CallStore: Received incoming call', { callId: callId.value, type, callerUserId })
    
    return true
  }

  /**
   * 接受通话
   */
  function acceptCall() {
    if (status.value !== 'incoming') {
      console.warn('CallStore: Cannot accept call, not in incoming state')
      return false
    }
    
    status.value = 'connecting'
    console.log('CallStore: Call accepted')
    return true
  }

  /**
   * 对方接受通话 (发起方收到)
   */
  function peerAccepted() {
    if (status.value !== 'calling') {
      console.warn('CallStore: Cannot set peer accepted, not in calling state')
      return false
    }
    
    status.value = 'connecting'
    console.log('CallStore: Peer accepted call')
    return true
  }

  /**
   * 通话连接成功
   */
  function setConnected() {
    if (status.value !== 'connecting') {
      console.warn('CallStore: Cannot set connected, not in connecting state')
      return
    }
    
    status.value = 'connected'
    callStartTime.value = Date.now()
    duration.value = 0
    
    // 开始计时
    durationTimer = setInterval(() => {
      duration.value++
    }, 1000)
    
    console.log('CallStore: Call connected')
  }

  /**
   * 结束通话
   */
  function endCall(reason = 'normal') {
    console.log('CallStore: Ending call', { reason, currentStatus: status.value })
    
    // 停止计时
    if (durationTimer) {
      clearInterval(durationTimer)
      durationTimer = null
    }
    
    endReason.value = reason
    status.value = 'ended'
    
    // 延迟重置状态
    setTimeout(() => {
      reset()
    }, 2000)
  }

  /**
   * 重置状态
   */
  function reset() {
    status.value = 'idle'
    callId.value = null
    callType.value = 'audio'
    isCaller.value = false
    peerInfo.value = {
      userId: null,
      nickname: '',
      avatar: ''
    }
    callStartTime.value = null
    duration.value = 0
    endReason.value = ''
    
    if (durationTimer) {
      clearInterval(durationTimer)
      durationTimer = null
    }
    
    console.log('CallStore: Reset')
  }

  /**
   * 获取结束原因文本
   */
  function getEndReasonText() {
    const reasons = {
      normal: '通话结束',
      cancelled: '已取消',
      rejected: '对方已拒绝',
      busy: '对方忙',
      timeout: '无人接听',
      network_error: '网络异常',
      peer_ended: '对方已挂断'
    }
    return reasons[endReason.value] || '通话结束'
  }

  return {
    // 状态
    status,
    callId,
    callType,
    isCaller,
    peerInfo,
    callStartTime,
    duration,
    endReason,
    
    // 计算属性
    isInCall,
    hasIncomingCall,
    formattedDuration,
    
    // 方法
    generateCallId,
    initiateCall,
    receiveCall,
    acceptCall,
    peerAccepted,
    setConnected,
    endCall,
    reset,
    getEndReasonText
  }
})
