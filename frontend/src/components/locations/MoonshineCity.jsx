import { Link } from 'react-router-dom'

export default function MoonshineCity() {
  return (
    <div className="moonshine-city">
      <div className="location-info">
        <h2>Welcome to Moonshine</h2>
        <p>This is the starting city. Explore the shops and locations around you.</p>
      </div>
      
      <div className="location-actions">
        <h3>Available Locations:</h3>
        <ul>
          <li><Link to="/locations/weapon_shop">Weapon Shop</Link></li>
          <li><Link to="/locations/craft_shop">Craft Shop</Link></li>
          <li><Link to="/locations/wayward_pines">Wayward Pines</Link></li>
        </ul>
      </div>
    </div>
  )
}

