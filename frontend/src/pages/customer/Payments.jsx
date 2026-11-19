import { useCallback, useEffect, useState } from "react";
import { paymentApi, reservationApi } from "../../services/api";
import { useAuth } from "../../context/AuthContext";

const API_ORIGIN = (
  import.meta.env.VITE_API_URL || "http://localhost:8080/api"
).replace(/\/api\/?$/, "");

function formatRupiah(value) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(Number(value) || 0);
}

function getProofUrl(path) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return `${API_ORIGIN}${path}`;
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
  const { user } = useAuth();

  const [reservations, setReservations] = useState([]);
  const [payments, setPayments] = useState([]);
  const [selectedId, setSelectedId] = useState("");

  const [paymentMethod, setPaymentMethod] = useState("transfer");
  const [paymentFile, setPaymentFile] = useState(null);
  const [preview, setPreview] = useState("");

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const loadData = useCallback(async () => {
    if (!user?.id) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError("");

    try {
      const reservationResult = await reservationApi.mine(user.id);
      const reservationList = Array.isArray(reservationResult)
        ? reservationResult
        : [];

      setReservations(reservationList);

      const paymentResults = await Promise.all(
        reservationList.map(async (reservation) => {
          try {
            const result = await paymentApi.byReservation(reservation.id);
            return Array.isArray(result) ? result : result ? [result] : [];
          } catch {
            return [];
          }
        }),
      );

      setPayments(paymentResults.flat());
    } catch (err) {
      setError(err.message || "Gagal mengambil data pembayaran.");
    } finally {
      setLoading(false);
    }
  }, [user?.id]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Membuat dan membersihkan URL preview gambar.
  useEffect(() => {
    if (!paymentFile) {
      setPreview("");
      return;
    }

    const objectUrl = URL.createObjectURL(paymentFile);
    setPreview(objectUrl);

    return () => URL.revokeObjectURL(objectUrl);
  }, [paymentFile]);

  const selectedReservation = reservations.find(
    (reservation) => String(reservation.id) === String(selectedId),
  );

  const selectedPayment = payments.find(
    (payment) =>
      String(payment.reservation_id) === String(selectedId) &&
      ["pending", "verified"].includes(payment.status),
  );

  const payableReservations = reservations.filter((reservation) => {
    const hasActivePayment = payments.some(
      (payment) =>
        Number(payment.reservation_id) === Number(reservation.id) &&
        ["pending", "verified"].includes(payment.status),
    );

    return reservation.status === "pending" && !hasActivePayment;
  });

  function handleFileChange(event) {
    const file = event.target.files?.[0];

    if (!file) {
      setPaymentFile(null);
      return;
    }

    const allowedTypes = ["image/jpeg", "image/png", "image/webp"];

    if (!allowedTypes.includes(file.type)) {
      setError("Format gambar harus JPG, PNG, atau WebP.");
      event.target.value = "";
      setPaymentFile(null);
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      setError("Ukuran gambar maksimal 5 MB.");
      event.target.value = "";
      setPaymentFile(null);
      return;
    }

    setError("");
    setPaymentFile(file);
  }

  async function handleSubmit(event) {
    event.preventDefault();

    setError("");
    setMessage("");

    if (!selectedReservation) {
      setError("Pilih reservasi terlebih dahulu.");
      return;
    }

    if (selectedPayment) {
      setError("Reservasi ini sudah memiliki pembayaran aktif.");
      return;
    }

    if (paymentMethod === "transfer" && !paymentFile) {
      setError("Silakan upload bukti transfer terlebih dahulu.");
      return;
    }

    setSubmitting(true);

    try {
      const formData = new FormData();

      formData.append("reservation_id", String(selectedId));
      formData.append("payment_method", paymentMethod);

      if (paymentFile) {
        formData.append("payment_proof", paymentFile);
      }

      const result = await paymentApi.create(formData);

      setMessage(result?.message || "Pembayaran berhasil dikirim.");

      setSelectedId("");
      setPaymentMethod("transfer");
      setPaymentFile(null);

      const fileInput = document.getElementById("payment-proof");
      if (fileInput) fileInput.value = "";

      await loadData();
    } catch (err) {
      setError(err.message || "Gagal mengirim pembayaran.");
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

        <h1 className="mt-2 text-3xl font-black text-slate-900">Pembayaran</h1>

        <p className="mt-2 text-slate-500">
          Bayar reservasi dan pantau status pembayaranmu.
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
        className="mt-6 max-w-2xl space-y-5 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
      >
        <h2 className="text-lg font-bold text-slate-900">Kirim Pembayaran</h2>

        <div>
          <label className="mb-2 block text-sm font-semibold text-slate-700">
            Pilih Reservasi
          </label>

          <select
            value={selectedId}
            onChange={(event) => setSelectedId(event.target.value)}
            required
            className="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:border-blue-500"
          >
            <option value="">Pilih reservasi yang belum dibayar</option>

            {payableReservations.map((reservation) => (
              <option key={reservation.id} value={reservation.id}>
                Reservasi #{reservation.id} —{" "}
                {formatRupiah(reservation.total_price)}
              </option>
            ))}
          </select>

          {payableReservations.length === 0 && !loading && (
            <p className="mt-2 text-sm text-slate-500">
              Tidak ada reservasi yang menunggu pembayaran.
            </p>
          )}
        </div>

        {selectedReservation && (
          <div className="rounded-xl bg-blue-50 p-4">
            <p className="text-sm text-slate-600">
              Detail tagihan reservasi #{selectedReservation.id}
            </p>

            <p className="mt-2 text-2xl font-black text-blue-700">
              {formatRupiah(selectedReservation.total_price)}
            </p>

            <p className="mt-1 text-xs text-slate-500">
              Total tagihan berdasarkan data reservasi dari backend.
            </p>
          </div>
        )}

        <div>
          <label className="mb-2 block text-sm font-semibold text-slate-700">
            Metode Pembayaran
          </label>

          <select
            value={paymentMethod}
            onChange={(event) => setPaymentMethod(event.target.value)}
            className="w-full rounded-xl border border-slate-200 px-4 py-3"
          >
            <option value="transfer">Transfer Bank</option>
            <option value="cash">Tunai</option>
          </select>
        </div>

        {paymentMethod === "transfer" && (
          <div>
            <label
              htmlFor="payment-proof"
              className="mb-2 block text-sm font-semibold text-slate-700"
            >
              Upload Bukti Pembayaran
            </label>

            <input
              id="payment-proof"
              type="file"
              accept="image/jpeg,image/png,image/webp"
              onChange={handleFileChange}
              required
              className="w-full rounded-xl border border-slate-200 p-3 text-sm file:mr-4 file:rounded-lg file:border-0 file:bg-blue-50 file:px-4 file:py-2 file:font-semibold file:text-blue-700"
            />

            <p className="mt-2 text-xs text-slate-500">
              Format JPG, PNG, atau WebP. Maksimal 5 MB.
            </p>

            {preview && (
              <div className="mt-4">
                <p className="mb-2 text-sm font-semibold text-slate-700">
                  Preview Bukti Pembayaran
                </p>

                <img
                  src={preview}
                  alt="Preview bukti pembayaran"
                  className="max-h-72 rounded-xl border border-slate-200 object-contain"
                />

                <button
                  type="button"
                  onClick={() => {
                    setPaymentFile(null);
                    const input = document.getElementById("payment-proof");
                    if (input) input.value = "";
                  }}
                  className="mt-2 text-sm font-semibold text-red-600 hover:text-red-700"
                >
                  Hapus gambar
                </button>
              </div>
            )}
          </div>
        )}

        {paymentMethod === "cash" && (
          <p className="rounded-xl bg-yellow-50 p-4 text-sm text-yellow-800">
            Pembayaran tunai akan dicatat dan menunggu konfirmasi admin.
          </p>
        )}

        <button
          type="submit"
          disabled={submitting || loading || !selectedId}
          className="rounded-xl bg-blue-600 px-6 py-3 font-bold text-white transition hover:bg-blue-700 disabled:opacity-50"
        >
          {submitting ? "Mengirim..." : "Kirim Pembayaran"}
        </button>
      </form>

      <div className="mt-9">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-black text-slate-900">
            Riwayat Pembayaran
          </h2>

          <button
            type="button"
            onClick={loadData}
            className="text-sm font-semibold text-blue-700"
          >
            Refresh
          </button>
        </div>

        {loading ? (
          <p className="mt-4 text-slate-500">Memuat pembayaran...</p>
        ) : payments.length === 0 ? (
          <p className="mt-4 rounded-xl border border-dashed border-slate-300 p-6 text-sm text-slate-500">
            Belum ada riwayat pembayaran.
          </p>
        ) : (
          <div className="mt-4 space-y-3">
            {payments.map((payment) => (
              <div
                key={payment.id}
                className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-slate-200 bg-white p-5"
              >
                <div>
                  <p className="font-bold text-slate-900">
                    Pembayaran #{payment.id}
                  </p>

                  <p className="mt-1 text-sm text-slate-500">
                    Reservasi #{payment.reservation_id}
                  </p>

                  <p className="mt-1 text-sm text-slate-500">
                    {payment.payment_method}
                  </p>

                  <p className="mt-2 font-bold text-blue-700">
                    {formatRupiah(payment.amount)}
                  </p>

                  {payment.payment_proof && (
                    <a
                      href={getProofUrl(payment.payment_proof)}
                      target="_blank"
                      rel="noreferrer"
                      className="mt-2 inline-block text-sm font-semibold text-blue-600 hover:underline"
                    >
                      Lihat bukti pembayaran
                    </a>
                  )}
                </div>

                <StatusBadge status={payment.status} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
