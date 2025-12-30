import MapGrid from './MapGrid'
import './WaywardPines.css'

export default function WaywardPines() {
  return (
    <div className="location-inner-content">
      <div className="wayward-pines-header">
        <img
          src="/assets/images/locations/wayward_pines/icon.png"
          alt="Wayward Pines"
          className="wayward-pines-icon"
        />
        <h2>Wayward Pines</h2>
      </div>
      <MapGrid locationSlug="wayward_pines" />
    </div>
  )
}