-- Seed Verisi: Golden State - Kitap Satış Demo Uygulaması
-- Türk Edebiyatı eserleri, TL fiyatlandırması

CREATE TABLE IF NOT EXISTS orders (
    id          BIGSERIAL       PRIMARY KEY,
    tenant_id   VARCHAR(50)     NOT NULL,
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
ALTER TABLE books ADD COLUMN IF NOT EXISTS image_url TEXT;

-- Mevcut kitapları temizle
DELETE FROM books;

-- DEMO_BLUEPRINT: Altın Şablon (Asla Değişmez)
INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock) VALUES
    ('demo_blueprint', 'Saatleri Ayarlama Enstitüsü', 'Ahmet Hamdi Tanpınar',        '978-9750718908', 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=800&auto=format&fit=crop&q=60', 185.00, 80),
    ('demo_blueprint', 'Benim Adım Kırmızı',          'Orhan Pamuk',                 '978-9750726439', 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800&auto=format&fit=crop&q=60', 210.00, 65),
    ('demo_blueprint', 'Beyaz Kale',                  'Orhan Pamuk',                 '978-9750726088', 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=800&auto=format&fit=crop&q=60', 195.50, 90),
    ('demo_blueprint', 'Tutunamayanlar',               'Oğuz Atay',                   '978-9750518256', 'https://images.unsplash.com/photo-1476275466078-4cdc8100eca9?w=800&auto=format&fit=crop&q=60', 220.00, 50),
    ('demo_blueprint', 'Çalıkuşu',                    'Reşat Nuri Güntekin',         '978-9750736131', 'https://images.unsplash.com/photo-1524985069026-dd778a71c7b4?w=800&auto=format&fit=crop&q=60', 175.00, 120),
    ('demo_blueprint', 'Kürk Mantolu Madonna',         'Sabahattin Ali',              '978-9750731501', 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=800&auto=format&fit=crop&q=60', 165.00, 100),
    ('demo_blueprint', 'İnce Memed',                  'Yaşar Kemal',                 '978-9750803352', 'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800&auto=format&fit=crop&q=60', 190.00, 70),
    ('demo_blueprint', 'Huzur',                       'Ahmet Hamdi Tanpınar',        '978-9750718823', 'https://images.unsplash.com/photo-1455390582262-044cdead277a?w=800&auto=format&fit=crop&q=60', 200.00, 55),
    ('demo_blueprint', 'Yaban',                       'Yakup Kadri Karaosmanoğlu',   '978-9750718243', 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=800&auto=format&fit=crop&q=60', 160.00, 85),
    ('demo_blueprint', 'Mavi Sürgün',                 'Cevat Şakir Kabaağaçlı',     '978-9750724862', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&auto=format&fit=crop&q=60', 155.00, 110);

-- DEMO_ACTIVE: Başlangıç Aktif Verisi (Golden State Reset ile yenilenir)
INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock) VALUES
    ('demo_active', 'Saatleri Ayarlama Enstitüsü', 'Ahmet Hamdi Tanpınar',        '978-9750718908', 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=800&auto=format&fit=crop&q=60', 185.00, 80),
    ('demo_active', 'Benim Adım Kırmızı',          'Orhan Pamuk',                 '978-9750726439', 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800&auto=format&fit=crop&q=60', 210.00, 65),
    ('demo_active', 'Beyaz Kale',                  'Orhan Pamuk',                 '978-9750726088', 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=800&auto=format&fit=crop&q=60', 195.50, 90),
    ('demo_active', 'Tutunamayanlar',               'Oğuz Atay',                   '978-9750518256', 'https://images.unsplash.com/photo-1476275466078-4cdc8100eca9?w=800&auto=format&fit=crop&q=60', 220.00, 50),
    ('demo_active', 'Çalıkuşu',                    'Reşat Nuri Güntekin',         '978-9750736131', 'https://images.unsplash.com/photo-1524985069026-dd778a71c7b4?w=800&auto=format&fit=crop&q=60', 175.00, 120),
    ('demo_active', 'Kürk Mantolu Madonna',         'Sabahattin Ali',              '978-9750731501', 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=800&auto=format&fit=crop&q=60', 165.00, 100),
    ('demo_active', 'İnce Memed',                  'Yaşar Kemal',                 '978-9750803352', 'https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800&auto=format&fit=crop&q=60', 190.00, 70),
    ('demo_active', 'Huzur',                       'Ahmet Hamdi Tanpınar',        '978-9750718823', 'https://images.unsplash.com/photo-1455390582262-044cdead277a?w=800&auto=format&fit=crop&q=60', 200.00, 55),
    ('demo_active', 'Yaban',                       'Yakup Kadri Karaosmanoğlu',   '978-9750718243', 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=800&auto=format&fit=crop&q=60', 160.00, 85),
    ('demo_active', 'Mavi Sürgün',                 'Cevat Şakir Kabaağaçlı',     '978-9750724862', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&auto=format&fit=crop&q=60', 155.00, 110);
