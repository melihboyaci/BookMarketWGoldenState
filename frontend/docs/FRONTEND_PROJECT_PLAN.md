# Frontend Project Plan

AI ajanı, bu fazları sırayla takip etmeli ve her faz bitiminde commit istemelidir.

## Phase 1: Vite & Tailwind Kurulumu
- `npm create vite@latest` ile React projesinin oluşturulması.
- Tailwind CSS'in kurulması ve yapılandırılması.
- Klasör yapısının (`components`, `pages`, `services`) oluşturulması.

## Phase 2: API Servisleri ve Context
- `services/api.js` dosyasının yazılması (Login, Kitapları Çek, Restore Golden State fonksiyonları).
- `context/AuthContext.jsx` oluşturulup global JWT ve Role state'inin yönetilmesi.

## Phase 3: Temel Arayüz (UI)
- `LoginPage` bileşeninin tasarlanması (Buyer, Seller, Admin giriş butonları).
- `DashboardPage` içine kitapları listeleyen Grid/Card yapısının kurulması.

## Phase 4: Golden State Entegrasyonu
- Sadece `ADMIN` rolüne görünen "Sistemi Sıfırla" butonunun eklenmesi.
- Butona basıldığında API isteği atılması, yükleniyor (spinner) gösterilmesi ve başarılı olunca kitap listesinin otomatik yenilenmesi.