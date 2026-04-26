-- Seed Verisi: Golden State - Kitap Satış Demo Uygulaması
-- Açıklama: Bu dosya hem 'demo_blueprint' (asla değişmeyen altın şablon)
-- hem de 'demo_active' (başlangıçta blueprint'in kopyası olan aktif veri) için
-- örnek kitap kayıtları içerir.
--
-- Bu betik idempotent'tir: Aynı isbn zaten varsa çakışmayı görmezden gelir
-- (ON CONFLICT DO NOTHING) ve güvenle defalarca çalıştırılabilir.

-- ============================================================
-- BÖLÜM 1: DEMO_BLUEPRINT — Altın Şablon (Asla Değiştirilmez)
-- ============================================================
INSERT INTO books (tenant_id, title, author, isbn, price, stock) VALUES
    ('demo_blueprint', 'The Pragmatic Programmer',         'David Thomas, Andrew Hunt',      '978-0135957059',  49.99, 120),
    ('demo_blueprint', 'Clean Code',                       'Robert C. Martin',               '978-0132350884',  44.99, 85),
    ('demo_blueprint', 'Designing Data-Intensive Apps',    'Martin Kleppmann',               '978-1449373320',  59.99, 60),
    ('demo_blueprint', 'The Go Programming Language',      'Alan Donovan, Brian Kernighan',  '978-0134190440',  39.99, 150),
    ('demo_blueprint', 'Domain-Driven Design',             'Eric Evans',                     '978-0321125217',  54.99, 45),
    ('demo_blueprint', 'Release It!',                      'Michael T. Nygard',              '978-1680502398',  42.99, 70),
    ('demo_blueprint', 'Site Reliability Engineering',     'Beyer, Jones, Petoff, Murphy',   '978-1491929124',  64.99, 35),
    ('demo_blueprint', 'Accelerate',                       'Nicole Forsgren, Jez Humble',    '978-1942788331',  34.99, 90),
    ('demo_blueprint', 'Building Microservices',           'Sam Newman',                     '978-1492034025',  57.99, 55),
    ('demo_blueprint', 'Staff Engineer',                   'Will Larson',                    '978-1736417911',  29.99, 200)
ON CONFLICT (isbn) DO NOTHING;

-- ============================================================
-- BÖLÜM 2: DEMO_ACTIVE — Başlangıç Aktif Verisi (Blueprint Kopyası)
-- Bu veri demo sırasında kirlenebilir; Golden State Reset ile sıfırlanır.
-- ============================================================
INSERT INTO books (tenant_id, title, author, isbn, price, stock) VALUES
    ('demo_active', 'The Pragmatic Programmer',         'David Thomas, Andrew Hunt',      '978-0135957059-a',  49.99, 120),
    ('demo_active', 'Clean Code',                       'Robert C. Martin',               '978-0132350884-a',  44.99, 85),
    ('demo_active', 'Designing Data-Intensive Apps',    'Martin Kleppmann',               '978-1449373320-a',  59.99, 60),
    ('demo_active', 'The Go Programming Language',      'Alan Donovan, Brian Kernighan',  '978-0134190440-a',  39.99, 150),
    ('demo_active', 'Domain-Driven Design',             'Eric Evans',                     '978-0321125217-a',  54.99, 45),
    ('demo_active', 'Release It!',                      'Michael T. Nygard',              '978-1680502398-a',  42.99, 70),
    ('demo_active', 'Site Reliability Engineering',     'Beyer, Jones, Petoff, Murphy',   '978-1491929124-a',  64.99, 35),
    ('demo_active', 'Accelerate',                       'Nicole Forsgren, Jez Humble',    '978-1942788331-a',  34.99, 90),
    ('demo_active', 'Building Microservices',           'Sam Newman',                     '978-1492034025-a',  57.99, 55),
    ('demo_active', 'Staff Engineer',                   'Will Larson',                    '978-1736417911-a',  29.99, 200)
ON CONFLICT (isbn) DO NOTHING;
