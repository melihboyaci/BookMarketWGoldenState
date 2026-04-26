-- Migrasyon: 001_create_books_table
-- Açıklama: Kitap satış uygulamasının ana tablosu.
-- tenant_id sütunu Golden State mekanizmasının temel taşıdır:
--   'demo_active'    -> Uygulamanın anlık kullandığı, kirlenmeye müsait veri
--   'demo_blueprint' -> Asla değiştirilmeyen Altın Şablon verisi

CREATE TABLE IF NOT EXISTS books (
    id          BIGSERIAL       PRIMARY KEY,
    tenant_id   VARCHAR(50)     NOT NULL,           -- 'demo_active' veya 'demo_blueprint'
    title       VARCHAR(255)    NOT NULL,
    author      VARCHAR(255)    NOT NULL,
    isbn        VARCHAR(20)     UNIQUE,
    price       NUMERIC(10, 2)  NOT NULL CHECK (price >= 0),
    stock       INTEGER         NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- tenant_id üzerinde index: Golden State reset sorgularını hızlandırır.
CREATE INDEX IF NOT EXISTS idx_books_tenant_id ON books (tenant_id);

COMMENT ON TABLE books IS 'Kitap satış demo uygulamasının ana tablosu. tenant_id ile Golden State mekanizması yönetilir.';
COMMENT ON COLUMN books.tenant_id IS 'demo_active: kullanılan veri | demo_blueprint: dokunulmaz altın şablon';
