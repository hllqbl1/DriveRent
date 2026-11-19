# DriveRent Frontend

Frontend starter untuk backend Go DriveRent. Stack: React, Vite, React Router, Tailwind CSS v4, dan Fetch API (tanpa Axios).

## Menjalankan

```bash
npm install
```

Buat `.env` dari `.env.example`:

```env
VITE_API_URL=http://localhost:8080/api
```

Lalu:

```bash
npm run dev
```

Build:

```bash
npm run build
```

## Catatan integrasi

- Login/register menggunakan `/auth/login` dan `/auth/register`.
- Daftar kendaraan menggunakan `GET /vehicles`.
- Detail kendaraan menggunakan `GET /vehicles/:id`.
- Reservasi menggunakan `POST /reservations`.
- Pembayaran menggunakan `POST /payments`.
- Endpoint dashboard admin menggunakan route yang tersedia di backend.
- Token disimpan di `localStorage` untuk starter ini. Untuk produksi, pertimbangkan cookie HttpOnly/Secure.
- Sesuaikan endpoint dan bentuk response di `src/services/api.js` atau halaman terkait jika implementasi backend berbeda.
- Backend harus mengizinkan CORS dari origin Vite (`http://localhost:5173`).
