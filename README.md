# Manufacture

Manufacture adalah platform manajemen manufaktur full-stack dengan frontend Vue 3 dan backend Go.
Dokumentasi ini membantu kamu menjalankan aplikasi, memahami fitur, serta menggunakan API yang tersedia.

## 👤 Tentang Saya
| Nama | Peran | Kontak |
| --- | --- | --- |
| Brian Restu Pratiwa | Hamba Allah | [Email](mailto:brestu352@gmail.com) • [BuySomeCoffee](https://teer.id/brianrstp) |

---

## 🚀 Fitur Utama

- Admin dashboard interaktif dengan ringkasan pesanan, produksi, pendapatan, dan pelanggan
- CRUD lengkap untuk:
  - Produk
  - Kategori
  - Unit
  - Gudang
  - Pelanggan
  - Pesanan
  - Produksi
  - Inventaris
  - Mutasi stok
- Filter, pencarian, dan pagination sederhana pada data admin
- Modal form dengan validasi frontend dan feedback sukses/gagal
- Token-based admin authentication via JWT
- Backend menggunakan database MySQL dengan soft delete

---

## 📌 Detail Fitur

### Fitur Admin
- Dashboard admin interaktif dengan ringkasan pesanan, produksi, pendapatan bulanan, dan total pelanggan.
- CRUD lengkap untuk produk, kategori, unit, gudang, pelanggan, pesanan, produksi, inventaris, dan mutasi stok.
- Manajemen pesanan dengan referensi pelanggan, status pesanan, tanggal, dan jumlah.
- Produksi dapat dicatat, diedit, dan dihapus dengan status serta jumlah barang yang diproduksi.
- Inventaris dan mutasi stok memudahkan admin memonitor pergerakan barang antar gudang.
- Filter dan pencarian di banyak halaman untuk menemukan data berdasarkan nama produk, status pesanan, tipe mutasi stok, dan gudang.
- Modal form dengan validasi frontend dan pesan sukses/gagal setelah proses CRUD.
- Admin authentication dengan JWT untuk melindungi semua endpoint admin.

### Fitur Customer
- Portal customer dengan halaman login khusus untuk pelanggan yang terdaftar.
- Customer dapat melihat riwayat pesanan mereka di halaman `Pesanan Saya`.
- Customer dapat membuka halaman profil untuk melihat informasi akun seperti nama dan email.
- Halaman `Bantuan` menyediakan navigasi portal customer serta tombol logout.
- Customer session dikelola di frontend menggunakan `localStorage` untuk token dan informasi nama/email.
- Permintaan customer memakai header `Authorization: Bearer <token>` untuk endpoint customer yang dilindungi.

---

## 🧩 Teknologi

- Frontend
  - Vue 3
  - Vue Router
  - Vite
  - Tailwind CSS
  - Chart.js
  - FontAwesome
- Backend
  - Go
  - MySQL
  - JWT authentication
  - `net/http`

---

## 📁 Struktur Proyek

- `frontend/`: aplikasi Vue 3
  - `src/pages/admin/`: semua halaman dashboard admin
  - `src/router/`: routing aplikasi
- `backend/`: API Go
  - `handler/`: HTTP route handlers
  - `database/`: model dan query database
  - `config/`: environment configuration

---

## ⚙️ Menjalankan Proyek

### Frontend

Masuk ke folder frontend:

```bash
cd frontend
npm install
npm run dev
```

Buka `http://localhost:5173`.

### Backend

Masuk ke folder backend:

```bash
cd backend
go run main.go
```

Default server berjalan di `http://localhost:8080`.

---

## 🛠️ Konfigurasi Environment Backend

Tambahkan variabel environment jika diperlukan:

```bash
SERVER_PORT=8080
MYSQL_USER=root
MYSQL_PASSWORD=yourpassword
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_DATABASE=manufacture
ADMIN_JWT_SECRET=supersecretkey
ADMIN_TOKEN_EXPIRY=1h
CUSTOMER_JWT_SECRET=customersecretkey
CUSTOMER_TOKEN_EXPIRY=24h
```

---

## Setup Database, Schema & Migrasi

1. Pastikan MySQL sudah berjalan dan database `manufacture` sudah dibuat:

```sql
CREATE DATABASE manufacture CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. Konfigurasi environment untuk koneksi MySQL di file `.env` atau environment variables.
3. Jalankan backend dari `backend/`:

```bash
cd backend
go run main.go
```

4. Backend akan otomatis membuat schema dan tabel yang dibutuhkan saat startup via `database.EnsureManufacturingSchema(db)` di `backend/main.go`.

> Tabel utama dibuat secara otomatis: `categories`, `units`, `products`, `warehouses`, `inventory`, `customers`, `orders`, `production_jobs`, `stock_movements`, `bills_of_materials`, dan tabel pendukung lain.

### Skema database penting

- `products`: menyimpan katalog produk dengan SKU, kategori, unit, harga, dan status lifecycle.
- `customers`: menyimpan data pelanggan, email, password hash, dan status.
- `orders`: menyimpan pesanan pelanggan dengan referensi `customer_id`.
- `production_jobs`: menyimpan pekerjaan produksi dengan kode, produk, tanggal mulai, durasi, dan status.
- `inventory`: menyimpan status stok per produk per gudang.
- `stock_movements`: merekam mutasi stok masuk/keluar di gudang.
- `warehouses`: menyimpan informasi gudang dan statusnya.
- `bills_of_materials`: menyimpan struktur BOM untuk produk manufaktur.

> Semua tabel utama mendukung `soft delete` dengan kolom `deleted_at`, sehingga data dapat dihapus secara logis tanpa hilang sepenuhnya.

---

## 📂 Struktur Folder Customer & Database

### Frontend customer
- `frontend/src/pages/customer/CustomerLoginPage.vue`: halaman login customer.
- `frontend/src/pages/customer/CustomerOrdersPage.vue`: halaman riwayat pesanan customer.
- `frontend/src/pages/customer/CustomerProfilePage.vue`: halaman profil customer.
- `frontend/src/pages/customer/CustomerHelpPage.vue`: halaman bantuan dan logout.

Route customer didefinisikan di `frontend/src/router/index.ts` dengan proteksi `requiresCustomerAuth` untuk route yang membutuhkan token.

### Backend database
- `backend/database/products.go`: schema dan query produk, unit, kategori, gudang, BOM, inventory, dan mutasi stok.
- `backend/database/customers.go`: query pelanggan, autentikasi, dan manajemen profil.
- `backend/database/orders.go`: query pesanan dan filter pesanan berdasarkan pelanggan atau status.
- `backend/database/production.go`: query produksi dan ringkasan status produksi.
- `backend/database/stock_movements.go`: query mutasi stok.
- `backend/database/mysql.go`: koneksi database MySQL dan konfigurasi koneksi.

---

## 🔐 Alur Autentikasi JWT

### Admin
- Admin login melalui `POST /api/admin/login`.
- Jika kredensial valid, backend mengembalikan JWT yang ditandatangani dengan `ADMIN_JWT_SECRET`.
- Semua endpoint admin menggunakan token ini di header:

```bash
Authorization: Bearer <ADMIN_TOKEN>
```

### Customer
- Customer login melalui `POST /api/customers/login`.
- Jika email dan password valid, backend mengembalikan JWT terpisah yang ditandatangani dengan `CUSTOMER_JWT_SECRET`.
- Customer menggunakan token ini untuk endpoint seperti:

```bash
GET /api/customer/orders
GET /api/customer/profile
```

### Perbedaan utama
- Admin token digunakan untuk operasi CRUD data master, dashboard, dan laporan.
- Customer token digunakan untuk portal customer, melihat profil, dan melihat pesanan miliknya.
- Secret JWT admin dan customer disimpan di konfigurasi environment yang berbeda.

---

## 🚀 Contoh Alur Fitur Utama

### Alur Admin
1. Login admin
2. Tambah produk atau kategori baru
3. Tambah gudang dan stok awal
4. Buat pesanan atau pekerjaan produksi
5. Pantau dashboard untuk ringkasan pesanan, produksi, dan pendapatan

Contoh `curl` membuat produk:

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"PRD-001","name":"Produk A","categoryId":1,"unitId":1,"standardPrice":10000,"costPrice":8000,"description":"Contoh produk"}'
```

Contoh `curl` membuat pesanan:

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"orderNumber":"ORD-001","customerId":1,"product":"Produk A","orderDate":"2026-04-14T08:00:00Z","amount":150000,"status":"pending"}'
```

Contoh `curl` membuat produksi:

```bash
curl -X POST http://localhost:8080/api/production \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"jobCode":"JOB-001","product":"Produk A","startDate":"2026-04-15T08:00:00Z","durationDays":7,"status":"pending"}'
```

### Alur Customer
1. Login customer di portal `Customer Login`
2. Buka halaman `Pesanan Saya` untuk melihat riwayat pesanan
3. Buka halaman `Profil` untuk melihat data akun
4. Klik logout di halaman `Bantuan` bila selesai

Contoh `curl` login customer:

```bash
curl -X POST http://localhost:8080/api/customers/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}'
```

Contoh `curl` melihat profil customer:

```bash
curl http://localhost:8080/api/customer/profile \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN_HERE"
```

---

## 🔌 API Utama

- `GET /api/health`
- `POST /api/admin/login`
- `GET /api/admin/overview`
- `GET|POST /api/products`
- `GET|PUT|DELETE /api/products/:id`
- `GET|POST /api/categories`
- `GET|PUT|DELETE /api/categories/:id`
- `GET|POST /api/units`
- `GET|PUT|DELETE /api/units/:id`
- `GET|POST /api/warehouses`
- `GET|PUT|DELETE /api/warehouses/:id`
- `GET|POST /api/orders`
- `GET|PUT|DELETE /api/orders/:id`
- `GET|POST /api/customers`
- `GET|PUT|DELETE /api/customers/:id`
- `POST /api/customers/login`
- `GET /api/customer/orders`
- `GET /api/customer/profile`
- `GET|POST /api/production`
- `GET|PUT|DELETE /api/production/:id`
- `GET|POST /api/inventory`
- `GET|PUT|DELETE /api/inventory/:id`
- `GET|POST /api/stock_movements`
- `GET|PUT|DELETE /api/stock_movements/:id`

> Semua endpoint admin memerlukan header `Authorization: Bearer <token>`.
> Endpoint customer yang terproteksi memakai customer JWT di header `Authorization: Bearer <token>`.

### Contoh `curl`

Login admin dan dapatkan token:

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Menggunakan token untuk request yang membutuhkan otentikasi:

```bash
curl http://localhost:8080/api/admin/overview \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

Contoh customer login:

```bash
curl -X POST http://localhost:8080/api/customers/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}'
```

Contoh customer mengambil profil:

```bash
curl http://localhost:8080/api/customer/profile \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN_HERE"
```

Contoh membuat produk baru:

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"sku":"PRD-001","name":"Produk A","categoryId":1,"unitId":1,"standardPrice":10000,"costPrice":8000,"description":"Contoh produk"}'
```

---

## 📌 Catatan

- Frontend mengkonsumsi API backend langsung melalui `fetch` dan menggunakan JWT admin untuk akses.
- Dashboard dan laporan menggunakan data real-time dari backend untuk production, order, dan revenue.
- Struktur modular memudahkan pengembangan fitur tambahan seperti laporan lebih lanjut, notifikasi, atau multi-gudang.

---

## ✅ Build Production

### Frontend

```bash
cd frontend
npm run build
```

### Backend

```bash
cd backend
go build ./...
```
