# Current Project State

*AI Ajanı bu dosyayı her faz veya görev bitiminde otomatik olarak güncellemelidir.*

## Status Check
- [x] **Phase 1:** Altyapı ve Veritabanı Hazırlığı ✅
- [x] **Phase 2:** Altın Veri (Seed) Yükleme ✅
- [x] **Phase 3:** Core API ve JWT Middleware ✅
- [x] **Phase 4:** Restore Golden State Mekanizması ✅

## Aktif Görev (Current Task)
Gerçek veritabanı destekli kimlik doğrulama refactoring'i yapıldı. `users` tablosu eklendi, hardcoded şifreler bcrypt'e geçirildi, `/register` endpoint'i eklendi ve Restore (Golden State) sıralaması FK güvenli hale getirildi.

## Phase 1–4 Sonrası — Güvenlik & Dayanıklılık & Testler ✅
- **[GÜVENLİK]** `system_handler.go`: iç hata detayı (SQL mesajı, tablo adı) artık API yanıtında gizleniyor, sunucu loguna yazılıyor
- **[RATE-LIMIT]** `restore_repository.go`: `sync.Mutex.TryLock()` ile eş zamanlı restore istekleri anında reddediliyor
- **[DAYANIKLILIK]** `db.go`: `SetConnMaxLifetime(30 * time.Minute)` eklendi — stale connection hatası önlendi

## Son Kararlar ve Notlar (Dev Notes)
- ORM kullanılmayacak, `database/sql` ile devam edilecek.
- Bağlantı havuzu: MaxOpenConns=25, MaxIdleConns=10, MaxConnLifetime=30m
- `tenant_id` üzerinde index oluşturuldu (Golden State reset sorgularını hızlandırmak için)
- **[SaaS Refactor]** Hardcoded auth yerine `users` tablosu ve **bcrypt** (cost=10) tabanlı bir sisteme geçildi. (Migration 005)
- **[Restore Sırası]** FK (Foreign Key) hatalarını engellemek için Golden State işlem sırası güncellendi: DELETE (orders -> users -> books), INSERT (books -> users).