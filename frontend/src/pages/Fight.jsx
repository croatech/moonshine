import { useState, useEffect } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { botAPI, userAPI } from '../lib/api'
import PlayerHeader from '../components/PlayerHeader'
import EquipmentDisplay from '../components/EquipmentDisplay'
import { useAuth } from '../context/AuthContext'
import './Profile.css'

export default function Fight() {
  const [searchParams] = useSearchParams()
  const botSlug = searchParams.get('bot')
  const navigate = useNavigate()
  const { logout } = useAuth()
  const [user, setUser] = useState(null)
  const [bot, setBot] = useState(null)
  const [equippedItems, setEquippedItems] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!botSlug) {
      setError('Bot slug is required')
      setLoading(false)
      return
    }

    setLoading(true)
    Promise.all([
      botAPI.attack(botSlug),
      userAPI.getEquippedItems()
    ])
      .then(([fightData, equipped]) => {
        setUser(fightData.user)
        setBot(fightData.bot)
        setEquippedItems(equipped)
        setLoading(false)
        navigate('/fight', { replace: true })
      })
      .catch((err) => {
        console.error('[Fight] Error loading fight data:', err)
        setError(err.message || 'Ошибка загрузки данных боя')
        setLoading(false)
      })
  }, [botSlug, navigate])

  const handleBack = () => {
    navigate(-1)
  }

  const handleLogout = () => {
    logout()
    localStorage.clear()
    navigate('/signin')
  }

  if (loading) {
    return (
      <div className="profile-container">
        <div className="profile-main-block">
          <div className="profile-header">
            <PlayerHeader fullWidth={true} />
            <div className="profile-header-actions">
              <button onClick={handleBack} className="profile-back-button">
                ← Назад
              </button>
              <button onClick={handleLogout} className="profile-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="profile-content">
            <p>Загрузка...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="profile-container">
        <div className="profile-main-block">
          <div className="profile-header">
            <PlayerHeader fullWidth={true} />
            <div className="profile-header-actions">
              <button onClick={handleBack} className="profile-back-button">
                ← Назад
              </button>
              <button onClick={handleLogout} className="profile-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="profile-content">
            <p className="fight-error">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (!user || !bot) {
    return (
      <div className="profile-container">
        <div className="profile-main-block">
          <div className="profile-header">
            <PlayerHeader fullWidth={true} />
            <div className="profile-header-actions">
              <button onClick={handleBack} className="profile-back-button">
                ← Назад
              </button>
              <button onClick={handleLogout} className="profile-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="profile-content">
            <p className="fight-error">Данные не найдены</p>
          </div>
        </div>
      </div>
    )
  }

  const botEquippedItems = {}

  return (
    <div className="profile-container">
      <div className="profile-main-block">
        <div className="profile-header">
          <PlayerHeader fullWidth={true} />
          <div className="profile-header-actions">
            <button onClick={handleBack} className="profile-back-button">
              ← Назад
            </button>
            <button onClick={handleLogout} className="profile-logout-button">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
              </svg>
            </button>
          </div>
        </div>
        <div className="profile-content">
          <div className="fight-profiles-container">
            <div className="fight-profile-section">
              <div className="fight-profile-title">
                <h2>{user.username}</h2>
                <div className="fight-profile-level">Уровень {user.level}</div>
              </div>
              <EquipmentDisplay 
                user={user} 
                equippedItems={equippedItems}
                readonly={true}
              />
            </div>

            <div className="fight-profile-section">
              <div className="fight-profile-title">
                <h2>{bot.name}</h2>
                <div className="fight-profile-level">Уровень {bot.level}</div>
              </div>
              <EquipmentDisplay 
                user={bot} 
                equippedItems={botEquippedItems}
                readonly={true}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
