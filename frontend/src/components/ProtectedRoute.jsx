import { Navigate, Outlet, useLocation } from 'react-router'
import { useAuth } from '../context/AuthContext'

export default function ProtectedRoute({ roles }) {
  const { user, loading } = useAuth()
  const location = useLocation()
  if (loading) return <div className="grid min-h-[50vh] place-items-center text-slate-500">Memuat...</div>
  if (!user) return <Navigate to="/login" replace state={{ from: location }} />
  if (roles?.length && !roles.includes(user.role)) {
    return <Navigate to={user.role === 'customer' ? '/customer/dashboard' : '/admin/dashboard'} replace />
  }
  return <Outlet />
}
