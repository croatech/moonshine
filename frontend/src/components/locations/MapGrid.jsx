import { useEffect, useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { locationAPI } from '../../lib/api'
import './MapGrid.css'

export default function MapGrid({ locationSlug }) {
  const { user } = useAuth()
  const [cells, setCells] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

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

  const handleCellClick = (e, cellSlug) => {
    e.preventDefault()
    window.location.href = `/locations/${cellSlug}/move`
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
              <a
                href={`/locations/${cell.slug}/move`}
                onClick={(e) => handleCellClick(e, cell.slug)}
                className={`map-cell ${cell.inactive ? 'map-cell-inactive' : ''}`}
                title={cell.name}
              >
                {cell.image && (
                  <img
                    src={`/assets/images/locations/${cell.image}`}
                    alt={cell.name}
                    className="map-cell-image"
                  />
                )}
              </a>
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

