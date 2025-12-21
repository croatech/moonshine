import { useState, useEffect } from 'react'
import { useAuth } from '../../context/AuthContext'
import config from '../../config'

export default function WaywardPines() {
  const { user } = useAuth()
  const [location, setLocation] = useState({})
  const [cells, setCells] = useState([])
  const [cellsLoaded, setCellsLoaded] = useState(false)

  useEffect(() => {
    getLocation()
  }, [])

  const getLocation = async () => {
    try {
      const response = await fetch(`${config.apiUrl}/locations/wayward_pines`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })

      if (!response.ok) {
        console.error('Failed to fetch location')
        return
      }

      const data = await response.json()
      setLocation(data)
      setCells(data.children || [])
      
      // Set cells loaded after a delay
      setTimeout(() => {
        setCellsLoaded(true)
        // TODO: Subscribe to movement stream if player has active movement
      }, 2000)
    } catch (error) {
      console.error('Error fetching location:', error)
    }
  }

  const checkForInactiveStyle = (cell) => {
    return cell.inactive ? 'inactive' : ''
  }

  const checkForCurrentCell = (cell) => {
    return cell.id === user?.location_id
  }

  const checkForMovementCell = (cell) => {
    if (!user?.active_movement) return false
    const pathArray = user.active_movement.path || []
    const indexOfCurrentCell = pathArray.indexOf(String(user.location_id))
    const indexOfNextCell = pathArray.indexOf(String(cell.id))
    return (
      pathArray.includes(String(cell.id)) &&
      user.location_id !== cell.id &&
      indexOfCurrentCell < indexOfNextCell
    )
  }

  const checkForLastMovementCell = (cell) => {
    if (!user?.active_movement) return false
    return String(cell.id) === user.active_movement.path[user.active_movement.path.length - 1]
  }

  const determineMapClass = () => {
    return !cellsLoaded ? 'invisible' : ''
  }

  const determineMovementCellClass = (cell) => {
    if (!user?.active_movement) return 'none'
    const pathArray = user.active_movement.path || []
    const indexOfCell = pathArray.indexOf(String(cell.id))
    const nextLocation = pathArray[indexOfCell + 1]
    if (nextLocation == null) {
      return 'none'
    } else {
      return cell.id - nextLocation
    }
  }

  const changeCell = async (cell) => {
    try {
      const response = await fetch(`${config.apiUrl}/cells/${cell.id}/move`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })

      if (!response.ok) {
        console.error('Failed to move to cell')
        return
      }

      const data = await response.json()
      // TODO: Update user data in context with new location
      // TODO: Subscribe to movement stream
    } catch (error) {
      console.error('Error moving to cell:', error)
    }
  }

  return (
    <div className="wayward-pines">
      <h1 className="text-center">{location.name}</h1>

      {user?.active_movement && (
        <div className="center">
          <i>Walking...</i>
        </div>
      )}

      {!cellsLoaded && (
        <img
          src="/assets/images/locations/spinner.gif"
          alt="spinner"
          className="spinner center"
        />
      )}

      <div className={`map center ${determineMapClass()}`}>
        <div className="cells">
          {cells.map((cell) => (
            <div
              key={cell.id}
              className={`cell ${checkForInactiveStyle(cell)}`}
              onClick={() => changeCell(cell)}
              style={{ cursor: 'pointer' }}
            >
              <img
                src={`/assets/images/locations/wayward_pines/cells/${cell.slug}.png`}
                alt={cell.slug}
              />

              {checkForCurrentCell(cell) && (
                <div className="current_position">
                  <img src="/assets/images/locations/point.png" alt="point" />
                </div>
              )}

              {checkForMovementCell(cell) && (
                <div className="arrow">
                  <img
                    src="/assets/images/locations/movements/arrow.png"
                    alt=""
                    className={`position-${determineMovementCellClass(cell)}`}
                  />
                  {checkForLastMovementCell(cell) && (
                    <img
                      src="/assets/images/locations/movements/target.png"
                      alt=""
                    />
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
