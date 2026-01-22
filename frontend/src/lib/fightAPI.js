const API_BASE_URL = 'http://localhost:8080/api'

function getAuthHeaders() {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  }
}

async function parseResponse(response) {
  const text = await response.text()
  const trimmed = text.trim()
  if (!trimmed) return null
  try {
    return JSON.parse(trimmed)
  } catch (e) {
    return null
  }
}

export const fightAPI = {
  getCurrentFight: async () => {
    const response = await fetch(`${API_BASE_URL}/fights/current`, {
      method: 'GET',
      headers: getAuthHeaders(),
    })

    const data = await parseResponse(response)
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('no active fight')
      }
      throw new Error(data?.error || 'Failed to get current fight')
    }
    return data
  },
}
