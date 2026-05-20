import { useState } from 'react'

export default function LoginForm({ onGoToRegister }) {
  const [form, setForm] = useState({ username: '', password: '' })
  const [error, setError] = useState('')

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  function handleSubmit(e) {
    e.preventDefault()
    if (!form.username || !form.password) {
      setError('Please fill in all fields.')
      return
    }
    setError('')
    // TODO: call backend API
    alert(`Logging in as: ${form.username}`)
  }

  return (
    <div className="card">
      <h1>Login</h1>

      <form onSubmit={handleSubmit}>
        <label htmlFor="username">Username or Email</label>
        <input
          id="username"
          name="username"
          type="text"
          placeholder="Enter username or email"
          value={form.username}
          onChange={handleChange}
          autoComplete="username"
        />

        <label htmlFor="password">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          placeholder="Enter password"
          value={form.password}
          onChange={handleChange}
          autoComplete="current-password"
        />

        {error && <p className="error">{error}</p>}

        <button type="submit" className="btn-primary">Login</button>
      </form>

      <p className="switch-text">
        Don&apos;t have an account?{' '}
        <button type="button" className="btn-link" onClick={onGoToRegister}>
          Register
        </button>
      </p>
    </div>
  )
}
