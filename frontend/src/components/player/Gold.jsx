import { useAuth } from '../../context/AuthContext'
import GoldIcon from '../GoldIcon'

export default function Gold() {
  const { user } = useAuth()
  
  if (!user) return null

  return (
    <div className="gold">
      <GoldIcon width={32} height={32} />
      {user.gold || 0}
    </div>
  )
}
















