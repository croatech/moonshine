import { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react'
import { userAPI } from '../lib/api'
import { useWebSocket } from '../hooks/useWebSocket'

const AuthContext = createContext()

const CACHE_DURATION = 30000
const cache = {
  user: null,
  timestamp: 0,
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(() => !!localStorage.getItem('token'))

  const handleWebSocketMessage = useCallback((message) => {
    if (message.type === 'hp_update' && message.data) {
      setUser((prevUser) => {
        if (!prevUser) return prevUser
        const updated = {
          ...prevUser,
          currentHp: message.data.currentHp,
          current_hp: message.data.currentHp,
          hp: message.data.hp,
        }
        cache.user = updated
        cache.timestamp = Date.now()
        return updated
      })
    }
  }, [])

  const shouldConnectWS = useMemo(() => {
    return !!token && !!user
  }, [token, user])

  useWebSocket(token, handleWebSocketMessage, shouldConnectWS)

  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    
    if (storedToken) {
      const now = Date.now()
      if (cache.user && (now - cache.timestamp) < CACHE_DURATION) {
        setUser(cache.user)
        setLoading(false)
        return
      }

      setLoading(true)
      userAPI.getCurrentUser()
        .then((userData) => {
          cache.user = userData
          cache.timestamp = Date.now()
          setUser(userData)
          setLoading(false)
        })
        .catch((err) => {
          console.error('[AuthContext] Error fetching current user:', err)
          const errorMsg = err.message || ''
          if (errorMsg.toLowerCase().includes('unauthorized')) {
            setToken(null)
            setUser(null)
            localStorage.removeItem('token')
          }
          setLoading(false)
        })
    } else {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token)
    } else {
      localStorage.removeItem('token')
      setUser(null)
      cache.user = null
      cache.timestamp = 0
    }
  }, [token])

  const login = useCallback((newToken, userData = null) => {
    setToken(newToken)
    if (userData) {
      cache.user = userData
      cache.timestamp = Date.now()
      setUser(userData)
    } else {
      setLoading(true)
      userAPI.getCurrentUser()
        .then((userData) => {
          cache.user = userData
          cache.timestamp = Date.now()
          setUser(userData)
          setLoading(false)
        })
        .catch((err) => {
          console.error('[AuthContext] Error fetching user after login:', err)
          setToken(null)
          setUser(null)
          localStorage.removeItem('token')
          setLoading(false)
        })
    }
  }, [])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
    localStorage.clear()
    cache.user = null
    cache.timestamp = 0
  }, [])

  const refetchUser = useCallback(() => {
    if (token) {
      return userAPI.getCurrentUser()
        .then((userData) => {
          cache.user = userData
          cache.timestamp = Date.now()
          setUser(userData)
          return userData
        })
        .catch((err) => {
          console.error('[AuthContext] Error refetching user:', err)
          throw err
        })
    }
    return Promise.resolve(null)
  }, [token])

  const value = useMemo(() => ({
    user,
    token,
    loading,
    isAuthenticated: !!token && !!user,
    login,
    logout,
    refetchUser,
  }), [user, token, loading, login, logout, refetchUser])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
