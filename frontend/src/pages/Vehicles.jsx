import { useEffect, useState } from 'react'
import VehicleCard from '../components/VehicleCard'
import { vehicleApi } from '../services/api'

export default function Vehicles() {
  const [vehicles, setVehicles] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  useEffect(() => { vehicleApi.all().then(data => setVehicles(Array.isArray(data) ? data : [])).catch(e => setError(e.message)).finally(() => setLoading(false)) }, [])
  return <section className="container-page py-12 sm:py-16"><div className="mb-8"><p className="text-sm font-bold uppercase tracking-wider text-blue-700">Katalog</p><h1 className="mt-2 text-3xl font-black sm:text-4xl">Daftar kendaraan</h1><p className="mt-3 text-slate-500">Pilih kendaraan yang sesuai dengan perjalananmu.</p></div>{loading ? <p className="py-12 text-center text-slate-500">Memuat kendaraan...</p> : error ? <p className="rounded-xl bg-red-50 p-4 text-red-700">{error}</p> : vehicles.length ? <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">{vehicles.map(v => <VehicleCard key={v.id} vehicle={v} />)}</div> : <div className="rounded-2xl border border-dashed p-12 text-center text-slate-500">Belum ada kendaraan tersedia.</div>}</section>
}
