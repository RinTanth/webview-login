import { useState } from 'react'
import { logout } from '../api/auth'

export default function Dashboard({ user, onLogout }) {
  const [loading, setLoading] = useState(false)
  const [error, setError]   = useState('')

  async function handleLogout() {
    setLoading(true)
    try {
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) await logout(refreshToken)
    } catch {
      // proceed with local logout even if server call fails
    } finally {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      onLogout()
    }
  }

  return (
    <div className="card">
      <h1>Welcome</h1>
      <p className="switch-text">Logged in as <strong>{user.username}</strong></p>

      {error && <p className="error">{error}</p>}

      <button
        className="btn-primary"
        onClick={handleLogout}
        disabled={loading}
        style={{ marginTop: '24px' }}
      >
        {loading ? 'Logging out...' : 'Logout'}
      </button>
    </div>
  )
}
