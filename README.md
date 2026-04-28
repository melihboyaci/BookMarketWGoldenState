# 🌟 Golden State: B2B SaaS Demo Platformu

> B2B yazılım demolarında bozulan verileri milisaniyeler içinde "İlk Günkü Satışa Hazır" durumuna getiren, yüksek performanslı Go & React uygulaması.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Vite](https://img.shields.io/badge/Vite-B73BFE?style=for-the-badge&logo=vite&logoColor=FFD62E)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)

---

## 🎯 The Problem & The Solution (Neden Yapıldı?)

### 🔴 Problem: Kirlenen Demo Verileri
B2B SaaS ürünlerinin satış toplantılarında (demolarda) karşılaşılan en büyük problem, **verilerin kirlenmesidir**. Satış temsilcileri ürünü müşteriye gösterirken yeni kayıtlar ekler, var olanları değiştirir veya siler. Ardından gelen bir sonraki toplantıda sistem "kullanılmış ve dağınık" bir halde bulunur. Her toplantı öncesi veritabanını manuel olarak temizlemek veya sıfırdan seed verisi yüklemek büyük bir zaman ve efor kaybıdır.

### 🟢 Çözüm: "Sihirli" Golden State Butonu
Golden State mimarisi, satış temsilcisine tek bir **"Sihirli Buton"** sunar. Toplantı bittiğinde veya yeni bir demo öncesinde Admin bu butona basarak tüm sistemi **saniyeler içinde ilk "Altın Durumuna" (Golden State) geri döndürür**. Müşteri verileri sıfırlanırken, sistemin ideal sunum hali korunur. Satış ekibi her demoya kusursuz, güvenilir ve tutarlı bir platform ile başlar.

---

## ⚙️ Under the Hood (Nasıl Çalışıyor?)

Sıfırlama işlemi sıradan bir tablonun silinip (`TRUNCATE`) baştan yaratılması veya dışarıdan yavaş bir SQL dosyasının (dump) içeri aktarılması **değildir**. Arka planda yüksek performanslı ve izole bir **Tenant Mimarisi** çalışır.

| Aşama | Teknik Detay |
| :--- | :--- |
| **İzolasyon** | Veritabanında aynı tabloda (`books`, `users`, `orders`) veriler `tenant_id` bazında izole edilir. Uygulamanın anlık kullandığı veriler `demo_active` tenant'ındadır. Asla değiştirilmeyen "Altın Şablon" verisi ise `demo_blueprint` tenant'ında yaşar. |
| **Concurrency (Eşzamanlılık)** | Go tarafında `sync.Mutex` ve `TryLock` kullanılarak aynı anda yalnızca bir sıfırlama işleminin çalışması garanti altına alınır. İşlem sırasında gelen diğer istekler anında reddedilir. |
| **Atomik İşlem (Transaction)** | ORM kullanılmamıştır. Go'nun `database/sql` paketi ile saf SQL kullanılarak işlemler başlatılır (`tx.Begin()`). FK sırasına göre: `orders → users → books` silinir, ardından blueprint'ten kopyalanır. |
| **INSERT ... SELECT** | `demo_active` kayıtları silindikten hemen sonra, veritabanı motoru seviyesinde `INSERT INTO ... SELECT ... FROM ... WHERE tenant_id = 'demo_blueprint'` operasyonu ile altın şablon klonlanır. |
| **Milisaniye Seviyesinde Hız** | İşlemlerin tümü PostgreSQL içinde gerçekleştiği için veri ağda gidip gelmez. Tek bir `COMMIT` ile kalıcı hale gelir. Sonuç: **milisaniyeler içinde atomik ve tutarlı bir reset işlemi.** |

---

## 🏗️ Architecture & Stack

Proje, katmanlı bir mimari yaklaşımıyla tasarlanmıştır. Backend tarafında ORM'siz, performansa odaklı bir yapı tercih edilirken, Frontend tarafında modern araçlar kullanılmıştır.

```mermaid
graph TD
    Client[Browser / React App]
    
    subgraph Frontend [React + Vite]
        UI[Components / Pages]
        State[AuthContext / CartContext]
        API_Layer[Axios API Services]
    end

    subgraph Backend [Go / Gin]
        Router[Gin Router]
        Middleware[JWT Auth & Role Middleware]
        Handlers[HTTP Handlers]
        Repo[Repository / Pure SQL]
    end

    subgraph Database [PostgreSQL]
        ActiveData[(tenant_id: demo_active)]
        BlueprintData[(tenant_id: demo_blueprint)]
    end

    Client <--> UI
    UI --> State
    State --> API_Layer
    API_Layer -- REST API --> Router
    Router --> Middleware
    Middleware --> Handlers
    Handlers --> Repo
    Repo -- SQL tx.Begin --> ActiveData
    Repo -- INSERT ... SELECT --> BlueprintData
```

### 🗂️ Katman Özeti
- **Frontend:** React, Vite (hızlı build), Vanilla CSS (styling), Context API (AuthContext + CartContext).
- **Backend:** Go, Gin (HTTP framework), `database/sql` (ORM'siz veri erişimi), JWT + bcrypt (kimlik doğrulama & şifreleme).
- **Veritabanı:** PostgreSQL (Docker üzerinde), çok kiracılı (`multi-tenant`) şema.
- **CI/CD:** GitHub Actions ile her `push`'ta otomatik Go testleri (Testcontainers ile gerçek PostgreSQL).

---

## 👥 Rol Tabanlı Yetki Sistemi

Uygulama üç farklı kullanıcı rolünü destekler:

| Rol | Yetkiler |
| :--- | :--- |
| **BUYER (Alıcı)** | Kitap kataloğunu görüntüleyebilir, sepete ekleyebilir ve sipariş verebilir. |
| **SELLER (Satıcı)** | Alıcı yetkilerine ek olarak kitap ekleyebilir, düzenleyebilir ve silebilir. Satış istatistiklerini görebilir. |
| **ADMIN (Yönetici)** | Tüm Satıcı yetkilerine ek olarak kullanıcı listesini yönetebilir ve **Golden State** sıfırlama işlemini gerçekleştirebilir. |

### 🔑 Demo Hesapları

Uygulamayı test etmek için aşağıdaki hazır hesapları kullanabilirsiniz:

| Rol | E-posta | Şifre |
| :--- | :--- | :--- |
| Admin | `admin@demo.com` | `admin123` |
| Satıcı | `seller@demo.com` | `seller123` |
| Alıcı | `buyer@demo.com` | `buyer123` |

> **Guest Checkout:** Giriş yapmadan da kitaplara göz atabilir ve sepete ekleyebilirsiniz. Ödeme adımında oturum açmanız istenir.

---

## 🔌 API Endpoints

| Method | Endpoint | Yetki | Açıklama |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Herkese Açık | Sunucu sağlık kontrolü |
| `POST` | `/api/v1/auth/register` | Herkese Açık | Yeni kullanıcı kaydı |
| `POST` | `/api/v1/auth/login` | Herkese Açık | Kullanıcı girişi, JWT döner |
| `GET` | `/api/v1/books` | Herkese Açık | Kitap kataloğunu listeler |
| `POST` | `/api/v1/books` | SELLER / ADMIN | Yeni kitap ekler |
| `PUT` | `/api/v1/books/:id` | SELLER / ADMIN | Kitap bilgilerini günceller |
| `DELETE` | `/api/v1/books/:id` | SELLER / ADMIN | Kitabı siler |
| `POST` | `/api/v1/checkout` | Giriş Gerekli | Sipariş tamamlar, stok düşer |
| `GET` | `/api/v1/sales` | Giriş Gerekli | Satış istatistiklerini getirir |
| `GET` | `/api/v1/admin/users` | ADMIN | Tüm kullanıcıları listeler |
| `POST` | `/api/v1/admin/system/restore` | ADMIN | Golden State sıfırlamasını tetikler |

---

## 🚀 Getting Started (Kurulum)

Projeyi kendi ortamınızda çalıştırmak için aşağıdaki adımları izleyin.

### Ön Koşullar
- [Docker & Docker Compose](https://www.docker.com/)
- [Go](https://go.dev/) (v1.22+)
- [Node.js](https://nodejs.org/) (v18+) & npm

### Adım 1: Altyapıyı Ayağa Kaldırın (Veritabanı)

```bash
# Proje dizinine gidin
cd BookMarketWGoldenState

# PostgreSQL veritabanını başlatın
docker-compose up -d
```

### Adım 2: Backend'i Çalıştırın (Go)

```bash
# Gerekli Go bağımlılıklarını indirin
go mod tidy

# API servisini başlatın
go run ./cmd/api/main.go
```
*Backend varsayılan olarak `http://localhost:8080` portunda çalışacaktır.*

### Adım 3: Frontend'i Çalıştırın (React)

Yeni bir terminal penceresi açın:

```bash
# Frontend dizinine gidin
cd frontend

# Bağımlılıkları yükleyin
npm install

# Geliştirme sunucusunu başlatın
npm run dev
```
*Frontend varsayılan olarak `http://localhost:5173` portunda çalışacaktır.*

---

## 🧪 Testleri Çalıştırma

```bash
# Tüm birim ve entegrasyon testlerini çalıştırır
# (Entegrasyon testleri için Docker gereklidir — Testcontainers otomatik PostgreSQL ayağa kaldırır)
go test ./...

# Sadece hızlı birim testleri (Docker gerektirmez)
go test -short ./...
```

Testler, her `git push`'ta **GitHub Actions** üzerinde de otomatik olarak çalışır.
