import React, { useState, useEffect } from 'react';
import { systemService } from '../../services/api';
import { Card, CardBody } from '../ui/Card';
import { RevenueChart } from './RevenueChart';
import { TrendingUp, Package, Loader2 } from 'lucide-react';

export const SalesStats = ({ refreshKey = 0 }) => {
  const [stats, setStats] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStats();
  }, [refreshKey]); // refreshKey değişince (yeni sipariş veya reset) yeniden çek

  const fetchStats = async () => {
    setLoading(true);
    try {
      const data = await systemService.getSales();
      setStats(data || []);
    } catch (error) {
      console.error("Satış verileri alınamadı", error);
      setStats([]);
    } finally {
      setLoading(false);
    }
  };

  const totalRevenue = stats.reduce((acc, stat) => acc + (stat.total_sales || 0), 0);
  const totalOrders = stats.reduce((acc, stat) => acc + (stat.total_orders || 0), 0);

  if (loading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
      </div>
    );
  }

  return (
    <div className="mb-8">
      {/* KPI Kartları */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <Card className="bg-gradient-to-br from-indigo-500 to-indigo-600 text-white border-none shadow-lg shadow-indigo-200">
          <CardBody className="flex items-center gap-4 p-6">
            <div className="bg-white/20 p-3 rounded-xl">
              <TrendingUp className="h-8 w-8 text-white" />
            </div>
            <div>
              <p className="text-indigo-100 font-medium">Toplam Gelir</p>
              <h3 className="text-3xl font-bold">
                ₺{totalRevenue.toLocaleString('tr-TR', { minimumFractionDigits: 2 })}
              </h3>
              {totalRevenue === 0 && (
                <p className="text-indigo-200 text-xs mt-0.5">Henüz satış yok</p>
              )}
            </div>
          </CardBody>
        </Card>
        <Card className="bg-gradient-to-br from-emerald-500 to-emerald-600 text-white border-none shadow-lg shadow-emerald-200">
          <CardBody className="flex items-center gap-4 p-6">
            <div className="bg-white/20 p-3 rounded-xl">
              <Package className="h-8 w-8 text-white" />
            </div>
            <div>
              <p className="text-emerald-100 font-medium">Toplam Sipariş</p>
              <h3 className="text-3xl font-bold">{totalOrders}</h3>
              {totalOrders === 0 && (
                <p className="text-emerald-200 text-xs mt-0.5">Henüz sipariş yok</p>
              )}
            </div>
          </CardBody>
        </Card>
      </div>

      {/* Aylık Gelir Grafiği — Nisan'ın değeri = default + gerçek sipariş geliri */}
      <RevenueChart currentMonthExtra={totalRevenue} />
    </div>
  );
};
