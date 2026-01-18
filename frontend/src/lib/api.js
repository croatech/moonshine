const API_BASE_URL = 'http://localhost:8080/api'

// Helper function to get auth headers
function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  }
}

// Auth API
export const authAPI = {
  signUp: async (username, email, password) => {
    const response = await fetch(`${API_BASE_URL}/auth/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, email, password }),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Sign up failed')
    }

    return response.json()
  },

  signIn: async (username, password) => {
    const response = await fetch(`${API_BASE_URL}/auth/signin`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Sign in failed')
    }

    return response.json()
  },
}

// User API
export const userAPI = {
  getCurrentUser: async () => {
    const response = await fetch(`${API_BASE_URL}/user/me`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to get current user')
    }

    return response.json()
  },

  getInventory: async () => {
    const response = await fetch(`${API_BASE_URL}/users/me/inventory`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        throw new Error('Unauthorized')
      }
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch inventory')
    }

    return response.json()
  },

  getEquippedItems: async () => {
    const response = await fetch(`${API_BASE_URL}/users/me/equipped`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        throw new Error('Unauthorized')
      }
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch equipped items')
    }

    return response.json()
  },

  updateProfile: async (data) => {
    const response = await fetch(`${API_BASE_URL}/user/me`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        throw new Error('Unauthorized')
      }
      const error = await response.json()
      throw new Error(error.error || 'Failed to update profile')
    }

    return response.json()
  },
}

// Equipment API
export const equipmentAPI = {
  getByCategory: async (category) => {
    const response = await fetch(`${API_BASE_URL}/equipment_items?category=${encodeURIComponent(category)}`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        throw new Error('Unauthorized')
      }
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch equipment items')
    }

    return response.json()
  },

  buy: async (itemSlug) => {
    const response = await fetch(`${API_BASE_URL}/equipment_items/${itemSlug}/buy`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to buy item')
    }

    return response.json()
  },

  sell: async (itemSlug) => {
    const response = await fetch(`${API_BASE_URL}/equipment_items/${itemSlug}/sell`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to sell item')
    }

    return response.json()
  },

  takeOn: async (itemSlug) => {
    const response = await fetch(`${API_BASE_URL}/equipment_items/${itemSlug}/take_on`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to equip item')
    }

    return response.json()
  },

  takeOff: async (slotName) => {
    const response = await fetch(`${API_BASE_URL}/equipment_items/take_off/${slotName}`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to remove item')
    }

    return response.json()
  },
}

// Avatar API
export const avatarAPI = {
  getAll: async () => {
    const response = await fetch(`${API_BASE_URL}/avatars`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token')
        throw new Error('Unauthorized')
      }
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch avatars')
    }

    return response.json()
  },
}

// Location API
export const locationAPI = {
  move: async (locationSlug) => {
    const response = await fetch(`${API_BASE_URL}/locations/${locationSlug}/move`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to move to location')
    }

    return response.json()
  },

  moveToCell: async (locationSlug, cellSlug) => {
    const response = await fetch(`${API_BASE_URL}/locations/${locationSlug}/cells/${cellSlug}/move`, {
      method: 'POST',
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to move to cell')
    }

    return response.json()
  },
  
  getCells: async (locationSlug) => {
    const response = await fetch(`${API_BASE_URL}/locations/${locationSlug}/cells`, {
      headers: getAuthHeaders(),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to get location cells')
    }

    return response.json()
  },
}

