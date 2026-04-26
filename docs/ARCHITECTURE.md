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