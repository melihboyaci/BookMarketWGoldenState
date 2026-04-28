import { Routes, Route, Navigate } from 'react-router-dom'
import { DashboardPage } from './pages/DashboardPage'

// Dashboard artık herkese açık — Lazy Auth akışı
function App() {
  return (
    <Routes>
      {/* Ana sayfa: misafir de görebilir */}
      <Route path="/" element={<DashboardPage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
      {/* Bilinmeyen yolları ana sayfaya yönlendir */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
