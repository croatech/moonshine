import { useEffect, useRef, useState } from 'react'
import { wsManager } from '../lib/websocket'

export function useWebSocket(token, onMessage, enabled = true) {
  const onMessageRef = useRef(onMessage)
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    onMessageRef.current = onMessage
  }, [onMessage])

  useEffect(() => {
    if (!enabled || !token) {
      return
    }

    const listener = (message) => {
      onMessageRef.current?.(message)
    }

    wsManager.addListener(listener)
    wsManager.connect(token)

    const checkConnection = setInterval(() => {
      setIsConnected(wsManager.isConnected)
    }, 1000)

    return () => {
      wsManager.removeListener(listener)
      clearInterval(checkConnection)
    }
  }, [token, enabled])

  return { isConnected }
}

export default useWebSocket
