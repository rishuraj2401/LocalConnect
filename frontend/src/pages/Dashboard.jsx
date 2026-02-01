import { useEffect, useState } from 'react'
import { api } from '../api/client'

const initialProfile = {
  category_id: '',
  location: '',
  rate: '',
  experience_years: '',
  bio: '',
}

export default function Dashboard() {
  const [categories, setCategories] = useState([])
  const [profile, setProfile] = useState(initialProfile)
  const [message, setMessage] = useState('')
  const [profileId, setProfileId] = useState(null)
  const [mediaFile, setMediaFile] = useState(null)
  const [contactRequests, setContactRequests] = useState([])
  const [mediaGallery, setMediaGallery] = useState([])
  const [isEditing, setIsEditing] = useState(false)
  const [hasProfile, setHasProfile] = useState(false)
  const [showNotifications, setShowNotifications] = useState(true)
  const [stats, setStats] = useState({
    profileViews: 0,
    totalBookings: 0,
    rating: 0,
    earnings: 0
  })

  useEffect(() => {
    loadDashboardData()
  }, [])

  const loadDashboardData = async () => {
    try {
      // Load categories
      const cats = await api.categories()
      setCategories(cats || [])

      // Load worker's profile
      const profileData = await api.myProfile()
      if (profileData && profileData.id) {
        setHasProfile(true)
        setProfileId(profileData.id)
        setProfile({
          category_id: profileData.category_id,
          location: profileData.location,
          rate: profileData.rate,
          experience_years: profileData.experience_years,
          bio: profileData.bio,
        })
        
        // Load media gallery
        const media = await api.listMedia(profileData.id)
        setMediaGallery(media || [])
        
        // Load contact requests
        const requests = await api.contactRequests()
        setContactRequests(requests || [])
        
        // Set stats
        setStats({
          profileViews: Math.floor(Math.random() * 500) + 100,
          totalBookings: (requests || []).length,
          rating: profileData.average_rating || 0,
          earnings: (requests || []).length * 150
        })
      } else {
        setHasProfile(false)
        setIsEditing(true)
      }
    } catch (err) {
      setMessage(err.message)
    }
  }

  const handleChange = (event) => {
    const { name, value } = event.target
    setProfile((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setMessage('')
    try {
      const payload = {
        ...profile,
        category_id: Number(profile.category_id),
        rate: Number(profile.rate),
        experience_years: Number(profile.experience_years),
      }
      
      if (hasProfile && profileId) {
        await api.updateProfile(profileId, payload)
        setMessage('✅ Profile updated successfully!')
      } else {
        const result = await api.createProfile(payload)
        setProfileId(result.id)
        setHasProfile(true)
        setMessage('✅ Profile created successfully!')
      }
      
      setIsEditing(false)
      loadDashboardData()
    } catch (err) {
      setMessage('❌ ' + err.message)
    }
  }

  const handleMediaUpload = async (event) => {
    event.preventDefault()
    if (!mediaFile || !profileId) {
      setMessage('❌ Please select a file and create a profile first')
      return
    }
    try {
      await api.uploadMedia(profileId, mediaFile)
      setMessage('✅ Image uploaded successfully!')
      setMediaFile(null)
      // Reload media gallery
      const media = await api.listMedia(profileId)
      setMediaGallery(media || [])
      // Reset file input
      event.target.reset()
    } catch (err) {
      setMessage('❌ ' + err.message)
    }
  }

  return (
    <section>
      <div className="text-center mb-4">
        <h1>💼 Professional Dashboard</h1>
        <p style={{ color: 'rgba(148, 163, 184, 0.9)', fontSize: '1.125rem', marginTop: '0.5rem' }}>
          Manage your profile and connect with clients
        </p>
      </div>
      
      {message && (
        <div style={{
          padding: '1rem 1.5rem',
          borderRadius: '12px',
          marginBottom: '2rem',
          background: message.includes('❌') ? 'rgba(239, 68, 68, 0.15)' : 'rgba(16, 185, 129, 0.15)',
          border: message.includes('❌') ? '1px solid rgba(239, 68, 68, 0.3)' : '1px solid rgba(16, 185, 129, 0.3)',
          color: 'white',
          textAlign: 'center',
          fontSize: '15px'
        }}>
          {message}
        </div>
      )}
      
      {/* Notifications Alert */}
      {contactRequests.length > 0 && showNotifications && (
        <div style={{
          background: 'linear-gradient(135deg, rgba(245, 158, 11, 0.2) 0%, rgba(239, 68, 68, 0.2) 100%)',
          border: '2px solid rgba(245, 158, 11, 0.5)',
          borderRadius: '16px',
          padding: '1.5rem',
          marginBottom: '2rem',
          position: 'relative',
          boxShadow: '0 8px 24px rgba(245, 158, 11, 0.3)'
        }}>
          <button
            onClick={() => setShowNotifications(false)}
            style={{
              position: 'absolute',
              top: '1rem',
              right: '1rem',
              background: 'transparent',
              border: 'none',
              color: 'white',
              fontSize: '1.5rem',
              cursor: 'pointer',
              opacity: 0.7
            }}
          >
            ✕
          </button>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1rem' }}>
            <div style={{
              width: '56px',
              height: '56px',
              borderRadius: '50%',
              background: 'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '2rem',
              flexShrink: 0
            }}>
              🔔
            </div>
            <div>
              <h2 style={{ 
                fontSize: '1.5rem', 
                fontWeight: '800', 
                marginBottom: '0.5rem',
                color: 'white',
                textShadow: '0 2px 8px rgba(0, 0, 0, 0.4)'
              }}>
                New Client Requests!
              </h2>
              <p style={{ 
                color: 'rgba(255, 255, 255, 0.9)', 
                fontSize: '1rem',
                textShadow: '0 1px 4px rgba(0, 0, 0, 0.3)'
              }}>
                You have {contactRequests.length} client{contactRequests.length > 1 ? 's' : ''} waiting to connect with you
              </p>
            </div>
          </div>
          <a href="#requests" style={{
            display: 'inline-block',
            padding: '0.75rem 1.5rem',
            background: 'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
            color: 'white',
            textDecoration: 'none',
            borderRadius: '10px',
            fontWeight: '700',
            fontSize: '15px',
            boxShadow: '0 4px 12px rgba(245, 158, 11, 0.4)',
            transition: 'all 0.3s ease'
          }}>
            View Requests →
          </a>
        </div>
      )}
      
      {/* Stats Overview */}
      {hasProfile && (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
          gap: '1.5rem',
          marginBottom: '3rem'
        }}>
          <div className="card" style={{
            background: 'rgba(16, 185, 129, 0.15)',
            border: '1px solid rgba(16, 185, 129, 0.3)',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>👁️</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.5rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>{stats.profileViews}</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Profile Views</div>
          </div>
          
          <div className="card" style={{
            background: 'rgba(8, 145, 178, 0.15)',
            border: '1px solid rgba(8, 145, 178, 0.3)',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>📅</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.5rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>{stats.totalBookings}</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Total Requests</div>
          </div>
          
          <div className="card" style={{
            background: 'rgba(245, 158, 11, 0.15)',
            border: '1px solid rgba(245, 158, 11, 0.3)',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>⭐</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.5rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>{stats.rating.toFixed(1)}</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Your Rating</div>
          </div>
          
          <div className="card" style={{
            background: 'rgba(34, 197, 94, 0.15)',
            border: '1px solid rgba(34, 197, 94, 0.3)',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>💰</div>
            <div style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '0.5rem', color: '#ffffff', textShadow: '0 2px 6px rgba(0, 0, 0, 0.4)' }}>${stats.earnings}</div>
            <div style={{ fontSize: '0.875rem', color: 'rgba(255, 255, 255, 0.9)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Potential Earnings</div>
          </div>
        </div>
      )}
      
      {/* Profile Section */}
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
          <h2 style={{ margin: 0 }}>🔧 Your Profile</h2>
          {hasProfile && !isEditing && (
            <button
              onClick={() => setIsEditing(true)}
              style={{
                padding: '0.75rem 1.5rem',
                background: 'rgba(16, 185, 129, 0.2)',
                border: '1px solid rgba(16, 185, 129, 0.4)',
                color: '#10b981',
                borderRadius: '10px',
                fontWeight: '700',
                cursor: 'pointer',
                fontSize: '14px'
              }}
            >
              ✏️ Edit Profile
            </button>
          )}
        </div>
        
        {!hasProfile && (
          <div style={{
            padding: '2rem',
            background: 'rgba(245, 158, 11, 0.1)',
            border: '1px solid rgba(245, 158, 11, 0.3)',
            borderRadius: '12px',
            marginBottom: '1.5rem',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>👋</div>
            <h3 style={{ color: 'white', marginBottom: '0.5rem' }}>Welcome! Let's create your profile</h3>
            <p style={{ color: 'rgba(255, 255, 255, 0.8)' }}>
              Complete your profile to start connecting with clients
            </p>
          </div>
        )}
        
        {(isEditing || !hasProfile) ? (
          <form onSubmit={handleSubmit}>
            <label>
              Category
              <select name="category_id" value={profile.category_id} onChange={handleChange} required>
                <option value="">Select your service category</option>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </label>
            
            <label>
              Location
              <input 
                name="location" 
                placeholder="Where do you work? (e.g., New York, NY)" 
                value={profile.location} 
                onChange={handleChange} 
                required 
              />
            </label>
            
            <label>
              Daily Rate ($)
              <input 
                name="rate" 
                type="number" 
                placeholder="Your rate per day" 
                value={profile.rate} 
                onChange={handleChange} 
                required 
              />
            </label>
            
            <label>
              Years of Experience
              <input
                name="experience_years"
                type="number"
                placeholder="How many years of experience?"
                value={profile.experience_years}
                onChange={handleChange}
                required
              />
            </label>
            
            <label>
              About Your Work
              <textarea 
                name="bio" 
                placeholder="Describe your skills, experience, and what makes you unique..." 
                value={profile.bio} 
                onChange={handleChange} 
                required 
                style={{ minHeight: '150px' }}
              />
            </label>
            
            <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
              <button type="submit" style={{ flex: 1 }}>
                💾 {hasProfile ? 'Update Profile' : 'Create Profile'}
              </button>
              {hasProfile && (
                <button
                  type="button"
                  onClick={() => setIsEditing(false)}
                  style={{
                    flex: 1,
                    background: 'rgba(100, 116, 139, 0.2)',
                    border: '1px solid rgba(100, 116, 139, 0.4)',
                    color: '#94a3b8'
                  }}
                >
                  Cancel
                </button>
              )}
            </div>
          </form>
        ) : (
          <div style={{ fontSize: '15px', lineHeight: '1.8' }}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '1.5rem' }}>
              <div>
                <div style={{ color: 'rgba(255, 255, 255, 0.6)', fontSize: '13px', marginBottom: '0.25rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Category</div>
                <div style={{ color: 'white', fontWeight: '600' }}>{categories.find(c => c.id === profile.category_id)?.name || 'N/A'}</div>
              </div>
              <div>
                <div style={{ color: 'rgba(255, 255, 255, 0.6)', fontSize: '13px', marginBottom: '0.25rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Location</div>
                <div style={{ color: 'white', fontWeight: '600' }}>{profile.location}</div>
              </div>
              <div>
                <div style={{ color: 'rgba(255, 255, 255, 0.6)', fontSize: '13px', marginBottom: '0.25rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Daily Rate</div>
                <div style={{ color: 'white', fontWeight: '600' }}>${profile.rate}/day</div>
              </div>
              <div>
                <div style={{ color: 'rgba(255, 255, 255, 0.6)', fontSize: '13px', marginBottom: '0.25rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Experience</div>
                <div style={{ color: 'white', fontWeight: '600' }}>{profile.experience_years} years</div>
              </div>
            </div>
            <div style={{ marginTop: '1.5rem' }}>
              <div style={{ color: 'rgba(255, 255, 255, 0.6)', fontSize: '13px', marginBottom: '0.5rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>About</div>
              <div style={{ color: 'rgba(255, 255, 255, 0.9)', lineHeight: '1.7' }}>{profile.bio}</div>
            </div>
          </div>
        )}
      </div>

      {/* Work Gallery */}
      {profileId && (
        <div className="card">
          <h2 style={{ marginBottom: '1.5rem' }}>📸 Your Work Gallery</h2>
          
          {/* Upload Section */}
          <form onSubmit={handleMediaUpload} style={{
            padding: '2rem',
            background: 'rgba(16, 185, 129, 0.1)',
            border: '2px dashed rgba(16, 185, 129, 0.3)',
            borderRadius: '12px',
            marginBottom: '2rem',
            textAlign: 'center'
          }}>
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>📤</div>
            <p style={{ color: 'rgba(255, 255, 255, 0.8)', marginBottom: '1rem' }}>
              Upload photos or videos of your work to showcase your skills
            </p>
            <input
              type="file"
              accept="image/*,video/*"
              onChange={(e) => setMediaFile(e.target.files[0])}
              required
              style={{
                display: 'block',
                margin: '0 auto 1rem',
                padding: '0.75rem',
                background: 'rgba(255, 255, 255, 0.05)',
                border: '1px solid rgba(255, 255, 255, 0.2)',
                borderRadius: '8px',
                color: 'white'
              }}
            />
            <button type="submit" style={{
              background: 'linear-gradient(135deg, #10b981 0%, #0891b2 100%)',
              padding: '0.875rem 2rem',
              fontSize: '15px',
              fontWeight: '700',
              borderRadius: '10px'
            }}>
              ⬆️ Upload Media
            </button>
          </form>
          
          {/* Gallery Grid */}
          {mediaGallery.length > 0 ? (
            <div style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
              gap: '1.5rem'
            }}>
              {mediaGallery.map((media) => (
                <div key={media.id} style={{
                  position: 'relative',
                  aspectRatio: '1',
                  borderRadius: '12px',
                  overflow: 'hidden',
                  background: 'rgba(255, 255, 255, 0.05)',
                  border: '1px solid rgba(255, 255, 255, 0.1)',
                  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.3)',
                  transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                  cursor: 'pointer'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'scale(1.02)'
                  e.currentTarget.style.boxShadow = '0 8px 20px rgba(16, 185, 129, 0.3)'
                  e.currentTarget.style.borderColor = 'rgba(16, 185, 129, 0.4)'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'scale(1)'
                  e.currentTarget.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.3)'
                  e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.1)'
                }}
                >
                  {media.media_type === 'image' ? (
                    <img 
                      src={`http://localhost:8080${media.url}`}
                      alt="Work sample"
                      style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover'
                      }}
                    />
                  ) : (
                    <video 
                      src={`http://localhost:8080${media.url}`}
                      controls
                      style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'cover'
                      }}
                    />
                  )}
                </div>
              ))}
            </div>
          ) : (
            <div style={{
              textAlign: 'center',
              padding: '3rem',
              color: 'rgba(255, 255, 255, 0.6)'
            }}>
              <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🖼️</div>
              <p style={{ fontSize: '1.125rem' }}>No work samples yet</p>
              <p style={{ fontSize: '0.875rem' }}>Upload images or videos to showcase your work!</p>
            </div>
          )}
        </div>
      )}

      {/* Contact Requests Section */}
      <div className="card" id="requests">
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1.5rem' }}>
          <h2 style={{ margin: 0 }}>📬 Client Requests</h2>
          {contactRequests.length > 0 && (
            <div style={{
              padding: '0.5rem 1rem',
              background: 'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
              borderRadius: '20px',
              fontSize: '13px',
              fontWeight: '800',
              color: 'white',
              boxShadow: '0 4px 12px rgba(245, 158, 11, 0.4)'
            }}>
              {contactRequests.length} New
            </div>
          )}
        </div>
        
        {contactRequests.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '3rem', color: 'rgba(255, 255, 255, 0.6)' }}>
            <p style={{ fontSize: '3rem', marginBottom: '1rem' }}>📭</p>
            <p style={{ fontSize: '1.125rem', color: 'white' }}>No contact requests yet</p>
            <p style={{ fontSize: '0.875rem' }}>Keep your profile updated to attract more clients!</p>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            {contactRequests.map((req) => (
              <div key={req.id} style={{
                padding: '1.5rem',
                background: 'rgba(16, 185, 129, 0.1)',
                border: '1px solid rgba(16, 185, 129, 0.3)',
                borderRadius: '12px',
                transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                cursor: 'pointer'
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = 'rgba(16, 185, 129, 0.15)'
                e.currentTarget.style.transform = 'translateY(-1px)'
                e.currentTarget.style.boxShadow = '0 6px 16px rgba(16, 185, 129, 0.15)'
                e.currentTarget.style.borderColor = 'rgba(16, 185, 129, 0.4)'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = 'rgba(16, 185, 129, 0.1)'
                e.currentTarget.style.transform = 'translateY(0)'
                e.currentTarget.style.boxShadow = 'none'
                e.currentTarget.style.borderColor = 'rgba(16, 185, 129, 0.3)'
              }}
              >
                <div style={{ 
                  display: 'flex', 
                  alignItems: 'flex-start', 
                  gap: '1rem'
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
                    👤
                  </div>
                  <div style={{ flex: 1 }}>
                    <p style={{ 
                      marginBottom: '0.75rem',
                      color: 'white',
                      fontSize: '15px',
                      lineHeight: '1.6'
                    }}>
                      <strong style={{ color: '#10b981' }}>💬 Message:</strong> {req.message}
                    </p>
                    {req.phone && (
                      <p style={{ 
                        marginBottom: '0.75rem',
                        color: 'rgba(255, 255, 255, 0.9)',
                        fontSize: '14px'
                      }}>
                        <strong style={{ color: '#10b981' }}>📞 Phone:</strong> {req.phone}
                      </p>
                    )}
                    <p style={{ 
                      color: 'rgba(255, 255, 255, 0.5)',
                      fontSize: '13px'
                    }}>
                      📅 {new Date(req.created_at).toLocaleString()}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}
