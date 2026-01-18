export default function ToolItem({ item, playerSkill }) {
  if (!item) return null

  const checkForNotAvailableSkill = (itemSkill, playerSkill) => {
    return itemSkill > playerSkill ? 'red' : ''
  }

  return (
    <div className="row top">
      <div className="col-md-3">
        <img src={item.image?.url} alt={item.name} />
      </div>

      <div className="col-md-9">
        <div className={`level ${checkForNotAvailableSkill(item.required_skill, playerSkill)}`}>
          [{item.required_skill}]
        </div>

        <h3>{item.name}</h3>
      </div>
    </div>
  )
}
















