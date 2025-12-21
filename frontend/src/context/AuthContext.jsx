import { createContext, useContext, useState, useEffect, useCallback, useRef } from 'react'
import { useLazyQuery, gql } from '@apollo/client'

const GET_CURRENT_USER = gql`
  query GetCurrentUser {
    currentUser {
      id
      username
      email
      hp
      level
      gold
      exp
    }
  }
`

const AuthContext = createContext()

export function AuthProvider({ children }) {
  const [token, setToken] = useState(localStorage.getItem('token'))
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(false)
  const hasInitialized = useRef(false)
  
  const [getCurrentUser] = useLazyQuery(GET_CURRENT_USER, {
    errorPolicy: 'ignore',
    onCompleted: (data) => {
      if (data?.currentUser) {
        setUser(data.currentUser)
      }
      setLoading(false)
    },
    onError: () => {
      // If query fails, clear token and user (invalid token)
      setToken(null)
      setUser(null)
      localStorage.removeItem('token')
      setLoading(false)
    },
  })

  // Load user on mount if token exists (only once)
  useEffect(() => {
    if (hasInitialized.current) return
    hasInitialized.current = true

    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setLoading(true)
      getCurrentUser()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []) // Only run on mount

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token)
    } else {
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
      getCurrentUser()
    }
  }, [getCurrentUser])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('token')
  }, [])

  const refetchUser = useCallback(() => {
    if (token) {
      setLoading(true)
      getCurrentUser()
    }
  }, [token, getCurrentUser])

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

