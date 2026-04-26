-- Migrasyon: 004_add_image_url_to_books
-- Açıklama: Kitaplara kapak görseli URL'si ekler.

ALTER TABLE books ADD COLUMN IF NOT EXISTS image_url TEXT;
