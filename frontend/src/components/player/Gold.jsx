import { useAuth } from '../../context/AuthContext'

export default function Gold() {
  const { user } = useAuth()
  
  if (!user) return null

  return (
    <div className="gold">
      <img src="/assets/images/gold.png" alt="gold" />
      {user.gold || 0}
    </div>
  )
}





