import { Navigate, Route, Routes } from 'react-router'
import MainLayout from './layouts/MainLayout'
import DashboardLayout from './layouts/DashboardLayout'
import ProtectedRoute from './components/ProtectedRoute'
import Home from './pages/Home'
import Vehicles from './pages/Vehicles'
import VehicleDetail from './pages/VehicleDetail'
import Login from './pages/Login'
import Register from './pages/Register'
import CustomerDashboard from './pages/customer/Dashboard'
import Reservations from './pages/customer/Reservations'
import Payments from './pages/customer/Payments'
import AdminDashboard from './pages/admin/Dashboard'
import Users from './pages/admin/Users'
import AdminVehicles from './pages/admin/Vehicles'
import AdminReservations from './pages/admin/Reservations'
import AdminPayments from './pages/admin/Payments'
import Rentals from './pages/admin/Rentals'
import Returns from './pages/admin/Returns'

function NotFound() {
  return <div className="container-page py-24 text-center"><p className="text-6xl font-black text-slate-200">404</p><h1 className="mt-3 text-2xl font-bold">Halaman tidak ditemukan</h1><a href="/" className="mt-5 inline-block font-semibold text-blue-700">Kembali ke beranda</a></div>
}

export default function App() {
  return <Routes>
    <Route element={<MainLayout />}>
      <Route path="/" element={<Home />} />
      <Route path="/vehicles" element={<Vehicles />} />
      <Route path="/vehicles/:id" element={<VehicleDetail />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
    </Route>
    <Route element={<ProtectedRoute roles={['customer']} />}>
      <Route element={<DashboardLayout />}>
        <Route path="/customer/dashboard" element={<CustomerDashboard />} />
        <Route path="/customer/reservations" element={<Reservations />} />
        <Route path="/customer/payments" element={<Payments />} />
      </Route>
    </Route>
    <Route element={<ProtectedRoute roles={['admin','staff']} />}>
      <Route element={<DashboardLayout />}>
        <Route path="/admin/dashboard" element={<AdminDashboard />} />
        <Route path="/admin/vehicles" element={<AdminVehicles />} />
        <Route path="/admin/reservations" element={<AdminReservations />} />
        <Route path="/admin/payments" element={<AdminPayments />} />
        <Route path="/admin/rentals" element={<Rentals />} />
        <Route path="/admin/returns" element={<Returns />} />
      </Route>
    </Route>
    <Route element={<ProtectedRoute roles={['admin']} />}>
      <Route element={<DashboardLayout />}>
        <Route path="/admin/users" element={<Users />} />
      </Route>
    </Route>
    <Route path="/404" element={<NotFound />} />
    <Route path="*" element={<NotFound />} />
  </Routes>
}
