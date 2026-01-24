const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/ws'

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
    if (this.token === token && this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    if (this.token === token && this.isConnecting) {
      return
    }

    this.token = token
    this.doConnect()
  }

  doConnect() {
    if (!this.token) return
    if (this.ws?.readyState === WebSocket.OPEN) return
    if (this.isConnecting) return

    this.isConnecting = true

    const url = `${WS_URL}?token=${encodeURIComponent(this.token)}`

    try {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        this.isConnecting = false
        this.reconnectAttempt = 0
      }

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          this.listeners.forEach((listener) => listener(message))
        } catch (e) {
        }
      }

      this.ws.onclose = (event) => {
        this.isConnecting = false
        this.ws = null

        if (this.token && !event.wasClean) {
          this.scheduleReconnect()
        }
      }

      this.ws.onerror = () => {
        this.isConnecting = false
      }
    } catch (e) {
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
