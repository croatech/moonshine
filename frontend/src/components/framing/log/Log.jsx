import { useState } from 'react'
import Events from './Events'
import Chat from './Chat'

export default function Log({ events = [], messages = [], recipient = null, onSetRecipient, onRemoveRecipient }) {
  const [inactiveComponent, setInactiveComponent] = useState('events')

  const switchComponent = () => {
    setInactiveComponent(inactiveComponent === 'events' ? 'chat' : 'events')
  }

  return (
    <div className="log col-md-8">
      <button onClick={switchComponent} className="btn btn-default switch">
        Switch to {inactiveComponent === 'events' ? 'chat' : 'events'}
      </button>
      {inactiveComponent !== 'events' && <Events events={events} />}
      {inactiveComponent !== 'chat' && (
        <Chat
          messages={messages}
          recipient={recipient}
          onSetRecipient={onSetRecipient}
          onRemoveRecipient={onRemoveRecipient}
        />
      )}
    </div>
  )
}











