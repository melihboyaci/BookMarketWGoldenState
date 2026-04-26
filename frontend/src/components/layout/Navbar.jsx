import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthContext } from '../../context/AuthContext';
import { Button } from '../ui/Button';
import { BookOpen, LogOut } from 'lucide-react';

export const Navbar = () => {
  const { user, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-white border-b border-slate-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <div className="flex-shrink-0 flex items-center gap-2 text-indigo-600">
              <BookOpen className="h-8 w-8" />
              <span className="font-bold text-xl tracking-tight text-slate-900">Golden State</span>
            </div>
          </div>
          <div className="flex items-center gap-4">
            {user ? (
              <>
                <div className="flex flex-col items-end hidden sm:flex">
                  <span className="text-sm font-medium text-slate-900">{user.email}</span>
                  <span className="text-xs text-slate-500 uppercase font-semibold tracking-wider">{user.role}</span>
                </div>
                <Button variant="ghost" size="sm" onClick={handleLogout} className="text-slate-500 hover:text-slate-700">
                  <LogOut className="h-5 w-5 mr-1" />
                  Çıkış Yap
                </Button>
              </>
            ) : null}
          </div>
        </div>
      </div>
    </nav>
  );
};
