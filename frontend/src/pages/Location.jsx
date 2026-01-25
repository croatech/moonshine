import { useParams } from 'react-router-dom'
import LocationView from '../components/locations/LocationView'
import MoonshineCity from '../components/locations/MoonshineCity'
import WeaponShop from '../components/locations/WeaponShop'
import ArtifactsShop from '../components/locations/ArtifactsShop'
import WaywardPines from '../components/locations/WaywardPines'

export default function Location() {
  const { slug } = useParams()

  const renderLocationContent = () => {
    switch (slug) {
      case 'moonshine':
        return <MoonshineCity />
      case 'weapon_shop':
        return <WeaponShop />
      case 'shop_of_artifacts':
        return <ArtifactsShop />
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

