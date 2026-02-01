import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../api/client'

export default function ProfileDetail() {
  const { id } = useParams()
  const [profile, setProfile] = useState(null)
  const [reviews, setReviews] = useState([])
  const [contactMessage, setContactMessage] = useState('')
  const [rating, setRating] = useState(5)
  const [comment, setComment] = useState('')
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')

  useEffect(() => {
    api.profile(id)
      .then(setProfile)
      .catch((err) => setError(err.message))
    api.reviews(id)
      .then(setReviews)
      .catch((err) => setError(err.message))
  }, [id])

  const handleUpvote = async () => {
    try {
      await api.upvote(id)
      setNotice('Upvoted successfully')
    } catch (err) {
      setError(err.message)
    }
  }

  const handleReview = async () => {
    try {
      await api.addReview(id, { rating, comment })
      const updated = await api.reviews(id)
      setReviews(updated)
      setNotice('Review added')
    } catch (err) {
      setError(err.message)
    }
  }

  const handleContact = async () => {
    try {
      await api.contact(id, { message: contactMessage, phone_shared: true })
      setNotice('Contact request sent')
      setContactMessage('')
    } catch (err) {
      setError(err.message)
    }
  }

  const handleMessage = async () => {
    if (!profile) return
    try {
      await api.sendMessage({
        receiver_id: profile.user_id,
        content: contactMessage,
      })
      setNotice('Message sent! Check your messages.')
      setContactMessage('')
    } catch (err) {
      setError(err.message)
    }
  }

  if (!profile) {
    return <p>Loading profile...</p>
  }

  return (
    <section>
      <div className="card" style={{ background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #d946ef 100%)', color: 'white' }}>
        <h1 style={{ color: 'white', marginBottom: '1rem' }}>👷 {profile.category_name}</h1>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.5rem' }}>
          <div>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>📍 Location</p>
            <p style={{ fontSize: '1.25rem', fontWeight: '700' }}>{profile.location}</p>
          </div>
          <div>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>💰 Rate</p>
            <p style={{ fontSize: '1.25rem', fontWeight: '700' }}>${profile.rate}/day</p>
          </div>
          <div>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>⭐ Rating</p>
            <p style={{ fontSize: '1.25rem', fontWeight: '700' }}>
              {profile.average_rating.toFixed(1)} ({profile.review_count} reviews)
            </p>
          </div>
          <div>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>💼 Experience</p>
            <p style={{ fontSize: '1.25rem', fontWeight: '700' }}>{profile.experience_years} years</p>
          </div>
        </div>
      </div>

      {error && <p className="error">{error}</p>}
      {notice && <p className="notice">{notice}</p>}

      <div className="card">
        <h3>📝 About</h3>
        <p style={{ lineHeight: '1.8', color: '#1f2937' }}>{profile.bio}</p>
        <div style={{ marginTop: '1.5rem' }}>
          <button onClick={handleUpvote} className="button-success">
            👍 Upvote This Worker
          </button>
        </div>
      </div>

      <div className="card">
        <h3>💬 Contact Worker</h3>
        <textarea 
          placeholder="Write your message here..."
          value={contactMessage} 
          onChange={(event) => setContactMessage(event.target.value)}
          style={{ minHeight: '120px' }}
        />
        <div className="button-group">
          <button onClick={handleContact} className="button-secondary">
            📧 Send Contact Request
          </button>
          <button onClick={handleMessage}>
            💬 Send Direct Message
          </button>
        </div>
      </div>

      <div className="card">
        <h3>⭐ Leave a Review</h3>
        <label>
          Rating
          <select value={rating} onChange={(event) => setRating(Number(event.target.value))}>
            {[1, 2, 3, 4, 5].map((value) => (
              <option key={value} value={value}>
                {'⭐'.repeat(value)} {value} Star{value !== 1 ? 's' : ''}
              </option>
            ))}
          </select>
        </label>
        <label>
          Your Review
          <textarea 
            placeholder="Share your experience..."
            value={comment} 
            onChange={(event) => setComment(event.target.value)}
            style={{ minHeight: '120px' }}
          />
        </label>
        <button onClick={handleReview}>📤 Submit Review</button>
      </div>

      <div className="card">
        <h3>📋 Reviews ({reviews.length})</h3>
        {reviews.length === 0 ? (
          <p style={{ textAlign: 'center', color: '#6b7280', padding: '2rem' }}>
            No reviews yet. Be the first to review!
          </p>
        ) : (
          reviews.map((review) => (
            <div key={review.id} className="review">
              <p style={{ fontSize: '1.125rem', marginBottom: '0.5rem' }}>
                {'⭐'.repeat(review.rating)} {review.rating}/5
              </p>
              <p style={{ color: '#1f2937', marginBottom: '0.5rem' }}>{review.comment}</p>
              <p className="muted">📅 {new Date(review.created_at).toLocaleDateString()}</p>
            </div>
          ))
        )}
      </div>
    </section>
  )
}
