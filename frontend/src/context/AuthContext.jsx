import { createContext, useContext, useEffect, useState } from 'react'
import { authApi } from '../services/api'

const AuthContext = createContext(null)

function readStoredUser() {
  try { return JSON.parse(localStorage.getItem('driverent_user') || 'null') }
  catch { return null }
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(readStoredUser)
  const [loading, setLoading] = useState(false)

  async function login(credentials) {
    const result = await authApi.login(credentials)
    if (!result?.token || !result?.user) throw new Error('Response login dari server tidak lengkap.')
    localStorage.setItem('driverent_token', result.token)
    localStorage.setItem('driverent_user', JSON.stringify(result.user))
    setUser(result.user)
    return result.user
  }

  async function register(payload) {
    return authApi.register(payload)
  }

  function logout() {
    localStorage.removeItem('driverent_token')
    localStorage.removeItem('driverent_user')
    setUser(null)
  }

  useEffect(() => {
    // Starter ini memakai data user dari response login. Validasi ulang sesi bisa ditambahkan
    // jika endpoint /auth/me sudah digunakan pada alur aplikasi.
    setLoading(false)
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const value = useContext(AuthContext)
  if (!value) throw new Error('useAuth harus digunakan di dalam AuthProvider')
  return value
}
