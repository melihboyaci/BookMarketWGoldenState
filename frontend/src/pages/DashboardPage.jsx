import React, { useState, useEffect, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { CartContext } from '../context/CartContext';
import { bookService, systemService } from '../services/api';
import { Card, CardBody, CardFooter } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { Navbar } from '../components/layout/Navbar';
import { SalesStats } from '../components/dashboard/SalesStats';
import { UsersPanel } from '../components/dashboard/UsersPanel';
import { Book, RefreshCw, ShoppingCart, Loader2, Plus, Edit, Trash2, X, RotateCcw } from 'lucide-react';

export const DashboardPage = () => {
  const { user } = useContext(AuthContext);
  const { addToCart } = useContext(CartContext);
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [restoring, setRestoring] = useState(false);
  const [statsRefreshKey, setStatsRefreshKey] = useState(0);

  // Sipariş tamamlandığında grafiği güncelle (Navbar'dan çağrılır)
  const handleOrderComplete = () => {
    setStatsRefreshKey(prev => prev + 1);
    fetchBooks(); // Stok güncellemesi için kitapları da yenile
  };

  // Seller Modals State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingBook, setEditingBook] = useState(null);
  const [formData, setFormData] = useState({ title: '', author: '', isbn: '', image_url: '', price: '', stock: '' });

  useEffect(() => {
    fetchBooks();
  }, []);

  const fetchBooks = async () => {
    setLoading(true);
    try {
      const data = await bookService.getBooks();
      setBooks(data || []);
    } catch (error) {
      console.error('Kitaplar çekilemedi', error);
      setBooks([]);
    } finally {
      setLoading(false);
    }
  };

  const handleRestore = async () => {
    if (!window.confirm("Sistemi ilk Altın Durumuna (Golden State) döndürmek istediğinize emin misiniz? Yapılan tüm siparişler ve eklenen kitaplar silinecek.")) return;
    
    setRestoring(true);
    try {
      await systemService.restoreGoldenState();
      // Grafik ve kitapları resetle — sayfa yenilemeye gerek yok
      setStatsRefreshKey(prev => prev + 1);
      await fetchBooks();
      alert("Sistem başarıyla sıfırlandı!");
    } catch (error) {
      alert("Sıfırlama başarısız oldu.");
    } finally {
      setRestoring(false);
    }
  };

  const openAddModal = () => {
    setEditingBook(null);
    setFormData({ title: '', author: '', isbn: '', image_url: '', price: '', stock: '' });
    setIsModalOpen(true);
  };

  const openEditModal = (book) => {
    setEditingBook(book);
    setFormData({ title: book.title, author: book.author, isbn: book.isbn, image_url: book.image_url || '', price: book.price, stock: book.stock });
    setIsModalOpen(true);
  };

  const handleSaveBook = async (e) => {
    e.preventDefault();
    try {
      const data = {
        title: formData.title,
        author: formData.author,
        isbn: formData.isbn,
        image_url: formData.image_url,
        price: parseFloat(formData.price),
        stock: parseInt(formData.stock, 10)
      };

      if (editingBook) {
        await bookService.updateBook(editingBook.id, data);
      } else {
        await bookService.createBook(data);
      }
      setIsModalOpen(false);
      fetchBooks();
    } catch (error) {
      alert("Kitap kaydedilirken hata oluştu.");
    }
  };

  const handleDeleteBook = async (id) => {
    if (window.confirm("Bu kitabı silmek istediğinize emin misiniz?")) {
      try {
        await bookService.deleteBook(id);
        fetchBooks();
      } catch (error) {
        alert("Kitap silinirken hata oluştu.");
      }
    }
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <Navbar onOrderComplete={handleOrderComplete} />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        
        {/* Admin Golden State Header */}
        {user?.role === 'ADMIN' && (
          <div className="mb-8 p-6 bg-gradient-to-r from-slate-900 to-slate-800 rounded-2xl shadow-lg flex flex-col sm:flex-row justify-between items-center gap-4 border border-slate-700">
            <div>
              <h2 className="text-xl font-bold text-white flex items-center gap-2">
                <span className="bg-rose-500 w-3 h-3 rounded-full animate-pulse"></span>
                Sistem Yönetimi (Golden State)
              </h2>
              <p className="text-slate-400 mt-1 text-sm">Demoyu ilk haline getirmek için sistemi sıfırlayın.</p>
            </div>
            <Button 
              variant="primary" 
              className="bg-rose-600 hover:bg-rose-700 text-white border-none shadow-rose-900/50" 
              onClick={handleRestore}
              disabled={restoring}
            >
              {restoring ? (
                <><Loader2 className="h-5 w-5 mr-2 animate-spin" /> Sıfırlanıyor...</>
              ) : (
                <><RotateCcw className="h-5 w-5 mr-2" /> Sistemi Sıfırla</>
              )}
            </Button>
          </div>
        )}

        {/* Stats: Seller ve Admin için */}
        {(user?.role === 'SELLER' || user?.role === 'ADMIN') && (
          <SalesStats refreshKey={statsRefreshKey} />
        )}

        {/* Kullanıcı paneli: yalnızca Admin */}
        {user?.role === 'ADMIN' && (
          <UsersPanel refreshKey={statsRefreshKey} />
        )}

        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Kitap Kataloğu</h1>
            <p className="text-sm text-slate-500 mt-1">Sistemdeki tüm kitapları inceleyin.</p>
          </div>
          <div className="flex gap-2">
            {(user?.role === 'SELLER' || user?.role === 'ADMIN') && (
              <Button variant="primary" onClick={openAddModal} className="gap-2">
                <Plus className="h-4 w-4" />
                Yeni Kitap
              </Button>
            )}
            <Button variant="secondary" onClick={fetchBooks} disabled={loading} className="gap-2">
              <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
              Yenile
            </Button>
          </div>
        </div>

        {loading ? (
          <div className="flex flex-col items-center justify-center h-64 text-slate-400">
            <Loader2 className="h-10 w-10 animate-spin mb-4 text-indigo-500" />
            <p>Kitaplar yükleniyor...</p>
          </div>
        ) : books.length === 0 ? (
          <div className="text-center py-20 text-slate-500">
            <Book className="h-12 w-12 mx-auto text-slate-300 mb-4" />
            <p>Sistemde hiç kitap bulunmuyor.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {books.map((book) => (
              <Card key={book.id} className="flex flex-col h-full group hover:-translate-y-1 transition-transform duration-300 overflow-hidden">
                <div className="h-48 bg-slate-100 relative overflow-hidden flex items-center justify-center">
                  {book.image_url ? (
                    <img 
                      src={book.image_url} 
                      alt={book.title} 
                      className="w-full h-full object-contain p-2 group-hover:scale-105 transition-transform duration-500"
                    />
                  ) : (
                    <Book className="h-16 w-16 text-slate-300" />
                  )}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
                </div>
                <CardBody className="flex-1 flex flex-col">
                  <h3 className="text-lg font-bold text-slate-900 line-clamp-1">{book.title}</h3>
                  <p className="text-sm text-slate-500 mb-4 line-clamp-1">{book.author}</p>
                  
                  <div className="mt-auto flex justify-between items-end">
                    <div className="text-2xl font-black text-indigo-600">
                      ₺{book.price?.toLocaleString('tr-TR', { minimumFractionDigits: 2 })}
                    </div>
                    <div className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-1 rounded-md">
                      Stok: {book.stock}
                    </div>
                  </div>
                </CardBody>
                
                {/* BUYER veya misafir → Sepete Ekle */}
                {(!user || user?.role === 'BUYER') && (
                  <CardFooter className="pt-0 pb-4 px-4 bg-white border-none">
                    <Button 
                      variant="primary" 
                      className="w-full gap-2 group-hover:bg-indigo-700 shadow-lg shadow-indigo-200"
                      onClick={() => addToCart(book)}
                      disabled={book.stock <= 0}
                    >
                      <ShoppingCart className="h-4 w-4" />
                      {book.stock <= 0 ? "Stokta Yok" : "Sepete Ekle"}
                    </Button>
                  </CardFooter>
                )}
                
                {(user?.role === 'SELLER' || user?.role === 'ADMIN') && (
                  <CardFooter className="pt-0 pb-4 px-4 bg-white border-none flex gap-2">
                    <Button variant="secondary" className="flex-1 gap-1 px-2" onClick={() => openEditModal(book)}>
                      <Edit className="h-4 w-4" />
                      Düzenle
                    </Button>
                    <Button variant="secondary" className="flex-1 gap-1 px-2 text-rose-600 hover:text-rose-700 hover:bg-rose-50" onClick={() => handleDeleteBook(book.id)}>
                      <Trash2 className="h-4 w-4" />
                      Sil
                    </Button>
                  </CardFooter>
                )}
              </Card>
            ))}
          </div>
        )}
      </main>

      {/* Add/Edit Book Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-slate-900/50 backdrop-blur-sm" onClick={() => setIsModalOpen(false)} />
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-200 flex justify-between items-center bg-slate-50">
              <h2 className="text-lg font-bold text-slate-900">
                {editingBook ? "Kitabı Düzenle" : "Yeni Kitap Ekle"}
              </h2>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-slate-600">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleSaveBook} className="p-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Kitap Adı</label>
                  <input type="text" required value={formData.title} onChange={e => setFormData({...formData, title: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Yazar</label>
                  <input type="text" required value={formData.author} onChange={e => setFormData({...formData, author: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Görsel URL</label>
                  <input type="url" value={formData.image_url} placeholder="https://..." onChange={e => setFormData({...formData, image_url: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">ISBN</label>
                  <input type="text" value={formData.isbn} onChange={e => setFormData({...formData, isbn: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Fiyat (₺)</label>
                    <input type="number" step="0.01" min="0" required value={formData.price} onChange={e => setFormData({...formData, price: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Stok</label>
                    <input type="number" min="0" required value={formData.stock} onChange={e => setFormData({...formData, stock: e.target.value})} className="w-full px-3 py-2 border border-slate-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500" />
                  </div>
                </div>
              </div>
              <div className="mt-8 flex justify-end gap-3">
                <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>İptal</Button>
                <Button type="submit" variant="primary">Kaydet</Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
