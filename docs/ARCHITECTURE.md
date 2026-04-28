# System Architecture

## 1. Klasör Yapısı (Project Layout)
Proje, standart Go backend dizin yapısına uygun olacaktır:
```text
├── cmd/
│   └── api/
│       └── main.go           # Uygulama başlangıç noktası
├── internal/
│   ├── db/                   # DB bağlantıları ve migrasyon betikleri
│   ├── handlers/             # HTTP endpoint fonksiyonları
│   ├── middleware/           # Auth ve Role kontrolleri
│   ├── models/               # Struct tanımlamaları
│   └── repository/           # Saf SQL sorgularının bulunduğu katman
├── docs/                     # PRD, Architecture, Plan ve State dosyaları
├── docker-compose.yml
├── Dockerfile
└── .env

## 2. Veritabanı Şeması (PostgreSQL)
Sistemde tüm tablolar `tenant_id` ('demo_active' veya 'demo_blueprint') kolonuna sahiptir.
- `users`: id (UUID, PK), username (VARCHAR UNIQUE), password_hash (VARCHAR), role (VARCHAR), tenant_id (VARCHAR)
- `books`: id (UUID, PK), title (VARCHAR), author (VARCHAR), price (DECIMAL), image_url (VARCHAR), tenant_id (VARCHAR)
- `orders`: id (UUID, PK), user_id (UUID, FK -> users.id), total_amount (DECIMAL), tenant_id (VARCHAR)

## 3. Temel API Endpoint'leri
- `POST /api/v1/auth/register` : Yeni kullanıcı kaydeder. Şifre **bcrypt** ile hashlenir. Eklenen kullanıcı `tenant_id='demo_active'` olarak kaydedilir.
- `POST /api/v1/auth/login` : DB'deki bcrypt hash ile gelen şifreyi kıyaslar, geçerliyse JWT döner.
- `POST /api/v1/system/restore` : Sadece ADMIN çalıştırabilir. İşlem sırası (Foreign Key çakışmasını önlemek için):
  1. DELETE (Sırayla): orders -> users -> books (Sadece demo_active olanlar)
  2. INSERT ... SELECT (Sırayla): books -> users (Sadece demo_blueprint olanları active olarak kopyala)