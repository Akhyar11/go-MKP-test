# E-Ticketing Transportasi Publik API

API untuk sistem E-Ticketing Transportasi Publik menggunakan Go, GORM, dan PostgreSQL.

## Persyaratan

- Go 1.19 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Postman atau tools API testing lainnya

## Setup Project

1. Clone repository

```bash
git clone <repository_url>
cd tester
```

2. Install dependencies

```bash
go mod tidy
```

3. Setup database

- Buat database PostgreSQL baru

```sql
CREATE DATABASE eticketing;
```

- Import struktur database

```bash
psql -U postgres -d eticketing -f import.sql
```

4. Konfigurasi environment

- Copy `.env.example` ke `.env`
- Sesuaikan konfigurasi database:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=eticketing
DB_SSLMODE=disable
JWT_SECRET=your_jwt_secret_key
```

5. Setup user admin

```bash
go run scripts/generate_hash.go
```

6. Jalankan server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Endpoints API

### 1. Login

Endpoint untuk mendapatkan token JWT.

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

Endpoint untuk membuat terminal baru. Membutuhkan autentikasi JWT.

- **URL**: `/api/terminals`
- **Method**: `POST`
- **Headers**:

```
Authorization: Bearer <token_jwt>
Content-Type: application/json
```

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

## Struktur Database

Database terdiri dari beberapa tabel utama:

- `users`: Menyimpan data pengguna
- `terminals`: Data terminal transportasi
- `gates`: Gate/pintu di setiap terminal
- `cards`: Kartu e-ticket
- `trips`: Riwayat perjalanan
- `balance_transactions`: Riwayat transaksi saldo

Detail lengkap struktur dapat dilihat di file `import.sql`.
