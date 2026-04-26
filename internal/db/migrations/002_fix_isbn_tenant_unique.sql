-- Migrasyon: 002_fix_isbn_tenant_unique_constraint
-- Açıklama: isbn sütunundaki global UNIQUE kısıtı, aynı ISBN'in
-- farklı tenant'larda (demo_blueprint ve demo_active) bulunmasına
-- izin vermez. Bu migrasyon:
--   1. Eski '-a' suffix'li aktif kayıtları temizler
--   2. Global kısıtı kaldırır
--   3. Yerine (isbn, tenant_id) bileşik UNIQUE index ekler
--
-- Bu migrasyon idempotent'tir: tekrar çalıştırılabilir.

-- Adım 1: Eski seed tarafından '-a' suffix'le oluşturulan active kayıtları sil
DELETE FROM books WHERE tenant_id = 'demo_active' AND isbn LIKE '%-a';

-- Adım 2: Global UNIQUE kısıtını kaldır
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_isbn_key;

-- Adım 3: Varsa hatalı oluşturulmuş partial indexi kaldır
DROP INDEX IF EXISTS idx_books_isbn_tenant;

-- Adım 4: Aynı ISBN farklı tenant'larda olabilir, aynı tenant'ta olamaz (Doğru kısıt)
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_isbn_tenant_key;
ALTER TABLE books ADD CONSTRAINT books_isbn_tenant_key UNIQUE (isbn, tenant_id);

