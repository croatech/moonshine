import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { userAPI } from '../lib/api'
import { fightAPI } from '../lib/fightAPI'
import EquipmentDisplay from '../components/EquipmentDisplay'
import { useAuth } from '../context/AuthContext'
import './Fight.css'

const BODY_PARTS = [
  { value: 'head', label: 'Голова' },
  { value: 'neck', label: 'Шея' },
  { value: 'chest', label: 'Грудь' },
  { value: 'belt', label: 'Пояс' },
  { value: 'legs', label: 'Ноги' },
]

export default function Fight() {
  const navigate = useNavigate()
  const { logout } = useAuth()
  const [user, setUser] = useState(null)
  const [bot, setBot] = useState(null)
  const [equippedItems, setEquippedItems] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [selectedAttack, setSelectedAttack] = useState(null)
  const [selectedDefense, setSelectedDefense] = useState(null)
  const [hitting, setHitting] = useState(false)

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

  const handleHit = async () => {
    if (!selectedAttack || !selectedDefense) {
      return
    }

    setHitting(true)
    try {
      await fightAPI.hit(selectedAttack, selectedDefense)
      const [fightData, equipped] = await Promise.all([
        fightAPI.getCurrentFight(),
        userAPI.getEquippedItems()
      ])
      setUser(fightData.user)
      setBot(fightData.bot)
      setEquippedItems(equipped)
    } catch (err) {
      console.error('[Fight] Error hitting:', err)
      setError(err.message || 'Ошибка при ударе')
    } finally {
      setHitting(false)
    }
  }

  if (loading) {
    return (
      <div className="fight-container">
        <div className="fight-main-block">
          <div className="fight-header">
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
            <div className="fight-header-actions">
              <button onClick={handleLogout} className="fight-logout-button">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 21V3h8v2H5v14h6v2H3zm13-4l-1.375-1.45 2.55-2.55H9v-2h8.175l-2.55-2.55L16 7l5 5-5 5z" fill="currentColor" />
                </svg>
              </button>
            </div>
          </div>
          <div className="fight-content">
            <p className="fight-error-text">Данные боя не найдены</p>
          </div>
        </div>
      </div>
    )
  }

  const getBotImageSrc = () => {
    if (!bot.avatar) {
      return '/assets/images/bots/rat.jpg'
    }
    
    let path = bot.avatar.trim()
    
    if (path.startsWith('/')) {
      path = path.slice(1)
    }
    
    path = path.replace(/^frontend\/assets\/images\//, '')
    path = path.replace(/^assets\/images\//, '')
    
    if (path.startsWith('images/')) {
      path = path.replace(/^images\//, '')
    }
    
    if (!path.includes('.')) {
      path = `${path}.jpg`
    }
    
    if (path.includes('rat') && !path.includes('bots/')) {
      path = `bots/${path}`
    } else if (!path.startsWith('bots/')) {
      path = `bots/${path}`
    }
    
    return `/assets/images/${path}`
  }
  
  const botImageSrc = getBotImageSrc()

  const botUserData = {
    username: bot.name,
    name: bot.name,
    level: bot.level,
    hp: bot.hp,
    currentHp: bot.currentHp ?? bot.current_hp ?? bot.hp,
    avatar: bot.avatar,
    attack: bot.attack,
    defense: bot.defense,
  }

  const currentHp = user.currentHp ?? user.current_hp ?? 0
  const maxHp = user.hp || 0
  const botCurrentHp = bot.currentHp ?? bot.current_hp ?? bot.hp
  const botMaxHp = bot.hp || 0

  return (
    <div className="fight-container">
      <div className="fight-main-block">
        <div className="fight-header">
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
            <div className="fight-player-info">
              <div className="fight-player-title">
                <h2>{user.username}</h2>
                <span className="fight-player-level">[{user.level}]</span>
              </div>
              <div className="fight-hp-bar">
                <div 
                  className="fight-hp-fill" 
                  style={{ 
                    width: `${maxHp > 0 ? Math.round((currentHp / maxHp) * 100) : 0}%`,
                    backgroundColor: '#dc3545'
                  }}
                ></div>
                <span className="fight-hp-text">
                  {currentHp}/{maxHp}
                </span>
              </div>
            </div>
            <EquipmentDisplay 
              user={user} 
              equippedItems={equippedItems}
              readonly={true}
            />
          </div>

          <div className="fight-arena-section">
            <div className="fight-vs-text">VS</div>
            
            <div className="fight-controls">
              <div className="fight-controls-column">
                <div className="fight-controls-title">Защита</div>
                <div className="fight-body-parts-list">
                  {BODY_PARTS.map((part) => (
                    <button
                      key={part.value}
                      className={`fight-body-part-button ${selectedDefense === part.value ? 'selected' : ''}`}
                      onClick={() => setSelectedDefense(part.value)}
                    >
                      {part.label}
                    </button>
                  ))}
                </div>
              </div>

              <div className="fight-controls-column">
                <div className="fight-controls-title">Атака</div>
                <div className="fight-body-parts-list">
                  {BODY_PARTS.map((part) => (
                    <button
                      key={part.value}
                      className={`fight-body-part-button ${selectedAttack === part.value ? 'selected' : ''}`}
                      onClick={() => setSelectedAttack(part.value)}
                    >
                      {part.label}
                    </button>
                  ))}
                </div>
              </div>
            </div>

            <button
              className="fight-hit-button"
              onClick={handleHit}
              disabled={!selectedAttack || !selectedDefense || hitting}
            >
              {hitting ? 'Удар...' : 'Ударить'}
            </button>
          </div>

          <div className="fight-bot-section">
            <div className="fight-bot-info">
              <div className="fight-bot-title">
                <h2>{bot.name}</h2>
                <span className="fight-bot-level">[{bot.level}]</span>
              </div>
              <div className="fight-hp-bar">
                <div 
                  className="fight-hp-fill fight-hp-fill-full" 
                  style={{ 
                    width: '100%',
                    backgroundColor: '#dc3545'
                  }}
                ></div>
                <span className="fight-hp-text">
                  {botMaxHp}/{botMaxHp}
                </span>
              </div>
            </div>
            <EquipmentDisplay 
              user={botUserData} 
              equippedItems={{}}
              readonly={true}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
