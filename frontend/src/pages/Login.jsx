import { useState } from 'react'
import { Link, Navigate, useLocation, useNavigate } from 'react-router'
import { useAuth } from '../context/AuthContext'

export default function Login() {
  const { login, user } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  if (user) return <Navigate to={user.role === 'customer' ? '/customer/dashboard' : '/admin/dashboard'} replace />
  async function submit(e) {
    e.preventDefault(); setError(''); setBusy(true)
    try {
      const logged = await login(form)
      navigate(location.state?.from?.pathname || (logged.role === 'customer' ? '/customer/dashboard' : '/admin/dashboard'), { replace: true })
    } catch (e) { setError(e.message) } finally { setBusy(false) }
  }
  return <div className="container-page grid min-h-[75vh] place-items-center py-12"><form onSubmit={submit} className="w-full max-w-md rounded-3xl border border-slate-200 bg-white p-7 shadow-sm sm:p-9"><p className="text-sm font-bold uppercase tracking-wider text-blue-700">Selamat datang kembali</p><h1 className="mt-2 text-3xl font-black">Masuk ke DriveRent</h1><p className="mt-2 text-sm text-slate-500">Gunakan akun yang sudah terdaftar.</p>{error && <p className="mt-5 rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p>}<label className="mt-6 block text-sm font-semibold">Email<input required type="email" value={form.email} onChange={e=>setForm({...form,email:e.target.value})} className="mt-2 w-full rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-100" placeholder="nama@email.com" /></label><label className="mt-4 block text-sm font-semibold">Password<input required type="password" value={form.password} onChange={e=>setForm({...form,password:e.target.value})} className="mt-2 w-full rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-100" placeholder="Masukkan password" /></label><button disabled={busy} className="mt-6 w-full rounded-xl bg-blue-700 px-5 py-3 font-bold text-white hover:bg-blue-800 disabled:opacity-60">{busy ? 'Memproses...' : 'Masuk'}</button><p className="mt-6 text-center text-sm text-slate-500">Belum punya akun? <Link to="/register" className="font-bold text-blue-700">Daftar</Link></p></form></div>
}
