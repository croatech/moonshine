import Chat from './Chat'

export default function Log({ messages = [], recipient = null, onSetRecipient, onRemoveRecipient }) {
  return (
    <div className="log col-md-8">
      <Chat
        messages={messages}
        recipient={recipient}
        onSetRecipient={onSetRecipient}
        onRemoveRecipient={onRemoveRecipient}
      />
    </div>
  )
}
















