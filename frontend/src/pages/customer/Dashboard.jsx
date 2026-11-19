import { Link } from "react-router";
import { useAuth } from "../../context/AuthContext";
export default function Dashboard() {
  const { user } = useAuth();
  return (
    <div>
      <p className="text-sm font-semibold text-blue-700">CUSTOMER</p>
      <h1 className="mt-1 text-3xl font-black">
        Halo, {user?.name || "Customer"}!
      </h1>
      <p className="mt-2 text-slate-500">
        Kelola reservasi dan pembayaran kendaraanmu dari sini.
      </p>
      <div className="mt-7 grid gap-4 sm:grid-cols-2">
        <Link
          to="/vehicles"
          className="rounded-2xl border bg-white p-6 shadow-sm hover:border-blue-300"
        >
          <span className="text-3xl">🚘</span>
          <h2 className="mt-4 text-lg font-bold">Cari kendaraan</h2>
          <p className="mt-1 text-sm text-slate-500">
            Lihat kendaraan yang tersedia.
          </p>
        </Link>
        <Link
          to="/customer/reservations"
          className="rounded-2xl border bg-white p-6 shadow-sm hover:border-blue-300"
        >
          <span className="text-3xl">📅</span>
          <h2 className="mt-4 text-lg font-bold">Reservasi saya</h2>
          <p className="mt-1 text-sm text-slate-500">
            Lihat proses penyewaanmu.
          </p>
        </Link>
      </div>
    </div>
  );
}
