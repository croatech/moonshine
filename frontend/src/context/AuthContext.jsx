import { createContext, useContext, useState, useEffect, useCallback } from 'react'
import { userAPI } from '../lib/api'

const AuthContext = createContext()

export function AuthProvider({ children }) {
  // Инициализируем токен из localStorage сразу
  const [token, setToken] = useState(() => localStorage.getItem('token'))
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(() => {
    const hasToken = !!localStorage.getItem('token')
    console.log('[AuthContext] Initial loading state:', hasToken)
    return hasToken
  })

  // Load user on mount if token exists
  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    console.log('[AuthContext] Mount effect, token exists:', !!storedToken)
    
    if (storedToken) {
      console.log('[AuthContext] Fetching current user...')
      setLoading(true)
      userAPI.getCurrentUser()
        .then((userData) => {
          console.log('[AuthContext] User data loaded:', userData)
          setUser(userData)
          setLoading(false)
        })
        .catch((err) => {
          console.error('[AuthContext] Error fetching current user:', err)
          // If query fails, clear token and user (invalid token)
          setToken(null)
          setUser(null)
          localStorage.removeItem('token')
          setLoading(false)
        })
    } else {
      setLoading(false)
    }
  }, []) // Only run on mount

  // Sync token with localStorage
  useEffect(() => {
    if (token) {
      console.log('[AuthContext] Setting token to localStorage')
      localStorage.setItem('token', token)
    } else {
      console.log('[AuthContext] Removing token from localStorage')
      localStorage.removeItem('token')
      setUser(null)
    }
  }, [token])

  const login = useCallback((newToken, userData = null) => {
    setToken(newToken)
    if (userData) {
      setUser(userData)
    } else {
      // If no user data provided, fetch it
      setLoading(true)
      userAPI.getCurrentUser()
        .then((userData) => {
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
    console.log('[AuthContext] Login called. New token:', newToken)
  }, [])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
    localStorage.clear()
    console.log('[AuthContext] Logout called. localStorage cleared.')
  }, [])

  const refetchUser = useCallback(() => {
    if (token) {
      setLoading(true)
      userAPI.getCurrentUser()
        .then((userData) => {
          setUser(userData)
          setLoading(false)
        })
        .catch((err) => {
          console.error('[AuthContext] Error refetching user:', err)
          setLoading(false)
        })
    }
  }, [token])

  const value = {
    user,
    token,
    loading,
    isAuthenticated: !!token && !!user,
    login,
    logout,
    refetchUser,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
