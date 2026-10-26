# Project Specification — Web Manajemen Sparepart

> **Tugas 1:** Spesifikasi ini harus selesai dan disetujui sebelum vibe coding dimulai.
> Item bertanda **[KEPUTUSAN]** adalah usulan default. Ubah jika tidak sesuai, lalu centang di bagian 10.

---

## 1. Informasi Aplikasi

**Nama aplikasi:** Web Manajemen Sparepart

**Deskripsi:** Aplikasi web untuk mencatat, menampilkan, memperbarui, dan menghapus data sparepart, serta mencatat barang masuk dan barang keluar beserta riwayat dan stok terkininya.

**Target pengguna:** Administrator bengkel / petugas manajemen gudang.

**Masalah yang diselesaikan:**
1. Pencarian sparepart manual (buku/Excel) lambat.
2. Stok sering tidak sinkron karena barang masuk-keluar tidak tercatat rapi.
3. Admin tidak tahu kapan stok menipis.

**Kriteria sukses (aplikasi dianggap selesai jika):**
- Admin dapat login, menambah sparepart, dan menemukannya lewat pencarian dalam kurang dari 10 detik.
- Setiap perubahan stok tercatat sebagai transaksi masuk/keluar dan stok tidak pernah negatif.
- Tidak ada kredensial di dalam kode maupun Git.

---

## 2. Kebutuhan Fungsional

### 2.1 Fitur utama

| # | Fitur | Prioritas |
|---|-------|-----------|
| 1 | CRUD sparepart (tambah, lihat, ubah, hapus) | Wajib |
| 2 | Pencatatan barang masuk dan barang keluar | Wajib |
| 3 | Klik kanan pada baris sparepart untuk memunculkan menu Ubah / Hapus | Wajib |
| 4 | Pencarian dan filter (nama, kode, kategori) | Wajib |
| 5 | Penanda stok menipis (stok <= stok minimum) | Tambahan |
| 6 | Riwayat transaksi per sparepart | Tambahan |
| 7 | Paginasi daftar sparepart | Tambahan |

**Catatan fitur klik kanan:** menu ini tidak bisa dipakai di layar sentuh. **[KEPUTUSAN]** Sediakan juga tombol "Ubah" dan "Hapus" pada setiap baris sebagai alternatif, dan tampilkan konfirmasi sebelum menghapus.

### 2.2 Halaman

- **Login:** form username dan password.
- **Register:** form pembuatan akun admin (lihat aturan di bagian 4).
- **Main Page:** daftar sparepart (tabel), pencarian, tombol tambah, form barang masuk/keluar, riwayat transaksi, tombol logout.

### 2.3 Data sparepart (usulan)

| Field | Tipe | Aturan |
|-------|------|--------|
| id | integer | primary key |
| kode | string | unik, wajib |
| nama | string | wajib |
| kategori | string | opsional |
| merek | string | opsional |
| harga | integer (rupiah) | >= 0 |
| stok | integer | >= 0, hanya berubah lewat transaksi |
| stok_minimum | integer | >= 0, default 0 |
| lokasi_rak | string | opsional |
| created_at, updated_at | timestamp | otomatis |
| deleted_at | timestamp | soft delete (bawaan GORM) |

### 2.4 Data transaksi stok (usulan)

| Field | Tipe | Aturan |
|-------|------|--------|
| id | integer | primary key |
| sparepart_id | integer | relasi ke sparepart |
| tipe | string | `masuk` atau `keluar` |
| jumlah | integer | > 0 |
| catatan | string | opsional (mis. no. nota/pemasok) |
| user_id | integer | admin yang mencatat |
| created_at | timestamp | otomatis |

### 2.5 Aturan bisnis

1. Barang keluar ditolak jika jumlah melebihi stok saat ini (pesan error jelas).
2. Perubahan stok dan pencatatan transaksi dilakukan dalam **satu database transaction** agar konsisten.
3. Stok tidak diedit langsung lewat form Ubah; stok berubah hanya lewat barang masuk/keluar. **[KEPUTUSAN]**
4. Hapus sparepart memakai soft delete agar riwayat transaksi tetap utuh.
5. Kode sparepart tidak boleh duplikat.

---

## 3. REST API (usulan)

Semua endpoint di bawah `/api`. Selain register dan login, semua endpoint wajib header `Authorization: Bearer <token>`.

| Method | Endpoint | Fungsi |
|--------|----------|--------|
| POST | `/api/auth/register` | Daftar akun admin |
| POST | `/api/auth/login` | Login, mengembalikan JWT |
| GET | `/api/auth/me` | Data admin yang sedang login |
| GET | `/api/spareparts?search=&kategori=&page=&limit=` | Daftar sparepart |
| POST | `/api/spareparts` | Tambah sparepart |
| GET | `/api/spareparts/:id` | Detail sparepart |
| PUT | `/api/spareparts/:id` | Ubah sparepart |
| DELETE | `/api/spareparts/:id` | Hapus sparepart |
| POST | `/api/spareparts/:id/transactions` | Catat barang masuk/keluar |
| GET | `/api/spareparts/:id/transactions` | Riwayat transaksi sparepart |

**Format respons error (seragam):**
```json
{ "error": "pesan yang jelas untuk pengguna" }
```
Kode status: 200/201 sukses, 400 input tidak valid, 401 belum login/token tidak valid, 404 tidak ditemukan, 409 konflik (mis. kode duplikat), 500 kesalahan server (tanpa detail internal).

---

## 4. Autentikasi dan Keamanan

**Registrasi:** Ya | **Login:** Ya | **Logout:** Ya | **Role:** Admin saja.

**Aturan akses:**
1. Pengguna harus login untuk mengakses Main Page dan semua endpoint data.
2. Pengguna yang belum login hanya dapat membuka Login dan Register.
3. Password disimpan sebagai hash **bcrypt**, tidak pernah sebagai teks biasa.
4. Password minimal 8 karakter; username unik.
5. Pesan error login tidak membedakan "username salah" dan "password salah".
6. JWT memiliki masa berlaku (**[KEPUTUSAN]** 24 jam). Token kedaluwarsa atau tidak valid menghasilkan 401 dan frontend mengarahkan ke Login.
7. **Logout:** frontend menghapus token dari localStorage dan mengarahkan ke Login.

**[KEPUTUSAN] Pembatasan Register.** Karena semua akun otomatis berstatus admin, Register yang terbuka untuk umum berarti siapa pun bisa menjadi admin. Pilih salah satu:
- [ ] A. Register dibuka hanya untuk akun pertama; setelah itu ditutup.
- [ ] B. Register membutuhkan kode undangan dari environment variable `REGISTER_INVITE_CODE`.
- [ ] C. Dibiarkan terbuka (hanya untuk keperluan belajar/lokal, tidak untuk produksi).

**Catatan keamanan JWT di localStorage:** localStorage dapat dibaca oleh JavaScript, sehingga rentan jika ada celah XSS. Mitigasinya: jangan pernah memasukkan data pengguna ke halaman lewat `innerHTML` (gunakan `textContent`), dan jangan memuat skrip pihak ketiga. Pilihan ini diterima untuk proyek belajar.

**Perlindungan tambahan:**
- Validasi semua input di backend (frontend hanya untuk kenyamanan).
- GORM dengan parameter query (tidak merangkai SQL dari input).
- CORS hanya mengizinkan origin frontend yang dikonfigurasi.
- **[KEPUTUSAN]** Pembatasan percobaan login (rate limit sederhana) boleh ditunda ke tahap akhir.

---

## 5. Pilihan Teknologi (final)

| Komponen | Pilihan | Ketentuan |
|----------|---------|-----------|
| Backend | **Gin** | Handler tipis, logika di service |
| ORM | **GORM** | Gunakan AutoMigrate untuk tahap belajar |
| Database | **PostgreSQL** | Koneksi dari `DATABASE_URL` |
| Frontend | **HTML, CSS, JavaScript murni** | Tanpa framework; akses API dengan `fetch()` |
| Autentikasi | **JWT** | Ditandatangani dengan `JWT_SECRET`, algoritma HS256 |
| Penyimpanan token | **localStorage** | Menyimpan **JWT** (bukan session ID) |
| API | **REST API** | Format JSON |
| Version control | **Git** | Commit kecil per tahap |

**Dependensi yang diperkirakan** (tetap perlu persetujuan saat dipasang): `gin-gonic/gin`, `gorm.io/gorm`, `gorm.io/driver/postgres`, `golang-jwt/jwt`, `golang.org/x/crypto` (bcrypt), `joho/godotenv`, `gin-contrib/cors`.

---

## 6. Konfigurasi Environment dan Database

**Ketentuan wajib:**
1. Semua konfigurasi sensitif memakai environment variable.
2. File `.env` untuk pengembangan lokal; `.env` masuk `.gitignore`.
3. Sediakan `.env.example` tanpa kredensial asli (boleh masuk Git).
4. `DATABASE_URL` dibaca dari environment variable, tidak pernah di-hardcode.
5. Jangan menulis URL database, username, password, secret, atau kredensial lain di kode.
6. Jika konfigurasi wajib hilang, aplikasi berhenti saat start dengan pesan jelas (menyebut nama variabel yang kurang) tanpa menampilkan nilai rahasia.
7. Nilai rahasia tidak boleh muncul di log maupun respons API.

**`.env.example`:**
```dotenv
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/sparepart_db?sslmode=disable
JWT_SECRET=replace-with-a-long-random-secret
JWT_EXPIRES_HOURS=24
ALLOWED_ORIGIN=http://localhost:5500
# REGISTER_INVITE_CODE=replace-with-invite-code
```

> `SESSION_SECRET` pada versi sebelumnya diganti `JWT_SECRET` karena sistem memakai JWT, bukan session. Untuk produksi, gunakan konfigurasi SSL sesuai penyedia database.

---

## 7. Struktur Folder (usulan awal, belum dibuat)

```
sparepart-app/
├── backend/
│   ├── main.go            # titik masuk, setup router
│   ├── config/            # baca dan validasi environment variable
│   ├── database/          # koneksi GORM dan migrasi
│   ├── models/            # struct User, Sparepart, StockTransaction
│   ├── handlers/          # handler HTTP (auth, sparepart, transaksi)
│   ├── middleware/        # verifikasi JWT, CORS
│   ├── services/          # logika bisnis (mis. update stok dalam transaksi)
│   └── *_test.go          # pengujian dasar
├── frontend/
│   ├── login.html
│   ├── register.html
│   ├── index.html         # Main Page
│   ├── css/style.css
│   └── js/                # api.js, auth.js, main.js
├── docs/prompts/00-project-spec.md
├── .env.example
├── .gitignore
└── README.md
```

---

## 8. Batasan Pengembangan

- [x] Struktur folder sederhana.
- [x] Kode mudah dipahami pemula (komentar singkat pada bagian penting).
- [x] Frontend dan backend dipisah.
- [x] Environment variable untuk konfigurasi.
- [x] Pengujian dasar (lihat bagian 9).
- [x] AI menjelaskan kode dan perubahan yang dibuat.
- [x] Tidak menambah teknologi/dependensi di luar pilihan tanpa persetujuan.

**Batasan tambahan:**
- Satu tahap = satu fitur kecil yang bisa diuji sebelum lanjut.
- Tidak ada refactor besar di luar tahap yang sedang dikerjakan.

**Fitur yang tidak dibutuhkan (di luar ruang lingkup):**
- Role selain Admin, manajemen pengguna.
- Reset password lewat email, login sosial.
- Laporan PDF/Excel, ekspor data, grafik dashboard.
- Upload foto sparepart.
- Multi-gudang, manajemen pemasok, pembelian/penjualan, aplikasi mobile.

---

## 9. Pengujian Dasar

- **Backend:** tes untuk hash/verifikasi password, pembuatan dan validasi JWT, aturan stok keluar (tidak boleh melebihi stok), dan validasi input sparepart.
- **API:** daftar langkah uji manual (curl atau REST client) per endpoint: kasus berhasil dan gagal.
- **Frontend:** daftar periksa manual: akses halaman tanpa login, login berhasil/gagal, CRUD, klik kanan, logout.

---

## 10. Pemeriksaan Sebelum Vibe Coding

- [ ] Nama dan tujuan aplikasi sudah jelas.
- [ ] Fitur utama dan aturan bisnis sudah ditentukan.
- [ ] Pilihan teknologi sudah final dan tidak saling bertentangan.
- [ ] Keputusan **[KEPUTUSAN]** sudah ditinjau, termasuk pembatasan Register.
- [ ] Kebutuhan autentikasi sudah ditentukan.
- [ ] Penggunaan `.env` dan `DATABASE_URL` sudah ditentukan.
- [ ] Batasan dan fitur di luar ruang lingkup sudah ditentukan.

**Status spesifikasi:** Siap diperiksa dan disetujui sebelum implementasi.

---

## 11. Instruksi untuk Agy

**Peran:** Kamu asisten pemrograman yang membimbing pemula membangun aplikasi Go (Gin + GORM + PostgreSQL) dengan frontend HTML/CSS/JavaScript murni. Jelaskan dengan bahasa Indonesia yang sederhana.

Baca file `docs/prompts/00-project-spec.md` secara lengkap dan gunakan sebagai acuan utama. Jika ada konflik antara permintaan lisan dan dokumen ini, tanyakan dulu.

### Fase 0 — Tinjauan (lakukan sekarang, belum menulis kode)
1. Periksa kelengkapan dan konsistensi spesifikasi; sebutkan jika ada yang bertentangan.
2. Ajukan pertanyaan jika ada informasi penting yang belum jelas (maksimal satu daftar, jangan bertanya berulang).
3. Usulkan struktur folder sederhana sesuai bagian 7 dan jelaskan fungsi tiap folder dan file utama.
4. **Berhenti dan tunggu persetujuan saya.** Jangan membuat, mengubah, atau menghapus file apa pun sebelum itu.

### Fase implementasi (setelah disetujui, satu tahap per giliran)
| Tahap | Isi | Cara uji |
|-------|-----|----------|
| 1 | Struktur proyek, `.gitignore`, `.env.example`, loader konfigurasi, koneksi database | Jalankan; coba tanpa `DATABASE_URL` dan lihat pesan error |
| 2 | Model dan migrasi | Cek tabel di PostgreSQL |
| 3 | Register, Login, middleware JWT | Uji endpoint dengan curl |
| 4 | CRUD sparepart + pencarian + paginasi | Uji endpoint |
| 5 | Barang masuk/keluar + riwayat | Uji stok tidak boleh negatif |
| 6 | Frontend Login/Register + penyimpanan token + logout | Uji manual di browser |
| 7 | Frontend Main Page, klik kanan, tombol alternatif | Daftar periksa manual |
| 8 | Pengujian dasar dan README | Jalankan `go test ./...` |

### Aturan kerja
1. Jangan membuat seluruh aplikasi sekaligus.
2. Di setiap tahap: jelaskan apa yang dibuat dan alasannya, tampilkan file yang berubah, berikan langkah menguji, lalu tunggu konfirmasi saya sebelum lanjut.
3. Jangan mengganti teknologi atau menambah dependensi tanpa menjelaskan alasannya dan meminta persetujuan.
4. `DATABASE_URL` dan `JWT_SECRET` hanya dibaca dari environment variable; jangan di-hardcode.
5. Pastikan `.env` ada di `.gitignore`; `.env.example` boleh masuk Git.
6. Jangan menampilkan nilai rahasia dari `.env` di log atau respons API.
7. Pesan error ke pengguna harus jelas tetapi tidak membocorkan detail internal.
8. Jika kamu tidak yakin, tanyakan; jangan menebak.
9. Di akhir tiap tahap, sarankan pesan commit Git yang singkat.
