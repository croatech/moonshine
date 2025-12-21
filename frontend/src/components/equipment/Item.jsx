import config from '../../config'

export default function EquipmentItem({ item, playerLevel }) {
  if (!item) return null

  const checkForNotAvailableLevel = (itemLevel, playerLevel) => {
    return itemLevel > playerLevel ? 'red' : ''
  }

  return (
    <div className="row top">
      <div className="col-md-3">
        <img src={item.image?.url} alt={item.name} />
      </div>

      <div className="col-md-9">
        <div className={`level ${checkForNotAvailableLevel(item.required_level, playerLevel)}`}>
          [{item.required_level}]
        </div>

        <h3>{item.name}</h3>

        {config.stats.map((stat) => (
          item[stat] && (
            <div key={stat} className={`${stat} stat`}>
              <img src={`/assets/images/${stat}.png`} alt={stat} />
              {item[stat]}
            </div>
          )
        ))}
      </div>
    </div>
  )
}

