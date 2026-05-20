import { useState } from 'react'

export default function RegisterForm({ onGoToLogin }) {
  const [form, setForm] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  })
  const [error, setError] = useState('')

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  function handleSubmit(e) {
    e.preventDefault()

    if (!form.username || !form.email || !form.password || !form.confirmPassword) {
      setError('Please fill in all fields.')
      return
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      setError('Please enter a valid email address.')
      return
    }
    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match.')
      return
    }
    if (form.password.length < 6) {
      setError('Password must be at least 6 characters.')
      return
    }

    setError('')
    // TODO: call backend API
    alert(`Registered: ${form.username} (${form.email})`)
  }

  return (
    <div className="card">
      <h1>Register</h1>

      <form onSubmit={handleSubmit}>
        <label htmlFor="reg-username">Username</label>
        <input
          id="reg-username"
          name="username"
          type="text"
          placeholder="Choose a username"
          value={form.username}
          onChange={handleChange}
          autoComplete="username"
        />

        <label htmlFor="reg-email">Email</label>
        <input
          id="reg-email"
          name="email"
          type="email"
          placeholder="Enter your email"
          value={form.email}
          onChange={handleChange}
          autoComplete="email"
        />

        <label htmlFor="reg-password">Password</label>
        <input
          id="reg-password"
          name="password"
          type="password"
          placeholder="Choose a password"
          value={form.password}
          onChange={handleChange}
          autoComplete="new-password"
        />

        <label htmlFor="reg-confirm">Confirm Password</label>
        <input
          id="reg-confirm"
          name="confirmPassword"
          type="password"
          placeholder="Repeat your password"
          value={form.confirmPassword}
          onChange={handleChange}
          autoComplete="new-password"
        />

        {error && <p className="error">{error}</p>}

        <button type="submit" className="btn-primary">Register</button>
      </form>

      <p className="switch-text">
        Already have an account?{' '}
        <button type="button" className="btn-link" onClick={onGoToLogin}>
          Login
        </button>
      </p>
    </div>
  )
}
