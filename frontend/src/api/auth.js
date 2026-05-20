const BASE = '/api'

async function request(path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Something went wrong')
  return data
}

export async function login(identifier, password) {
  return request('/login', { username: identifier, password })
}

export async function register(username, email, password, confirmPassword) {
  return request('/register', { username, email, password, confirmPassword })
}

export async function logout(refreshToken) {
  return request('/auth/logout', { refresh_token: refreshToken })
}

export async function refresh(refreshToken) {
  return request('/auth/refresh', { refresh_token: refreshToken })
}
