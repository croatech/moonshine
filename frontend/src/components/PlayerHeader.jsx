import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import './PlayerHeader.css'

export default function PlayerHeader({ fullWidth = false }) {
  const { user } = useAuth()
  const currentHp = user?.currentHp ?? user?.current_hp ?? 0
  const maxHp = user?.hp || 20

  return (
    <div className={`player-header-info ${fullWidth ? 'player-header-info-full' : ''}`}>
      <div className="player-name-level">
        <span className="player-name">{user?.username || 'Игрок'}</span>
        <span className="player-level">[<span className="player-level-number">{user?.level || 1}</span>]</span>
      </div>
      <div className="player-hp-bar">
        <div 
          className="player-hp-fill" 
          style={{ 
            width: `${maxHp > 0 ? Math.round((currentHp / maxHp) * 100) : 0}%`,
            backgroundColor: '#dc3545'
          }}
        ></div>
        <span className="player-hp-text">
          {currentHp}/{maxHp}
        </span>
      </div>
      <Link 
        to="/profile" 
        className="profile-link-button" 
        title="Профиль персонажа"
      >
        <svg 
          width="20" 
          height="20" 
          viewBox="0 0 24 24" 
          fill="none" 
          xmlns="http://www.w3.org/2000/svg"
        >
          <path 
            d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" 
            fill="currentColor"
          />
        </svg>
      </Link>
    </div>
  )
}
