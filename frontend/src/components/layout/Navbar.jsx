import React, { useContext, useState } from 'react';
import { AuthContext } from '../../context/AuthContext';
import { CartContext } from '../../context/CartContext';
import { Button } from '../ui/Button';
import { AuthModal } from '../auth/AuthModal';
import { BookOpen, LogOut, ShoppingCart, X, CreditCard, Loader2 } from 'lucide-react';
import { cartService } from '../../services/api';

export const Navbar = ({ onOrderComplete }) => {
  const { user, logout } = useContext(AuthContext);
  const { cartItems, getCartTotal, removeFromCart, clearCart } = useContext(CartContext);
  const [isCartOpen, setIsCartOpen] = useState(false);
  const [checkoutLoading, setCheckoutLoading] = useState(false);
  const [checkoutSuccess, setCheckoutSuccess] = useState(false);
  const [showAuthModal, setShowAuthModal] = useState(false);

  const handleLogout = () => {
    logout();
    // Yönlendirme yok — misafir olarak kalmaya devam eder
  };

  // Gerçek checkout işlemi — token zaten var
  const performCheckout = async () => {
    if (cartItems.length === 0) return;
    setCheckoutLoading(true);
    try {
      const items = cartItems.map(item => ({ book_id: item.book_id, quantity: item.quantity }));
      await cartService.checkout(items);
      setCheckoutSuccess(true);
      clearCart();
      if (onOrderComplete) onOrderComplete();
      setTimeout(() => {
        setCheckoutSuccess(false);
        setIsCartOpen(false);
      }, 3000);
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Sipariş tamamlanırken hata oluştu.';
      if (errorMsg.includes('not found')) {
        alert('Sepetinizdeki bazı kitaplar sistemde bulunamadı (sistem sıfırlanmış olabilir). Sepetiniz temizlendi, lütfen yeniden ekleyin.');
        clearCart();
        setIsCartOpen(false);
      } else {
        alert(errorMsg);
      }
    } finally {
      setCheckoutLoading(false);
    }
  };

  // "Siparişi Tamamla" butonuna basıldığında:
  // Kullanıcı giriş yapmamışsa önce modal, giriş sonrası checkout
  const handleCheckoutClick = () => {
    if (!user) {
      setShowAuthModal(true);
    } else {
      performCheckout();
    }
  };

  // AuthModal'dan giriş/kayıt başarılı olunca çağrılır
  const handleAuthSuccess = () => {
    setShowAuthModal(false);
    // Token artık localStorage'da → checkout isteği gönder
    performCheckout();
  };

  return (
    <>
      <nav className="bg-white border-b border-slate-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0 flex items-center gap-2 text-indigo-600">
                <BookOpen className="h-8 w-8" />
                <span className="font-bold text-xl tracking-tight text-slate-900">Golden State</span>
              </div>
            </div>
            <div className="flex items-center gap-4">
              {/* Sepet ikonu: giriş yapmış BUYER veya misafir görebilir */}
              <button
                onClick={() => setIsCartOpen(true)}
                className="relative p-2 text-slate-500 hover:text-indigo-600 transition-colors"
              >
                <ShoppingCart className="h-6 w-6" />
                {cartItems.length > 0 && (
                  <span className="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/4 -translate-y-1/4 bg-rose-600 rounded-full">
                    {cartItems.reduce((acc, item) => acc + item.quantity, 0)}
                  </span>
                )}
              </button>

              {user ? (
                <>
                  <div className="hidden sm:flex flex-col items-end border-l border-slate-200 pl-4 ml-2">
                    <span className="text-sm font-medium text-slate-900">{user.email}</span>
                    <span className="text-xs text-slate-500 uppercase font-semibold tracking-wider">{user.role}</span>
                  </div>
                  <Button variant="ghost" size="sm" onClick={handleLogout} className="text-slate-500 hover:text-slate-700">
                    <LogOut className="h-5 w-5 mr-1" />
                    Çıkış Yap
                  </Button>
                </>
              ) : (
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setShowAuthModal(true)}
                  className="text-slate-600"
                >
                  Giriş Yap
                </Button>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Cart Drawer */}
      {isCartOpen && (
        <div className="fixed inset-0 z-50 overflow-hidden">
          <div className="absolute inset-0 bg-slate-900/50 backdrop-blur-sm transition-opacity" onClick={() => setIsCartOpen(false)} />
          <div className="fixed inset-y-0 right-0 pl-10 max-w-full flex">
            <div className="w-screen max-w-md transform transition-transform ease-in-out duration-500">
              <div className="h-full flex flex-col bg-white shadow-xl overflow-y-scroll">
                <div className="flex-1 py-6 overflow-y-auto px-4 sm:px-6">
                  <div className="flex items-start justify-between">
                    <h2 className="text-lg font-medium text-slate-900">Alışveriş Sepeti</h2>
                    <button onClick={() => setIsCartOpen(false)} className="text-slate-400 hover:text-slate-500">
                      <X className="h-6 w-6" />
                    </button>
                  </div>

                  <div className="mt-8">
                    {checkoutSuccess ? (
                      <div className="text-center py-12">
                        <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-emerald-100 mb-4">
                          <svg className="h-6 w-6 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                          </svg>
                        </div>
                        <h3 className="text-lg font-medium text-slate-900">Sipariş Başarılı!</h3>
                        <p className="mt-2 text-sm text-slate-500">Satın alma işleminiz tamamlandı. Stoklar güncellendi.</p>
                      </div>
                    ) : cartItems.length === 0 ? (
                      <p className="text-center text-slate-500 py-8">Sepetiniz şu an boş.</p>
                    ) : (
                      <div className="flow-root">
                        <ul className="-my-6 divide-y divide-slate-200">
                          {cartItems.map((item) => (
                            <li key={item.book_id} className="py-6 flex">
                              <div className="flex-1 flex flex-col">
                                <div>
                                  <div className="flex justify-between text-base font-medium text-slate-900">
                                    <h3>{item.title}</h3>
                                    <p className="ml-4">₺{(item.price * item.quantity).toLocaleString('tr-TR', { minimumFractionDigits: 2 })}</p>
                                  </div>
                                </div>
                                <div className="flex-1 flex items-end justify-between text-sm">
                                  <p className="text-slate-500">Adet: {item.quantity}</p>
                                  <button onClick={() => removeFromCart(item.book_id)} className="font-medium text-rose-600 hover:text-rose-500">
                                    Kaldır
                                  </button>
                                </div>
                              </div>
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                </div>

                {cartItems.length > 0 && !checkoutSuccess && (
                  <div className="border-t border-slate-200 py-6 px-4 sm:px-6 bg-slate-50">
                    <div className="flex justify-between text-base font-medium text-slate-900 mb-4">
                      <p>Ara Toplam</p>
                      <p>₺{getCartTotal().toLocaleString('tr-TR', { minimumFractionDigits: 2 })}</p>
                    </div>
                    <Button
                      variant="primary"
                      className="w-full h-12 text-lg flex items-center justify-center gap-2"
                      onClick={handleCheckoutClick}
                      disabled={checkoutLoading}
                    >
                      {checkoutLoading ? (
                        <>
                          <Loader2 className="h-5 w-5 animate-spin" />
                          İşleniyor...
                        </>
                      ) : (
                        <>
                          <CreditCard className="h-5 w-5" />
                          {user ? 'Siparişi Tamamla' : 'Giriş Yap'}
                        </>
                      )}
                    </Button>
                    {!user && (
                      <p className="text-xs text-center text-slate-400 mt-3">
                        Ödeme için hızlı giriş veya ücretsiz kayıt gereklidir.
                      </p>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Auth Modal — Lazy: sadece checkout tetiklendiğinde açılır */}
      {showAuthModal && (
        <AuthModal
          onSuccess={handleAuthSuccess}
          onClose={() => setShowAuthModal(false)}
        />
      )}
    </>
  );
};
