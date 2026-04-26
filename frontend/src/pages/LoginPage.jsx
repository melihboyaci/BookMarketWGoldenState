import React, { useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { Card, CardBody, CardHeader } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { BookOpen, Shield, ShoppingCart, User } from 'lucide-react';

export const LoginPage = () => {
  const { login, user } = useContext(AuthContext);
  const navigate = useNavigate();
  const [loadingRole, setLoadingRole] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    if (user) {
      navigate('/dashboard');
    }
  }, [user, navigate]);

  const handleDemoLogin = async (email, password, role) => {
    setLoadingRole(role);
    setError('');
    
    const result = await login(email, password);
    if (!result.success) {
      setError(result.error);
    }
    setLoadingRole(null);
  };

  const demoAccounts = [
    {
      role: 'ADMIN',
      email: 'admin@demo.com',
      password: 'admin123',
      icon: <Shield className="h-6 w-6 text-indigo-600" />,
      color: 'bg-indigo-50 border-indigo-200 hover:border-indigo-300',
      description: 'Tam yetki, sistemi sıfırlama (Golden State) erişimi.',
    },
    {
      role: 'SELLER',
      email: 'seller@demo.com',
      password: 'seller123',
      icon: <BookOpen className="h-6 w-6 text-emerald-600" />,
      color: 'bg-emerald-50 border-emerald-200 hover:border-emerald-300',
      description: 'Kendi kitaplarını listeleme ve yönetme yetkisi.',
    },
    {
      role: 'BUYER',
      email: 'buyer@demo.com',
      password: 'buyer123',
      icon: <ShoppingCart className="h-6 w-6 text-amber-600" />,
      color: 'bg-amber-50 border-amber-200 hover:border-amber-300',
      description: 'Kitapları görüntüleme ve satın alma yetkisi.',
    }
  ];

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 z-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-[30%] -right-[10%] w-[70%] h-[70%] rounded-full bg-gradient-to-br from-indigo-100 to-purple-50 opacity-50 blur-3xl" />
        <div className="absolute -bottom-[20%] -left-[10%] w-[60%] h-[60%] rounded-full bg-gradient-to-tr from-emerald-50 to-teal-50 opacity-50 blur-3xl" />
      </div>

      <div className="sm:mx-auto sm:w-full sm:max-w-md relative z-10 text-center mb-8">
        <div className="mx-auto w-16 h-16 bg-white rounded-2xl shadow-sm border border-slate-100 flex items-center justify-center mb-4">
          <BookOpen className="h-8 w-8 text-indigo-600" />
        </div>
        <h2 className="text-3xl font-extrabold text-slate-900 tracking-tight">Golden State</h2>
        <p className="mt-2 text-sm text-slate-600">
          Kitap Satış Demo Uygulamasına Hoş Geldiniz
        </p>
      </div>

      <div className="sm:mx-auto sm:w-full sm:max-w-xl relative z-10">
        <Card className="shadow-xl shadow-slate-200/50 border-white/50 backdrop-blur-sm bg-white/90">
          <CardHeader className="text-center pb-2 border-none">
            <h3 className="text-lg font-medium leading-6 text-slate-900">
              Demo Hesabı Seçin
            </h3>
            <p className="mt-1 text-sm text-slate-500">
              Uygulamayı test etmek için aşağıdaki rollerden biriyle giriş yapın.
            </p>
          </CardHeader>
          <CardBody className="pt-4">
            {error && (
              <div className="mb-4 bg-rose-50 border border-rose-200 text-rose-700 px-4 py-3 rounded-lg text-sm">
                {error}
              </div>
            )}
            <div className="space-y-4">
              {demoAccounts.map((account) => (
                <div 
                  key={account.role}
                  className={`relative rounded-xl border p-5 flex items-center justify-between cursor-pointer transition-all duration-200 group ${account.color}`}
                  onClick={() => handleDemoLogin(account.email, account.password, account.role)}
                >
                  <div className="flex items-center gap-4">
                    <div className="p-3 bg-white rounded-lg shadow-sm">
                      {account.icon}
                    </div>
                    <div>
                      <h4 className="text-base font-semibold text-slate-900 group-hover:text-slate-800">
                        {account.role} Olarak Gir
                      </h4>
                      <p className="text-sm text-slate-600 mt-0.5">
                        {account.description}
                      </p>
                    </div>
                  </div>
                  <Button 
                    variant="primary" 
                    size="sm" 
                    className="ml-4 opacity-0 group-hover:opacity-100 transition-opacity"
                    isLoading={loadingRole === account.role}
                  >
                    Giriş
                  </Button>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
};
