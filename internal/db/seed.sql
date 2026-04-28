-- Seed Verisi: Golden State - Kitap Satış Demo Uygulaması
-- Türk Edebiyatı eserleri, TL fiyatlandırması

-- ── Güvenli Temizlik ───────────────────────────────────────────────────────
-- Foreign Key sırasına göre: önce orders, sonra users, sonra books
-- (orders.user_id -> users.id olduğunda çakışma oluşmaz)
DELETE FROM orders WHERE tenant_id IN ('demo_active', 'demo_blueprint');
DELETE FROM users  WHERE tenant_id IN ('demo_active', 'demo_blueprint');
DELETE FROM books  WHERE tenant_id IN ('demo_active', 'demo_blueprint');

-- ── DEMO_BLUEPRINT: Kullanıcı Şablonları (Asla Değişmez) ──────────────────
-- Şifreler bcrypt cost=10 ile hashlenmiştir:
--   admin@demo.com  -> admin123
--   seller@demo.com -> seller123
--   buyer@demo.com  -> buyer123
INSERT INTO users (username, email, password_hash, role, tenant_id) VALUES
    ('Admin Kullanıcı',  'admin@demo.com',  '$2a$10$vxvwRa9m.WRu7Tn/Dp.EMe2v8.TLmJepgE01ynlPBbMdk6enwRpHG', 'ADMIN',  'demo_blueprint'),
    ('Demo Satıcı',      'seller@demo.com', '$2a$10$IWxD7MYLXDPtgkX9yrzSBehjYByrk57W1xlT2hhdw1ln/5qMjT4qe', 'SELLER', 'demo_blueprint'),
    ('Demo Alıcı',       'buyer@demo.com',  '$2a$10$izanH5NRG6L1Ko7.pi9HgOhbrmHtmKUD.sAgYo139pmlvpFy5ntES', 'BUYER',  'demo_blueprint');

-- ── DEMO_ACTIVE: Başlangıç Aktif Kullanıcılar (Golden State Reset ile yenilenir) ──
INSERT INTO users (username, email, password_hash, role, tenant_id) VALUES
    ('Admin Kullanıcı',  'admin@demo.com',  '$2a$10$vxvwRa9m.WRu7Tn/Dp.EMe2v8.TLmJepgE01ynlPBbMdk6enwRpHG', 'ADMIN',  'demo_active'),
    ('Demo Satıcı',      'seller@demo.com', '$2a$10$IWxD7MYLXDPtgkX9yrzSBehjYByrk57W1xlT2hhdw1ln/5qMjT4qe', 'SELLER', 'demo_active'),
    ('Demo Alıcı',       'buyer@demo.com',  '$2a$10$izanH5NRG6L1Ko7.pi9HgOhbrmHtmKUD.sAgYo139pmlvpFy5ntES', 'BUYER',  'demo_active');

-- ── DEMO_BLUEPRINT: Kitap Şablonları (Asla Değişmez) ──────────────────────
INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock) VALUES
    ('demo_blueprint', 'Saatleri Ayarlama Enstitüsü', 'Ahmet Hamdi Tanpınar',        '978-9750718908', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000163351-1.jpg', 185.00, 80),
    ('demo_blueprint', 'Benim Adım Kırmızı',          'Orhan Pamuk',                 '978-9750726439', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000064074-1.jpg', 210.00, 65),
    ('demo_blueprint', 'Beyaz Kale',                  'Orhan Pamuk',                 '978-9750726088', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000064010-1.jpg', 195.50, 90),
    ('demo_blueprint', 'Tutunamayanlar',               'Oğuz Atay',                   '978-9750518256', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000061214-1.jpg', 220.00, 50),
    ('demo_blueprint', 'Çalıkuşu',                    'Reşat Nuri Güntekin',         '978-9750736131', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000057973-1.jpg', 175.00, 120),
    ('demo_blueprint', 'Kürk Mantolu Madonna',         'Sabahattin Ali',              '978-9750731501', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000058317-1.jpg', 165.00, 100),
    ('demo_blueprint', 'İnce Memed',                  'Yaşar Kemal',                 '978-9750803352', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000063538-1.jpg', 190.00, 70),
    ('demo_blueprint', 'Huzur',                       'Ahmet Hamdi Tanpınar',        '978-9750718823', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000163350-1.jpg', 200.00, 55),
    ('demo_blueprint', 'Yaban',                       'Yakup Kadri Karaosmanoğlu',   '978-9750718243', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000061617-1.jpg', 160.00, 85),
    ('demo_blueprint', 'Mavi Sürgün',                 'Cevat Şakir Kabaağaçlı',     '978-9750724862', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000062758-1.jpg', 155.00, 110);

-- ── DEMO_ACTIVE: Başlangıç Aktif Kitaplar (Golden State Reset ile yenilenir) ──
INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock) VALUES
    ('demo_active', 'Saatleri Ayarlama Enstitüsü', 'Ahmet Hamdi Tanpınar',        '978-9750718908', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000163351-1.jpg', 185.00, 80),
    ('demo_active', 'Benim Adım Kırmızı',          'Orhan Pamuk',                 '978-9750726439', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000064074-1.jpg', 210.00, 65),
    ('demo_active', 'Beyaz Kale',                  'Orhan Pamuk',                 '978-9750726088', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000064010-1.jpg', 195.50, 90),
    ('demo_active', 'Tutunamayanlar',               'Oğuz Atay',                   '978-9750518256', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000061214-1.jpg', 220.00, 50),
    ('demo_active', 'Çalıkuşu',                    'Reşat Nuri Güntekin',         '978-9750736131', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000057973-1.jpg', 175.00, 120),
    ('demo_active', 'Kürk Mantolu Madonna',         'Sabahattin Ali',              '978-9750731501', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000058317-1.jpg', 165.00, 100),
    ('demo_active', 'İnce Memed',                  'Yaşar Kemal',                 '978-9750803352', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000063538-1.jpg', 190.00, 70),
    ('demo_active', 'Huzur',                       'Ahmet Hamdi Tanpınar',        '978-9750718823', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000163350-1.jpg', 200.00, 55),
    ('demo_active', 'Yaban',                       'Yakup Kadri Karaosmanoğlu',   '978-9750718243', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000061617-1.jpg', 160.00, 85),
    ('demo_active', 'Mavi Sürgün',                 'Cevat Şakir Kabaağaçlı',     '978-9750724862', 'https://i.dr.com.tr/cache/600x600-0/originals/0000000062758-1.jpg', 155.00, 110);
