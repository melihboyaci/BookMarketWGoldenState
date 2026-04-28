# Frontend Architecture

## 1. Klasör Yapısı (Vite Standardı)
```text
src/
 ├── assets/          # Resimler, logolar
 ├── components/      # Tekrar kullanılabilir UI parçaları (Button, Modal, Card)
 ├── pages/           # Ana görünümler (Login, Dashboard)
 ├── services/        # Backend API istekleri (api.js)
 ├── context/         # AuthContext (Kullanıcı oturum yönetimi)
 └── App.jsx          # Ana Router

## 4. Guest Checkout (Misafir Alışverişi) ve Auth Akışı
- **Ana Sayfa / Katalog:** Herkese açıktır. Kullanıcılar üye olmadan (AuthContext null iken) kitapları görebilir ve sepete ekleyebilir. Sepet verisi (CartContext) sadece tarayıcıda yaşar.
- **Lazy Registration (Gecikmeli Kayıt):** Arayüzde doğrudan bir 'Kayıt Ol' duvarı yoktur. Kullanıcı ancak sepetindeki ürünleri almak için "Siparişi Tamamla" butonuna bastığında karşısına bir Login/Register Modal'ı (açılır pencere) veya sayfası gelir.
- **Kayıt Sonrası:** Başarıyla kayıt olup/giriş yapıp JWT alındığında, sepetteki veriler backend `orders` endpoint'ine gönderilir.