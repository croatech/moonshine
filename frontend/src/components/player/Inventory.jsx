import Profile from './Profile'
import Items from './Items'

export default function Inventory({ player }) {
  return (
    <div className="show inventory">
      <div className="col-md-5 grid col-md-offset-1">
        <Profile player={player} />
      </div>

      <div className="col-md-5 col-md-offset-1">
        <Items player={player} />
      </div>
    </div>
  )
}
















