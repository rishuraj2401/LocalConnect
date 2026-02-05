const API_URL = import.meta.env.VITE_API_URL || 'https://localconnect-ghti.onrender.com'

export const tokenStore = {
  get() {
    return localStorage.getItem('lc_token')
  },
  set(token) {
    localStorage.setItem('lc_token', token)
  },
  clear() {
    localStorage.removeItem('lc_token')
  },
}

async function request(path, options = {}) {
  const headers = new Headers(options.headers || {})
  headers.set('Content-Type', 'application/json')
  const token = tokenStore.get()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    throw new Error(error.error || 'Request failed')
  }
  if (response.status === 204) {
    return null
  }
  return response.json()
}

async function uploadMedia(path, file) {
  const formData = new FormData()
  formData.append('file', file)
  const token = tokenStore.get()
  const headers = new Headers()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  const response = await fetch(`${API_URL}${path}`, {
    method: 'POST',
    headers,
    body: formData,
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    throw new Error(error.error || 'Upload failed')
  }
  return response.json()
}

export const api = {
  login: (payload) => request('/auth/login', { method: 'POST', body: JSON.stringify(payload) }),
  register: (payload) => request('/auth/register', { method: 'POST', body: JSON.stringify(payload) }),
  categories: () => request('/categories'),
  profiles: (params) => {
    const search = new URLSearchParams(params)
    return request(`/profiles?${search.toString()}`)
  },
  profile: (id) => request(`/profiles/${id}`),
  myProfile: () => request('/profiles/me'),
  reviews: (id) => request(`/profiles/${id}/reviews`),
  addReview: (id, payload) => request(`/profiles/${id}/reviews`, { method: 'POST', body: JSON.stringify(payload) }),
  upvote: (id) => request(`/profiles/${id}/upvote`, { method: 'POST' }),
  removeUpvote: (id) => request(`/profiles/${id}/upvote`, { method: 'DELETE' }),
  createProfile: (payload) => request('/profiles', { method: 'POST', body: JSON.stringify(payload) }),
  updateProfile: (id, payload) => request(`/profiles/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
  uploadMedia: (id, file) => uploadMedia(`/profiles/${id}/media`, file),
  listMedia: (id) => request(`/profiles/${id}/media`),
  contact: (id, payload) => request(`/profiles/${id}/contact-requests`, { method: 'POST', body: JSON.stringify(payload) }),
  contactRequests: () => request('/contact-requests'),
  sendMessage: (payload) => request('/messages', { method: 'POST', body: JSON.stringify(payload) }),
  conversations: () => request('/conversations'),
  messages: (conversationId) => request(`/conversations/${conversationId}/messages`),
}
