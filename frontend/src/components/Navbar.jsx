import { useState } from 'react'
import { Link, NavLink } from 'react-router'
import { useAuth } from '../context/AuthContext'

const linkClass = ({ isActive }) => `text-sm font-medium transition ${isActive ? 'text-blue-700' : 'text-slate-600 hover:text-blue-700'}`

export default function Navbar() {
  const [open, setOpen] = useState(false)
  const { user, logout } = useAuth()
  return (
    <header className="sticky top-0 z-40 border-b border-slate-200 bg-white/95 backdrop-blur">
      <div className="container-page flex h-16 items-center justify-between gap-4">
        <Link to="/" className="flex items-center gap-2 font-black tracking-tight text-slate-900">
          <span className="grid size-9 place-items-center rounded-xl bg-blue-700 text-lg text-white">D</span>
          <span className="text-xl">Drive<span className="text-blue-700">Rent</span></span>
        </Link>
        <button className="rounded-lg border border-slate-200 px-3 py-2 text-sm md:hidden" onClick={() => setOpen(!open)} aria-label="Buka menu">Menu</button>
        <nav className={`${open ? 'absolute left-0 right-0 top-16 flex flex-col border-b border-slate-200 bg-white p-5 shadow-lg' : 'hidden'} gap-5 md:static md:flex md:flex-row md:items-center md:border-0 md:bg-transparent md:p-0 md:shadow-none`}>
          <NavLink to="/" className={linkClass} onClick={() => setOpen(false)}>Beranda</NavLink>
          <NavLink to="/vehicles" className={linkClass} onClick={() => setOpen(false)}>Kendaraan</NavLink>
          {user ? (
            <>
              <NavLink to={user.role === 'customer' ? '/customer/dashboard' : '/admin/dashboard'} className={linkClass}>Dashboard</NavLink>
              <button onClick={logout} className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-700">Keluar</button>
            </>
          ) : (
            <div className="flex items-center gap-3">
              <Link to="/login" className="text-sm font-semibold text-slate-700">Masuk</Link>
              <Link to="/register" className="rounded-lg bg-blue-700 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-800">Daftar</Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  )
}
