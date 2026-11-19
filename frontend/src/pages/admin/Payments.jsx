import { useCallback, useEffect, useState } from "react";
import { paymentApi } from "../../services/api";

const API_ORIGIN = (
  import.meta.env.VITE_API_URL || "http://localhost:8080/api"
).replace(/\/api\/?$/, "");

function formatRupiah(amount) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount || 0);
}

function formatDate(date) {
  if (!date) return "-";

  return new Date(date).toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

function getProofUrl(proof) {
  if (!proof) return null;

  if (proof.startsWith("http://") || proof.startsWith("https://")) {
    return proof;
  }

  return `${API_ORIGIN}${proof.startsWith("/") ? "" : "/"}${proof}`;
}

function StatusBadge({ status }) {
  const styles = {
    pending: "bg-yellow-100 text-yellow-700",
    verified: "bg-green-100 text-green-700",
    rejected: "bg-red-100 text-red-700",
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

export default function Payments() {
  const [payments, setPayments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [actionId, setActionId] = useState(null);

  const loadPayments = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const result = await paymentApi.all();
      setPayments(Array.isArray(result) ? result : []);
    } catch (err) {
      setError(err.message || "Gagal mengambil data pembayaran");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadPayments();
  }, [loadPayments]);

  async function handleVerify(payment) {
    const confirmed = window.confirm(
      `Verifikasi pembayaran untuk reservasi #${payment.reservation_id}?`,
    );

    if (!confirmed) return;

    setActionId(payment.id);
    setError("");
    setNotice("");

    try {
      await paymentApi.verify(payment.id);
      setNotice(`Pembayaran #${payment.id} berhasil diverifikasi.`);
      await loadPayments();
    } catch (err) {
      setError(err.message || "Gagal memverifikasi pembayaran");
    } finally {
      setActionId(null);
    }
  }

  async function handleReject(payment) {
    const confirmed = window.confirm(
      `Yakin ingin menolak pembayaran untuk reservasi #${payment.reservation_id}?`,
    );

    if (!confirmed) return;

    setActionId(payment.id);
    setError("");
    setNotice("");

    try {
      await paymentApi.reject(payment.id);
      setNotice(`Pembayaran #${payment.id} berhasil ditolak.`);
      await loadPayments();
    } catch (err) {
      setError(err.message || "Gagal menolak pembayaran");
    } finally {
      setActionId(null);
    }
  }

  const pendingCount = payments.filter(
    (payment) => payment.status === "pending",
  ).length;

  const verifiedCount = payments.filter(
    (payment) => payment.status === "verified",
  ).length;

  const rejectedCount = payments.filter(
    (payment) => payment.status === "rejected",
  ).length;

  return (
    <div>
      {/* Header */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <p className="font-bold uppercase tracking-wide text-blue-600">
            Manajemen Transaksi
          </p>
          <h1 className="mt-2 text-3xl font-black text-slate-900">
            Pembayaran
          </h1>
          <p className="mt-2 text-slate-500">
            Kelola dan verifikasi pembayaran pelanggan DriveRent.
          </p>
        </div>

        <button
          onClick={loadPayments}
          disabled={loading}
          className="rounded-xl border border-slate-200 bg-white px-5 py-3 font-semibold text-slate-700 transition hover:bg-slate-50 disabled:opacity-50"
        >
          {loading ? "Memuat..." : "↻ Refresh"}
        </button>
      </div>

      {/* Notifications */}
      {error && (
        <div className="mt-6 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {notice && (
        <div className="mt-6 rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {notice}
        </div>
      )}

      {/* Statistics */}
      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <p className="text-sm text-slate-500">Total Pembayaran</p>
          <h2 className="mt-2 text-2xl font-black text-slate-900">
            {payments.length}
          </h2>
        </div>

        <div className="rounded-2xl border border-yellow-200 bg-white p-5">
          <p className="text-sm text-slate-500">Menunggu Verifikasi</p>
          <h2 className="mt-2 text-2xl font-black text-yellow-600">
            {pendingCount}
          </h2>
        </div>

        <div className="rounded-2xl border border-green-200 bg-white p-5">
          <p className="text-sm text-slate-500">Terverifikasi</p>
          <h2 className="mt-2 text-2xl font-black text-green-600">
            {verifiedCount}
          </h2>
        </div>

        <div className="rounded-2xl border border-red-200 bg-white p-5">
          <p className="text-sm text-slate-500">Ditolak</p>
          <h2 className="mt-2 text-2xl font-black text-red-600">
            {rejectedCount}
          </h2>
        </div>
      </div>

      {/* Payment Table */}
      <div className="mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5">
          <h2 className="font-bold text-slate-900">
            Daftar Pembayaran ({payments.length})
          </h2>
        </div>

        <div className="overflow-x-auto">
          {loading ? (
            <p className="p-6 text-slate-500">Memuat data pembayaran...</p>
          ) : error && payments.length === 0 ? (
            <p className="p-6 text-slate-500">
              Data pembayaran tidak dapat ditampilkan.
            </p>
          ) : payments.length === 0 ? (
            <p className="p-6 text-slate-500">Belum ada data pembayaran.</p>
          ) : (
            <table className="w-full min-w-[950px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-5 py-4">ID</th>
                  <th className="px-5 py-4">Reservasi</th>
                  <th className="px-5 py-4">Nominal</th>
                  <th className="px-5 py-4">Metode</th>
                  <th className="px-5 py-4">Tanggal</th>
                  <th className="px-5 py-4">Bukti</th>
                  <th className="px-5 py-4">Status</th>
                  <th className="px-5 py-4">Aksi</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {payments.map((payment) => {
                  const proofUrl = getProofUrl(payment.payment_proof);
                  const isProcessing = actionId === payment.id;

                  return (
                    <tr key={payment.id} className="hover:bg-slate-50">
                      <td className="px-5 py-4 font-semibold text-slate-900">
                        #{payment.id}
                      </td>

                      <td className="px-5 py-4">#{payment.reservation_id}</td>

                      <td className="px-5 py-4 font-semibold text-slate-900">
                        {formatRupiah(payment.amount)}
                      </td>

                      <td className="px-5 py-4 capitalize">
                        {payment.payment_method || "-"}
                      </td>

                      <td className="px-5 py-4 text-slate-500">
                        {formatDate(payment.created_at)}
                      </td>

                      <td className="px-5 py-4">
                        {proofUrl ? (
                          <a
                            href={proofUrl}
                            target="_blank"
                            rel="noreferrer"
                            className="font-semibold text-blue-600 hover:underline"
                          >
                            Lihat Bukti
                          </a>
                        ) : (
                          <span className="text-slate-400">Tidak ada</span>
                        )}
                      </td>

                      <td className="px-5 py-4">
                        <StatusBadge status={payment.status} />
                      </td>

                      <td className="px-5 py-4">
                        {payment.status === "pending" ? (
                          <div className="flex gap-2">
                            <button
                              onClick={() => handleVerify(payment)}
                              disabled={isProcessing || actionId !== null}
                              className="rounded-lg bg-green-50 px-3 py-2 font-semibold text-green-700 hover:bg-green-100 disabled:opacity-50"
                            >
                              {isProcessing ? "Proses..." : "Verifikasi"}
                            </button>

                            <button
                              onClick={() => handleReject(payment)}
                              disabled={isProcessing || actionId !== null}
                              className="rounded-lg bg-red-50 px-3 py-2 font-semibold text-red-700 hover:bg-red-100 disabled:opacity-50"
                            >
                              Tolak
                            </button>
                          </div>
                        ) : (
                          <span className="text-slate-400">Sudah diproses</span>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
