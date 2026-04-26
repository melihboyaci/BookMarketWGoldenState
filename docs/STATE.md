# Current Project State

*AI Ajanı bu dosyayı her faz veya görev bitiminde otomatik olarak güncellemelidir.*

## Status Check
- [x] **Phase 1:** Altyapı ve Veritabanı Hazırlığı ✅
- [x] **Phase 2:** Altın Veri (Seed) Yükleme ✅
- [x] **Phase 3:** Core API ve JWT Middleware ✅
- [x] **Phase 4:** Restore Golden State Mekanizması ✅

## Aktif Görev (Current Task)
Unit ve Integration testleri tamamlandı (Tüm testler PASS). Test altyapısı hazırlandı. Sonraki adım: Frontend (React SPA) geliştirme süreci.

## Phase 1–4 Sonrası — Güvenlik & Dayanıklılık & Testler ✅
- **[GÜVENLİK]** `system_handler.go`: iç hata detayı (SQL mesajı, tablo adı) artık API yanıtında gizleniyor, sunucu loguna yazılıyor
- **[RATE-LIMIT]** `restore_repository.go`: `sync.Mutex.TryLock()` ile eş zamanlı restore istekleri anında reddediliyor
- **[DAYANIKLILIK]** `db.go`: `SetConnMaxLifetime(30 * time.Minute)` eklendi — stale connection hatası önlendi

## Son Kararlar ve Notlar (Dev Notes)
- ORM kullanılmayacak, `database/sql` ile devam edilecek.
- Bağlantı havuzu: MaxOpenConns=25, MaxIdleConns=10, MaxConnLifetime=30m
- `tenant_id` üzerinde index oluşturuldu (Golden State reset sorgularını hızlandırmak için)
- Sonraki aşama: Unit/Integration testler → Frontend geliştirme