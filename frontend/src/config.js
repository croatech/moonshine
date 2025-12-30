export default {
  stats: ["attack", "defense", "hp"],
  baseUrl: process.env.VITE_API_URL || 'http://localhost:8080',
  apiUrl: process.env.VITE_API_URL || 'http://localhost:8080/api',
  cableUrl: process.env.VITE_CABLE_URL || 'ws://localhost:8080/cable'
}








