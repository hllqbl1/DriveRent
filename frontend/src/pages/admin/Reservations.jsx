import { useEffect, useState } from "react";
import { reservationApi } from "../../services/api";

export default function Reservations() {
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [filter, setFilter] = useState("all");
  const [processing, setProcessing] = useState(null);
  const [selected, setSelected] = useState(null);

  async function loadReservations() {
    setLoading(true);
    setError("");

    try {
      const result = await reservationApi.all();
      setReservations(Array.isArray(result) ? result : []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadReservations();
  }, []);

  async function handleStatus(id, status) {
    const action = status === "confirmed" ? "menyetujui" : "menolak";

    if (!window.confirm(`Yakin ingin ${action} reservasi #${id}?`)) {
      return;
    }

    setProcessing(id);
    setError("");
    setSuccess("");

    try {
      await reservationApi.updateStatus(id, status);

      setSuccess(
        `Reservasi #${id} berhasil ${
          status === "confirmed" ? "dikonfirmasi" : "ditolak"
        }.`,
      );

      setSelected(null);
      await loadReservations();
    } catch (err) {
      setError(err.message);
    } finally {
      setProcessing(null);
    }
  }

  async function handleCancel(id) {
    if (!window.confirm(`Yakin ingin membatalkan reservasi #${id}?`)) {
      return;
    }

    setProcessing(id);
    setError("");
    setSuccess("");

    try {
      await reservationApi.cancel(id);
      setSuccess(`Reservasi #${id} berhasil dibatalkan.`);
      setSelected(null);
      await loadReservations();
    } catch (err) {
      setError(err.message);
    } finally {
      setProcessing(null);
    }
  }

  const filteredReservations = reservations.filter((item) => {
    return filter === "all" || item.status === filter;
  });

  function formatDate(date) {
    if (!date) return "-";

    return new Date(date).toLocaleDateString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      timeZone: "UTC",
    });
  }

  function formatPrice(price) {
    return Number(price || 0).toLocaleString("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    });
  }

  function statusStyle(status) {
    const styles = {
      pending: "bg-amber-50 text-amber-700",
      confirmed: "bg-green-50 text-green-700",
      rejected: "bg-red-50 text-red-700",
      cancelled: "bg-slate-100 text-slate-600",
      completed: "bg-blue-50 text-blue-700",
    };

    return styles[status] || "bg-slate-100 text-slate-600";
  }

  const pendingCount = reservations.filter(
    (item) => item.status === "pending",
  ).length;

  const confirmedCount = reservations.filter(
    (item) => item.status === "confirmed",
  ).length;

  const rejectedCount = reservations.filter(
    (item) => item.status === "rejected",
  ).length;

  return (
    <div className="space-y-6">
      {/* HEADER */}
      <div>
        <p className="text-sm font-bold uppercase tracking-wider text-blue-700">
          Manajemen Pemesanan
        </p>

        <h1 className="mt-2 text-3xl font-black text-slate-900">Reservasi</h1>

        <p className="mt-2 text-sm text-slate-500">
          Pantau dan proses reservasi kendaraan customer.
        </p>
      </div>

      {/* STATISTICS */}
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard title="Total Reservasi" value={reservations.length} />
        <StatCard title="Menunggu" value={pendingCount} color="amber" />
        <StatCard title="Dikonfirmasi" value={confirmedCount} color="green" />
        <StatCard title="Ditolak" value={rejectedCount} color="red" />
      </div>

      {/* NOTIFICATIONS */}
      {error && (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {success && (
        <div className="rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {success}
        </div>
      )}

      {/* FILTER */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="font-bold text-slate-900">Daftar Reservasi</h2>

          <p className="mt-1 text-sm text-slate-500">
            {filteredReservations.length} reservasi ditemukan
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="rounded-xl border border-slate-300 bg-white px-4 py-2.5 text-sm outline-none focus:border-blue-600"
          >
            <option value="all">Semua Status</option>
            <option value="pending">Pending</option>
            <option value="confirmed">Confirmed</option>
            <option value="rejected">Rejected</option>
            <option value="cancelled">Cancelled</option>
            <option value="completed">Completed</option>
          </select>

          <button
            onClick={loadReservations}
            className="rounded-xl border border-slate-300 px-4 py-2.5 text-sm font-semibold text-slate-700 hover:bg-slate-50"
          >
            Refresh
          </button>
        </div>
      </div>

      {/* TABLE */}
      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        {loading ? (
          <p className="p-10 text-center text-slate-500">Memuat reservasi...</p>
        ) : filteredReservations.length === 0 ? (
          <p className="p-10 text-center text-slate-500">
            Tidak ada reservasi ditemukan.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[950px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-5 py-4">ID</th>
                  <th className="px-5 py-4">Customer</th>
                  <th className="px-5 py-4">Kendaraan</th>
                  <th className="px-5 py-4">Tanggal Sewa</th>
                  <th className="px-5 py-4">Total Harga</th>
                  <th className="px-5 py-4">Status</th>
                  <th className="px-5 py-4">Aksi</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {filteredReservations.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-50">
                    <td className="px-5 py-4 font-bold">#{item.id}</td>

                    <td className="px-5 py-4">Customer #{item.user_id}</td>

                    <td className="px-5 py-4">Kendaraan #{item.vehicle_id}</td>

                    <td className="px-5 py-4">
                      <p>{formatDate(item.start_date)}</p>
                      <p className="mt-1 text-xs text-slate-500">
                        sampai {formatDate(item.end_date)}
                      </p>
                    </td>

                    <td className="px-5 py-4 font-semibold">
                      {formatPrice(item.total_price)}
                    </td>

                    <td className="px-5 py-4">
                      <span
                        className={`rounded-full px-3 py-1 text-xs font-bold ${statusStyle(item.status)}`}
                      >
                        {item.status}
                      </span>
                    </td>

                    <td className="px-5 py-4">
                      <div className="flex flex-wrap gap-2">
                        <button
                          onClick={() => setSelected(item)}
                          className="rounded-lg bg-slate-100 px-3 py-2 text-xs font-bold text-slate-700 hover:bg-slate-200"
                        >
                          Detail
                        </button>

                        {item.status === "pending" && (
                          <>
                            <button
                              disabled={processing === item.id}
                              onClick={() => handleStatus(item.id, "confirmed")}
                              className="rounded-lg bg-green-50 px-3 py-2 text-xs font-bold text-green-700 hover:bg-green-100 disabled:opacity-50"
                            >
                              Setujui
                            </button>

                            <button
                              disabled={processing === item.id}
                              onClick={() => handleStatus(item.id, "rejected")}
                              className="rounded-lg bg-red-50 px-3 py-2 text-xs font-bold text-red-700 hover:bg-red-100 disabled:opacity-50"
                            >
                              Tolak
                            </button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* DETAIL MODAL */}
      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 p-4">
          <div className="w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-black text-slate-900">
                Detail Reservasi #{selected.id}
              </h2>

              <button
                onClick={() => setSelected(null)}
                className="rounded-lg px-3 py-2 text-slate-500 hover:bg-slate-100"
              >
                Tutup
              </button>
            </div>

            <div className="mt-6 space-y-4">
              <DetailRow label="Customer" value={`#${selected.user_id}`} />
              <DetailRow label="Kendaraan" value={`#${selected.vehicle_id}`} />
              <DetailRow
                label="Tanggal Mulai"
                value={formatDate(selected.start_date)}
              />
              <DetailRow
                label="Tanggal Selesai"
                value={formatDate(selected.end_date)}
              />
              <DetailRow label="Durasi" value={`${selected.total_days} hari`} />
              <DetailRow
                label="Total Harga"
                value={formatPrice(selected.total_price)}
              />
              <DetailRow label="Status" value={selected.status} />
            </div>

            {selected.status === "pending" && (
              <div className="mt-7 flex flex-wrap gap-3">
                <button
                  disabled={processing === selected.id}
                  onClick={() => handleStatus(selected.id, "confirmed")}
                  className="flex-1 rounded-xl bg-green-600 px-4 py-3 font-bold text-white hover:bg-green-700 disabled:opacity-50"
                >
                  Setujui
                </button>

                <button
                  disabled={processing === selected.id}
                  onClick={() => handleStatus(selected.id, "rejected")}
                  className="flex-1 rounded-xl bg-red-600 px-4 py-3 font-bold text-white hover:bg-red-700 disabled:opacity-50"
                >
                  Tolak
                </button>
              </div>
            )}

            {["pending", "confirmed"].includes(selected.status) && (
              <button
                disabled={processing === selected.id}
                onClick={() => handleCancel(selected.id)}
                className="mt-3 w-full rounded-xl border border-red-200 px-4 py-3 font-bold text-red-600 hover:bg-red-50 disabled:opacity-50"
              >
                Batalkan Reservasi
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

function StatCard({ title, value, color = "blue" }) {
  const styles = {
    blue: "bg-blue-50 text-blue-700",
    amber: "bg-amber-50 text-amber-700",
    green: "bg-green-50 text-green-700",
    red: "bg-red-50 text-red-700",
  };

  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      <p className="text-sm text-slate-500">{title}</p>
      <p className="mt-3 text-3xl font-black text-slate-900">{value}</p>
      <span
        className={`mt-3 inline-block rounded-lg px-2.5 py-1 text-xs font-bold ${styles[color]}`}
      >
        Data reservasi
      </span>
    </div>
  );
}

function DetailRow({ label, value }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-slate-100 pb-3">
      <span className="text-sm text-slate-500">{label}</span>
      <span className="text-right text-sm font-semibold text-slate-900">
        {value}
      </span>
    </div>
  );
}
