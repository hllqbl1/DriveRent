import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { useAuth } from '../context/AuthContext'

export default function Register() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ name: '', email: '', phone: '', password: '' })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [busy, setBusy] = useState(false)
  async function submit(e) {
    e.preventDefault(); setError(''); setSuccess(''); setBusy(true)
    try { await register(form); setSuccess('Akun berhasil dibuat. Silakan masuk.'); setTimeout(() => navigate('/login'), 900) }
    catch (e) { setError(e.message) } finally { setBusy(false) }
  }
  return <div className="container-page grid min-h-[75vh] place-items-center py-12"><form onSubmit={submit} className="w-full max-w-lg rounded-3xl border border-slate-200 bg-white p-7 shadow-sm sm:p-9"><p className="text-sm font-bold uppercase tracking-wider text-blue-700">Mulai bersama DriveRent</p><h1 className="mt-2 text-3xl font-black">Buat akun</h1><p className="mt-2 text-sm text-slate-500">Daftar sebagai customer untuk mulai menyewa kendaraan.</p>{error && <p className="mt-5 rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p>}{success && <p className="mt-5 rounded-xl bg-emerald-50 p-3 text-sm text-emerald-700">{success}</p>}{[['name','Nama lengkap','text'],['email','Email','email'],['phone','Nomor telepon','tel'],['password','Password','password']].map(([key,label,type])=><label key={key} className="mt-4 block text-sm font-semibold">{label}<input required type={type} value={form[key]} onChange={e=>setForm({...form,[key]:e.target.value})} className="mt-2 w-full rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-100" /></label>)}<button disabled={busy} className="mt-6 w-full rounded-xl bg-blue-700 px-5 py-3 font-bold text-white hover:bg-blue-800 disabled:opacity-60">{busy ? 'Memproses...' : 'Daftar'}</button><p className="mt-6 text-center text-sm text-slate-500">Sudah punya akun? <Link to="/login" className="font-bold text-blue-700">Masuk</Link></p></form></div>
}
