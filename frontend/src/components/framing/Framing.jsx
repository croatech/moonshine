import { useState } from 'react'
import Log from './log/Log'
import OnlinePlayers from './OnlinePlayers'

export default function Framing({ events = [], messages = [] }) {
  const [recipient, setRecipient] = useState(null)

  const handleSetRecipient = (player) => {
    setRecipient(player)
  }

  const handleRemoveRecipient = () => {
    setRecipient(null)
  }

  return (
    <div className="frame">
      <Log
        events={events}
        messages={messages}
        recipient={recipient}
        onSetRecipient={handleSetRecipient}
        onRemoveRecipient={handleRemoveRecipient}
      />
      <OnlinePlayers onSetRecipient={handleSetRecipient} />
    </div>
  )
}

