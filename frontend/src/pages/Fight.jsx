import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { userAPI } from '../lib/api'
import { fightAPI } from '../lib/fightAPI'
import PlayerHeader from '../components/PlayerHeader'
import EquipmentDisplay from '../components/EquipmentDisplay'
import { useAuth } from '../context/AuthContext'
import './Fight.css'

export default function Fight() {
  const navigate = useNavigate()
  const { logout } = useAuth()
  const [user, setUser] = useState(null)
  const [bot, setBot] = useState(null)
  const [equippedItems, setEquippedItems] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    setLoading(true)
    Promise.all([
      fightAPI.getCurrentFight(),
      userAPI.getEquippedItems()
    ])
      .then(([fightData, equipped]) => {
        setUser(fightData.user)
        setBot(fightData.bot)
        setEquippedItems(equipped)
        setLoading(false)
      })
      .catch((err) => {
        console.error('[Fight] Error loading fight data:', err)
        if (err.message.includes('no active fight')) {
          navigate('/', { replace: true })
        } else if (err.message.includes('user not found') || err.message.includes('Unauthorized')) {
          localStorage.removeItem('token')
          navigate('/signin', { replace: true })
        } else {
          setError(err.message || 'Ошибка загрузки данных боя')
        }
        setLoading(false)
      })
  }, [navigate])

  const handleLogout = () => {
    logout()
    localStorage.clear()
    navigate('/signin')
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
      <div className="fight-container">
        <div className="fight-main-block">
          <div className="fight-header">
            <PlayerHeader fullWidth={true} />
            <div className="fight-header-actions">
              <button onClick={handleLogout} className="fight-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="fight-content">
            <p>Загрузка боя...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="fight-container">
        <div className="fight-main-block">
          <div className="fight-header">
            <PlayerHeader fullWidth={true} />
            <div className="fight-header-actions">
              <button onClick={handleLogout} className="fight-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="fight-content">
            <p className="fight-error-text">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (!user || !bot) {
    return (
      <div className="fight-container">
        <div className="fight-main-block">
          <div className="fight-header">
            <PlayerHeader fullWidth={true} />
          </div>
          <div className="fight-content">
            <p className="fight-error-text">Данные боя не найдены</p>
          </div>
        </div>
      </div>
    )
  }

  const botAvatarSrc = bot.avatar ? normalizeImagePath(bot.avatar) : '/assets/images/equipment_items/grid/head.png'

  return (
    <div className="fight-container">
      <div className="fight-main-block">
        <div className="fight-header">
          <PlayerHeader fullWidth={true} />
          <div className="fight-header-actions">
            <button onClick={handleLogout} className="fight-logout-button">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
              </svg>
            </button>
          </div>
        </div>
        
        <div className="fight-content">
          <div className="fight-player-section">
            <div className="fight-player-title">
              <h2>{user.username}</h2>
              <span className="fight-player-level">[{user.level}]</span>
            </div>
            <EquipmentDisplay 
              user={user} 
              equippedItems={equippedItems}
              readonly={true}
            />
          </div>

          <div className="fight-arena-section">
            <div className="fight-vs-text">VS</div>
          </div>

          <div className="fight-bot-section">
            <div className="fight-bot-title">
              <h2>{bot.name}</h2>
              <span className="fight-bot-level">[{bot.level}]</span>
            </div>
            <div className="fight-bot-avatar">
              <img 
                src={botAvatarSrc} 
                alt={bot.name}
                onError={(e) => {
                  e.target.src = '/assets/images/equipment_items/grid/head.png'
                }}
              />
            </div>
            <div className="fight-bot-stats">
              <div className="fight-bot-stat">
                <img src="/assets/images/attack.png" alt="Attack" className="fight-stat-icon" />
                <span>{bot.attack || 0}</span>
              </div>
              <div className="fight-bot-stat">
                <img src="/assets/images/defense.png" alt="Defense" className="fight-stat-icon" />
                <span>{bot.defense || 0}</span>
              </div>
              <div className="fight-bot-stat">
                <img src="/assets/images/hp.png" alt="HP" className="fight-stat-icon" />
                <span>{bot.currentHp ?? bot.current_hp ?? bot.hp}/{bot.hp || 0}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
