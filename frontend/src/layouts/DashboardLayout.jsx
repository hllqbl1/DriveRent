import { Link, NavLink, Outlet } from 'react-router'
import { useAuth } from '../context/AuthContext'

export default function DashboardLayout() {
  const { user, logout } = useAuth()
  const adminLinks = [['/admin/dashboard', 'Ringkasan'], ['/admin/users', 'Pengguna'], ['/admin/vehicles', 'Kendaraan'], ['/admin/reservations', 'Reservasi'], ['/admin/payments', 'Pembayaran'], ['/admin/rentals', 'Rental'], ['/admin/returns', 'Pengembalian']]
  const customerLinks = [['/customer/dashboard', 'Ringkasan'], ['/customer/reservations', 'Reservasi saya'], ['/customer/payments', 'Pembayaran']]
  const links = user?.role === 'customer' ? customerLinks : adminLinks
  return (
    <div className="min-h-screen bg-slate-50">
      <header className="border-b border-slate-200 bg-white"><div className="container-page flex h-16 items-center justify-between"><Link to="/" className="text-xl font-black">Drive<span className="text-blue-700">Rent</span></Link><div className="flex items-center gap-3"><span className="hidden text-sm text-slate-600 sm:block">{user?.name} · {user?.role}</span><button onClick={logout} className="rounded-lg border px-3 py-2 text-sm font-semibold">Keluar</button></div></div></header>
      <div className="container-page grid gap-6 py-6 lg:grid-cols-[230px_minmax(0,1fr)] lg:py-8">
        <aside className="rounded-2xl border border-slate-200 bg-white p-3 lg:h-fit">
          <p className="px-3 py-2 text-xs font-bold uppercase tracking-wider text-slate-400">Menu {user?.role}</p>
          <nav className="flex gap-2 overflow-x-auto lg:flex-col">{links.map(([to, label]) => <NavLink key={to} to={to} end className={({isActive}) => `shrink-0 rounded-xl px-3 py-2.5 text-sm font-semibold ${isActive ? 'bg-blue-50 text-blue-700' : 'text-slate-600 hover:bg-slate-50'}`}>{label}</NavLink>)}</nav>
        </aside>
        <main className="min-w-0"><Outlet /></main>
      </div>
    </div>
  )
}
