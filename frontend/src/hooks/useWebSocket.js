import { useEffect, useRef, useCallback, useState } from 'react'

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/ws'

const RECONNECT_DELAYS = [1000, 2000, 4000, 8000, 16000, 30000]

export function useWebSocket(token, onMessage, enabled = true) {
  const wsRef = useRef(null)
  const reconnectAttemptRef = useRef(0)
  const reconnectTimeoutRef = useRef(null)
  const onMessageRef = useRef(onMessage)
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    onMessageRef.current = onMessage
  }, [onMessage])

  const scheduleReconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }

    const attempt = reconnectAttemptRef.current
    const delay = RECONNECT_DELAYS[Math.min(attempt, RECONNECT_DELAYS.length - 1)]

    reconnectTimeoutRef.current = setTimeout(() => {
      reconnectAttemptRef.current++
      connect()
    }, delay)
  }, [])

  const connect = useCallback(() => {
    if (!token || !enabled) {
      return
    }

    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return
    }

    if (wsRef.current?.readyState === WebSocket.CONNECTING) {
      return
    }

    const url = `${WS_URL}?token=${encodeURIComponent(token)}`
    
    try {
      const ws = new WebSocket(url)

      ws.onopen = () => {
        reconnectAttemptRef.current = 0
        setIsConnected(true)
      }

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          onMessageRef.current?.(message)
        } catch (e) {
        }
      }

      ws.onclose = (event) => {
        setIsConnected(false)
        wsRef.current = null

        if (enabled && token && !event.wasClean) {
          scheduleReconnect()
        }
      }

      ws.onerror = () => {
        setIsConnected(false)
      }

      wsRef.current = ws
    } catch (e) {
      scheduleReconnect()
    }
  }, [token, enabled, scheduleReconnect])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnect')
      wsRef.current = null
    }

    setIsConnected(false)
    reconnectAttemptRef.current = 0
  }, [])

  useEffect(() => {
    if (enabled && token) {
      connect()
    } else {
      disconnect()
    }

    return () => {
      disconnect()
    }
  }, [token, enabled])

  return { isConnected, disconnect }
}

export default useWebSocket
