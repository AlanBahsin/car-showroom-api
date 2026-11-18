# AutoDrive Showroom API (Backend Go + Gin + Supabase)

Backend REST API untuk sistem informasi & katalog showroom mobil **AutoDrive Showroom**, dilengkapi sistem rekomendasi mobil berbasis algoritma **Simple Additive Weighting (SAW)**, perbandingan spesifikasi, simulasi kredit, serta manajemen transaksi 3 role (`admin`, `sales`, `buyer`).

## Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin Web Framework
- **ORM & Database:** GORM + PostgreSQL (Supabase Cloud)
- **Authentication:** JWT (`golang-jwt/jwt/v5`) & Bcrypt

## Cara Menjalankan di Lokal
1. Clone repository ini:
   ```bash
   git clone <URL_REPO_GITHUB_KAMU>
   cd car-showroom-api