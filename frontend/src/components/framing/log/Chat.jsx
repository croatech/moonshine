import { useState } from 'react'
import { useAuth } from '../../../context/AuthContext'
import config from '../../../config'

export default function Chat({ messages = [], recipient = null, onSetRecipient, onRemoveRecipient }) {
  const { user } = useAuth()
  const [messageText, setMessageText] = useState('')

  if (!user) return null

  const submitMessage = async () => {
    if (!messageText.trim()) return

    try {
      const response = await fetch(`${config.apiUrl}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({
          message: {
            text: messageText,
            recipient_id: recipient?.id || null,
          },
        }),
      })

      if (!response.ok) {
        console.error('Failed to send message')
        return
      }

      setMessageText('')
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      submitMessage()
    }
  }

  const setRecipient = (player) => {
    if (user.id === player.id) return
    onSetRecipient(player)
  }

  const isPrivateMessage = (message) => {
    if (!message.recipient) return false
    return message.recipient.id === user.id || message.player.id === user.id
  }

  const messageBelongsToRecipient = (message) => {
    if (!message.recipient || isPrivateMessage(message)) {
      return true
    }
    return false
  }

  const setPrivateClass = (message) => {
    if (isPrivateMessage(message)) {
      return 'private'
    }
    return ''
  }

  return (
    <div className="chat">
      <div className="input row">
        <input
          type="text"
          value={messageText}
          onChange={(e) => setMessageText(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="your message here"
          className="message-input form-control"
          required
        />
        {recipient && (
          <div className="recipient">
            <p className="for">Message for</p>
            <b>{recipient.name}</b>
            <p className="remove" onClick={onRemoveRecipient}>x</p>
          </div>
        )}
      </div>

      <div className="messages row">
        {messages
          .filter(messageBelongsToRecipient)
          .map((message, index) => (
            <div key={index} className="message">
              <b className="name" onClick={() => setRecipient(message.player)}>
                {message.player.name}:
              </b>
              {isPrivateMessage(message) && (
                <i className="recipient" onClick={() => setRecipient(message.recipient)}>
                  to {message.recipient.name}
                </i>
              )}
              <p className={`text ${setPrivateClass(message)}`}>{message.text}</p>
            </div>
          ))}
      </div>
    </div>
  )
}














