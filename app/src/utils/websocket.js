let ws = null
let reconnectTimer = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5
let currentToken = null

// message event listeners for chat functionality
const messageListeners = new Set()

// call signaling listeners
const callListeners = new Set()

export function addMessageListener(callback) {
  messageListeners.add(callback)
}

export function removeMessageListener(callback) {
  messageListeners.delete(callback)
}

// 添加通话信令监听器
export function addCallListener(callback) {
  callListeners.add(callback)
}

// 移除通话信令监听器
export function removeCallListener(callback) {
  callListeners.delete(callback)
}

export function connectWebSocket(token) {
  // save token for reconnect
  if (token) {
    currentToken = token
  }
  
  if (!currentToken) {
    console.log('WebSocket: No token, skip connect')
    return
  }

  // close existing connection
  if (ws) {
    ws.close()
    ws = null
  }

  // use relative path, vite will proxy to backend
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const url = `${protocol}//${host}/ws?token=${currentToken}`

  console.log('WebSocket: Connecting to', url)

  try {
    ws = new WebSocket(url)
  } catch (e) {
    console.error('WebSocket: Failed to create connection', e)
    return
  }

  ws.onopen = () => {
    console.log('WebSocket: Connected successfully')
    reconnectAttempts = 0
  }

  ws.onmessage = (event) => {
    console.log('WebSocket: Received message', event.data)
    try {
      const message = JSON.parse(event.data)
      handleMessage(message)
    } catch (e) {
      console.error('WebSocket: Message parse error:', e)
    }
  }

  ws.onclose = (event) => {
    console.log('WebSocket: Disconnected', event.code, event.reason)
    ws = null
    attemptReconnect()
  }

  ws.onerror = (error) => {
    console.error('WebSocket: Error:', error)
  }
}

function attemptReconnect() {
  if (reconnectAttempts >= maxReconnectAttempts) {
    console.log('Max reconnect attempts reached')
    return
  }

  reconnectTimer = setTimeout(() => {
    reconnectAttempts++
    console.log(`Reconnecting... (${reconnectAttempts}/${maxReconnectAttempts})`)
    connectWebSocket()
  }, 3000 * reconnectAttempts)
}

// callback for handling app store actions
let onMessageCallback = null

export function setOnMessageCallback(callback) {
  onMessageCallback = callback
}

function handleMessage(message) {
  switch (message.type) {
    case 'match_found':
    case 'match_success':
    case 'match_rejected':
    case 'trip_grabbed':
      if (onMessageCallback) {
        onMessageCallback(message)
      }
      break
    case 'new_message':
      // notify all listeners about new chat message
      messageListeners.forEach(callback => callback(message.data))
      if (onMessageCallback) {
        onMessageCallback(message)
      }
      break
    // 通话信令消息
    case 'call_invite':
    case 'call_answer':
    case 'call_end':
    case 'webrtc_offer':
    case 'webrtc_answer':
    case 'ice_candidate':
      console.log('WebSocket: Call signaling received:', message.type)
      callListeners.forEach(callback => callback(message))
      break
    default:
      console.log('Unknown message type:', message.type)
  }
}

export function disconnectWebSocket() {
  currentToken = null
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (ws) {
    ws.close()
    ws = null
  }
}

/**
 * 发送通话信令消息
 * @param {string} type - 信令类型
 * @param {object} data - 信令数据
 */
export function sendCallSignaling(type, data) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.error('WebSocket: Cannot send signaling, connection not ready')
    return false
  }

  const message = JSON.stringify({
    type,
    data
  })

  console.log('WebSocket: Sending call signaling:', type)
  ws.send(message)
  return true
}

/**
 * 发起通话邀请
 */
export function sendCallInvite(targetOpenId, callId, callType, callerInfo) {
  return sendCallSignaling('call_invite', {
    target_open_id: targetOpenId,
    call_id: callId,
    call_type: callType,
    caller_info: callerInfo
  })
}

/**
 * 发送通话应答
 */
export function sendCallAnswer(targetOpenId, callId, accept) {
  return sendCallSignaling('call_answer', {
    target_open_id: targetOpenId,
    call_id: callId,
    accept
  })
}

/**
 * 发送通话结束信令
 */
export function sendCallEnd(targetOpenId, callId, reason = 'normal') {
  return sendCallSignaling('call_end', {
    target_open_id: targetOpenId,
    call_id: callId,
    reason
  })
}

/**
 * 发送 WebRTC Offer
 */
export function sendWebRTCOffer(targetOpenId, callId, sdp) {
  return sendCallSignaling('webrtc_offer', {
    target_open_id: targetOpenId,
    call_id: callId,
    sdp
  })
}

/**
 * 发送 WebRTC Answer
 */
export function sendWebRTCAnswer(targetOpenId, callId, sdp) {
  return sendCallSignaling('webrtc_answer', {
    target_open_id: targetOpenId,
    call_id: callId,
    sdp
  })
}

/**
 * 发送 ICE Candidate
 */
export function sendIceCandidate(targetOpenId, callId, candidate) {
  return sendCallSignaling('ice_candidate', {
    target_open_id: targetOpenId,
    call_id: callId,
    candidate
  })
}
