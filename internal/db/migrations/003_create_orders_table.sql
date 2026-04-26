-- Migrasyon: 003_create_orders_table
-- Açıklama: E-Ticaret senaryosu için siparişler ve sipariş detayları.
-- Golden State restore işleminden etkilenmesi için tenant_id içerir.

CREATE TABLE IF NOT EXISTS orders (
    id          BIGSERIAL       PRIMARY KEY,
    tenant_id   VARCHAR(50)     NOT NULL,           -- 'demo_active'
    total_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id          BIGSERIAL       PRIMARY KEY,
    order_id    BIGINT          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    book_id     BIGINT          NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    quantity    INTEGER         NOT NULL CHECK (quantity > 0),
    price       NUMERIC(10, 2)  NOT NULL CHECK (price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders (tenant_id);
