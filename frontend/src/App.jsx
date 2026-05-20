import { useState } from 'react'
import LoginForm    from './components/LoginForm'
import RegisterForm from './components/RegisterForm'
import Dashboard    from './components/Dashboard'
import './index.css'

function getStoredUser() {
  const token = localStorage.getItem('access_token')
  if (!token) return null
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    if (payload.exp * 1000 < Date.now()) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      return null
    }
    return { username: payload.username, user_id: payload.sub }
  } catch {
    return null
  }
}

export default function App() {
  const [page, setPage] = useState('login')
  const [user, setUser] = useState(getStoredUser)

  function handleLoginSuccess(userData) {
    setUser(userData)
  }

  function handleLogout() {
    setUser(null)
    setPage('login')
  }

  if (user) {
    return (
      <div className="page">
        <Dashboard user={user} onLogout={handleLogout} />
      </div>
    )
  }

  return (
    <div className="page">
      {page === 'login' ? (
        <LoginForm
          onGoToRegister={() => setPage('register')}
          onLoginSuccess={handleLoginSuccess}
        />
      ) : (
        <RegisterForm onGoToLogin={() => setPage('login')} />
      )}
    </div>
  )
}
