# DriveRent – Vehicle Rental Management System

DriveRent adalah aplikasi berbasis web untuk mengelola penyewaan kendaraan secara digital. Aplikasi ini dirancang untuk membantu pengelola rental dalam mengelola kendaraan, pengguna, reservasi, pembayaran, penyewaan, dan pengembalian kendaraan melalui satu sistem terintegrasi.

Project ini dikembangkan sebagai portfolio project Fullstack Web Development menggunakan React dan Golang.

## ✨ Features

### Authentication & Authorization
- Register dan login pengguna.
- Autentikasi menggunakan JSON Web Token (JWT).
- Password disimpan menggunakan bcrypt.
- Pembagian hak akses berdasarkan role:
  - **Admin:** Mengelola seluruh operasional rental.
  - **Customer:** Melihat kendaraan, membuat reservasi, dan melakukan pembayaran.

### Vehicle Management
- Menampilkan daftar kendaraan.
- Menambahkan kendaraan baru.
- Mengubah informasi kendaraan.
- Mengunggah foto kendaraan.
- Mengelola harga sewa.
- Mengelola status kendaraan.

### Reservation Management
- Customer dapat memilih kendaraan.
- Customer dapat menentukan tanggal penyewaan.
- Sistem mencatat reservasi.
- Customer dapat melihat riwayat reservasi.
- Admin dapat mengelola dan memperbarui status reservasi.

### Payment Management
- Customer dapat mengirim pembayaran berdasarkan reservasi.
- Mendukung metode pembayaran transfer bank dan tunai.
- Upload bukti pembayaran.
- Admin dapat melihat dan memverifikasi pembayaran.
- Admin dapat menolak pembayaran.
- Status pembayaran terintegrasi dengan status reservasi.

### Rental Management
- Admin dapat membuat data rental dari reservasi yang telah dikonfirmasi.
- Mencatat tanggal pickup kendaraan.
- Mencatat kondisi bahan bakar awal.
- Mencatat odometer awal.
- Memantau status rental.

### Return Management
- Admin dapat mencatat pengembalian kendaraan.
- Mencatat kondisi bahan bakar akhir.
- Mencatat odometer akhir.
- Mencatat denda keterlambatan.
- Mencatat biaya kerusakan.
- Mengatur status kendaraan setelah pengembalian.

### Dashboard
- Ringkasan data operasional.
- Statistik kendaraan.
- Statistik reservasi.
- Statistik pembayaran.
- Statistik rental dan pengembalian.

## 🛠️ Tech Stack

### Frontend
- React.js
- Vite
- React Router
- Tailwind CSS v4
- JavaScript
- Fetch API

### Backend
- Golang
- Go Standard Library (`net/http`)
- PostgreSQL
- pgx
- JWT Authentication
- bcrypt
- godotenv
- golang-migrate

## 📁 Project Structure

```text
DriveRentNew/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── middleware/
│   │   ├── modules/
│   │   │   ├── users/
│   │   │   ├── vehicles/
│   │   │   ├── reservations/
│   │   │   ├── payments/
│   │   │   ├── rentals/
│   │   │   └── returns/
│   │   └── routes/
│   ├── migrations/
│   ├── uploads/
│   ├── .env
│   └── go.mod
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── context/
│   │   ├── pages/
│   │   ├── services/
│   │   └── App.jsx
│   ├── package.json
│   └── vite.config.js
│
└── README.md
```

## ⚙️ Installation & Setup

### Prerequisites

Pastikan perangkat sudah memiliki:
- Go
- Node.js dan npm
- PostgreSQL
- Git

### 1. Clone Repository

```bash
git clone https://github.com/USERNAME/DriveRentNew.git
```

Masuk ke folder project:

```bash
cd DriveRentNew
```

### 2. Setup Backend

Masuk ke folder backend:

```bash
cd backend
```

Install dependencies:

```bash
go mod download
```

Buat file `.env` berdasarkan konfigurasi environment yang digunakan backend.

Contoh konfigurasi:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=driverent
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_secret_key
```

Sesuaikan nama variabel dengan konfigurasi yang digunakan di source code backend.

Buat database PostgreSQL dengan nama yang sesuai dengan konfigurasi.

Jalankan migrasi database menggunakan konfigurasi golang-migrate yang tersedia di project.

Jalankan backend dari folder `backend`:

```bash
go run ./cmd/api
```

Backend berjalan pada:

```text
http://localhost:8080
```

### 3. Setup Frontend

Buka terminal baru dan masuk ke folder frontend:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Jalankan development server:

```bash
npm run dev
```

Frontend biasanya berjalan pada:

```text
http://localhost:5173
```

## 🔄 Business Workflow

Alur utama penyewaan kendaraan pada DriveRent:

1. **Register & Login**  
   Customer membuat akun dan masuk ke aplikasi.

2. **Vehicle Selection**  
   Customer melihat daftar kendaraan dan memilih kendaraan yang tersedia.

3. **Reservation**  
   Customer membuat reservasi dengan menentukan tanggal penyewaan.

4. **Payment**  
   Customer mengirim pembayaran dan bukti pembayaran.

5. **Payment Verification**  
   Admin memeriksa dan memverifikasi pembayaran.

6. **Rental Process**  
   Admin mencatat penyerahan kendaraan dan informasi kondisi awal kendaraan.

7. **Vehicle Return**  
   Admin mencatat pengembalian kendaraan, odometer akhir, bahan bakar, dan biaya tambahan jika ada.

8. **Transaction Completed**  
   Rental diselesaikan dan status kendaraan diperbarui menjadi tersedia atau maintenance.

## 🔐 User Roles

| Role | Permissions |
|---|---|
| Admin | Mengelola pengguna, kendaraan, reservasi, pembayaran, rental, dan pengembalian |
| Customer | Melihat kendaraan, membuat reservasi, melakukan pembayaran, dan melihat riwayat transaksi sendiri |

## 🎯 Project Goals

Project DriveRent dikembangkan untuk menerapkan konsep Fullstack Web Development, meliputi:

- RESTful API Development.
- Frontend and Backend Integration.
- Authentication and Authorization.
- Database Design and Management.
- CRUD Operations.
- File Upload Management.
- Transaction Management.
- Role-Based Access Control.

## 👨‍💻 Developer

**DriveRent – Vehicle Rental Management System**

Developed as a Fullstack Web Development Portfolio Project.

**Technologies:** React.js, Golang, PostgreSQL, and Tailwind CSS.
