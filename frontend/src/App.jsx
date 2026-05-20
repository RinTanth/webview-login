import { useState } from 'react'
import LoginForm from './components/LoginForm'
import RegisterForm from './components/RegisterForm'
import './index.css'

export default function App() {
  const [page, setPage] = useState('login')

  return (
    <div className="page">
      {page === 'login' ? (
        <LoginForm onGoToRegister={() => setPage('register')} />
      ) : (
        <RegisterForm onGoToLogin={() => setPage('login')} />
      )}
    </div>
  )
}
