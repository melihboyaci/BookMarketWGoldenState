# Project Plan & Phases

AI ajanı, bu fazları sırayla takip etmeli ve her faz bitiminde onay istemelidir.

## Phase 1: Altyapı ve Veritabanı Hazırlığı
- Projenin `git init` ile başlatılması ve `.gitignore` (Go ve ortam değişkenleri için) dosyasının oluşturulması.
- `docker-compose.yml` dosyasının oluşturulması (Sadece PostgreSQL içerecek).
- Go modüllerinin başlatılması (`go mod init`).
- `internal/db` içinde PostgreSQL bağlantı kodlarının yazılması.
- `books` tablosunu oluşturacak SQL migrasyon betiğinin yazılması.

## Phase 2: Altın Veri (Seed) Yükleme
- Test ve Blueprint verilerini içeren bir `seed.sql` dosyası oluşturulması.
- Uygulama ayağa kalkarken `demo_blueprint` verilerinin veritabanına otomatik basılması.

## Phase 3: Core API ve JWT Middleware
- Gin router'ın `cmd/api/main.go` içinde ayağa kaldırılması.
- Basit bir JWT oluşturma mantığının yazılması.
- Sadece "ADMIN" rolünü kontrol eden `RequireAdminRole` middleware'inin `internal/middleware` içine yazılması.

## Phase 4: Restore Golden State Mekanizması
- `internal/repository` içinde veritabanı Transaction (Tx) mantığının kurulması.
- `POST /api/v1/system/restore` endpoint'inin yazılması. (demo_active sil, demo_blueprint'i kopyala, commit et).