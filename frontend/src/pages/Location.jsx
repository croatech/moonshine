import { useParams } from 'react-router-dom'
import LocationView from '../components/locations/LocationView'
import MoonshineCity from '../components/locations/MoonshineCity'
import WeaponShop from '../components/locations/WeaponShop'
import CraftShop from '../components/locations/CraftShop'
import WaywardPines from '../components/locations/WaywardPines'

export default function Location() {
  const { slug } = useParams()

  const renderLocationContent = () => {
    switch (slug) {
      case 'moonshine':
        return <MoonshineCity />
      case 'weapon_shop':
        return <WeaponShop />
      case 'craft_shop':
        return <CraftShop />
      case 'wayward_pines':
        return <WaywardPines />
      default:
        return <div>Unknown location: {slug}</div>
    }
  }

  return (
    <LocationView slug={slug}>
      {renderLocationContent()}
    </LocationView>
  )
}

