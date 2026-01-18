import { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react'
import { userAPI } from '../lib/api'

const AuthContext = createContext()

const CACHE_DURATION = 30000
const cache = {
  user: null,
  timestamp: 0,
  inventory: null,
  inventoryTimestamp: 0,
  equipped: null,
  equippedTimestamp: 0,
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(() => !!localStorage.getItem('token'))

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
          setToken(null)
          setUser(null)
          localStorage.removeItem('token')
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

  useEffect(() => {
    if (!user || !user.hp || !token) return

    const currentHp = user.current_hp || user.currentHp || 0
    const maxHp = user.hp || 0
    const isHpLessThan100 = currentHp < maxHp

    if (!isHpLessThan100) {
      return
    }

    const intervalId = setInterval(() => {
      userAPI.getCurrentUser()
        .then((userData) => {
          const newCurrentHp = userData.current_hp || userData.currentHp || 0
          const newMaxHp = userData.hp || 0
          
          cache.user = userData
          cache.timestamp = Date.now()
          setUser(userData)
          
          if (newCurrentHp >= newMaxHp) {
            clearInterval(intervalId)
          }
        })
        .catch((err) => {
          console.error('[AuthContext] Error polling user data:', err)
        })
    }, 5000)

    return () => {
      clearInterval(intervalId)
    }
  }, [user, token])

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
    cache.inventory = null
    cache.inventoryTimestamp = 0
    cache.equipped = null
    cache.equippedTimestamp = 0
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
