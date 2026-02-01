import { Link, useNavigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { tokenStore } from '../api/client'
import AuthDebug from './AuthDebug'

export default function Layout({ children }) {
  const navigate = useNavigate()
  const [user, setUser] = useState(null)
  const [showUserMenu, setShowUserMenu] = useState(false)

  const checkAuth = () => {
    const token = tokenStore.get()
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        setUser({
          id: payload.user_id,
          email: payload.email,
          role: payload.role,
          name: payload.name || payload.email.split('@')[0]
        })
      } catch (e) {
        console.error('Failed to decode token', e)
        setUser(null)
      }
    } else {
      setUser(null)
    }
  }

  useEffect(() => {
    checkAuth()
    
    // Listen for storage changes (when user logs in from another tab or Auth component)
    const handleStorageChange = () => {
      checkAuth()
    }
    
    // Close menu when clicking outside
    const handleClickOutside = (event) => {
      if (showUserMenu && !event.target.closest('.user-menu-container')) {
        setShowUserMenu(false)
      }
    }
    
    window.addEventListener('storage', handleStorageChange)
    window.addEventListener('auth-change', checkAuth)
    document.addEventListener('click', handleClickOutside)
    
    return () => {
      window.removeEventListener('storage', handleStorageChange)
      window.removeEventListener('auth-change', checkAuth)
      document.removeEventListener('click', handleClickOutside)
    }
  }, [showUserMenu])

  const handleLogout = () => {
    tokenStore.clear()
    setUser(null)
    setShowUserMenu(false)
    navigate('/auth')
  }

  return (
    <div className="app">
      <header className="header">
        <div className="brand">
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            LocalConnect
          </Link>
        </div>
        <nav className="nav">
          <Link to="/">Find Professionals</Link>
          {user && <Link to="/messages">Messages</Link>}
          {user && user.role === 'worker' && <Link to="/dashboard">Dashboard</Link>}
          
          {!user ? (
            <Link to="/auth" style={{
              background: 'linear-gradient(135deg, #10b981 0%, #0891b2 100%)',
              padding: '0.75rem 1.5rem',
              borderRadius: '12px',
              fontWeight: '700'
            }}>
              Login
            </Link>
          ) : (
            <div className="user-menu-container" style={{ position: 'relative', display: 'flex', alignItems: 'center', gap: '1rem' }}>
              {/* User Name and Avatar */}
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.75rem',
                padding: '0.5rem 1rem',
                background: 'rgba(16, 185, 129, 0.15)',
                borderRadius: '12px',
                border: '1px solid rgba(16, 185, 129, 0.3)'
              }}>
                <div style={{
                  width: '36px',
                  height: '36px',
                  borderRadius: '50%',
                  background: 'linear-gradient(135deg, #10b981 0%, #0891b2 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '1.25rem',
                  cursor: 'pointer'
                }}
                onClick={() => setShowUserMenu(!showUserMenu)}
                >
                  {user.role === 'worker' ? '💼' : '👤'}
                </div>
                <div style={{ 
                  display: 'flex', 
                  flexDirection: 'column',
                  cursor: 'pointer'
                }}
                onClick={() => setShowUserMenu(!showUserMenu)}
                >
                  <span style={{ 
                    fontSize: '14px', 
                    fontWeight: '600',
                    color: 'white'
                  }}>
                    {user.name}
                  </span>
                  <span style={{ 
                    fontSize: '11px', 
                    color: 'rgba(255, 255, 255, 0.6)',
                    textTransform: 'capitalize'
                  }}>
                    {user.role === 'worker' ? 'Professional' : 'Client'}
                  </span>
                </div>
              </div>
              
              {/* Logout Button */}
              <button
                onClick={handleLogout}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '0.5rem',
                  padding: '0.75rem 1.5rem',
                  background: 'rgba(239, 68, 68, 0.2)',
                  border: '1px solid rgba(239, 68, 68, 0.4)',
                  color: '#fca5a5',
                  fontSize: '15px',
                  fontWeight: '700',
                  borderRadius: '12px',
                  cursor: 'pointer',
                  transition: 'all 0.3s ease'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(239, 68, 68, 0.3)'
                  e.currentTarget.style.color = 'white'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'rgba(239, 68, 68, 0.2)'
                  e.currentTarget.style.color = '#fca5a5'
                }}
              >
                <span>🚪</span>
                <span>Logout</span>
              </button>
              
              {showUserMenu && (
                <div style={{
                  position: 'absolute',
                  top: 'calc(100% + 0.75rem)',
                  right: 0,
                  minWidth: '280px',
                  background: 'rgba(10, 15, 30, 0.98)',
                  backdropFilter: 'blur(24px)',
                  border: '1px solid rgba(16, 185, 129, 0.3)',
                  borderRadius: '16px',
                  boxShadow: '0 12px 40px rgba(0, 0, 0, 0.7)',
                  padding: '1rem',
                  zIndex: 1000,
                  animation: 'slideDown 0.2s ease'
                }}>
                  {/* User Info Header */}
                  <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '1rem',
                    padding: '1rem',
                    background: 'rgba(16, 185, 129, 0.1)',
                    borderRadius: '12px',
                    marginBottom: '1rem',
                    border: '1px solid rgba(16, 185, 129, 0.2)'
                  }}>
                    <div style={{
                      width: '48px',
                      height: '48px',
                      borderRadius: '50%',
                      background: 'linear-gradient(135deg, #10b981 0%, #0891b2 100%)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontSize: '1.5rem',
                      flexShrink: 0
                    }}>
                      {user.role === 'worker' ? '💼' : '👤'}
                    </div>
                    <div style={{ flex: 1, minWidth: 0 }}>
                      <div style={{ 
                        fontWeight: '700', 
                        color: 'white',
                        marginBottom: '0.25rem',
                        fontSize: '15px'
                      }}>
                        {user.name}
                      </div>
                      <div style={{ 
                        fontSize: '13px', 
                        color: 'rgba(255, 255, 255, 0.6)',
                        marginBottom: '0.5rem',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                      }}>
                        {user.email}
                      </div>
                      <div style={{
                        display: 'inline-block',
                        padding: '0.25rem 0.75rem',
                        background: user.role === 'worker' ? 'rgba(16, 185, 129, 0.3)' : 'rgba(8, 145, 178, 0.3)',
                        border: `1px solid ${user.role === 'worker' ? 'rgba(16, 185, 129, 0.5)' : 'rgba(8, 145, 178, 0.5)'}`,
                        borderRadius: '8px',
                        fontSize: '11px',
                        color: 'white',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em',
                        fontWeight: '700'
                      }}>
                        {user.role === 'worker' ? 'Professional' : 'Client'}
                      </div>
                    </div>
                  </div>
                  
                  {/* Quick Links */}
                  <div style={{ marginBottom: '0.5rem' }}>
                    {user.role === 'worker' && (
                      <Link
                        to="/dashboard"
                        onClick={() => setShowUserMenu(false)}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '0.75rem',
                          padding: '0.875rem 1rem',
                          color: 'white',
                          textDecoration: 'none',
                          borderRadius: '10px',
                          transition: 'all 0.2s',
                          fontSize: '14px',
                          fontWeight: '600',
                          marginBottom: '0.5rem'
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.background = 'rgba(16, 185, 129, 0.15)'
                          e.currentTarget.style.transform = 'translateX(4px)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.background = 'transparent'
                          e.currentTarget.style.transform = 'translateX(0)'
                        }}
                      >
                        <span style={{ fontSize: '1.25rem' }}>📊</span>
                        <span>Dashboard</span>
                      </Link>
                    )}
                    
                    <Link
                      to="/messages"
                      onClick={() => setShowUserMenu(false)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.75rem',
                        padding: '0.875rem 1rem',
                        color: 'white',
                        textDecoration: 'none',
                        borderRadius: '10px',
                        transition: 'all 0.2s',
                        fontSize: '14px',
                        fontWeight: '600'
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.background = 'rgba(16, 185, 129, 0.15)'
                        e.currentTarget.style.transform = 'translateX(4px)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.background = 'transparent'
                        e.currentTarget.style.transform = 'translateX(0)'
                      }}
                    >
                      <span style={{ fontSize: '1.25rem' }}>💬</span>
                      <span>Messages</span>
                    </Link>
                  </div>
                </div>
              )}
            </div>
          )}
        </nav>
      </header>
      <main className="content">{children}</main>
      <AuthDebug />
    </div>
  )
}
