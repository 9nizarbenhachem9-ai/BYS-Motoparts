# Web Manajemen Sparepart

Aplikasi web modern untuk manajemen inventaris sparepart bengkel dan gudang, dibangun dengan Go (Gin, GORM, PostgreSQL) dan frontend modern HTML/CSS/JavaScript murni.

---

## 🌟 Fitur Utama

1. **Autentikasi & Keamanan:**
   - Registrasi dan Login akun Administrator.
   - Enkripsi password menggunakan **bcrypt**.
   - Keamanan API berbasis **JSON Web Token (JWT)**.
   - Mencegah celah XSS dengan sanitasi rendering DOM.
2. **Manajemen Sparepart (CRUD):**
   - Tambah, lihat detail, perbarui data, dan hapus sparepart.
   - Menggunakan **Soft Delete** agar riwayat transaksi tetap terjaga.
   - Kode sparepart unik dan tervalidasi.
3. **Pencatatan Stok Masuk & Keluar:**
   - Perubahan stok hanya bisa terjadi melalui transaksi barang masuk atau keluar.
   - Transaksi database bersifat **atomic** (ACID).
   - Validasi ketat: transaksi barang keluar ditolak jika melebihi stok yang tersedia (stok tidak pernah negatif).
4. **Pencarian, Filter & Paginasi:**
   - Pencarian real-time (debounce) berdasarkan nama, kode, atau merek sparepart.
   - Filter berdasarkan kategori barang.
   - Paginasi data yang rapi dan responsif.
5. **Antarmuka Modern & Ramah Pengguna:**
   - **Penanda Stok Menipis (Warning Badge):** Otomatis menyala jika stok saat ini $\le$ stok minimum.
   - **Menu Klik Kanan (Context Menu):** Klik kanan pada baris sparepart untuk membuka menu cepat (Barang Masuk/Keluar, Riwayat, Ubah, Hapus).
   - **Tombol Alternatif:** Tersedia tombol aksi langsung pada setiap baris untuk perangkat layar sentuh / mobile.
   - Notifikasi toast melayang yang halus dan informatif.

---

## 🚀 Cara Menjalankan Aplikasi

### 1. Prasyarat
- **Go** (versi 1.20+)
- **Node.js** 18+ dan **npm** (untuk tooling frontend/Vite)
- **PostgreSQL** aktif

### 2. Konfigurasi Variabel Lingkungan
Salin `.env.example` menjadi `.env` di folder root proyek:
```bash
cp .env.example .env
```
Sesuaikan isi file `.env`:
```dotenv
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/sparepart_db?sslmode=disable
JWT_SECRET=kunci-rahasia-jwt-acak-yang-panjang
JWT_EXPIRES_HOURS=24
ALLOWED_ORIGIN=http://localhost:5500
```

### 3. Menjalankan Backend & Frontend

Frontend ditulis dengan **HTML/CSS/JavaScript murni** lalu dikelola oleh **npm + Vite** (dev server, proxy API, dan bundling untuk produksi).

**Mode development — buka dua terminal:**

Terminal 1 (backend Go):
```powershell
go run .
```

Terminal 2 (frontend Vite):
```powershell
cd frontend
npm install   # cukup sekali di awal
npm run dev
```

- Buka **http://localhost:5173** — semua panggilan `/api` otomatis di-proxy ke backend Go (port 8080), jadi tanpa masalah CORS.
- Alternatif tanpa Vite: buka [http://localhost:8080](http://localhost:8080) — backend menyajikan sumber frontend mentah (tanpa bundling/minify).

**Mode produksi — build frontend lalu jalankan backend:**
```powershell
cd frontend
npm run build    # menghasilkan frontend/dist/
cd ..
go run .         # backend menyajikan isi frontend/dist
```
Buka [http://localhost:8080](http://localhost:8080).

> Deploy ke Hostless tidak perlu langkah tambahan: Dockerfile otomatis menjalankan `npm ci` + `npm run build` pada tahap `frontend-builder`, lalu menyalin hasilnya ke `/app/frontend/` yang disajikan biner Go.

---

## 🧪 Menjalankan Pengujian Otomatis

Untuk menjalankan seluruh unit test backend:
```powershell
cd backend
go test -v ./...
```

