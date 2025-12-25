import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { userAPI, equipmentAPI, avatarAPI } from '../lib/api'
import PlayerHeader from '../components/PlayerHeader'
import './Profile.css'
import './EquipmentItems.css'

export default function Profile() {
  const { user: authUser, logout } = useAuth()
  const navigate = useNavigate()
  const [user, setUser] = useState(authUser)
  const [loading, setLoading] = useState(!authUser)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('inventory') // 'inventory' or 'settings'
  const [inventory, setInventory] = useState([])
  const [inventoryLoading, setInventoryLoading] = useState(false)
  const [equippedItems, setEquippedItems] = useState({}) // Map of slot -> item
  const [avatars, setAvatars] = useState([])
  const [avatarsLoading, setAvatarsLoading] = useState(false)
  const [notification, setNotification] = useState(null) // { message, type: 'error' | 'success' }

  // Auto-hide notification after 3 seconds
  useEffect(() => {
    if (notification) {
      const timer = setTimeout(() => {
        setNotification(null)
      }, 3000)
      return () => clearTimeout(timer)
    }
  }, [notification])

  const showNotification = (message, type = 'error') => {
    setNotification({ message, type })
  }

  useEffect(() => {
    // Always fetch fresh user data for profile page
    setLoading(true)
    userAPI.getCurrentUser()
      .then((userData) => {
        console.log('[Profile] User data loaded:', userData)
        setUser(userData)
        setLoading(false)
        
        // Load equipped items separately (non-blocking)
        userAPI.getEquippedItems()
          .then((equipped) => {
            console.log('[Profile] Equipped items loaded:', equipped)
            setEquippedItems(equipped)
          })
          .catch((err) => {
            console.error('[Profile] Error loading equipped items:', err)
            // Continue without equipped items
            setEquippedItems({})
          })
      })
      .catch((err) => {
        console.error('[Profile] Error loading profile:', err)
        setError('Ошибка загрузки профиля')
        setLoading(false)
      })
  }, [])

  useEffect(() => {
    // Load inventory when inventory tab is active
    if (activeTab === 'inventory') {
      setInventoryLoading(true)
      userAPI.getInventory()
        .then((items) => {
          console.log('[Profile] Inventory loaded:', items)
          setInventory(items)
          setInventoryLoading(false)
        })
        .catch((err) => {
          console.error('[Profile] Error loading inventory:', err)
          setInventory([])
          setInventoryLoading(false)
        })
    } else if (activeTab === 'settings') {
      // Load avatars when settings tab is active
      if (avatars.length === 0) {
        setAvatarsLoading(true)
        avatarAPI.getAll()
          .then((avatarsList) => {
            console.log('[Profile] Avatars loaded:', avatarsList)
            setAvatars(avatarsList)
            setAvatarsLoading(false)
          })
          .catch((err) => {
            console.error('[Profile] Error loading avatars:', err)
            setAvatars([])
            setAvatarsLoading(false)
          })
      }
    }
  }, [activeTab])

  if (loading) return <div className="profile-container">Загрузка...</div>
  if (error) return <div className="profile-container">{error}</div>
  if (!user) return <div className="profile-container">Пользователь не найден</div>

  const handleLogout = () => {
    logout()
    localStorage.clear()
    navigate('/signin')
  }

  const handleBack = () => {
    navigate(-1)
  }

  const handleTakeOn = async (item) => {
    if (!item.slug) {
      showNotification('Slug предмета не определен', 'error')
      return
    }

    // Check level requirement
    if (user.level < item.requiredLevel) {
      showNotification(`Недостаточный уровень! Требуется уровень ${item.requiredLevel}`, 'error')
      return
    }

    try {
      await equipmentAPI.takeOn(item.slug)
      // Refresh user data, inventory, and equipped items
      const [updatedUser, items, equipped] = await Promise.all([
        userAPI.getCurrentUser(),
        userAPI.getInventory(),
        userAPI.getEquippedItems()
      ])
      setUser(updatedUser)
      setInventory(items)
      setEquippedItems(equipped)
      // Removed success alert
    } catch (error) {
      console.error('[Profile] Error equipping item:', error)
      let errorMessage = 'Неизвестная ошибка'
      if (error.message.includes('not in inventory')) {
        errorMessage = 'Предмет не в инвентаре'
      } else if (error.message.includes('insufficient level')) {
        errorMessage = 'Недостаточный уровень'
      } else if (error.message.includes('not found')) {
        errorMessage = 'Предмет не найден'
      } else if (error.message.includes('invalid equipment type')) {
        errorMessage = 'Неверный тип экипировки'
      } else {
        errorMessage = error.message
      }
      showNotification(errorMessage, 'error')
    }
  }

  const handleTakeOff = async (slotName) => {
    // Check if there's an item in this slot
    const equippedItem = equippedItems[slotName] || equippedItems[slotName.toLowerCase()]
    if (!equippedItem) {
      return // No item to remove
    }

    try {
      await equipmentAPI.takeOff(slotName)
      // Refresh user data, inventory, and equipped items
      const [updatedUser, items, equipped] = await Promise.all([
        userAPI.getCurrentUser(),
        userAPI.getInventory(),
        userAPI.getEquippedItems()
      ])
      setUser(updatedUser)
      setInventory(items)
      setEquippedItems(equipped)
      // Removed success alert
    } catch (error) {
      console.error('[Profile] Error removing item:', error)
      let errorMessage = 'Неизвестная ошибка'
      if (error.message.includes('no item equipped')) {
        errorMessage = 'В этом слоте нет предмета'
      } else if (error.message.includes('invalid slot')) {
        errorMessage = 'Неверный слот'
      } else {
        errorMessage = error.message
      }
      showNotification(errorMessage, 'error')
    }
  }

  const handleSelectAvatar = async (avatarId) => {
    try {
      const updatedUser = await userAPI.updateProfile({ avatarId })
      setUser(updatedUser)
      // Removed success alert
    } catch (error) {
      console.error('[Profile] Error updating avatar:', error)
      let errorMessage = 'Неизвестная ошибка'
      if (error.message.includes('not found')) {
        errorMessage = 'Аватар не найден'
      } else if (error.message.includes('Unauthorized')) {
        errorMessage = 'Необходимо войти в систему'
      } else {
        errorMessage = error.message
      }
      showNotification(errorMessage, 'error')
    }
  }

  const handleSell = async (item) => {
    if (!item.slug) {
      showNotification('Slug предмета не определен', 'error')
      return
    }

    try {
      await equipmentAPI.sell(item.slug)
      // Refresh user data and inventory
      const [updatedUser, items] = await Promise.all([
        userAPI.getCurrentUser(),
        userAPI.getInventory()
      ])
      setUser(updatedUser)
      setInventory(items)
      showNotification(`Продано за ${item.price} золота`, 'success')
    } catch (error) {
      console.error('[Profile] Error selling item:', error)
      let errorMessage = 'Неизвестная ошибка'
      if (error.message.includes('not owned')) {
        errorMessage = 'У вас нет этого предмета'
      } else if (error.message.includes('not found')) {
        errorMessage = 'Предмет не найден'
      } else {
        errorMessage = error.message
      }
      showNotification(errorMessage, 'error')
    }
  }

  // Normalize image path helper
  const normalizeImagePath = (img) => {
    if (!img) return null
    let p = img
    if (p.startsWith('/')) p = p.slice(1)
    p = p.replace(/^frontend\/assets\/images\//, '')
    if (p.startsWith('assets/images/')) p = p.replace(/^assets\/images\//, '')
    return `/assets/images/${p}`
  }

  // Normalize image path for equipment items (same as in EquipmentItems.jsx)
  const normalizeEquipmentImagePath = (img) => {
    if (!img) return null
    let p = img
    if (p.startsWith('/')) p = p.slice(1)
    p = p.replace(/^frontend\/assets\/images\//, '')
    if (p.startsWith('assets/images/')) p = p.replace(/^assets\/images\//, '')
    return `/assets/images/${p}`
  }

  // Get equipment slot image or placeholder
  const getEquipmentSlotImage = (slotName) => {
    // Check if there's an equipped item for this slot
    const equippedItem = equippedItems[slotName] || equippedItems[slotName.toLowerCase()]
    
    if (equippedItem && equippedItem.image) {
      return normalizeEquipmentImagePath(equippedItem.image)
    }
    
    // Return placeholder from grid
    // For rings (ring1, ring2, etc), use "ring" placeholder
    const placeholderName = slotName.startsWith('ring') ? 'ring' : slotName
    return `/assets/images/equipment_items/grid/${placeholderName}.png`
  }

  // Check if slot has an item equipped
  const hasEquippedItem = (slotName) => {
    const equippedItem = equippedItems[slotName] || equippedItems[slotName.toLowerCase()]
    return equippedItem && equippedItem.image
  }

  // Render equipment slot with click handler
  const renderEquipmentSlot = (slotName, alt) => {
    const hasItem = hasEquippedItem(slotName)
    
    return (
      <div 
        className={`equipment-slot ${hasItem ? 'has-item' : ''}`}
        onClick={() => hasItem && handleTakeOff(slotName)}
        style={{ cursor: hasItem ? 'pointer' : 'default' }}
        title={hasItem ? 'Нажмите чтобы снять' : ''}
      >
        <img src={getEquipmentSlotImage(slotName)} alt={alt} />
      </div>
    )
  }

  // Get avatar image
  const avatarImage = user.avatar?.image
  const avatarSrc = normalizeImagePath(avatarImage) || getEquipmentSlotImage('head') // fallback

  return (
    <div className="profile-container">
      {/* Notification toast */}
      {notification && (
        <div className={`notification-toast notification-${notification.type}`}>
          {notification.message}
        </div>
      )}

      <div className="profile-main-block">
        <div className="profile-header">
          <PlayerHeader 
            user={user} 
            fullWidth={true}
          />
          <div className="profile-header-actions">
            <button 
              onClick={handleBack}
              className="profile-back-button"
              title="Назад"
            >
              ← Назад
            </button>
            <button 
              onClick={handleLogout} 
              className="profile-logout-button"
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
        <div className="profile-content">
          {/* Equipment grid and stats wrapper */}
          <div className="profile-equipment-wrapper">
              {/* Equipment grid - всегда показываем */}
              <div className="equipment-grid">
            {/* Left column: head, neck, weapon, legs, feet */}
            <div className="equipment-column-left">
              {renderEquipmentSlot('head', 'head')}
              {renderEquipmentSlot('neck', 'neck')}
              {renderEquipmentSlot('weapon', 'weapon')}
              {renderEquipmentSlot('legs', 'legs')}
              {renderEquipmentSlot('feet', 'feet')}
            </div>

            {/* Center: Avatar with rings below */}
            <div className="equipment-column-center">
              <div className="equipment-avatar">
                {avatarSrc ? (
                  <img 
                    src={avatarSrc} 
                    alt={user.username}
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

            {/* Right: bag and throw horizontally, then arms, hands, shield, chest, belt, box vertically */}
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

              {/* Stats - вертикально */}
              <div className="profile-stats-simple">
                <div className="stat-row">
                  <img src="/assets/images/attack.png" alt="Attack" className="stat-icon-simple" />
                  <span>{user.attack || 0}</span>
                </div>
                <div className="stat-row">
                  <img src="/assets/images/defense.png" alt="Defense" className="stat-icon-simple" />
                  <span>{user.defense || 0}</span>
                </div>
                <div className="stat-row">
                  <img src="/assets/images/hp.png" alt="HP" className="stat-icon-simple" />
                  <span>{user.hp || 0}</span>
                </div>
                <div className="stat-row">
                  <span>Свободных статов: {user.freeStats || 0}</span>
                </div>
              </div>
            </div>

          {/* Right sidebar with tabs */}
        <div className="profile-sidebar">
          <div className="profile-tabs">
            <button
              className={`profile-tab ${activeTab === 'inventory' ? 'active' : ''}`}
              onClick={() => setActiveTab('inventory')}
            >
              Инвентарь
            </button>
            <button
              className={`profile-tab ${activeTab === 'settings' ? 'active' : ''}`}
              onClick={() => setActiveTab('settings')}
            >
              Настройки
            </button>
          </div>

          {/* Inventory tab content */}
          {activeTab === 'inventory' && (
            <div className="profile-inventory-content">
              {inventoryLoading ? (
                <div>Загрузка инвентаря...</div>
              ) : (
                <div className="equipment-items-list">
                  {inventory.length === 0 ? (
                    <p>Инвентарь пуст</p>
                  ) : (
                    inventory.map((item) => (
                      <div key={item.id} className="equipment-item-card">
                        {item.image && (
                          <img 
                            src={normalizeEquipmentImagePath(item.image)} 
                            alt={item.name}
                            className="equipment-item-image"
                          />
                        )}
                        <div className="equipment-item-info">
                          <h3>{item.name}</h3>
                          <div className="equipment-item-stats">
                            <div>Уровень: {item.requiredLevel}</div>
                            {item.attack > 0 && <div>Атака: {item.attack}</div>}
                            {item.defense > 0 && <div>Защита: {item.defense}</div>}
                            {item.hp > 0 && <div>HP: {item.hp}</div>}
                          </div>
                          <div className="equipment-item-buttons">
                            <button 
                              className="equipment-item-equip-button"
                              onClick={() => handleTakeOn(item)}
                              disabled={user.level < item.requiredLevel}
                            >
                              {user.level < item.requiredLevel ? `Нужен ${item.requiredLevel} ур.` : 'Надеть'}
                            </button>
                            <button 
                              className="equipment-item-sell-button"
                              onClick={() => handleSell(item)}
                            >
                              Продать ({item.price}💰)
                            </button>
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>
          )}

          {/* Settings tab content */}
          {activeTab === 'settings' && (
            <div className="profile-settings-content">
              <h2>Выбор аватара</h2>
              {avatarsLoading ? (
                <div>Загрузка аватаров...</div>
              ) : (
                <div className="avatars-list">
                  {avatars.length === 0 ? (
                    <p>Аватары не найдены</p>
                  ) : (
                    avatars.map((avatar) => {
                      const isSelected = user.avatar?.id === avatar.id
                      return (
                        <div
                          key={avatar.id}
                          className={`avatar-item ${isSelected ? 'selected' : ''}`}
                          onClick={() => handleSelectAvatar(avatar.id)}
                        >
                          <img 
                            src={normalizeImagePath(avatar.image)} 
                            alt={`Avatar ${avatar.id}`}
                            className="avatar-image-select"
                            onError={(e) => {
                              e.target.style.display = 'none'
                            }}
                          />
                          {isSelected && <div className="avatar-selected-badge">✓</div>}
                        </div>
                      )
                    })
                  )}
                </div>
              )}
            </div>
          )}
        </div>
        </div>
      </div>
    </div>
  )
}
