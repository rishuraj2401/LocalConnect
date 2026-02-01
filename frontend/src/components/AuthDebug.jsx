import { useState, useEffect } from 'react'
import { tokenStore } from '../api/client'

export default function AuthDebug() {
  const [tokenData, setTokenData] = useState(null)
  
  useEffect(() => {
    const checkToken = () => {
      const token = tokenStore.get()
      if (token) {
        try {
          const payload = JSON.parse(atob(token.split('.')[1]))
          setTokenData(payload)
        } catch (e) {
          setTokenData({ error: e.message })
        }
      } else {
        setTokenData(null)
      }
    }
    
    checkToken()
    window.addEventListener('auth-change', checkToken)
    
    return () => {
      window.removeEventListener('auth-change', checkToken)
    }
  }, [])
  
  if (!tokenData) {
    return (
      <div style={{
        position: 'fixed',
        bottom: '1rem',
        right: '1rem',
        background: 'rgba(239, 68, 68, 0.9)',
        color: 'white',
        padding: '1rem',
        borderRadius: '8px',
        fontSize: '12px',
        maxWidth: '300px',
        zIndex: 9999,
        boxShadow: '0 4px 12px rgba(0,0,0,0.3)'
      }}>
        ❌ Not logged in
      </div>
    )
  }
  
  return (
    <div style={{
      position: 'fixed',
      bottom: '1rem',
      right: '1rem',
      background: 'rgba(16, 185, 129, 0.9)',
      color: 'white',
      padding: '1rem',
      borderRadius: '8px',
      fontSize: '12px',
      maxWidth: '300px',
      zIndex: 9999,
      boxShadow: '0 4px 12px rgba(0,0,0,0.3)'
    }}>
      <div><strong>✅ Logged in as:</strong></div>
      <div>Name: {tokenData.name || 'N/A'}</div>
      <div>Email: {tokenData.email || 'N/A'}</div>
      <div>Role: {tokenData.role || 'N/A'}</div>
      <div>User ID: {tokenData.user_id?.substring(0, 8) || 'N/A'}...</div>
    </div>
  )
}
