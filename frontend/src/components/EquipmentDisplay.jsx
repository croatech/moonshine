import './EquipmentDisplay.css'

export default function EquipmentDisplay({ user, equippedItems, readonly = false }) {
  const normalizeImagePath = (img) => {
    if (!img) return null
    let p = img
    if (p.startsWith('/')) p = p.slice(1)
    p = p.replace(/^frontend\/assets\/images\//, '')
    if (p.startsWith('assets/images/')) p = p.replace(/^assets\/images\//, '')
    if (!p.startsWith('images/')) p = `images/${p}`
    return `/assets/${p}`
  }

  const getEquipmentSlotImage = (slotName) => {
    const placeholderName = slotName.startsWith('ring') ? 'ring' : slotName
    return `/assets/images/equipment_items/grid/${placeholderName}.png`
  }

  const hasEquippedItem = (slotName) => {
    const equippedItem = equippedItems?.[slotName] || equippedItems?.[slotName.toLowerCase()]
    return equippedItem && equippedItem.image
  }

  const renderEquipmentSlot = (slotName, alt) => {
    const hasItem = hasEquippedItem(slotName)
    
    return (
      <div 
        className={`equipment-slot ${hasItem ? 'has-item' : ''} ${readonly ? 'readonly' : ''}`}
        title={hasItem && !readonly ? 'Нажмите чтобы снять' : ''}
      >
        <img src={getEquipmentSlotImage(slotName)} alt={alt} />
      </div>
    )
  }

  const avatarImage = typeof user?.avatar === 'string' 
    ? user.avatar 
    : user?.avatar?.image
  const avatarSrc = avatarImage ? normalizeImagePath(avatarImage) : getEquipmentSlotImage('head')

  return (
    <div className="equipment-display-wrapper">
      <div className="equipment-grid">
        <div className="equipment-column-left">
          {renderEquipmentSlot('head', 'head')}
          {renderEquipmentSlot('neck', 'neck')}
          {renderEquipmentSlot('weapon', 'weapon')}
          {renderEquipmentSlot('legs', 'legs')}
          {renderEquipmentSlot('feet', 'feet')}
        </div>

        <div className="equipment-column-center">
          <div className="equipment-avatar">
            {avatarSrc ? (
              <img 
                src={avatarSrc} 
                alt={user?.username || user?.name || 'Avatar'}
                className="avatar-image"
                onError={(e) => {
                  e.target.src = getEquipmentSlotImage('head')
                }}
              />
            ) : (
              <img src={getEquipmentSlotImage('head')} alt="avatar placeholder" />
            )}
          </div>
          <div className="equipment-rings">
            {renderEquipmentSlot('ring1', 'ring1')}
            {renderEquipmentSlot('ring2', 'ring2')}
            {renderEquipmentSlot('ring3', 'ring3')}
            {renderEquipmentSlot('ring4', 'ring4')}
          </div>
        </div>

        <div className="equipment-column-right">
          <div className="equipment-row-top">
            {renderEquipmentSlot('bag', 'bag')}
            {renderEquipmentSlot('throw', 'throw')}
          </div>
          <div className="equipment-column-right-items">
            {renderEquipmentSlot('arms', 'arms')}
            {renderEquipmentSlot('hands', 'hands')}
            {renderEquipmentSlot('shield', 'shield')}
            {renderEquipmentSlot('chest', 'chest')}
            {renderEquipmentSlot('belt', 'belt')}
            {renderEquipmentSlot('box', 'box')}
          </div>
        </div>
      </div>

      <div className="profile-stats-simple">
        {user?.gold !== undefined && (
          <div className="stat-row">
            <span>{user.gold || 0} зол.</span>
          </div>
        )}
        <div className="stat-row">
          <img src="/assets/images/attack.png" alt="Attack" className="stat-icon-simple" />
          <span>{user?.attack || 0}</span>
        </div>
        <div className="stat-row">
          <img src="/assets/images/defense.png" alt="Defense" className="stat-icon-simple" />
          <span>{user?.defense || 0}</span>
        </div>
        <div className="stat-row">
          <img src="/assets/images/hp.png" alt="HP" className="stat-icon-simple" />
          <span>{(user?.currentHp ?? user?.current_hp ?? 0)}/{user?.hp || 0}</span>
        </div>
        {user?.freeStats !== undefined && (
          <div className="stat-row">
            <span>Свободных статов: {user.freeStats || 0}</span>
          </div>
        )}
      </div>
    </div>
  )
}
