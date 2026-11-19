import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { vehicleApi } from '../services/api'
import { useAuth } from '../context/AuthContext'

export default function VehicleDetail() {
  const { id } = useParams()
  const { user } = useAuth()
  const [vehicle, setVehicle] = useState(null)
  const [error, setError] = useState('')
  useEffect(() => { vehicleApi.byId(id).then(setVehicle).catch(e => setError(e.message)) }, [id])
  if (error) return <div className="container-page py-16 text-red-700">{error}</div>
  if (!vehicle) return <div className="container-page py-16 text-slate-500">Memuat detail kendaraan...</div>
  return <section className="container-page py-12"><Link to="/vehicles" className="text-sm font-semibold text-blue-700">← Kembali ke kendaraan</Link><div className="mt-6 grid gap-8 lg:grid-cols-2"><div className="grid min-h-72 place-items-center overflow-hidden rounded-3xl bg-gradient-to-br from-slate-100 to-blue-100">{vehicle.image ? <img src={vehicle.image} alt={vehicle.name} className="h-full max-h-[460px] w-full object-cover" /> : <span className="text-9xl">🚙</span>}</div><div className="py-2"><p className="font-bold uppercase tracking-wider text-blue-700">{vehicle.brand} · {vehicle.category}</p><h1 className="mt-3 text-4xl font-black">{vehicle.name}</h1><p className="mt-4 leading-7 text-slate-600">{vehicle.description || 'Kendaraan nyaman untuk perjalananmu.'}</p><div className="mt-7 grid grid-cols-2 gap-4">{[['Tahun', vehicle.year], ['Transmisi', vehicle.transmission], ['Bahan bakar', vehicle.fuel_type], ['Kapasitas', `${vehicle.capacity} orang`]].map(([k,v]) => <div key={k} className="rounded-xl bg-slate-100 p-4"><p className="text-xs text-slate-500">{k}</p><p className="mt-1 font-bold capitalize">{v || '-'}</p></div>)}</div><p className="mt-8 text-2xl font-black">Rp{Number(vehicle.price_per_day).toLocaleString('id-ID')} <span className="text-sm font-medium text-slate-500">/ hari</span></p>{user?.role === 'customer' ? <Link to={`/customer/reservations?vehicle=${vehicle.id}`} className="mt-5 inline-flex rounded-xl bg-blue-700 px-6 py-3 font-bold text-white">Ajukan reservasi</Link> : <p className="mt-5 text-sm text-slate-500">Masuk sebagai customer untuk mengajukan reservasi.</p>}</div></div></section>
}
