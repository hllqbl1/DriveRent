import { useEffect, useState } from 'react'
import { apiRequest } from '../../services/api'

export default function ResourcePage({ title, endpoint }) {
  const [data,setData] = useState([])
  const [error,setError] = useState('')
  const [loading,setLoading] = useState(true)
  useEffect(()=>{apiRequest(endpoint).then(result=>setData(Array.isArray(result)?result:[])).catch(e=>setError(e.message)).finally(()=>setLoading(false))},[endpoint])
  return <div><h1 className="text-3xl font-black">{title}</h1><p className="mt-2 text-slate-500">Data {title.toLowerCase()} dari backend DriveRent.</p><div className="mt-6 overflow-hidden rounded-2xl border bg-white"><div className="overflow-x-auto">{loading?<p className="p-6 text-slate-500">Memuat data...</p>:error?<p className="p-6 text-red-700">{error}</p>:data.length?<table className="w-full min-w-[620px] text-left text-sm"><thead className="bg-slate-50 text-xs uppercase tracking-wider text-slate-500"><tr><th className="px-5 py-4">ID</th><th className="px-5 py-4">Data</th><th className="px-5 py-4">Status</th></tr></thead><tbody className="divide-y divide-slate-100">{data.map((item,i)=><tr key={item.id||i}><td className="px-5 py-4 font-semibold">#{item.id||'-'}</td><td className="px-5 py-4">{item.name||item.email||`Reservation / Rental #${item.reservation_id||item.rental_id||item.id}`}</td><td className="px-5 py-4"><span className="rounded-full bg-blue-50 px-3 py-1 text-xs font-bold text-blue-700">{item.status||'-'}</span></td></tr>)}</tbody></table>:<p className="p-6 text-slate-500">Belum ada data.</p>}</div></div></div>
}
