# **NutriSnap Backend Roadmap**

Dokumen ini menjadi blueprint pengembangan backend **NutriSnap**, sebuah platform yang memproses foto nutrition facts atau barcode untuk menghasilkan informasi nutrisi, skor kesehatan, insight, dan perbandingan antar produk.
Roadmap disusun berbasis EPIC untuk mempermudah eksekusi bertahap oleh solo developer.

---

# **EPIC 1 — Core Platform Foundation**

Tujuan: membangun fondasi backend yang stabil, modular, aman, dan siap mensupport pipeline AI (OCR + parser).

### Checklist

* [x] Setup repository & struktur proyek (Go + Fiber)
* [x] Setup readme & license
* [x] Setup Docker Compose (App, PostgreSQL, Adminer, MinIO, Tesseract runtime)
* [x] Tambahkan `.env.example` dan config loader
* [x] Setup structured logging (slog)
* [x] Tambahkan rate limiting middleware
* [x] Setup API response envelope (success, error, metadata)
* [x] Inisialisasi database (GORM) + AutoMigrate (dev environment)
* [x] Setup API versioning (`/api/v1`)
* [x] Setup folder structure: routes, controllers, services, repositories
* [x] Setup health check endpoint (`/api/v1/healthz`)
* [x] Integrasi Swagger/OpenAPI (`/docs`)
* [x] Integrasi Prometheus & Grafana (monitoring dasar)
* [x] Setup API Status constants
* [x] Timezone WIB (Asia/Jakarta) standardization

---

# **EPIC 2 — Authentication & Authorization**

Tujuan: menyediakan autentikasi aman dan modern, sekelas aplikasi production.

### Checklist

* [x] Registrasi user (email + password hash)
* [x] Login user (email + password)
* [x] Generate JWT access tokens
* [x] JWT middleware untuk protected routes
* [x] Endpoint "Get Current User"
* [x] Endpoint "Update Profile"
* [x] Endpoint "Change Password"
* [x] Implementasi Google OAuth2 login
* [x] Model OAuthAccount + relasi ke User
* [x] Role-based access (user, admin)
* [x] Refresh token

---

# **EPIC 3 — Image Ingestion & Storage**

Tujuan: menangani upload gambar nutrition facts secara efisien dan aman.

### Checklist

* [x] Buat MinIO client
* [x] Endpoint upload gambar nutrition (`POST /v1/scan`)
* [x] Validasi file (ukuran, tipe MIME, dimensi minimal)
* [x] Opsi `store_image=true|false` (opt-in)
* [x] Upload file ke MinIO dan simpan `image_ref`
* [x] Endpoint presigned URL untuk akses gambar
* [x] Endpoint delete image
* [x] Implementasi retention policy (cleanup rutin)

---

# **EPIC 4 — Barcode Fast-Path Pipeline**

Tujuan: memproses barcode untuk hasil yang cepat & akurat, sebelum OCR.

### Checklist

* [x] Integrasi API OpenFoodFacts
* [x] Terima `barcode` pada endpoint upload
* [x] Mapping nutriments OF → canonical nutrient schema
* [x] Simpan ke tabel `products` sebagai cache
* [x] Endpoint `GET /v1/product/:barcode`
* [x] Fallback otomatis ke OCR jika tidak ditemukan

---

# **EPIC 5 — OCR Processing Pipeline**

Tujuan: mengubah foto nutrition facts menjadi teks mentah untuk parsing.

### Checklist

* [x] Image preprocessing (handled by OCR lib)
* [x] Eksekusi Tesseract OCR (gosseract + CGO)
* [x] Simpan hasil raw OCR (`ocr_raw`)
* [x] Background Worker (Queue System)
* [x] Error handling untuk kasus foto blur atau kondisi buruk
* [x] Bahasa Indonesia + English support

---

# **EPIC 6 — Nutrition Fact Parser & Normalizer**

Tujuan: mengekstrak tabel gizi dari teks OCR menjadi format terstruktur.

### Checklist

* [x] Regex untuk nama nutrisi (Indonesia + English)
* [x] Regex ekstraksi angka + digit cleanup
* [x] Normalisasi unit (Dot/Comma handling)
* [x] Simpan `nutrients` ke tabel `products`
* [x] Serving Size detection (`Takaran Saji`)
* [x] OCR Typo Correction (`cleanupOCRText`)
* [x] Bilingual text handling (separator `/`)

---

# **EPIC 7 — Scoring, Highlights, Insights**

Tujuan: menghasilkan interpretasi nutrisi dalam bahasa sederhana untuk pengguna awam.

### Checklist

* [x] Implementasi NutriScore-style scoring (A–E)
* [x] Nutrient highlight (low, medium, high)
* [x] Indikator bahaya (high sugar, high sodium, etc)
* [x] Integrasi scoring & insights ke `/v1/scan/:id` response

---

# **EPIC 8 — Product Comparison Engine**

Tujuan: memungkinkan pengguna membandingkan dua produk secara objektif.

### Checklist

* [ ] Endpoint `POST /v1/compare`
* [ ] Ambil dua produk via barcode atau scan_id
* [ ] Normalize kedua produk ke basis per-100g
* [ ] Hitung delta absolut & persentase
* [ ] Generate verdict (human readable)
* [ ] Return JSON structured comparison

---

# **EPIC 9 — Scan History & User Corrections**

Tujuan: menyediakan riwayat scan dan kemampuan koreksi manual dari pengguna.

### Checklist

* [x] Endpoint `GET /v1/scan` (paginated history)
* [x] Simpan metadata scan (barcode, score, timestamps, image_ref)
* [x] Endpoint `POST /v1/scan/:id/correct` untuk koreksi data nutrisi
* [x] Simpan koreksi ke tabel `corrections`
* [ ] Tampilkan hasil koreksi dalam endpoint detail scan

---

# **EPIC 10 — Reporting, Monitoring & Observability**

Tujuan: memantau performa sistem dan mendukung debugging.

### Checklist

* [x] Implementasi metrics Prometheus
* [x] Dashboard Grafana
* [x] Logging terstruktur (slog)
* [x] Endpoint `/metrics`

---

# **EPIC 11 — Documentation & Developer Experience**

Tujuan: memastikan maintainability & onboarding developer berjalan baik.

### Checklist

* [x] Dokumentasi folder structure
* [x] Dokumentasi API via Swagger/OpenAPI
* [x] Panduan local development
