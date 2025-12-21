import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useMutation, gql } from '@apollo/client'
import { useAuth } from '../context/AuthContext'
import './Auth.css'

const SIGN_IN = gql`
  mutation SignIn($input: SignInInput!) {
    signIn(input: $input) {
      token
      user {
        id
        username
        email
        hp
        level
        gold
      }
    }
  }
`

export default function SignIn() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  })
  const [error, setError] = useState('')

  const [signIn, { loading }] = useMutation(SIGN_IN, {
    onCompleted: (data) => {
      login(data.signIn.token, data.signIn.user)
      navigate('/locations/moonshine')
    },
    onError: (err) => {
      const graphQLError = err.graphQLErrors?.[0]
      const errorMessage = graphQLError?.message || err.message || ''
      const lowerMessage = errorMessage.toLowerCase()
      
      // Handle validation errors - show specific messages
      if (lowerMessage === 'invalid credentials') {
        setError('Invalid username or password. Please try again.')
      } else if (lowerMessage === 'invalid input') {
        setError('Please check your input. Username and password must be 3-20 characters.')
      } else {
        // All other errors (including internal server errors) - show generic message
        setError('Something went wrong. Please try again later.')
      }
    },
  })

  const handleSubmit = (e) => {
    e.preventDefault()
    setError('')
    signIn({ variables: { input: formData } })
  }

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    })
  }

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1>Sign In</h1>
        <form onSubmit={handleSubmit}>
          {error && <div className="error-message">{error}</div>}
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              required
              minLength={3}
              maxLength={20}
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              minLength={3}
              maxLength={20}
            />
          </div>
          <button type="submit" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
        <p className="auth-link">
          Don't have an account? <Link to="/signup">Sign Up</Link>
        </p>
      </div>
    </div>
  )
}
