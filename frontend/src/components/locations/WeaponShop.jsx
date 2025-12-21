import { useState, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import EquipmentItem from '../equipment/Item'
import config from '../../config'

export default function WeaponShop() {
  const { user } = useAuth()
  const location = useLocation()
  const [categories, setCategories] = useState([])
  const [currentCategory, setCurrentCategory] = useState(null)
  const [items, setItems] = useState([])
  const [boughtItemId, setBoughtItemId] = useState(null)
  const [errorMessage, setErrorMessage] = useState('')
  const [successMessage, setSuccessMessage] = useState('')
  const [resourceName, setResourceName] = useState(null)

  useEffect(() => {
    getCategoriesList()
  }, [])

  const getCategoriesList = async () => {
    const pathname = location.pathname
    const match = pathname.match(/weapon_shop|shop_of_artifacts/gi)
    const resourceName = match ? match[0] : 'weapon_shop'
    setResourceName(resourceName)

    const categoryType = resourceName === 'weapon_shop' ? 'equipment' : 'artifacts'

    try {
      const response = await fetch(`${config.apiUrl}/stuff/categories?category_type=${categoryType}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })

      if (!response.ok) {
        console.error('Failed to fetch categories')
        return
      }

      const data = await response.json()
      setCategories(data || [])
    } catch (error) {
      console.error('Error fetching categories:', error)
    }
  }

  const showCategory = (index) => {
    setCurrentCategory(index)
    setItems(categories[index]?.items || [])
  }

  const buyItem = async (itemId) => {
    setBoughtItemId(itemId)
    setErrorMessage('')
    setSuccessMessage('')

    try {
      const response = await fetch(`${config.apiUrl}/stuff/items/${itemId}/buy`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ item_type: 'equipment' }),
      })

      if (!response.ok) {
        const error = await response.text()
        setErrorMessage(error)
        return
      }

      const data = await response.json()
      setSuccessMessage(`Congrats! You have bought ${data.name}`)
      // TODO: Update user gold in context
    } catch (error) {
      setErrorMessage('Failed to buy item')
      console.error('Error buying item:', error)
    }
  }

  const resolveCategoryClass = (index) => {
    return currentCategory === index ? 'active' : ''
  }

  return (
    <div className="shop">
      <div className="item-categories">
        <div className="col-md-3">
          {categories.map((category, i) => (
            <div key={category.id || i} className="category">
              <a
                onClick={() => showCategory(i)}
                className={`category-btn btn btn-default ${resolveCategoryClass(i)}`}
                style={{ cursor: 'pointer' }}
              >
                {category.name}
              </a>
            </div>
          ))}
        </div>

        <div className="col-md-9">
          {currentCategory !== null && items.length > 0 && (
            <div className="items">
              {items.map((item) => (
                <div key={item.id} className="item row">
                  <EquipmentItem item={item} playerLevel={user?.level || 0} />
                  <a
                    onClick={() => buyItem(item.id)}
                    className="buy-button btn btn-success"
                    style={{ cursor: 'pointer' }}
                  >
                    Buy for {item.price} gold
                  </a>
                  {errorMessage && boughtItemId === item.id && (
                    <div className="alert alert-danger" role="alert">
                      {errorMessage}
                    </div>
                  )}
                  {successMessage && boughtItemId === item.id && (
                    <div className="alert alert-success" role="alert">
                      {successMessage}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}

          {currentCategory === null && resourceName && (
            <img
              src={`/assets/images/locations/cities/moonshine/${resourceName}/bg.jpg`}
              alt="equipment"
              className="center background"
            />
          )}
        </div>
      </div>
    </div>
  )
}

