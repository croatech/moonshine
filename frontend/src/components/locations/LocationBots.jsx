import { useAuth } from '../../context/AuthContext'
import './LocationBots.css'

export default function LocationBots() {
  const { user } = useAuth()
  
  // Получаем bots из user.location.bots
  const bots = user?.location?.bots || []

  if (!bots || bots.length === 0) {
    return null
  }

  const handleAttack = (botId) => {
    // TODO: Реализовать логику атаки
    console.log('Attack bot:', botId)
  }

  return (
    <div className="location-bots">
      <h3 className="location-bots-title">Боты</h3>
      <div className="location-bots-list">
        {bots.map((bot) => (
          <div key={bot.id} className="location-bot-item">
            <span className="location-bot-name">
              {bot.name} [{bot.level}]
            </span>
            <a
              href="#"
              onClick={(e) => {
                e.preventDefault()
                handleAttack(bot.id)
              }}
              className="location-bot-attack-link"
            >
              напасть
            </a>
          </div>
        ))}
      </div>
    </div>
  )
}

