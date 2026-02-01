import { useEffect, useState } from 'react'
import { api } from '../api/client'
import './Messages.css'

export default function Messages() {
  const [conversations, setConversations] = useState([])
  const [selectedConversation, setSelectedConversation] = useState(null)
  const [messages, setMessages] = useState([])
  const [newMessage, setNewMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    loadConversations()
  }, [])

  useEffect(() => {
    if (selectedConversation) {
      loadMessages(selectedConversation.id)
    }
  }, [selectedConversation])

  const loadConversations = async () => {
    try {
      const data = await api.conversations()
      setConversations(data || [])
    } catch (err) {
      setError(err.message)
    }
  }

  const loadMessages = async (conversationId) => {
    try {
      const data = await api.messages(conversationId)
      setMessages(data || [])
    } catch (err) {
      setError(err.message)
    }
  }

  const handleSendMessage = async (e) => {
    e.preventDefault()
    if (!newMessage.trim() || !selectedConversation) return

    try {
      await api.sendMessage({
        receiver_id: selectedConversation.other_user_id,
        content: newMessage,
      })
      setNewMessage('')
      // Reload messages
      setTimeout(() => loadMessages(selectedConversation.id), 500)
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <section className="messages-container">
      <h1>Messages</h1>
      {error && <p className="error">{error}</p>}
      
      <div className="messages-layout">
        <div className="conversations-list">
          <h2>Conversations</h2>
          {conversations.length === 0 ? (
            <p className="empty">No conversations yet</p>
          ) : (
            conversations.map((conv) => (
              <div
                key={conv.id}
                className={`conversation-item ${selectedConversation?.id === conv.id ? 'active' : ''}`}
                onClick={() => setSelectedConversation(conv)}
              >
                <div className="conversation-header">
                  <strong>{conv.other_user_name}</strong>
                  {conv.unread_count > 0 && (
                    <span className="unread-badge">{conv.unread_count}</span>
                  )}
                </div>
                <p className="last-message">{conv.last_message}</p>
                <small>{new Date(conv.last_message_at).toLocaleString()}</small>
              </div>
            ))
          )}
        </div>

        <div className="messages-panel">
          {selectedConversation ? (
            <>
              <div className="messages-header">
                <h2>{selectedConversation.other_user_name}</h2>
              </div>
              
              <div className="messages-list">
                {messages.map((msg) => (
                  <div
                    key={msg.id}
                    className={`message ${msg.sender_id !== selectedConversation.other_user_id ? 'sent' : 'received'}`}
                  >
                    <p>{msg.content}</p>
                    <small>{new Date(msg.created_at).toLocaleString()}</small>
                  </div>
                ))}
              </div>

              <form className="message-input" onSubmit={handleSendMessage}>
                <input
                  type="text"
                  placeholder="Type a message..."
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                />
                <button type="submit">Send</button>
              </form>
            </>
          ) : (
            <div className="no-selection">
              <p>Select a conversation to view messages</p>
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
