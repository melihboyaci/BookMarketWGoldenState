import React, { useState, useEffect, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { bookService } from '../services/api';
import { Card, CardBody, CardFooter, CardHeader } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { Navbar } from '../components/layout/Navbar';
import { Book, RefreshCw, ShoppingCart, Loader2 } from 'lucide-react';

// Mock data as fallback
const MOCK_BOOKS = [
  { id: 1, title: 'Clean Code', author: 'Robert C. Martin', price: 45.99, stock: 12, cover: 'bg-indigo-100' },
  { id: 2, title: 'Pragmatic Programmer', author: 'Andrew Hunt', price: 55.00, stock: 8, cover: 'bg-emerald-100' },
  { id: 3, title: 'Design Patterns', author: 'Erich Gamma', price: 60.50, stock: 5, cover: 'bg-rose-100' },
  { id: 4, title: 'Refactoring', author: 'Martin Fowler', price: 49.99, stock: 15, cover: 'bg-amber-100' },
];

export const DashboardPage = () => {
  const { user } = useContext(AuthContext);
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBooks();
  }, []);

  const fetchBooks = async () => {
    setLoading(true);
    try {
      const data = await bookService.getBooks();
      if (data && data.length > 0) {
        setBooks(data);
      } else {
        // Fallback if backend not implemented yet
        setBooks(MOCK_BOOKS);
      }
    } catch (error) {
      console.error('Kitaplar çekilemedi', error);
      setBooks(MOCK_BOOKS);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <Navbar />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Kitap Kataloğu</h1>
            <p className="text-sm text-slate-500 mt-1">Sistemdeki tüm kitapları inceleyin.</p>
          </div>
          <Button variant="secondary" onClick={fetchBooks} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Yenile
          </Button>
        </div>

        {loading ? (
          <div className="flex flex-col items-center justify-center h-64 text-slate-400">
            <Loader2 className="h-10 w-10 animate-spin mb-4 text-indigo-500" />
            <p>Kitaplar yükleniyor...</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {books.map((book) => (
              <Card key={book.id} className="flex flex-col h-full group hover:-translate-y-1 transition-transform duration-300">
                <div className={`h-48 ${book.cover || 'bg-slate-200'} relative flex items-center justify-center p-6`}>
                  <Book className="h-20 w-20 text-slate-900/10 drop-shadow-sm" />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />
                </div>
                <CardBody className="flex-1 flex flex-col">
                  <h3 className="text-lg font-bold text-slate-900 line-clamp-1">{book.title}</h3>
                  <p className="text-sm text-slate-500 mb-4 line-clamp-1">{book.author}</p>
                  
                  <div className="mt-auto flex justify-between items-end">
                    <div className="text-2xl font-black text-indigo-600">
                      ${book.price?.toFixed(2)}
                    </div>
                    <div className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-1 rounded-md">
                      Stok: {book.stock}
                    </div>
                  </div>
                </CardBody>
                
                {user?.role === 'BUYER' && (
                  <CardFooter className="pt-0 pb-4 px-4 bg-white border-none">
                    <Button variant="primary" className="w-full gap-2 group-hover:bg-indigo-700">
                      <ShoppingCart className="h-4 w-4" />
                      Sepete Ekle
                    </Button>
                  </CardFooter>
                )}
                
                {user?.role === 'SELLER' && (
                  <CardFooter className="pt-0 pb-4 px-4 bg-white border-none">
                    <Button variant="secondary" className="w-full gap-2">
                      Düzenle
                    </Button>
                  </CardFooter>
                )}
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  );
};
