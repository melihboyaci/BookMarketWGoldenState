import React, { useState, useContext } from 'react';
import { AuthContext } from '../../context/AuthContext';
import { Button } from '../ui/Button';
import { X, Mail, Lock, User, Loader2, BookOpen, Zap } from 'lucide-react';

const DEMO_ACCOUNTS = [
  { label: 'Alıcı',  role: 'BUYER',  email: 'buyer@demo.com',  password: 'buyer123',  color: 'bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100' },
  { label: 'Satıcı', role: 'SELLER', email: 'seller@demo.com', password: 'seller123', color: 'bg-indigo-50 text-indigo-700 border-indigo-200 hover:bg-indigo-100' },
  { label: 'Admin',  role: 'ADMIN',  email: 'admin@demo.com',  password: 'admin123',  color: 'bg-rose-50 text-rose-700 border-rose-200 hover:bg-rose-100' },
];

export const AuthModal = ({ onSuccess, onClose }) => {
  const { login, register, loading } = useContext(AuthContext);
  const [mode, setMode] = useState('login'); // 'login' | 'register'
  const [formData, setFormData] = useState({ username: '', email: '', password: '' });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    let result;
    if (mode === 'login') {
      result = await login(formData.email, formData.password);
    } else {
      result = await register(formData.username, formData.email, formData.password);
    }

    if (result.success) {
      onSuccess(result.role);
    } else {
      setError(result.error);
    }
  };

  // Demo hesabını forma otomatik doldur ve hemen giriş yap
  const handleDemoLogin = async (account) => {
    setError('');
    const result = await login(account.email, account.password);
    if (result.success) {
      onSuccess(result.role);
    } else {
      setError(result.error);
    }
  };

  const switchMode = () => {
    setMode(prev => prev === 'login' ? 'register' : 'login');
    setError('');
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" onClick={onClose} />

      {/* Modal */}
      <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden z-10">

        {/* Header */}
        <div className="bg-gradient-to-r from-indigo-600 to-indigo-700 px-8 py-6 text-white">
          <button onClick={onClose} className="absolute top-4 right-4 text-indigo-200 hover:text-white transition-colors">
            <X className="h-5 w-5" />
          </button>
          <div className="flex items-center gap-3 mb-1">
            <BookOpen className="h-7 w-7" />
            <span className="font-bold text-xl">Golden State</span>
          </div>
          <p className="text-indigo-200 text-sm">
            {mode === 'login'
              ? 'Siparişinizi tamamlamak için giriş yapın.'
              : 'Ücretsiz hesap oluşturun ve alışverişe başlayın.'}
          </p>
        </div>

        {/* Demo Accounts — sadece login modunda */}
        {mode === 'login' && (
          <div className="px-8 pt-5 pb-1">
            <div className="flex items-center gap-2 mb-3">
              <Zap className="h-4 w-4 text-amber-500" />
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Demo Hesapları</span>
            </div>
            <div className="grid grid-cols-3 gap-2">
              {DEMO_ACCOUNTS.map(account => (
                <button
                  key={account.role}
                  type="button"
                  onClick={() => handleDemoLogin(account)}
                  disabled={loading}
                  className={`flex flex-col items-center gap-1 px-2 py-2.5 rounded-xl border text-xs font-semibold transition-all ${account.color} disabled:opacity-50`}
                >
                  <span className="text-base">
                    {account.role === 'BUYER' ? '🛒' : account.role === 'SELLER' ? '🏪' : '👑'}
                  </span>
                  <span>{account.label}</span>
                </button>
              ))}
            </div>
            <div className="relative mt-5 mb-1">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-slate-200" />
              </div>
              <div className="relative flex justify-center">
                <span className="bg-white px-3 text-xs text-slate-400">ya da kendin giriş yap</span>
              </div>
            </div>
          </div>
        )}

        {/* Tab Switcher */}
        <div className="flex border-b border-slate-200">
          <button
            className={`flex-1 py-3 text-sm font-semibold transition-colors ${
              mode === 'login' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-slate-500 hover:text-slate-700'
            }`}
            onClick={() => { setMode('login'); setError(''); }}
          >
            Giriş Yap
          </button>
          <button
            className={`flex-1 py-3 text-sm font-semibold transition-colors ${
              mode === 'register' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-slate-500 hover:text-slate-700'
            }`}
            onClick={() => { setMode('register'); setError(''); }}
          >
            Kayıt Ol
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="px-8 py-5 space-y-4">
          {mode === 'register' && (
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Kullanıcı Adı</label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                <input
                  type="text"
                  required
                  minLength={2}
                  value={formData.username}
                  onChange={e => setFormData({ ...formData, username: e.target.value })}
                  placeholder="adınız"
                  className="w-full pl-10 pr-4 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
              </div>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">E-posta</label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
              <input
                type="email"
                required
                value={formData.email}
                onChange={e => setFormData({ ...formData, email: e.target.value })}
                placeholder="ornek@email.com"
                className="w-full pl-10 pr-4 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Şifre</label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
              <input
                type="password"
                required
                minLength={6}
                value={formData.password}
                onChange={e => setFormData({ ...formData, password: e.target.value })}
                placeholder={mode === 'register' ? 'En az 6 karakter' : '••••••••'}
                className="w-full pl-10 pr-4 py-2.5 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
            </div>
          </div>

          {error && (
            <div className="bg-rose-50 border border-rose-200 text-rose-700 text-sm rounded-lg px-4 py-3">
              {error}
            </div>
          )}

          <Button
            type="submit"
            variant="primary"
            className="w-full h-11 text-base flex items-center justify-center gap-2"
            disabled={loading}
          >
            {loading ? (
              <><Loader2 className="h-5 w-5 animate-spin" /> İşleniyor...</>
            ) : mode === 'login' ? (
              'Giriş Yap'
            ) : (
              'Kayıt Ol'
            )}
          </Button>

          <p className="text-center text-sm text-slate-500 pb-1">
            {mode === 'login' ? 'Hesabın yok mu? ' : 'Zaten üye misin? '}
            <button type="button" onClick={switchMode} className="text-indigo-600 hover:underline font-medium">
              {mode === 'login' ? 'Kayıt Ol' : 'Giriş Yap'}
            </button>
          </p>
        </form>
      </div>
    </div>
  );
};
