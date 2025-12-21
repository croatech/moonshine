import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import SignUp from './pages/SignUp'
import SignIn from './pages/SignIn'
import Location from './pages/Location'
import ProtectedRoute from './components/ProtectedRoute'

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/signup" element={<SignUp />} />
        <Route path="/signin" element={<SignIn />} />
        <Route
          path="/locations/:slug"
          element={
            <ProtectedRoute>
              <Location />
            </ProtectedRoute>
          }
        />
        <Route
          path="/"
          element={<Navigate to="/locations/moonshine" replace />}
        />
        <Route path="*" element={<Navigate to="/signin" replace />} />
      </Routes>
    </AuthProvider>
  )
}

export default App

