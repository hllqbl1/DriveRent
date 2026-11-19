import { Link } from "react-router";
import { useEffect, useState } from "react";
import { vehicleApi } from "../services/api";
import VehicleCard from "../components/VehicleCard";

export default function Home() {
  const [vehicles, setVehicles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    vehicleApi
      .all()
      .then((data) => {
        setVehicles(Array.isArray(data) ? data.slice(0, 3) : []);
      })
      .catch((err) => {
        console.error("Gagal mengambil kendaraan:", err);
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <>
      {/* HERO SECTION */}
      <section className="relative overflow-hidden bg-slate-950 text-white">
        <div className="container-page grid min-h-[600px] items-center gap-12 py-16 lg:grid-cols-2 lg:py-20">
          {/* Hero Content */}
          <div className="relative z-10">
            <span className="inline-flex items-center gap-2 rounded-full border border-blue-400/30 bg-blue-400/10 px-4 py-2 text-sm font-semibold text-blue-300">
              <span className="h-2 w-2 rounded-full bg-blue-400" />
              Rental Mobil Terpercaya
            </span>

            <h1 className="mt-6 max-w-xl text-4xl font-black leading-tight tracking-tight sm:text-5xl lg:text-6xl">
              Perjalanan Lebih Nyaman Bersama{" "}
              <span className="text-blue-400">DriveRent.</span>
            </h1>

            <p className="mt-6 max-w-lg text-base leading-8 text-slate-300">
              Temukan mobil yang sesuai dengan kebutuhan perjalananmu. Pilihan
              kendaraan berkualitas, harga transparan, dan proses pemesanan yang
              mudah.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <Link
                to="/vehicles"
                className="rounded-xl bg-blue-600 px-7 py-3.5 font-bold text-white transition hover:bg-blue-500"
              >
                Cari Kendaraan
              </Link>

              <Link
                to="/register"
                className="rounded-xl border border-white/30 px-7 py-3.5 font-bold text-white transition hover:bg-white/10"
              >
                Daftar Sekarang
              </Link>
            </div>

            {/* Highlights */}
            <div className="mt-12 grid max-w-md grid-cols-3 gap-5 border-t border-white/10 pt-7">
              <div>
                <p className="text-2xl font-black">Mudah</p>
                <p className="mt-1 text-xs text-slate-400">Pemesanan</p>
              </div>

              <div>
                <p className="text-2xl font-black">Aman</p>
                <p className="mt-1 text-xs text-slate-400">Transaksi</p>
              </div>

              <div>
                <p className="text-2xl font-black">Fleksibel</p>
                <p className="mt-1 text-xs text-slate-400">Pilihan mobil</p>
              </div>
            </div>
          </div>

          {/* Hero Image */}
          <div className="relative mx-auto w-full max-w-xl">
            <div className="absolute -inset-5 rounded-full bg-blue-600/20 blur-3xl" />

            <div className="relative overflow-hidden rounded-3xl border border-white/10 bg-slate-800 shadow-2xl">
              <img
                src="/images/hero-car.jpg"
                alt="Kendaraan DriveRent"
                className="h-[300px] w-full object-cover sm:h-[400px] lg:h-[470px]"
                onError={(event) => {
                  event.currentTarget.style.display = "none";
                }}
              />

              <div className="absolute inset-0 bg-gradient-to-t from-slate-950/80 via-transparent to-transparent" />

              <div className="absolute bottom-0 left-0 p-6 sm:p-8">
                <p className="text-sm font-bold tracking-widest text-blue-300">
                  DRIVE YOUR JOURNEY
                </p>

                <h2 className="mt-2 text-2xl font-black sm:text-3xl">
                  Ready for your next trip?
                </h2>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* SERVICE BENEFITS */}
      <section className="border-b border-slate-200 bg-white">
        <div className="container-page grid gap-8 py-10 sm:grid-cols-3">
          <div>
            <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-700">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="23"
                height="23"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M5 17h14l1.5-5H3.5L5 17Z" />
                <path d="M5 12l2-5h10l2 5" />
                <circle cx="7.5" cy="17.5" r="1.5" />
                <circle cx="16.5" cy="17.5" r="1.5" />
              </svg>
            </div>

            <h3 className="font-bold text-slate-900">Pilihan Kendaraan</h3>

            <p className="mt-2 text-sm leading-6 text-slate-500">
              Berbagai jenis kendaraan untuk kebutuhan perjalananmu.
            </p>
          </div>

          <div>
            <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-700">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="23"
                height="23"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M12 2v20" />
                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
              </svg>
            </div>

            <h3 className="font-bold text-slate-900">Harga Transparan</h3>

            <p className="mt-2 text-sm leading-6 text-slate-500">
              Informasi biaya sewa ditampilkan dengan jelas.
            </p>
          </div>

          <div>
            <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-700">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="23"
                height="23"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <rect x="3" y="4" width="18" height="17" rx="2" />
                <path d="M16 2v4M8 2v4M3 10h18" />
                <path d="m9 15 2 2 4-4" />
              </svg>
            </div>

            <h3 className="font-bold text-slate-900">Reservasi Praktis</h3>

            <p className="mt-2 text-sm leading-6 text-slate-500">
              Pesan kendaraan melalui sistem reservasi online.
            </p>
          </div>
        </div>
      </section>

      {/* VEHICLE CATALOG */}
      <section className="container-page py-16 sm:py-20">
        <div className="mb-9 flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-sm font-bold uppercase tracking-wider text-blue-700">
              Koleksi Kendaraan
            </p>

            <h2 className="mt-2 text-3xl font-black text-slate-900 sm:text-4xl">
              Pilihan Mobil untukmu
            </h2>

            <p className="mt-3 text-slate-500">
              Temukan kendaraan yang cocok untuk perjalananmu.
            </p>
          </div>

          <Link
            to="/vehicles"
            className="font-bold text-blue-700 transition hover:text-blue-900"
          >
            Lihat semua kendaraan →
          </Link>
        </div>

        {loading ? (
          <div className="py-16 text-center">
            <p className="text-slate-500">Memuat kendaraan...</p>
          </div>
        ) : error ? (
          <div className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center">
            <p className="font-semibold text-red-700">
              Gagal mengambil data kendaraan.
            </p>
            <p className="mt-2 text-sm text-red-600">{error}</p>
          </div>
        ) : vehicles.length > 0 ? (
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {vehicles.map((vehicle) => (
              <VehicleCard key={vehicle.id} vehicle={vehicle} />
            ))}
          </div>
        ) : (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-12 text-center">
            <h3 className="font-bold text-slate-800">
              Belum ada kendaraan tersedia
            </h3>

            <p className="mt-2 text-sm text-slate-500">
              Kendaraan akan muncul di sini setelah ditambahkan ke database.
            </p>
          </div>
        )}
      </section>

      {/* CALL TO ACTION */}
      <section className="bg-blue-700">
        <div className="container-page flex flex-col items-start justify-between gap-6 py-12 text-white sm:flex-row sm:items-center">
          <div>
            <h2 className="text-2xl font-black sm:text-3xl">
              Siap memulai perjalanan?
            </h2>

            <p className="mt-2 text-blue-100">
              Temukan mobil yang tepat untuk perjalanan berikutnya.
            </p>
          </div>

          <Link
            to="/vehicles"
            className="rounded-xl bg-white px-6 py-3 font-bold text-blue-700 transition hover:bg-blue-50"
          >
            Lihat Kendaraan
          </Link>
        </div>
      </section>
    </>
  );
}
