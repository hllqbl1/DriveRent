import { useCallback, useEffect, useState } from "react";
import { rentalApi, reservationApi } from "../../services/api";
import { useAuth } from "../../context/AuthContext";

function formatDate(date) {
  if (!date) return "-";

  return new Date(date).toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

function StatusBadge({ status }) {
  const styles = {
    active: "bg-blue-100 text-blue-700",
    completed: "bg-green-100 text-green-700",
    cancelled: "bg-red-100 text-red-700",
  };

  return (
    <span
      className={`rounded-full px-3 py-1 text-xs font-bold ${
        styles[status] || "bg-slate-100 text-slate-600"
      }`}
    >
      {status || "-"}
    </span>
  );
}

function DetailItem({ label, value }) {
  return (
    <div className="flex items-center justify-between gap-4">
      <span className="text-sm text-slate-500">{label}</span>
      <span className="text-right text-sm font-semibold text-slate-900">
        {value || "-"}
      </span>
    </div>
  );
}

const initialForm = {
  reservation_id: "",
  pickup_date: "",
  initial_fuel: "Full",
  initial_odometer: "",
};

export default function Rentals() {
  const { user } = useAuth();

  const [rentals, setRentals] = useState([]);
  const [reservations, setReservations] = useState([]);

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const [showModal, setShowModal] = useState(false);
  const [selectedRental, setSelectedRental] = useState(null);
  const [form, setForm] = useState(initialForm);

  const loadData = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const [rentalResult, reservationResult] = await Promise.all([
        rentalApi.all(),
        reservationApi.all(),
      ]);

      setRentals(Array.isArray(rentalResult) ? rentalResult : []);
      setReservations(
        Array.isArray(reservationResult) ? reservationResult : [],
      );
    } catch (err) {
      setError(err.message || "Gagal mengambil data rental.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const activeCount = rentals.filter(
    (rental) => rental.status === "active",
  ).length;

  const completedCount = rentals.filter(
    (rental) => rental.status === "completed",
  ).length;

  const otherCount = rentals.length - activeCount - completedCount;

  // Reservasi confirmed yang belum memiliki rental.
  const availableReservations = reservations.filter((reservation) => {
    const hasRental = rentals.some(
      (rental) => Number(rental.reservation_id) === Number(reservation.id),
    );

    return reservation.status === "confirmed" && !hasRental;
  });

  function updateForm(event) {
    const { name, value } = event.target;

    setForm((previous) => ({
      ...previous,
      [name]: value,
    }));
  }

  function openCreateModal() {
    setError("");
    setMessage("");
    setForm(initialForm);
    setShowModal(true);
  }

  async function handleCreateRental(event) {
    event.preventDefault();

    setError("");
    setMessage("");

    if (!form.reservation_id) {
      setError("Pilih reservasi terlebih dahulu.");
      return;
    }

    if (!form.initial_odometer || Number(form.initial_odometer) < 0) {
      setError("Odometer awal harus diisi dengan angka yang valid.");
      return;
    }

    if (!user?.id) {
      setError("Identitas admin tidak ditemukan. Silakan login ulang.");
      return;
    }

    setSubmitting(true);

    try {
      const payload = {
        reservation_id: Number(form.reservation_id),

        // Backend saat ini masih menggunakan field staff_id.
        // ID admin dikirim sebagai pencatat rental sementara.
        staff_id: Number(user.id),

        pickup_date: form.pickup_date
          ? new Date(form.pickup_date).toISOString()
          : new Date().toISOString(),

        initial_fuel: form.initial_fuel,
        initial_odometer: Number(form.initial_odometer),
      };

      const result = await rentalApi.create(payload);

      setMessage(result?.message || "Rental berhasil dibuat.");

      setShowModal(false);
      setForm(initialForm);

      await loadData();
    } catch (err) {
      setError(err.message || "Gagal membuat rental.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div>
      {/* Header */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <p className="font-bold uppercase tracking-wide text-blue-600">
            Manajemen Transaksi
          </p>

          <h1 className="mt-2 text-3xl font-black text-slate-900">Rental</h1>

          <p className="mt-2 text-slate-500">
            Kelola proses penyerahan kendaraan kepada customer.
          </p>
        </div>

        <div className="flex flex-wrap gap-3">
          <button
            type="button"
            onClick={loadData}
            disabled={loading}
            className="rounded-xl border border-slate-200 bg-white px-5 py-3 font-semibold text-slate-700 transition hover:bg-slate-50 disabled:opacity-50"
          >
            {loading ? "Memuat..." : "↻ Refresh"}
          </button>

          <button
            type="button"
            onClick={openCreateModal}
            className="rounded-xl bg-blue-600 px-5 py-3 font-bold text-white transition hover:bg-blue-700"
          >
            + Buat Rental
          </button>
        </div>
      </div>

      {/* Notification */}
      {error && (
        <div className="mt-6 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {message && (
        <div className="mt-6 rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {message}
        </div>
      )}

      {/* Statistics */}
      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <p className="text-sm text-slate-500">Total Rental</p>
          <h2 className="mt-2 text-2xl font-black text-slate-900">
            {rentals.length}
          </h2>
        </div>

        <div className="rounded-2xl border border-blue-200 bg-white p-5">
          <p className="text-sm text-slate-500">Rental Aktif</p>
          <h2 className="mt-2 text-2xl font-black text-blue-600">
            {activeCount}
          </h2>
        </div>

        <div className="rounded-2xl border border-green-200 bg-white p-5">
          <p className="text-sm text-slate-500">Rental Selesai</p>
          <h2 className="mt-2 text-2xl font-black text-green-600">
            {completedCount}
          </h2>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <p className="text-sm text-slate-500">Status Lainnya</p>
          <h2 className="mt-2 text-2xl font-black text-slate-600">
            {otherCount}
          </h2>
        </div>
      </div>

      {/* Rental Table */}
      <div className="mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5">
          <h2 className="font-bold text-slate-900">
            Daftar Rental ({rentals.length})
          </h2>
        </div>

        <div className="overflow-x-auto">
          {loading ? (
            <p className="p-6 text-slate-500">Memuat data rental...</p>
          ) : rentals.length === 0 ? (
            <p className="p-6 text-slate-500">Belum ada data rental.</p>
          ) : (
            <table className="w-full min-w-[900px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-5 py-4">ID Rental</th>
                  <th className="px-5 py-4">Reservasi</th>
                  <th className="px-5 py-4">Admin Pencatat</th>
                  <th className="px-5 py-4">Tanggal Pickup</th>
                  <th className="px-5 py-4">Bahan Bakar</th>
                  <th className="px-5 py-4">Odometer Awal</th>
                  <th className="px-5 py-4">Status</th>
                  <th className="px-5 py-4">Aksi</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {rentals.map((rental) => (
                  <tr key={rental.id} className="transition hover:bg-slate-50">
                    <td className="px-5 py-4 font-bold text-slate-900">
                      #{rental.id}
                    </td>

                    <td className="px-5 py-4">#{rental.reservation_id}</td>

                    <td className="px-5 py-4">#{rental.staff_id}</td>

                    <td className="px-5 py-4 text-slate-600">
                      {formatDate(rental.pickup_date)}
                    </td>

                    <td className="px-5 py-4 capitalize">
                      {rental.initial_fuel || "-"}
                    </td>

                    <td className="px-5 py-4">
                      {Number(rental.initial_odometer || 0).toLocaleString(
                        "id-ID",
                      )}{" "}
                      km
                    </td>

                    <td className="px-5 py-4">
                      <StatusBadge status={rental.status} />
                    </td>

                    <td className="px-5 py-4">
                      <button
                        type="button"
                        onClick={() => setSelectedRental(rental)}
                        className="rounded-lg bg-blue-50 px-3 py-2 font-semibold text-blue-700 transition hover:bg-blue-100"
                      >
                        Detail
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {/* Create Rental Modal */}
      {showModal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-900/50 p-4"
          onClick={() => !submitting && setShowModal(false)}
        >
          <div
            className="my-8 w-full max-w-xl rounded-2xl bg-white p-6 shadow-xl"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-xl font-black text-slate-900">
                  Buat Rental Baru
                </h2>

                <p className="mt-1 text-sm text-slate-500">
                  Catat penyerahan kendaraan kepada customer.
                </p>
              </div>

              <button
                type="button"
                disabled={submitting}
                onClick={() => setShowModal(false)}
                className="rounded-lg px-3 py-1 text-2xl text-slate-400 hover:bg-slate-100"
              >
                ×
              </button>
            </div>

            {availableReservations.length === 0 ? (
              <div className="mt-6 rounded-xl border border-yellow-200 bg-yellow-50 p-4">
                <p className="font-semibold text-yellow-800">
                  Tidak ada reservasi yang bisa diproses.
                </p>

                <p className="mt-1 text-sm text-yellow-700">
                  Pastikan ada reservasi dengan status confirmed dan belum
                  memiliki rental.
                </p>
              </div>
            ) : (
              <form onSubmit={handleCreateRental} className="mt-6 space-y-5">
                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Reservasi
                  </label>

                  <select
                    name="reservation_id"
                    value={form.reservation_id}
                    onChange={updateForm}
                    required
                    className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
                  >
                    <option value="">Pilih reservasi</option>

                    {availableReservations.map((reservation) => (
                      <option key={reservation.id} value={reservation.id}>
                        Reservasi #{reservation.id} —{" "}
                        {reservation.customer_name ||
                          `Customer #${reservation.user_id || "-"}`}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Tanggal Pickup
                  </label>

                  <input
                    type="datetime-local"
                    name="pickup_date"
                    value={form.pickup_date}
                    onChange={updateForm}
                    className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
                  />

                  <p className="mt-1 text-xs text-slate-500">
                    Kosongkan jika menggunakan waktu saat ini.
                  </p>
                </div>

                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Bahan Bakar Awal
                  </label>

                  <select
                    name="initial_fuel"
                    value={form.initial_fuel}
                    onChange={updateForm}
                    required
                    className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
                  >
                    <option value="Empty">Kosong</option>
                    <option value="1/4">1/4</option>
                    <option value="1/2">1/2</option>
                    <option value="3/4">3/4</option>
                    <option value="Full">Penuh</option>
                  </select>
                </div>

                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Odometer Awal (km)
                  </label>

                  <input
                    type="number"
                    name="initial_odometer"
                    value={form.initial_odometer}
                    onChange={updateForm}
                    min="0"
                    step="1"
                    required
                    placeholder="Contoh: 15000"
                    className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
                  />
                </div>

                <div className="flex justify-end gap-3 border-t border-slate-100 pt-5">
                  <button
                    type="button"
                    disabled={submitting}
                    onClick={() => setShowModal(false)}
                    className="rounded-xl border border-slate-200 px-5 py-3 font-semibold text-slate-600 hover:bg-slate-50 disabled:opacity-50"
                  >
                    Batal
                  </button>

                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-blue-600 px-5 py-3 font-bold text-white hover:bg-blue-700 disabled:opacity-50"
                  >
                    {submitting ? "Menyimpan..." : "Buat Rental"}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}

      {/* Detail Modal */}
      {selectedRental && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4"
          onClick={() => setSelectedRental(null)}
        >
          <div
            className="w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-black text-slate-900">
                Detail Rental #{selectedRental.id}
              </h2>

              <button
                type="button"
                onClick={() => setSelectedRental(null)}
                className="rounded-lg px-3 py-1 text-xl text-slate-400 hover:bg-slate-100"
              >
                ×
              </button>
            </div>

            <div className="mt-6 space-y-4">
              <DetailItem
                label="ID Reservasi"
                value={`#${selectedRental.reservation_id}`}
              />

              <DetailItem
                label="ID Admin Pencatat"
                value={`#${selectedRental.staff_id}`}
              />

              <DetailItem
                label="Tanggal Pickup"
                value={formatDate(selectedRental.pickup_date)}
              />

              <DetailItem
                label="Bahan Bakar Awal"
                value={selectedRental.initial_fuel}
              />

              <DetailItem
                label="Odometer Awal"
                value={`${Number(
                  selectedRental.initial_odometer || 0,
                ).toLocaleString("id-ID")} km`}
              />

              <div className="flex items-center justify-between border-t border-slate-100 pt-4">
                <span className="text-sm text-slate-500">Status</span>
                <StatusBadge status={selectedRental.status} />
              </div>

              <DetailItem
                label="Dibuat"
                value={formatDate(selectedRental.created_at)}
              />
            </div>

            <button
              type="button"
              onClick={() => setSelectedRental(null)}
              className="mt-6 w-full rounded-xl bg-blue-600 py-3 font-bold text-white transition hover:bg-blue-700"
            >
              Tutup
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
