# E-Ticketing Transportasi Publik API

API untuk sistem E-Ticketing Transportasi Publik menggunakan Go, GORM, dan PostgreSQL.

---

## Persyaratan

- Go 1.19 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Postman/cURL/HTTP client lain

---

## Setup Project

1. **Clone repository**
   ```bash
   git clone <repository_url>
   cd tester
   ```
2. **Install dependencies**
   ```bash
   go mod tidy
   ```
3. **Setup database**
   - Buat database PostgreSQL baru:
     ```sql
     CREATE DATABASE eticketing;
     ```
   - Import struktur database:
     ```bash
     psql -U postgres -d eticketing -f import.sql
     ```
4. **Konfigurasi environment**
   - Buat file `.env` (atau copy dari `.env.example`)
   - Isi dengan:
     ```env
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=postgres
     DB_PASSWORD=yourpassword
     DB_NAME=eticketing
     DB_SSLMODE=disable
     JWT_SECRET=your_jwt_secret_key
     ```
5. **Setup user admin**
   ```bash
   go run scripts/generate_hash.go
   ```
6. **Jalankan server**
   ```bash
   go run main.go
   ```
   Server berjalan di `http://localhost:8080`

---

## Rancangan Alur Kerja Sistem

Dokumen ini merincikan alur kerja operasional sistem E-Ticketing dalam dua kondisi krusial: saat terhubung ke internet (Online) dan saat koneksi terputus (Offline).

### 1. Rancangan Saat Ada Jaringan Internet (Mode Online)

Ini adalah kondisi operasional standar di mana setiap gate terhubung secara real-time ke Server Pusat. Keuntungannya adalah data selalu akurat dan terpusat.

#### Alur Kerja Check-in (Masuk)

1. **Tap Kartu**: Penumpang menempelkan Kartu Prepaid pada reader di gate masuk.
2. **Kirim Data ke Server**: Gate membaca ID_Kartu yang unik, lalu secara instan mengirimkan permintaan check-in ke Server Pusat. Permintaan ini berisi: ID_Kartu, ID_Gate, dan Waktu_Checkin.
3. **Validasi oleh Server Pusat**:
   - Memverifikasi apakah ID_Kartu valid dan berstatus 'aktif'.
   - Memastikan saldo di database mencukupi syarat minimum untuk memulai perjalanan.
   - Memastikan kartu tidak sedang dalam status 'check-in' di tempat lain (mencegah double entry).
4. **Pencatatan & Respon**: Jika valid, Server Pusat membuat catatan perjalanan baru di database dengan status "berjalan". Kemudian, server mengirimkan respon "OK" kembali ke gate. **Gate juga menulis data perjalanan (check-in) ke dalam kartu sebagai backup jika sewaktu-waktu terjadi offline di gate tujuan.**
5. **Aksi Gate**: Gate menerima respon "OK", palang gerbang terbuka, dan layar menampilkan pesan sambutan beserta sisa saldo resmi yang diterima dari server. Seluruh proses ini terjadi dalam sepersekian detik.

#### Alur Kerja Check-out (Keluar)

1. **Tap Kartu**: Penumpang menempelkan kartu yang sama pada reader di gate keluar.
2. **Kirim Data ke Server**: Gate mengirimkan permintaan check-out ke Server Pusat berisi: ID_Kartu, ID_Gate, dan Waktu_Checkout.
3. **Kalkulasi & Transaksi oleh Server**:
   - Mengambil data check-in yang tersimpan untuk kartu tersebut.
   - Menghitung tarif perjalanan berdasarkan terminal masuk dan terminal keluar.
   - Memotong saldo pengguna di database pusat sesuai dengan tarif.
   - Memperbarui catatan perjalanan di database menjadi "selesai", lengkap dengan rincian tarif dan waktu.
4. **Respon & Sinkronisasi**: Server mengirimkan respon "OK" kembali ke gate, beserta rincian transaksi seperti jumlah tarif yang dipotong dan sisa saldo akhir.
5. **Aksi Gate**: Gate menerima respon, palang gerbang terbuka, dan layar menampilkan rincian perjalanan. Gate juga bisa memperbarui data saldo di kartu agar tetap sinkron dengan server.

---

### 2. Solusi Saat Tidak Ada Jaringan Internet (Mode Offline)

Ini adalah solusi krusial untuk memastikan layanan tetap berjalan 24 jam tanpa henti. Saat koneksi putus, tanggung jawab validasi dan pencatatan sementara dialihkan dari Server Pusat ke Gate dan Kartu Prepaid itu sendiri.

#### Alur Kerja Check-in (Offline)

1. **Deteksi Offline**: Penumpang menempelkan kartu. Gate mencoba menghubungi server, namun gagal. Gate secara otomatis beralih ke Mode Offline.
2. **Validasi Lokal pada Kartu**: Gate kini mengandalkan data yang tersimpan di dalam kartu:
   - Membaca Saldo_Lokal dan memastikan nilainya di atas saldo minimum.
   - Membaca Status_Perjalanan dan memastikan statusnya "checked-out".
3. **TULIS Data ke Kartu**: Jika validasi lokal berhasil, Gate menulis informasi check-in ke dalam memori kartu. Data yang ditulis meliputi ID_Terminal_Masuk dan Waktu_Checkin. Status perjalanan di kartu juga diubah menjadi "checked-in".
4. **Catat Log Lokal**: Gate menyimpan catatan transaksi check-in ini di dalam penyimpanan lokalnya (seperti SD card atau flash memory). Log ini akan menjadi antrean untuk sinkronisasi nanti.
5. **Aksi Gate**: Palang gerbang terbuka. Pengguna bisa masuk tanpa menyadari adanya gangguan jaringan.

#### Alur Kerja Check-out (Offline)

1. **Deteksi Offline**: Penumpang menempelkan kartu di gate tujuan. Gate kembali beroperasi dalam mode offline.
2. **BACA Data Perjalanan dari Kartu**: Gate membaca informasi check-in yang sebelumnya ditulis pada langkah check-in offline langsung dari kartu.
3. **Kalkulasi Tarif Lokal**: Gate menggunakan cache matriks tarif yang tersimpan di memorinya untuk menghitung tarif perjalanan berdasarkan terminal masuk (dari kartu) dan terminal keluar (lokasi gate saat ini).
4. **TULIS Ulang ke Kartu**: Gate memperbarui data di kartu:
   - Mengurangi Saldo_Lokal sesuai tarif yang dihitung.
   - Mengubah Status_Perjalanan kembali menjadi "checked-out".
   - Menghapus data check-in yang sudah tidak relevan.
5. **Catat Log Lokal**: Gate menyimpan log transaksi lengkap (check-in dan check-out) di penyimpanan lokalnya.
6. **Aksi Gate**: Palang gerbang terbuka. Informasi tarif dan sisa saldo ditampilkan berdasarkan perhitungan lokal.

#### Proses Rekonsiliasi (Saat Kembali Online)

1. **Koneksi Pulih**: Ketika koneksi internet di sebuah gate kembali pulih, gate akan secara otomatis memulai proses sinkronisasi.
2. **Kirim Antrean Log**: Gate mengirimkan semua log transaksi yang tersimpan di memori lokalnya ke Server Pusat.
3. **Pembaruan Database Pusat**: Server Pusat menerima log ini dan memperbarui database utamanya satu per satu, seolah-olah transaksi tersebut terjadi secara real-time. Proses ini memastikan bahwa meskipun transaksi terjadi secara offline, pada akhirnya semua data akan tercatat secara akurat di pusat data untuk keperluan audit dan pelaporan.

---

## Endpoints API

### 1. Login

- **URL**: `/api/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "admin@example.com",
    "password": "password123"
  }
  ```
- **Response Success** (200):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```
- **Response Error** (401):
  ```json
  {
    "message": "Email atau password salah"
  }
  ```

### 2. Create Terminal

- **URL**: `/api/terminals`
- **Method**: `POST`
- **Headers**:
  - `Authorization: Bearer <token_jwt>`
  - `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "nama_terminal": "Terminal Pusat",
    "lokasi": "Jl. Raya Utama No. 1"
  }
  ```
- **Response Success** (201):
  ```json
  {
    "terminal_id": 1,
    "nama_terminal": "Terminal Pusat",
    "lokasi": "Jl. Raya Utama No. 1",
    "created_at": "2024-03-23T10:00:00Z"
  }
  ```
- **Response Error** (401):
  ```json
  {
    "message": "Unauthorized"
  }
  ```

---

## Testing API

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
```

### Create Terminal

```bash
curl -X POST http://localhost:8080/api/terminals \
  -H "Authorization: Bearer <token_dari_login>" \
  -H "Content-Type: application/json" \
  -d '{"nama_terminal":"Terminal Pusat","lokasi":"Jl. Raya Utama No. 1"}'
```

---

## Struktur Database

Database terdiri dari beberapa tabel utama:

- `users`: Menyimpan data pengguna
- `terminals`: Data terminal transportasi
- `gates`: Gate/pintu di setiap terminal
- `cards`: Kartu e-ticket
- `trips`: Riwayat perjalanan
- `balance_transactions`: Riwayat transaksi saldo

Detail lengkap struktur dapat dilihat di file `import.sql`.
