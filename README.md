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
Buka terminal di direktori `backend`:
```powershell
cd backend
go run .
```

Aplikasi backend akan langsung berjalan dan menyajikan frontend secara otomatis:
- **Buka Browser:** Kunjungi [http://localhost:8080](http://localhost:8080)
- Masuk ke menu **Daftar di sini** untuk membuat akun admin pertama Anda, lalu login untuk masuk ke dashboard utama.

---

## 🧪 Menjalankan Pengujian Otomatis

Untuk menjalankan seluruh unit test backend:
```powershell
cd backend
go test -v ./...
```
