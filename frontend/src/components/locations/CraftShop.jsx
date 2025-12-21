import { useState, useEffect } from 'react'
import { useAuth } from '../../context/AuthContext'
import ToolItem from '../tool/Item'
import config from '../../config'

export default function CraftShop() {
  const { user } = useAuth()
  const [categories, setCategories] = useState([])
  const [currentCategoryIndex, setCurrentCategoryIndex] = useState(null)
  const [items, setItems] = useState([])
  const [playerSkill, setPlayerSkill] = useState(null)
  const [boughtItemId, setBoughtItemId] = useState(null)
  const [errorMessage, setErrorMessage] = useState('')
  const [successMessage, setSuccessMessage] = useState('')

  useEffect(() => {
    getCategoriesList()
  }, [])

  const getCategoriesList = async () => {
    try {
      const response = await fetch(`${config.apiUrl}/stuff/categories?category_type=tool`, {
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
    setCurrentCategoryIndex(index)
    const categoryName = categories[index]?.name?.toLowerCase()
    const skill = user?.[`${categoryName}_skill`] || 0
    setPlayerSkill(skill)
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
        body: JSON.stringify({ item_type: 'tool' }),
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
    return currentCategoryIndex === index ? 'active' : ''
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
          {currentCategoryIndex !== null && items.length > 0 && (
            <div className="items">
              {items.map((item) => (
                <div key={item.id} className="item row">
                  <ToolItem item={item} playerSkill={playerSkill} />
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

          {currentCategoryIndex === null && (
            <img
              src="/assets/images/locations/cities/moonshine/craft_shop/bg.jpg"
              alt="equipment"
              className="center background"
            />
          )}

          {currentCategoryIndex !== null && items.length === 0 && (
            <h3>Items will be soon</h3>
          )}
        </div>
      </div>
    </div>
  )
}

