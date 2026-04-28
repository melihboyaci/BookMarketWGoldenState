-- Migrasyon: 005_create_users_table
-- Açıklama: SaaS kimlik doğrulama için kullanıcılar tablosu.
-- tenant_id ile Golden State mekanizmasına entegre edilmiştir:
--   'demo_blueprint' -> Asla değiştirilmeyen şablon kullanıcılar
--   'demo_active'    -> Demo sırasında kayıt olan / değişen kullanıcılar

CREATE TABLE IF NOT EXISTS users (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(100)  NOT NULL,
    email         VARCHAR(255)  NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,
    role          VARCHAR(20)   NOT NULL CHECK (role IN ('ADMIN', 'SELLER', 'BUYER')),
    tenant_id     VARCHAR(50)   NOT NULL,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- (email, tenant_id) çifti unique: aynı e-posta farklı tenant'larda var olabilir
-- (blueprint kopyası vs active kaydı)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_tenant
    ON users (email, tenant_id);

-- tenant_id üzerinde index: restore sorgularını ve kimlik doğrulamayı hızlandırır
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id);

COMMENT ON TABLE users IS 'Demo kullanıcıları. tenant_id ile Golden State mekanizması yönetilir.';
COMMENT ON COLUMN users.password_hash IS 'bcrypt (cost=10) ile hashlenmiş şifre. Ham şifre hiçbir zaman saklanmaz.';
COMMENT ON COLUMN users.tenant_id IS 'demo_active: aktif demo | demo_blueprint: dokunulmaz altın şablon';
