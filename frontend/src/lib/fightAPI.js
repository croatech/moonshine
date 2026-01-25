const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

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
  if (!trimmed) {
    console.error('[fightAPI] Empty response body')
    return null
  }
  try {
    const parsed = JSON.parse(trimmed)
    console.log('[fightAPI] Parsed response:', parsed)
    return parsed
  } catch (e) {
    console.error('[fightAPI] JSON parse error:', e, 'Response text:', trimmed)
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

  hit: async (attack, defense) => {
    const payload = { attack, defense }
    console.log('[fightAPI] Sending hit request:', payload)
    
    const response = await fetch(`${API_BASE_URL}/fights/current/hit`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(payload),
    })

    const text = await response.text()
    console.log('[fightAPI] Raw response text:', text)
    console.log('[fightAPI] Response status:', response.status, 'ok:', response.ok)
    
    let data
    try {
      data = JSON.parse(text.trim())
      console.log('[fightAPI] Parsed response:', data)
    } catch (e) {
      console.error('[fightAPI] JSON parse error:', e)
      throw new Error('Failed to parse response')
    }
    
    if (!response.ok) {
      console.error('[fightAPI] Hit failed:', response.status, data)
      throw new Error(data?.error || 'Failed to hit')
    }
    console.log('[fightAPI] Hit response data:', data)
    return data
  },
}
