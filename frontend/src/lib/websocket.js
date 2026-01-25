function getWsUrl() {
  if (import.meta.env.VITE_WS_URL) return import.meta.env.VITE_WS_URL
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${location.host}/api/ws`
}
const WS_URL = getWsUrl()

const RECONNECT_DELAYS = [1000, 2000, 4000, 8000, 16000, 30000]

class WebSocketManager {
  constructor() {
    this.ws = null
    this.token = null
    this.listeners = new Set()
    this.reconnectAttempt = 0
    this.reconnectTimeout = null
    this.isConnecting = false
  }

  connect(token) {
    console.log('[WSManager] connect() called, token:', token?.substring(0, 20), 'current readyState:', this.ws?.readyState)
    
    if (this.token === token && this.ws?.readyState === WebSocket.OPEN) {
      console.log('[WSManager] Already connected with same token, skipping')
      return
    }

    if (this.token === token && this.isConnecting) {
      console.log('[WSManager] Already connecting with same token, skipping')
      return
    }

    this.token = token
    this.doConnect()
  }

  doConnect() {
    if (!this.token) return
    if (this.ws?.readyState === WebSocket.OPEN) return
    if (this.isConnecting) return

    console.log('[WSManager] Starting new connection')
    this.isConnecting = true

    const url = `${WS_URL}?token=${encodeURIComponent(this.token)}`

    try {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('[WSManager] Connection opened')
        this.isConnecting = false
        this.reconnectAttempt = 0
      }

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          console.log('[WSManager] Message received:', message.type, message)
          this.listeners.forEach((listener) => listener(message))
        } catch (e) {
          console.error('[WSManager] Error parsing message:', e)
        }
      }

      this.ws.onclose = (event) => {
        console.log('[WSManager] Connection closed, code:', event.code, 'wasClean:', event.wasClean, 'hasToken:', !!this.token)
        this.isConnecting = false
        
        if (!this.token) {
          this.ws = null
          return
        }
        
        if (event.wasClean && event.code === 1000) {
          console.log('[WSManager] Clean close, not reconnecting')
          this.ws = null
          return
        }

        this.ws = null
        this.scheduleReconnect()
      }

      this.ws.onerror = () => {
        console.log('[WSManager] Connection error')
        this.isConnecting = false
      }
    } catch (e) {
      console.log('[WSManager] Exception during connect:', e)
      this.isConnecting = false
      this.scheduleReconnect()
    }
  }

  scheduleReconnect() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
    }

    const delay = RECONNECT_DELAYS[Math.min(this.reconnectAttempt, RECONNECT_DELAYS.length - 1)]

    this.reconnectTimeout = setTimeout(() => {
      this.reconnectAttempt++
      this.doConnect()
    }, delay)
  }

  disconnect() {
    this.token = null
    
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }

    if (this.ws) {
      this.ws.close(1000, 'User disconnect')
      this.ws = null
    }

    this.reconnectAttempt = 0
    this.isConnecting = false
  }

  addListener(listener) {
    this.listeners.add(listener)
  }

  removeListener(listener) {
    this.listeners.delete(listener)
  }

  get isConnected() {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

export const wsManager = new WebSocketManager()
