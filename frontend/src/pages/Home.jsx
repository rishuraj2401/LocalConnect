import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export default function Home() {
  const [categories, setCategories] = useState([])
  const [profiles, setProfiles] = useState([])
  const [filters, setFilters] = useState({ category: '', location: '' })
  const [error, setError] = useState('')

  useEffect(() => {
    api.categories()
      .then((data) => setCategories(data || []))
      .catch((err) => {
        setError(err.message)
        setCategories([])
      })
  }, [])

  useEffect(() => {
    api.profiles(filters)
      .then((data) => setProfiles(data || []))
      .catch((err) => {
        setError(err.message)
        setProfiles([])
      })
  }, [filters])

  return (
    <section>
      {/* Hero Section - Premium Glass with Background */}
      <div style={{
        backgroundImage: 'linear-gradient(rgba(10, 15, 30, 0.90), rgba(10, 15, 30, 0.85)), url(https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=1920&q=80)',
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backdropFilter: 'blur(30px) saturate(180%)',
        WebkitBackdropFilter: 'blur(30px) saturate(180%)',
        borderRadius: 'var(--radius-lg)',
        padding: '4rem 3rem',
        marginBottom: '4rem',
        color: 'white',
        textAlign: 'center',
        boxShadow: '0 8px 32px 0 rgba(0, 0, 0, 0.6)',
        border: '1px solid rgba(16, 185, 129, 0.3)',
        position: 'relative',
        overflow: 'hidden'
      }}>
        <div style={{
          position: 'absolute',
          top: '-50%',
          left: '-50%',
          width: '200%',
          height: '200%',
          background: 'radial-gradient(circle, rgba(16, 185, 129, 0.15) 0%, transparent 70%)',
          animation: 'rotate 20s linear infinite',
          pointerEvents: 'none'
        }}></div>
        <h1 style={{ 
          fontSize: '3.5rem', 
          marginBottom: '1.5rem', 
          color: '#ffffff',
          fontWeight: '900',
          letterSpacing: '-0.02em',
          position: 'relative',
          zIndex: 1,
          textShadow: '0 4px 20px rgba(0, 0, 0, 0.8), 0 2px 8px rgba(0, 0, 0, 0.6)'
        }}>
          Find Trusted Local Professionals
        </h1>
        <p style={{ 
          fontSize: '1.5rem', 
          marginBottom: '3rem', 
          fontWeight: '400',
          color: '#ffffff',
          position: 'relative',
          zIndex: 1,
          textShadow: '0 2px 12px rgba(0, 0, 0, 0.7), 0 1px 4px rgba(0, 0, 0, 0.5)'
        }}>
          Connect with verified professionals in your area
        </p>
        <div style={{ 
          display: 'flex', 
          gap: '3rem', 
          justifyContent: 'center', 
          flexWrap: 'wrap', 
          marginTop: '2rem',
          position: 'relative',
          zIndex: 1
        }}>
          <div style={{ 
            textAlign: 'center',
            background: 'rgba(10, 15, 30, 0.85)',
            padding: '1.5rem 2rem',
            borderRadius: '16px',
            border: '1px solid rgba(16, 185, 129, 0.4)',
            backdropFilter: 'blur(15px)',
            minWidth: '140px',
            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.6)'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>💼</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.25rem', color: '#ffffff', textShadow: '0 2px 8px rgba(0, 0, 0, 0.5)' }}>{(profiles || []).length}+</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em', textShadow: '0 1px 4px rgba(0, 0, 0, 0.5)' }}>Professionals</div>
          </div>
          <div style={{ 
            textAlign: 'center',
            background: 'rgba(10, 15, 30, 0.85)',
            padding: '1.5rem 2rem',
            borderRadius: '16px',
            border: '1px solid rgba(8, 145, 178, 0.4)',
            backdropFilter: 'blur(15px)',
            minWidth: '140px',
            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.6)'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>📂</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.25rem', color: '#ffffff', textShadow: '0 2px 8px rgba(0, 0, 0, 0.5)' }}>{(categories || []).length}</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em', textShadow: '0 1px 4px rgba(0, 0, 0, 0.5)' }}>Categories</div>
          </div>
          <div style={{ 
            textAlign: 'center',
            background: 'rgba(10, 15, 30, 0.85)',
            padding: '1.5rem 2rem',
            borderRadius: '16px',
            border: '1px solid rgba(245, 158, 11, 0.4)',
            backdropFilter: 'blur(15px)',
            minWidth: '140px',
            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.6)'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>⭐</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.25rem', color: '#ffffff', textShadow: '0 2px 8px rgba(0, 0, 0, 0.5)' }}>4.8</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em', textShadow: '0 1px 4px rgba(0, 0, 0, 0.5)' }}>Avg Rating</div>
          </div>
        </div>
      </div>
      
      {error && <p className="error">{error}</p>}
      
      <h2 style={{ marginBottom: '1.5rem', fontSize: '1.875rem' }}>🔍 Find the Perfect Professional</h2>
      
      <div className="filters">
        <select
          value={filters.category}
          onChange={(event) => setFilters((prev) => ({ ...prev, category: event.target.value }))}
        >
          <option value="">🔍 All Categories</option>
          {categories.map((cat) => (
            <option key={cat.id} value={cat.name}>
              {cat.name}
            </option>
          ))}
        </select>
        <input
          placeholder="📍 Search by location"
          value={filters.location}
          onChange={(event) => setFilters((prev) => ({ ...prev, location: event.target.value }))}
        />
      </div>
      
      {!profiles || profiles.length === 0 ? (
        <div className="card text-center" style={{ padding: '4rem 2rem' }}>
          <div style={{ fontSize: '4rem', marginBottom: '1rem' }}>🔍</div>
          <h3 style={{ marginBottom: '1rem', fontSize: '1.5rem' }}>No Professionals Found</h3>
          <p style={{ fontSize: '1.125rem', color: '#6b7280', marginBottom: '1.5rem' }}>
            Try adjusting your filters or browse all categories
          </p>
          <button onClick={() => setFilters({ category: '', location: '' })} style={{ marginTop: '1rem' }}>
            🔄 Clear Filters
          </button>
        </div>
      ) : (
        <>
          <div style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            marginBottom: '1.5rem'
          }}>
            <h3 style={{ fontSize: '1.25rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.5)' }}>
              📋 {(profiles || []).length} Professional{(profiles || []).length !== 1 ? 's' : ''} Available
            </h3>
          </div>
          <div className="cards">
            {(profiles || []).map((profile) => (
            <Link key={profile.id} className="card" to={`/profiles/${profile.id}`}>
              <h3>{profile.category_name}</h3>
              <p style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span>📍</span> {profile.location || 'Location not specified'}
              </p>
              <p style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: '#1f2937', fontWeight: '600' }}>
                <span>💰</span> ${profile.rate || 0}/day
              </p>
              <p style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span>⭐</span> {(profile.average_rating || 0).toFixed(1)} ({profile.review_count || 0} reviews)
              </p>
              <p style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span>👍</span> {profile.upvote_count || 0} upvotes
              </p>
            </Link>
          ))}
          </div>
        </>
      )}
      
      {/* Features Section */}
      <div style={{ marginTop: '4rem', padding: '3rem 0' }}>
        <h2 style={{ textAlign: 'center', marginBottom: '3rem', fontSize: '2rem' }}>
          Why Choose LocalConnect?
        </h2>
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
          gap: '2rem'
        }}>
          <div className="card" style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>✅</div>
            <h3 style={{ marginBottom: '0.75rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>Verified Professionals</h3>
            <p style={{ color: 'rgba(255, 255, 255, 0.85)' }}>
              All professionals are verified with reviews and ratings from real clients
            </p>
          </div>
          <div className="card" style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>💬</div>
            <h3 style={{ marginBottom: '0.75rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>Direct Communication</h3>
            <p style={{ color: 'rgba(255, 255, 255, 0.85)' }}>
              Message professionals directly to discuss your project needs
            </p>
          </div>
          <div className="card" style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🛡️</div>
            <h3 style={{ marginBottom: '0.75rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>Safe & Secure</h3>
            <p style={{ color: 'rgba(255, 255, 255, 0.85)' }}>
              Your information is protected with industry-standard security
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
