import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router'
import { reservationApi } from '../../services/api'
import { useAuth } from '../../context/AuthContext'

export default function Reservations() {
  const { user } = useAuth()
  const [params] = useSearchParams()
  const [form, setForm] = useState({ vehicle_id: params.get('vehicle') || '', start_date: '', end_date: '' })
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [items, setItems] = useState([])
  async function load() {
    if (!user?.id) return
    try { const data = await reservationApi.mine(user.id); setItems(Array.isArray(data) ? data : []) } catch { /* route backend bisa berbeda; form tetap bisa digunakan */ }
  }
  useEffect(() => { load() }, [user?.id])
  async function submit(e) {
    e.preventDefault(); setMessage(''); setError('')
    try { const result = await reservationApi.create({ vehicle_id: Number(form.vehicle_id), start_date: form.start_date, end_date: form.end_date }); setMessage(result?.message || 'Reservasi berhasil dibuat.'); setForm({...form,start_date:'',end_date:''}); load() }
    catch (e) { setError(e.message) }
  }
  return <div><h1 className="text-3xl font-black">Reservasi saya</h1><p className="mt-2 text-slate-500">Ajukan reservasi dan pantau statusnya.</p><form onSubmit={submit} className="mt-6 grid gap-4 rounded-2xl border bg-white p-5 sm:grid-cols-3"><label className="text-sm font-semibold">ID kendaraan<input required type="number" min="1" value={form.vehicle_id} onChange={e=>setForm({...form,vehicle_id:e.target.value})} className="mt-2 w-full rounded-xl border px-3 py-3" /></label><label className="text-sm font-semibold">Tanggal mulai<input required type="date" value={form.start_date} onChange={e=>setForm({...form,start_date:e.target.value})} className="mt-2 w-full rounded-xl border px-3 py-3" /></label><label className="text-sm font-semibold">Tanggal selesai<input required type="date" value={form.end_date} onChange={e=>setForm({...form,end_date:e.target.value})} className="mt-2 w-full rounded-xl border px-3 py-3" /></label><div className="sm:col-span-3"><button className="rounded-xl bg-blue-700 px-5 py-3 font-bold text-white">Kirim reservasi</button></div>{message && <p className="text-sm text-emerald-700 sm:col-span-3">{message}</p>}{error && <p className="text-sm text-red-700 sm:col-span-3">{error}</p>}</form><h2 className="mt-8 text-xl font-bold">Riwayat reservasi</h2><div className="mt-3 space-y-3">{items.length ? items.map(item=><div key={item.id} className="flex flex-wrap justify-between gap-3 rounded-xl border bg-white p-4"><div><p className="font-bold">Reservasi #{item.id}</p><p className="text-sm text-slate-500">Kendaraan #{item.vehicle_id} · {String(item.start_date).slice(0,10)} – {String(item.end_date).slice(0,10)}</p></div><span className="h-fit rounded-full bg-blue-50 px-3 py-1 text-xs font-bold text-blue-700">{item.status}</span></div>) : <p className="rounded-xl border border-dashed p-6 text-sm text-slate-500">Belum ada data reservasi atau endpoint daftar customer perlu disesuaikan.</p>}</div></div>
}
