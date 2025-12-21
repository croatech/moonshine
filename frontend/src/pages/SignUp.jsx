import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useMutation, gql } from '@apollo/client'
import { useAuth } from '../context/AuthContext'
import './Auth.css'

const SIGN_UP = gql`
  mutation SignUp($input: SignUpInput!) {
    signUp(input: $input) {
      token
      user {
        id
        username
        email
        hp
        level
        gold
        exp
      }
    }
  }
`

export default function SignUp() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  })
  const [error, setError] = useState('')

  const [signUp, { loading }] = useMutation(SIGN_UP, {
    onCompleted: (data) => {
      login(data.signUp.token, data.signUp.user)
      navigate('/locations/moonshine')
    },
    onError: (err) => {
      const graphQLError = err.graphQLErrors?.[0]
      const errorMessage = graphQLError?.message || err.message || ''
      const lowerMessage = errorMessage.toLowerCase()
      
      // Handle validation errors - show specific messages
      if (lowerMessage === 'user already exists') {
        setError('Username or email already exists. Please try a different one.')
      } else if (lowerMessage === 'invalid input') {
        setError('Please check your input. Username and password must be 3-20 characters.')
      } else if (lowerMessage === 'invalid credentials') {
        setError('Invalid credentials. Please check your input.')
      } else {
        // All other errors (including internal server errors) - show generic message
        setError('Something went wrong. Please try again later.')
      }
    },
  })

  const handleSubmit = (e) => {
    e.preventDefault()
    setError('')
    signUp({ variables: { input: formData } })
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
        <h1>Sign Up</h1>
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
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
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
            {loading ? 'Signing up...' : 'Sign Up'}
          </button>
        </form>
        <p className="auth-link">
          Already have an account? <Link to="/signin">Sign In</Link>
        </p>
      </div>
    </div>
  )
}
