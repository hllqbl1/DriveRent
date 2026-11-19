import { Link } from "react-router";
import { useAuth } from "../../context/AuthContext";
const cards = [
  ["Kendaraan", "Kelola armada", "/admin/vehicles", "🚘"],
  ["Reservasi", "Pantau permintaan sewa", "/admin/reservations", "📅"],
  ["Pembayaran", "Verifikasi pembayaran", "/admin/payments", "💳"],
  ["Rental", "Proses serah terima", "/admin/rentals", "🔑"],
  ["Pengembalian", "Catat kendaraan kembali", "/admin/returns", "↩️"],
  ["Pengguna", "Data akun pengguna", "/admin/users", "👥"],
];
export default function Dashboard() {
  const { user } = useAuth();
  return (
    <div>
      <p className="text-sm font-semibold text-blue-700">ADMIN / STAFF</p>
      <h1 className="mt-1 text-3xl font-black">Dashboard</h1>
      <p className="mt-2 text-slate-500">
        Selamat datang, {user?.name || "Admin"}. Pilih modul yang ingin
        dikelola.
      </p>
      <div className="mt-7 grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {cards
          .filter((c) => user?.role === "admin" || c[0] !== "Pengguna")
          .map((c) => (
            <Link
              key={c[0]}
              to={c[2]}
              className="rounded-2xl border bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:border-blue-300 hover:shadow-md"
            >
              <span className="text-3xl">{c[3]}</span>
              <h2 className="mt-4 font-bold">{c[0]}</h2>
              <p className="mt-1 text-sm text-slate-500">{c[1]}</p>
            </Link>
          ))}
      </div>
    </div>
  );
}
