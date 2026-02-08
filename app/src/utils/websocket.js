let ws = null
let reconnectTimer = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5
let currentToken = null

// message event listeners for chat functionality
const messageListeners = new Set()

export function addMessageListener(callback) {
  messageListeners.add(callback)
}

export function removeMessageListener(callback) {
  messageListeners.delete(callback)
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
