import { useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router";
import { vehicleApi, reservationApi } from "../../services/api";
import { useAuth } from "../../context/AuthContext";

const today = new Date().toLocaleDateString("en-CA");

function formatRupiah(value) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(Number(value) || 0);
}

function getDuration(start, end) {
  if (!start || !end || end < start) return 0;

  const startDate = new Date(`${start}T00:00:00`);
  const endDate = new Date(`${end}T00:00:00`);

  return Math.floor((endDate - startDate) / 86400000);
}

function StatusBadge({ status }) {
  const styles = {
    pending: "bg-yellow-100 text-yellow-700",
    confirmed: "bg-green-100 text-green-700",
    cancelled: "bg-red-100 text-red-700",
    completed: "bg-blue-100 text-blue-700",
  };

  return (
    <span
      className={`rounded-full px-3 py-1 text-xs font-bold ${styles[status] || "bg-slate-100 text-slate-600"}`}
    >
      {status || "-"}
    </span>
  );
}

export default function Reservations() {
  const { user } = useAuth();
  const [params] = useSearchParams();

  const [vehicles, setVehicles] = useState([]);
  const [items, setItems] = useState([]);
  const [form, setForm] = useState({
    vehicle_id: params.get("vehicle") || "",
    start_date: "",
    end_date: "",
  });

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const loadData = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const vehicleResult = await vehicleApi.all();
      setVehicles(Array.isArray(vehicleResult) ? vehicleResult : []);

      if (user?.id) {
        const reservationResult = await reservationApi.mine(user.id);
        setItems(Array.isArray(reservationResult) ? reservationResult : []);
      }
    } catch (err) {
      setError(err.message || "Gagal mengambil data.");
    } finally {
      setLoading(false);
    }
  }, [user?.id]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const selectedVehicle = vehicles.find(
    (vehicle) => String(vehicle.id) === String(form.vehicle_id),
  );

  const duration = getDuration(form.start_date, form.end_date);

  const totalPrice = useMemo(() => {
    return duration * Number(selectedVehicle?.price_per_day || 0);
  }, [duration, selectedVehicle]);

  const availableVehicles = vehicles.filter(
    (vehicle) => vehicle.status === "available",
  );

  function handleChange(event) {
    const { name, value } = event.target;

    setForm((previous) => ({
      ...previous,
      [name]: value,
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setMessage("");

    if (!selectedVehicle) {
      setError("Silakan pilih kendaraan.");
      return;
    }

    if (selectedVehicle.status !== "available") {
      setError("Kendaraan ini sedang tidak tersedia.");
      return;
    }

    if (duration <= 0) {
      setError("Tanggal selesai harus setelah tanggal mulai.");
      return;
    }

    setSubmitting(true);

    try {
      const result = await reservationApi.create({
        vehicle_id: Number(form.vehicle_id),
        start_date: form.start_date,
        end_date: form.end_date,
      });

      setMessage(result?.message || "Reservasi berhasil diajukan.");

      setForm({
        vehicle_id: "",
        start_date: "",
        end_date: "",
      });

      await loadData();
    } catch (err) {
      setError(err.message || "Gagal membuat reservasi.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div>
      <div>
        <p className="font-bold uppercase tracking-wide text-blue-600">
          Transaksi Customer
        </p>
        <h1 className="mt-2 text-3xl font-black text-slate-900">
          Reservasi Saya
        </h1>
        <p className="mt-2 text-slate-500">
          Pilih kendaraan dan tanggal sewa yang kamu inginkan.
        </p>
      </div>

      {error && (
        <div className="mt-5 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {message && (
        <div className="mt-5 rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {message}
        </div>
      )}

      <form
        onSubmit={handleSubmit}
        className="mt-6 space-y-5 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
      >
        <h2 className="text-lg font-bold text-slate-900">
          Buat Reservasi Baru
        </h2>

        <div>
          <label className="mb-2 block text-sm font-semibold text-slate-700">
            Pilih Kendaraan
          </label>

          <select
            name="vehicle_id"
            value={form.vehicle_id}
            onChange={handleChange}
            required
            className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
          >
            <option value="">Pilih kendaraan</option>

            {availableVehicles.map((vehicle) => (
              <option key={vehicle.id} value={vehicle.id}>
                {vehicle.name} — {formatRupiah(vehicle.price_per_day)}/hari
              </option>
            ))}
          </select>

          {availableVehicles.length === 0 && !loading && (
            <p className="mt-2 text-sm text-orange-600">
              Belum ada kendaraan berstatus available.
            </p>
          )}
        </div>

        {selectedVehicle && (
          <div className="rounded-xl bg-blue-50 p-4">
            <p className="font-bold text-slate-900">{selectedVehicle.name}</p>
            <p className="mt-1 text-sm text-slate-600">
              {selectedVehicle.brand} · {selectedVehicle.year} ·{" "}
              {selectedVehicle.transmission}
            </p>
            <p className="mt-2 font-bold text-blue-700">
              {formatRupiah(selectedVehicle.price_per_day)} / hari
            </p>
          </div>
        )}

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label className="mb-2 block text-sm font-semibold text-slate-700">
              Tanggal Mulai
            </label>
            <input
              type="date"
              name="start_date"
              min={today}
              value={form.start_date}
              onChange={handleChange}
              required
              className="w-full rounded-xl border border-slate-200 px-4 py-3"
            />
          </div>

          <div>
            <label className="mb-2 block text-sm font-semibold text-slate-700">
              Tanggal Selesai
            </label>
            <input
              type="date"
              name="end_date"
              min={form.start_date || today}
              value={form.end_date}
              onChange={handleChange}
              required
              className="w-full rounded-xl border border-slate-200 px-4 py-3"
            />
          </div>
        </div>

        {duration > 0 && selectedVehicle && (
          <div className="rounded-xl border border-blue-100 bg-slate-50 p-4">
            <div className="flex justify-between text-sm text-slate-600">
              <span>Durasi sewa</span>
              <span>{duration} hari</span>
            </div>
            <div className="mt-2 flex justify-between text-sm text-slate-600">
              <span>Harga per hari</span>
              <span>{formatRupiah(selectedVehicle.price_per_day)}</span>
            </div>
            <div className="mt-3 flex justify-between border-t border-slate-200 pt-3">
              <span className="font-bold text-slate-900">Estimasi total</span>
              <span className="text-lg font-black text-blue-700">
                {formatRupiah(totalPrice)}
              </span>
            </div>
            <p className="mt-2 text-xs text-slate-500">
              Total final mengikuti perhitungan backend.
            </p>
          </div>
        )}

        <button
          type="submit"
          disabled={submitting || loading || availableVehicles.length === 0}
          className="rounded-xl bg-blue-600 px-6 py-3 font-bold text-white transition hover:bg-blue-700 disabled:opacity-50"
        >
          {submitting ? "Mengirim..." : "Kirim Reservasi"}
        </button>
      </form>

      <div className="mt-9">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-black text-slate-900">
            Riwayat Reservasi
          </h2>
          <button
            onClick={loadData}
            className="text-sm font-semibold text-blue-700"
          >
            Refresh
          </button>
        </div>

        {loading ? (
          <p className="mt-4 text-slate-500">Memuat riwayat reservasi...</p>
        ) : items.length === 0 ? (
          <p className="mt-4 rounded-xl border border-dashed border-slate-300 p-6 text-sm text-slate-500">
            Belum ada reservasi. Silakan pilih kendaraan untuk mulai memesan.
          </p>
        ) : (
          <div className="mt-4 space-y-3">
            {items.map((item) => {
              const vehicle = vehicles.find(
                (v) => Number(v.id) === Number(item.vehicle_id),
              );

              return (
                <div
                  key={item.id}
                  className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-slate-200 bg-white p-5"
                >
                  <div>
                    <p className="font-bold text-slate-900">
                      {vehicle?.name || `Kendaraan #${item.vehicle_id}`}
                    </p>
                    <p className="mt-1 text-sm text-slate-500">
                      Reservasi #{item.id}
                    </p>
                    <p className="mt-1 text-sm text-slate-500">
                      {String(item.start_date).slice(0, 10)} –{" "}
                      {String(item.end_date).slice(0, 10)}
                    </p>
                    {item.total_price != null && (
                      <p className="mt-2 font-bold text-blue-700">
                        {formatRupiah(item.total_price)}
                      </p>
                    )}
                  </div>

                  <StatusBadge status={item.status} />
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
