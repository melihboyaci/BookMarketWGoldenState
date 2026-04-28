import React, { useState, useEffect } from 'react';
import { adminService } from '../../services/api';
import { Users, Loader2, ShieldCheck, Store, ShoppingBag, RefreshCw } from 'lucide-react';

const ROLE_CONFIG = {
  ADMIN:  { label: 'Admin',  icon: ShieldCheck, cls: 'bg-rose-100 text-rose-700' },
  SELLER: { label: 'Satıcı', icon: Store,        cls: 'bg-indigo-100 text-indigo-700' },
  BUYER:  { label: 'Alıcı',  icon: ShoppingBag,  cls: 'bg-emerald-100 text-emerald-700' },
};

const formatDate = (iso) =>
  new Date(iso).toLocaleString('tr-TR', { dateStyle: 'medium', timeStyle: 'short' });

export const UsersPanel = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await adminService.getUsers();
      setUsers(data || []);
    } catch (err) {
      setError('Kullanıcılar yüklenirken hata oluştu.');
    } finally {
      setLoading(false);
    }
  };

  const counts = users.reduce((acc, u) => {
    acc[u.role] = (acc[u.role] || 0) + 1;
    return acc;
  }, {});

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm mb-8">
      {/* Başlık */}
      <div className="flex items-center justify-between px-6 py-4 border-b border-slate-100">
        <div className="flex items-center gap-3">
          <div className="bg-slate-100 p-2 rounded-xl">
            <Users className="h-5 w-5 text-slate-600" />
          </div>
          <div>
            <h2 className="text-base font-bold text-slate-900">Kullanıcı Yönetimi</h2>
            <p className="text-xs text-slate-500">demo_active kiracısındaki tüm kullanıcılar</p>
          </div>
        </div>
        <button
          onClick={fetchUsers}
          className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
          title="Yenile"
        >
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
        </button>
      </div>

      {/* Özet sayaçlar */}
      <div className="grid grid-cols-3 divide-x divide-slate-100 border-b border-slate-100">
        {Object.entries(ROLE_CONFIG).map(([role, cfg]) => (
          <div key={role} className="px-6 py-3 flex items-center gap-3">
            <div className={`p-1.5 rounded-lg ${cfg.cls}`}>
              <cfg.icon className="h-4 w-4" />
            </div>
            <div>
              <p className="text-xl font-bold text-slate-900">{counts[role] ?? 0}</p>
              <p className="text-xs text-slate-500">{cfg.label}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Tablo */}
      {loading ? (
        <div className="flex justify-center py-10">
          <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
        </div>
      ) : error ? (
        <p className="text-center text-rose-500 py-8 text-sm">{error}</p>
      ) : users.length === 0 ? (
        <p className="text-center text-slate-400 py-8 text-sm">Kayıtlı kullanıcı bulunamadı.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full">
            <thead>
              <tr className="bg-slate-50 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                <th className="px-6 py-3 text-left">Kullanıcı</th>
                <th className="px-6 py-3 text-left">E-posta</th>
                <th className="px-6 py-3 text-left">Rol</th>
                <th className="px-6 py-3 text-left">Kayıt Tarihi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {users.map((user) => {
                const cfg = ROLE_CONFIG[user.role] ?? { label: user.role, cls: 'bg-slate-100 text-slate-600' };
                return (
                  <tr key={user.id} className="hover:bg-slate-50 transition-colors">
                    <td className="px-6 py-3">
                      <div className="flex items-center gap-3">
                        <div className="h-8 w-8 rounded-full bg-gradient-to-br from-indigo-400 to-indigo-600 flex items-center justify-center text-white text-sm font-bold">
                          {user.username?.charAt(0).toUpperCase()}
                        </div>
                        <span className="text-sm font-medium text-slate-800">{user.username}</span>
                      </div>
                    </td>
                    <td className="px-6 py-3 text-sm text-slate-600">{user.email}</td>
                    <td className="px-6 py-3">
                      <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${cfg.cls}`}>
                        <cfg.icon className="h-3 w-3" />
                        {cfg.label}
                      </span>
                    </td>
                    <td className="px-6 py-3 text-sm text-slate-500">{formatDate(user.created_at)}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};
