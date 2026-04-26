# Product Requirements Document (PRD)
**Proje Adı:** Golden State - Kitap Satış Demo Uygulaması

## 1. Projenin Amacı
B2B satış toplantılarında müşterilere gösterilmek üzere hazırlanan bir kitap satış uygulamasının, testler veya müşteri etkileşimi sırasında bozulan verilerini milisaniyeler içinde ilk "Satışa Hazır" (Golden State) durumuna getiren bir backend altyapısı kurmak.

## 2. Kullanıcı Rolleri (Personas)
- **Buyer (Alıcı):** Kitapları listeler, arama yapar, satın alma işlemi simüle eder.
- **Seller (Satıcı):** Kitap ekler, stok günceller. (Demo sırasında sistemi en çok kirleten roldür).
- **Admin (Yönetici):** Sistemi izler ve canlı demo sırasında `Restore Golden State` işlemini tetikleme yetkisine sahip tek roldür.

## 3. Temel İş Akışı (The Golden Reset)
Sistemde aynı tablo içinde (`books` vb.) iki farklı veri kümesi yaşar:
1. `tenant_id = 'demo_active'`: Uygulamanın anlık kullandığı, kirlenmeye müsait veri.
2. `tenant_id = 'demo_blueprint'`: Geliştirme aşamasında üretilmiş, asla değiştirilmeyen "Altın Şablon" verisi.

Admin sıfırlama butonuna bastığında, sistem `demo_active` verilerini siler ve `demo_blueprint` verilerini kopyalayarak sistemi saniyeler içinde ilk haline getirir.