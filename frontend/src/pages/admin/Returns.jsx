import { useCallback, useEffect, useState } from "react";
import { returnsApi, rentalApi } from "../../services/api";

function formatDate(date) {
  if (!date) return "-";

  return new Date(date).toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

function formatRupiah(value) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(Number(value) || 0);
}

function StatusBadge({ status }) {
  return (
    <span className="rounded-full bg-green-100 px-3 py-1 text-xs font-bold text-green-700">
      {status || "-"}
    </span>
  );
}

const initialForm = {
  rental_id: "",
  final_fuel: "Full",
  final_odometer: "",
  late_fee: "0",
  damage_fee: "0",
  notes: "",
  vehicle_status: "available",
};

export default function Returns() {
  const [returns, setReturns] = useState([]);
  const [rentals, setRentals] = useState([]);
  const [form, setForm] = useState(initialForm);
  const [showForm, setShowForm] = useState(false);
  const [selectedReturn, setSelectedReturn] = useState(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const loadData = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const [returnData, rentalData] = await Promise.all([
        returnsApi.all(),
        rentalApi.all(),
      ]);

      setReturns(Array.isArray(returnData) ? returnData : []);
      setRentals(Array.isArray(rentalData) ? rentalData : []);
    } catch (err) {
      setError(err.message || "Gagal mengambil data pengembalian");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const activeRentals = rentals.filter((rental) => rental.status === "active");

  const totalLateFee = returns.reduce(
    (total, item) => total + Number(item.late_fee || 0),
    0,
  );

  const totalDamageFee = returns.reduce(
    (total, item) => total + Number(item.damage_fee || 0),
    0,
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
    setSuccess("");

    if (!form.rental_id) {
      setError("Pilih rental terlebih dahulu.");
      return;
    }

    const payload = {
      rental_id: Number(form.rental_id),
      final_fuel: form.final_fuel,
      final_odometer: Number(form.final_odometer),
      late_fee: Number(form.late_fee),
      damage_fee: Number(form.damage_fee),
      notes: form.notes.trim() || null,
      vehicle_status: form.vehicle_status,
    };

    setSubmitting(true);

    try {
      await returnsApi.create(payload);

      setSuccess("Pengembalian berhasil dicatat.");
      setForm(initialForm);
      setShowForm(false);

      await loadData();
    } catch (err) {
      setError(err.message || "Gagal mencatat pengembalian.");
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

          <h1 className="mt-2 text-3xl font-black text-slate-900">
            Pengembalian
          </h1>

          <p className="mt-2 text-slate-500">
            Kelola proses pengembalian kendaraan DriveRent.
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          <button
            onClick={loadData}
            disabled={loading}
            className="rounded-xl border border-slate-200 bg-white px-4 py-3 font-semibold text-slate-700 hover:bg-slate-50 disabled:opacity-50"
          >
            Refresh
          </button>

          <button
            onClick={() => {
              setError("");
              setSuccess("");
              setShowForm(true);
            }}
            className="rounded-xl bg-blue-600 px-5 py-3 font-bold text-white hover:bg-blue-700"
          >
            + Catat Pengembalian
          </button>
        </div>
      </div>

      {/* Notifications */}
      {error && (
        <div className="mt-5 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {success && (
        <div className="mt-5 rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {success}
        </div>
      )}

      {/* Statistics */}
      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <p className="text-sm text-slate-500">Total Pengembalian</p>
          <h2 className="mt-2 text-2xl font-black text-slate-900">
            {returns.length}
          </h2>
        </div>

        <div className="rounded-2xl border border-orange-200 bg-white p-5">
          <p className="text-sm text-slate-500">Total Denda Keterlambatan</p>
          <h2 className="mt-2 text-2xl font-black text-orange-600">
            {formatRupiah(totalLateFee)}
          </h2>
        </div>

        <div className="rounded-2xl border border-red-200 bg-white p-5">
          <p className="text-sm text-slate-500">Total Biaya Kerusakan</p>
          <h2 className="mt-2 text-2xl font-black text-red-600">
            {formatRupiah(totalDamageFee)}
          </h2>
        </div>
      </div>

      {/* Return Table */}
      <div className="mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5">
          <h2 className="font-bold text-slate-900">
            Riwayat Pengembalian ({returns.length})
          </h2>
        </div>

        <div className="overflow-x-auto">
          {loading ? (
            <p className="p-6 text-slate-500">Memuat data...</p>
          ) : returns.length === 0 ? (
            <p className="p-6 text-slate-500">Belum ada data pengembalian.</p>
          ) : (
            <table className="w-full min-w-[850px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                <tr>
                  <th className="px-5 py-4">ID</th>
                  <th className="px-5 py-4">ID Rental</th>
                  <th className="px-5 py-4">Tanggal Kembali</th>
                  <th className="px-5 py-4">Odometer Akhir</th>
                  <th className="px-5 py-4">Total Biaya</th>
                  <th className="px-5 py-4">Status</th>
                  <th className="px-5 py-4">Aksi</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {returns.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-50">
                    <td className="px-5 py-4 font-bold">#{item.id}</td>
                    <td className="px-5 py-4">#{item.rental_id}</td>
                    <td className="px-5 py-4 text-slate-600">
                      {formatDate(item.return_date)}
                    </td>
                    <td className="px-5 py-4">
                      {Number(item.final_odometer || 0).toLocaleString("id-ID")}{" "}
                      km
                    </td>
                    <td className="px-5 py-4 font-semibold">
                      {formatRupiah(
                        Number(item.late_fee || 0) +
                          Number(item.damage_fee || 0),
                      )}
                    </td>
                    <td className="px-5 py-4">
                      <StatusBadge status={item.status} />
                    </td>
                    <td className="px-5 py-4">
                      <button
                        onClick={() => setSelectedReturn(item)}
                        className="rounded-lg bg-blue-50 px-3 py-2 font-semibold text-blue-700 hover:bg-blue-100"
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

      {/* Create Return Modal */}
      {showForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-900/50 p-4">
          <form
            onSubmit={handleSubmit}
            className="my-8 w-full max-w-xl rounded-2xl bg-white p-6 shadow-xl"
          >
            <div className="mb-6 flex items-center justify-between">
              <h2 className="text-xl font-black text-slate-900">
                Catat Pengembalian
              </h2>

              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="text-2xl text-slate-400 hover:text-slate-700"
              >
                ×
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-semibold text-slate-700">
                  Rental Aktif
                </label>

                <select
                  name="rental_id"
                  value={form.rental_id}
                  onChange={handleChange}
                  required
                  className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
                >
                  <option value="">Pilih rental</option>
                  {activeRentals.map((rental) => (
                    <option key={rental.id} value={rental.id}>
                      Rental #{rental.id} — Reservasi #{rental.reservation_id}
                    </option>
                  ))}
                </select>

                {activeRentals.length === 0 && (
                  <p className="mt-2 text-xs text-orange-600">
                    Tidak ada rental aktif yang bisa dikembalikan.
                  </p>
                )}
              </div>

              <div>
                <label className="mb-2 block text-sm font-semibold text-slate-700">
                  Bahan Bakar Akhir
                </label>

                <select
                  name="final_fuel"
                  value={form.final_fuel}
                  onChange={handleChange}
                  required
                  className="w-full rounded-xl border border-slate-200 px-4 py-3"
                >
                  <option value="Full">Full</option>
                  <option value="3/4">3/4</option>
                  <option value="1/2">1/2</option>
                  <option value="1/4">1/4</option>
                  <option value="Empty">Empty</option>
                </select>
              </div>

              <div>
                <label className="mb-2 block text-sm font-semibold text-slate-700">
                  Odometer Akhir (km)
                </label>

                <input
                  type="number"
                  name="final_odometer"
                  min="0"
                  value={form.final_odometer}
                  onChange={handleChange}
                  required
                  className="w-full rounded-xl border border-slate-200 px-4 py-3"
                  placeholder="Contoh: 15000"
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Denda Keterlambatan (Rp)
                  </label>

                  <input
                    type="number"
                    name="late_fee"
                    min="0"
                    value={form.late_fee}
                    onChange={handleChange}
                    required
                    className="w-full rounded-xl border border-slate-200 px-4 py-3"
                  />
                </div>

                <div>
                  <label className="mb-2 block text-sm font-semibold text-slate-700">
                    Biaya Kerusakan (Rp)
                  </label>

                  <input
                    type="number"
                    name="damage_fee"
                    min="0"
                    value={form.damage_fee}
                    onChange={handleChange}
                    required
                    className="w-full rounded-xl border border-slate-200 px-4 py-3"
                  />
                </div>
              </div>

              <div>
                <label className="mb-2 block text-sm font-semibold text-slate-700">
                  Kondisi Kendaraan
                </label>

                <select
                  name="vehicle_status"
                  value={form.vehicle_status}
                  onChange={handleChange}
                  required
                  className="w-full rounded-xl border border-slate-200 px-4 py-3"
                >
                  <option value="available">
                    Available — Kendaraan siap disewakan
                  </option>
                  <option value="maintenance">
                    Maintenance — Perlu perawatan
                  </option>
                </select>
              </div>

              <div>
                <label className="mb-2 block text-sm font-semibold text-slate-700">
                  Catatan
                </label>

                <textarea
                  name="notes"
                  value={form.notes}
                  onChange={handleChange}
                  rows="3"
                  className="w-full rounded-xl border border-slate-200 px-4 py-3"
                  placeholder="Catatan kondisi kendaraan..."
                />
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-3">
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="rounded-xl border border-slate-200 px-5 py-3 font-semibold text-slate-600 hover:bg-slate-50"
              >
                Batal
              </button>

              <button
                type="submit"
                disabled={submitting || activeRentals.length === 0}
                className="rounded-xl bg-blue-600 px-5 py-3 font-bold text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {submitting ? "Menyimpan..." : "Simpan Pengembalian"}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Detail Modal */}
      {selectedReturn && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4"
          onClick={() => setSelectedReturn(null)}
        >
          <div
            className="w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-black text-slate-900">
                Detail Pengembalian #{selectedReturn.id}
              </h2>

              <button
                onClick={() => setSelectedReturn(null)}
                className="text-2xl text-slate-400"
              >
                ×
              </button>
            </div>

            <div className="mt-6 space-y-4">
              <DetailItem
                label="ID Rental"
                value={`#${selectedReturn.rental_id}`}
              />
              <DetailItem
                label="Tanggal Kembali"
                value={formatDate(selectedReturn.return_date)}
              />
              <DetailItem
                label="Bahan Bakar Akhir"
                value={selectedReturn.final_fuel}
              />
              <DetailItem
                label="Odometer Akhir"
                value={`${Number(selectedReturn.final_odometer || 0).toLocaleString("id-ID")} km`}
              />
              <DetailItem
                label="Denda Keterlambatan"
                value={formatRupiah(selectedReturn.late_fee)}
              />
              <DetailItem
                label="Biaya Kerusakan"
                value={formatRupiah(selectedReturn.damage_fee)}
              />
              <DetailItem label="Catatan" value={selectedReturn.notes || "-"} />

              <div className="flex items-center justify-between border-t border-slate-100 pt-4">
                <span className="text-sm text-slate-500">Status</span>
                <StatusBadge status={selectedReturn.status} />
              </div>
            </div>

            <button
              onClick={() => setSelectedReturn(null)}
              className="mt-6 w-full rounded-xl bg-blue-600 py-3 font-bold text-white hover:bg-blue-700"
            >
              Tutup
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

function DetailItem({ label, value }) {
  return (
    <div className="flex items-start justify-between gap-4">
      <span className="text-sm text-slate-500">{label}</span>
      <span className="max-w-[60%] text-right text-sm font-semibold text-slate-900">
        {value || "-"}
      </span>
    </div>
  );
}
