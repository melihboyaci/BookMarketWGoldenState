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