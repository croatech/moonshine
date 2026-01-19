import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { botAPI } from '../lib/api'
import PlayerHeader from '../components/PlayerHeader'
import { useAuth } from '../context/AuthContext'
import './Bots.css'

export default function Bots() {
  const { location_slug } = useParams()
  const { logout } = useAuth()
  const navigate = useNavigate()
  const [bots, setBots] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!location_slug) {
      setError('Location slug is required')
      setLoading(false)
      return
    }

    setLoading(true)
    botAPI.getBots(location_slug)
      .then((data) => {
        setBots(data)
        setLoading(false)
      })
      .catch((err) => {
        console.error('[Bots] Error loading bots:', err)
        setError('Ошибка загрузки ботов')
        setLoading(false)
      })
  }, [location_slug])

  const handleBack = () => {
    navigate(`/locations/${location_slug}`)
  }

  const handleLogout = () => {
    logout()
    localStorage.clear()
    navigate('/signin')
  }

  const handleAttack = (botSlug) => {
    if (!botSlug) {
      return
    }
    navigate(`/fight?bot=${botSlug}`)
  }

  const normalizeImagePath = (img) => {
    if (!img) return null
    let p = img
    if (p.startsWith('/')) p = p.slice(1)
    p = p.replace(/^frontend\/assets\/images\//, '')
    if (p.startsWith('assets/images/')) p = p.replace(/^assets\/images\//, '')
    if (!p.startsWith('images/')) p = `images/${p}`
    return `/assets/${p}`
  }

  if (loading) {
    return (
      <div className="bots-container">
        <div className="bots-main-block">
          <div className="bots-header">
            <PlayerHeader />
          </div>
          <div className="bots-content">Загрузка...</div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bots-container">
        <div className="bots-main-block">
          <div className="bots-header">
            <PlayerHeader />
          </div>
          <div className="bots-content">{error}</div>
        </div>
      </div>
    )
  }

  return (
    <div className="bots-container">
      <div className="bots-main-block">
        <div className="bots-header">
          <PlayerHeader />
          <div className="bots-header-actions">
            <button 
              onClick={handleBack}
              className="bots-back-button"
              title="Назад"
            >
              ← Назад
            </button>
            <button 
              onClick={handleLogout} 
              className="bots-logout-button"
              title="Выйти из игры"
            >
              <svg 
                width="24" 
                height="24" 
                viewBox="0 0 24 24" 
                fill="none" 
                xmlns="http://www.w3.org/2000/svg"
              >
                <path 
                  d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" 
                  fill="currentColor"
                />
              </svg>
            </button>
          </div>
        </div>
        <div className="bots-content">
          <h2>Боты</h2>
          <div className="bots-list">
            {bots.length === 0 ? (
              <p>Боты не найдены</p>
            ) : (
              bots.map((bot) => (
                <div key={bot.id} className="bot-card">
                  {bot.avatar && (
                    <img 
                      src={normalizeImagePath(bot.avatar)} 
                      alt={bot.name}
                      className="bot-image"
                    />
                  )}
                  <div className="bot-info">
                    <h3>{bot.name}</h3>
                    <div className="bot-stats">
                      <div>Уровень: {bot.level}</div>
                      <div>Атака: {bot.attack}</div>
                      <div>Защита: {bot.defense}</div>
                      <div>HP: {bot.hp}</div>
                    </div>
                    <button 
                      className="bot-attack-button"
                      onClick={() => handleAttack(bot.slug)}
                    >
                      Атаковать
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
