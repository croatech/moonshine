import './StatsDisplay.css'

export default function StatsDisplay({ attack, defense, hp, currentHp, gold, showGold = false }) {
  return (
    <div className="equipment-stats-compact">
      <div className="stats-column">
        <div className="stat-row-compact">
          <img src="/assets/images/attack.png" alt="Attack" className="stat-icon-compact" />
          <span>{attack || 0}</span>
        </div>
        <div className="stat-row-compact">
          <img src="/assets/images/defense.png" alt="Defense" className="stat-icon-compact" />
          <span>{defense || 0}</span>
        </div>
      </div>
      <div className="stats-right-column">
        {showGold && gold !== undefined && (
          <div className="stat-row-compact stat-gold">
            <span>{gold || 0} зол.</span>
          </div>
        )}
        <div className="stat-row-compact stat-hp">
          <img src="/assets/images/hp.png" alt="HP" className="stat-icon-compact" />
          <span>{currentHp || 0}/{hp || 0}</span>
        </div>
      </div>
    </div>
  )
}
