import { useEffect, useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { locationAPI } from '../../lib/api'
import './MapGrid.css'

export default function MapGrid({ locationSlug }) {
  const { user, refetchUser } = useAuth()
  const [cells, setCells] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [moving, setMoving] = useState(false)

  useEffect(() => {
    const loadCells = async () => {
      try {
        setLoading(true)
        const data = await locationAPI.getCells(locationSlug)
        setCells(data.cells || [])
        setError(null)
      } catch (err) {
        console.error('[MapGrid] Error loading cells:', err)
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }

    if (locationSlug) {
      loadCells()
    }
  }, [locationSlug])

  useEffect(() => {
    if (!locationSlug) return

    const intervalId = setInterval(() => {
      refetchUser()
    }, 2000)

    return () => clearInterval(intervalId)
  }, [locationSlug, refetchUser])

  const handleCellClick = async (e, cellSlug) => {
    e.preventDefault()
    
    if (moving || user?.locationSlug === cellSlug) {
      return
    }

    try {
      setMoving(true)
      await locationAPI.moveToCell(locationSlug, cellSlug)
      await refetchUser()
    } catch (err) {
      console.error('[MapGrid] Error moving to cell:', err)
      alert(err.message || 'Ошибка при перемещении')
    } finally {
      setMoving(false)
    }
  }

  if (loading) {
    return <div className="map-grid-loading">Загрузка карты...</div>
  }

  if (error) {
    return <div className="map-grid-error">Ошибка загрузки: {error}</div>
  }

  const gridSize = 8
  const totalCells = gridSize * gridSize
  const cellMap = new Map()
  cells.forEach((cell) => {
    if (!cell || !cell.slug) {
      return
    }
    const match = cell.slug.match(/^(\d+)cell$/)
    if (match) {
      const cellNum = parseInt(match[1], 10)
      cellMap.set(cellNum, cell)
    }
  })

  return (
    <div className="map-grid">
      <div className="map-grid-container">
        {Array.from({ length: totalCells }, (_, index) => {
          const cellNum = index + 1
          const cell = cellMap.get(cellNum)
          const isPlayerHere = user?.locationSlug === cell?.slug

          if (!cell) {
            return (
              <div key={index} className="map-cell map-cell-empty" />
            )
          }

          return (
            <div key={cell.id} className="map-cell-wrapper">
              <div
                onClick={(e) => handleCellClick(e, cell.slug)}
                className={`map-cell ${cell.inactive ? 'map-cell-inactive' : ''} ${moving ? 'map-cell-moving' : ''}`}
                title={cell.name}
                style={{ cursor: moving ? 'wait' : 'pointer' }}
              >
                {cell.image && (
                  <img
                    src={`/assets/images/locations/${cell.image}`}
                    alt={cell.name}
                    className="map-cell-image"
                  />
                )}
              </div>
              {isPlayerHere && (
                <img
                  src="/assets/images/warrior.png"
                  alt="Персонаж"
                  className="map-cell-player-icon"
                />
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

