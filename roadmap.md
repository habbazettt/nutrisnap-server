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

---

# **EPIC 2 — Authentication & Authorization**

Tujuan: menyediakan autentikasi aman dan modern, sekelas aplikasi production.

### Checklist

* [x] Registrasi user (email + password hash)
* [x] Login user (email + password)
* [x] Generate JWT access tokens
* [x] JWT middleware untuk protected routes
* [x] Endpoint “Get Current User”
* [x] Endpoint "Update Profile"
* [x] Endpoint "Change Password"
* [x] Implementasi Google OAuth2 login
* [x] Model OAuthAccount + relasi ke User
* [x] Role-based access (user, admin)

---

# **EPIC 3 — Image Ingestion & Storage**

Tujuan: menangani upload gambar nutrition facts secara efisien dan aman.

### Checklist

* [ ] Endpoint upload gambar nutrition (`POST /v1/scan`)
* [ ] Validasi file (ukuran, tipe MIME, dimensi minimal)
* [ ] Opsi `store_image=true|false` (opt-in)
* [ ] Upload file ke MinIO dan simpan `image_ref`
* [ ] Endpoint presigned URL untuk akses gambar
* [ ] Endpoint delete image
* [ ] Implementasi retention policy (cleanup rutin)

---

# **EPIC 4 — Barcode Fast-Path Pipeline**

Tujuan: memproses barcode untuk hasil yang cepat & akurat, sebelum OCR.

### Checklist

* [ ] Terima `barcode` pada endpoint upload
* [ ] Integrasi API OpenFoodFacts
* [ ] Mapping nutriments OF → canonical nutrient schema
* [ ] Simpan ke tabel `products` sebagai cache
* [ ] Endpoint `GET /v1/product/:barcode`
* [ ] Fallback otomatis ke OCR jika tidak ditemukan
* [ ] Script sync subset produk populer ke DB

---

# **EPIC 5 — OCR Processing Pipeline**

Tujuan: mengubah foto nutrition facts menjadi teks mentah untuk parsing.

### Checklist

* [ ] Image preprocessing (rotate, grayscale, autocontrast, resize)
* [ ] Eksekusi Tesseract OCR via CLI dalam container
* [ ] Simpan hasil raw OCR (`ocr_raw`)
* [ ] Heuristik confidence dasar (hitungan digit, panjang teks, dsb.)
* [ ] Error handling untuk kasus foto blur atau kondisi buruk

---

# **EPIC 6 — Nutrition Fact Parser & Normalizer**

Tujuan: mengekstrak tabel gizi dari teks OCR menjadi format terstruktur.

### Checklist

* [ ] Fuzzy mapping nama nutrisi (English + Indonesia)
* [ ] Regex ekstraksi angka + digit cleanup (O→0, l→1, comma→dot)
* [ ] Normalisasi unit (kJ→kcal, mg<->g, dsb.)
* [ ] Normalisasi per-100g
* [ ] Normalisasi per-serving (jika tersedia)
* [ ] Konsistensi energi vs makronutrien
* [ ] Simpan `parsed_json` & `normalized_json` ke tabel `scans`

---

# **EPIC 7 — Scoring, Highlights, Insights**

Tujuan: menghasilkan interpretasi nutrisi dalam bahasa sederhana untuk pengguna awam.

### Checklist

* [ ] Implementasi NutriScore-style scoring (A–E)
* [ ] Nutrient highlight (low, medium, high)
* [ ] Indikator bahaya (high sugar, high sodium, high calorie)
* [ ] Template-based insight generation
* [ ] Integrasi scoring & insights ke `/v1/scan/:id` response

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

* [ ] Endpoint `GET /v1/history?user_id=&page=`
* [ ] Simpan metadata scan (barcode, score, timestamps, image_ref)
* [ ] Endpoint `POST /v1/scan/:id/correct` untuk koreksi data nutrisi
* [ ] Simpan koreksi ke tabel `corrections`
* [ ] Tampilkan hasil koreksi dalam endpoint detail scan

---

# **EPIC 10 — Reporting, Monitoring & Observability**

Tujuan: memantau performa sistem dan mendukung debugging.

### Checklist

* [ ] Implementasi metrics Prometheus (scan_count, barcode_hit_rate, ocr_latency)
* [ ] Dashboard Grafana
* [ ] Logging terstruktur + correlation ID
* [ ] Error tracking (Sentry / middleware Fiber)
* [ ] Endpoint `/metrics`

---

# **EPIC 11 — Documentation & Developer Experience**

Tujuan: memastikan maintainability & onboarding developer berjalan baik.

### Checklist

* [ ] `architecture.md` (flow pipeline, komponen, data alur)
* [ ] Dokumentasi folder structure
* [ ] Dokumentasi API via Swagger/OpenAPI
* [ ] Panduan local development
* [ ] Troubleshooting OCR/MinIO/Parser
* [ ] Developer runbook (backup, restore, maintenance)

---

# **NutriSnap Roadmap Finalized**

Roadmap ini sengaja disederhanakan, fokus pada fitur inti yang benar-benar memberikan value kepada user awam:
**foto → nutrition facts → insight → score → compare → history.**
