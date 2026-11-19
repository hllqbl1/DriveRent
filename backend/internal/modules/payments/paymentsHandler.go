package payments

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxPaymentProofSize = 5 << 20 // 5 MB

type PaymentHandler interface {
	CreatePayment(w http.ResponseWriter, r *http.Request)
	GetAllPayments(w http.ResponseWriter, r *http.Request)
	GetPaymentByID(w http.ResponseWriter, r *http.Request)
	GetPaymentsByReservationID(w http.ResponseWriter, r *http.Request)
	VerifyPayment(w http.ResponseWriter, r *http.Request)
	RejectPayment(w http.ResponseWriter, r *http.Request)
}

type paymentHandler struct {
	service PaymentService
}

func NewPaymentHandler(service PaymentService) PaymentHandler {
	return &paymentHandler{
		service: service,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// Menyimpan bukti pembayaran dan mengembalikan path-nya.
func savePaymentProof(file io.Reader) (string, error) {
	// Pastikan folder upload tersedia.
	uploadDir := filepath.Join("uploads", "payments")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Baca sebagian data untuk mendeteksi tipe file.
	buffer := make([]byte, 512)
	n, err := io.ReadFull(file, buffer)

	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	var extension string

	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/webp":
		extension = ".webp"
	default:
		return "", errors.New("format bukti pembayaran harus JPG, PNG, atau WebP")
	}

	// Buat nama file acak.
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	filename := hex.EncodeToString(randomBytes) + extension
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)

	if err != nil {
		return "", err
	}

	// Hapus file jika proses penyimpanan gagal.
	success := false
	defer func() {
		_ = dst.Close()
		if !success {
			_ = os.Remove(filePath)
		}
	}()

	if _, err := dst.Write(buffer[:n]); err != nil {
		return "", err
	}

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	success = true

	return "/uploads/payments/" + filename, nil
}

// CreatePayment membuat pembayaran baru.
func (h *paymentHandler) CreatePayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Batasi seluruh request termasuk file upload.
	r.Body = http.MaxBytesReader(w, r.Body, maxPaymentProofSize+(1<<20))

	if err := r.ParseMultipartForm(maxPaymentProofSize); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Form pembayaran tidak valid atau ukuran file terlalu besar",
		)
		return
	}

	reservationID, err := strconv.ParseInt(
		r.FormValue("reservation_id"),
		10,
		64,
	)

	if err != nil || reservationID <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID reservasi tidak valid",
		)
		return
	}

	paymentMethod := strings.TrimSpace(r.FormValue("payment_method"))

	var paymentProof *string

	file, _, fileErr := r.FormFile("payment_proof")

	if fileErr == nil {
		defer file.Close()

		path, err := savePaymentProof(file)

		if err != nil {
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		paymentProof = &path
	} else if !errors.Is(fileErr, http.ErrMissingFile) {
		writeError(
			w,
			http.StatusBadRequest,
			"Gagal membaca file bukti pembayaran",
		)
		return
	}

	err = h.service.CreatePayment(
		r.Context(),
		reservationID,
		paymentMethod,
		paymentProof,
	)

	if err != nil {
		// Hapus file jika pembayaran gagal dibuat.
		if paymentProof != nil {
			_ = os.Remove("." + *paymentProof)
		}

		switch {
		case errors.Is(err, ErrReservationNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrReservationNotPending),
			errors.Is(err, ErrPaymentAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())

		case errors.Is(err, ErrInvalidPayment),
			errors.Is(err, ErrInvalidMethod):
			writeError(w, http.StatusBadRequest, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal membuat pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Pembayaran berhasil dibuat dan menunggu verifikasi",
	})
}

// GetAllPayments mengambil seluruh pembayaran.
func (h *paymentHandler) GetAllPayments(
	w http.ResponseWriter,
	r *http.Request,
) {
	paymentList, err := h.service.GetAllPayments(r.Context())

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil data pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, paymentList)
}

// GetPaymentByID mengambil detail pembayaran.
func (h *paymentHandler) GetPaymentByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID pembayaran tidak valid")
		return
	}

	payment, err := h.service.GetPaymentByID(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

// GetPaymentsByReservationID mengambil riwayat pembayaran reservasi.
func (h *paymentHandler) GetPaymentsByReservationID(
	w http.ResponseWriter,
	r *http.Request,
) {
	reservationID, err := strconv.ParseInt(
		r.PathValue("reservation_id"),
		10,
		64,
	)

	if err != nil || reservationID <= 0 {
		writeError(w, http.StatusBadRequest, "ID reservasi tidak valid")
		return
	}

	paymentList, err := h.service.GetPaymentsByReservationID(
		r.Context(),
		reservationID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil riwayat pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, paymentList)
}

// VerifyPayment memverifikasi pembayaran.
func (h *paymentHandler) VerifyPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID pembayaran tidak valid")
		return
	}

	err = h.service.VerifyPayment(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, ErrPaymentNotFound),
			errors.Is(err, ErrReservationNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrPaymentAlreadyProcessed),
			errors.Is(err, ErrReservationNotPending):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal memverifikasi pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Pembayaran berhasil diverifikasi dan reservasi dikonfirmasi",
	})
}

// RejectPayment menolak pembayaran.
func (h *paymentHandler) RejectPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID pembayaran tidak valid")
		return
	}

	err = h.service.RejectPayment(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, ErrPaymentNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrPaymentAlreadyProcessed):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal menolak pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Pembayaran berhasil ditolak",
	})
}

// Menjaga agar fmt tetap tidak diperlukan jika kode diubah.
var _ = fmt.Sprintf
