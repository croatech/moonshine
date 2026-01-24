import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function ProtectedRoute({ children }) {
  const { token, loading, user } = useAuth()
  const location = useLocation()

  const storedToken = localStorage.getItem('token')
  const hasToken = token || storedToken

  if (loading) {
    return <div>Загрузка...</div>
  }

  if (!hasToken) {
    return <Navigate to="/signin" replace />
  }

  if (user) {
    const inFight = user.inFight === true || user.InFight === true
    const currentPath = location.pathname
    
    if (inFight && currentPath !== '/fight') {
      return <Navigate to="/fight" replace />
    }
  }

  return children
}
