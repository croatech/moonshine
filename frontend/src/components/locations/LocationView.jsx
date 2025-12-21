// Placeholder component - Location query will be added to GraphQL schema later
// For now, just render children with location name based on slug
export default function LocationView({ slug, children }) {
  const locationNames = {
    moonshine: 'Moonshine',
    weapon_shop: 'Weapon Shop',
    craft_shop: 'Craft Shop',
    wayward_pines: 'Wayward Pines',
  }

  const locationName = locationNames[slug] || slug

  // Background images mapping
  const locationBgs = {
    moonshine: 'cities/moonshine/bg.jpg',
    weapon_shop: 'cities/moonshine/weapon_shop/bg.jpg',
    craft_shop: 'cities/moonshine/craft_shop/bg.jpg',
    wayward_pines: 'wayward_pines/bg.png',
  }

  const bgImage = locationBgs[slug]

  return (
    <div className="location-view">
      <h1>{locationName}</h1>
      {bgImage && (
        <div 
          className="location-bg"
          style={{
            backgroundImage: `url(/assets/assets/images/locations/${bgImage})`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            minHeight: '400px',
            borderRadius: '8px',
            marginBottom: '20px',
          }}
        />
      )}
      <div className="location-content">
        {children}
      </div>
    </div>
  )
}

