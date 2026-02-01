import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api, tokenStore } from '../api/client'

const initialForm = {
  name: '',
  email: '',
  phone: '',
  password: '',
  role: 'client',
}

export default function Auth() {
  const navigate = useNavigate()
  const [form, setForm] = useState(initialForm)
  const [mode, setMode] = useState('login')
  const [message, setMessage] = useState('')

  const handleChange = (event) => {
    const { name, value } = event.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setMessage('')
    try {
      const response =
        mode === 'login'
          ? await api.login({ email: form.email, password: form.password })
          : await api.register(form)
      tokenStore.set(response.token)
      setMessage('Success! Redirecting...')
      
      // Trigger auth change event for Layout component
      window.dispatchEvent(new Event('auth-change'))
      
      // Redirect based on role
      setTimeout(() => {
        if (response.role === 'worker') {
          navigate('/dashboard')
        } else {
          navigate('/')
        }
      }, 1000)
    } catch (err) {
      setMessage(err.message)
    }
  }

  return (
    <section>
      <div style={{ maxWidth: '500px', margin: '0 auto' }}>
        <div className="text-center mb-4">
          <h1>{mode === 'login' ? 'Welcome Back!' : 'Join LocalConnect'}</h1>
          <p style={{ color: '#6b7280', marginTop: '0.5rem' }}>
            {mode === 'login' 
              ? 'Sign in to your account to continue' 
              : 'Create an account to get started'}
          </p>
        </div>
        
        {message && <p className={message.includes('success') ? 'notice' : 'error'}>{message}</p>}
        
        <form className="card" onSubmit={handleSubmit}>
          {mode === 'register' && (
            <>
              <label>
                Full Name
                <input 
                  name="name" 
                  placeholder="Enter your full name" 
                  value={form.name} 
                  onChange={handleChange}
                  required 
                />
              </label>
              <label>
                Phone Number
                <input 
                  name="phone" 
                  placeholder="Enter your phone number" 
                  value={form.phone} 
                  onChange={handleChange}
                  required 
                />
              </label>
              <label>
                Account Type
                <select name="role" value={form.role} onChange={handleChange}>
                  <option value="client">👤 Client - Looking for professionals</option>
                  <option value="worker">💼 Professional - Offering services</option>
                </select>
              </label>
            </>
          )}
          <label>
            Email Address
            <input 
              name="email" 
              type="email"
              placeholder="Enter your email" 
              value={form.email} 
              onChange={handleChange}
              required 
            />
          </label>
          <label>
            Password
            <input
              name="password"
              type="password"
              placeholder="Enter your password"
              value={form.password}
              onChange={handleChange}
              required
            />
          </label>
          <button type="submit" style={{ width: '100%', marginTop: '1rem' }}>
            {mode === 'login' ? '🔐 Login' : '🚀 Create Account'}
          </button>
        </form>
        
        <div className="text-center">
          <button className="link" onClick={() => setMode(mode === 'login' ? 'register' : 'login')}>
            {mode === 'login' ? '➡️ Need an account? Register here' : '⬅️ Have an account? Login here'}
          </button>
        </div>
      </div>
    </section>
  )
}
