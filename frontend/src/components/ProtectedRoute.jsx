import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function ProtectedRoute({ children }) {
  const { token, loading } = useAuth()

  console.log('[ProtectedRoute] token:', !!token, 'loading:', loading)

  // Показываем loading пока идет загрузка
  if (loading) {
    return <div>Загрузка...</div>
  }

  // Редиректим только если нет токена
  if (!token) {
    console.log('[ProtectedRoute] No token, redirecting to signin')
    return <Navigate to="/signin" replace />
  }

  return children
}


