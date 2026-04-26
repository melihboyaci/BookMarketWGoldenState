import React from 'react';
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';

// Ocak – Nisan için default Golden State değerleri (₺ cinsinden)
const BASE_MONTHLY_DATA = [
  { ay: 'Oca', gelir: 12500 },
  { ay: 'Şub', gelir: 18200 },
  { ay: 'Mar', gelir: 15800 },
  { ay: 'Nis', gelir: 24000 },
];

const CustomTooltip = ({ active, payload, label }) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-white border border-slate-200 rounded-xl shadow-xl px-4 py-3">
        <p className="text-sm font-semibold text-slate-700 mb-1">{label}</p>
        <p className="text-base font-bold text-indigo-600">
          ₺{payload[0].value.toLocaleString('tr-TR', { minimumFractionDigits: 2 })}
        </p>
      </div>
    );
  }
  return null;
};

export const RevenueChart = ({ currentMonthExtra = 0 }) => {
  // Nisan'ın değerine gerçek sipariş gelirini (₺) ekle
  const data = BASE_MONTHLY_DATA.map((d, i) => {
    if (i === BASE_MONTHLY_DATA.length - 1) {
      return { ...d, gelir: d.gelir + currentMonthExtra };
    }
    return d;
  });

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 mb-8">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-lg font-bold text-slate-900">Aylık Gelir</h2>
          <p className="text-sm text-slate-500 mt-0.5">Ocak – Nisan 2026 satış performansı</p>
        </div>
        <div className="flex items-center gap-2 bg-indigo-50 text-indigo-700 text-xs font-semibold px-3 py-1.5 rounded-full">
          <span className="w-2 h-2 rounded-full bg-indigo-500 inline-block"></span>
          Aylık Gelir (₺)
        </div>
      </div>

      <ResponsiveContainer width="100%" height={240}>
        <AreaChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="gelirGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#6366f1" stopOpacity={0.18} />
              <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" vertical={false} />
          <XAxis
            dataKey="ay"
            tick={{ fill: '#94a3b8', fontSize: 12 }}
            axisLine={false}
            tickLine={false}
          />
          <YAxis
            tick={{ fill: '#94a3b8', fontSize: 12 }}
            axisLine={false}
            tickLine={false}
            tickFormatter={(v) => `₺${(v / 1000).toFixed(0)}k`}
            width={50}
          />
          <Tooltip content={<CustomTooltip />} />
          <Area
            type="monotone"
            dataKey="gelir"
            stroke="#6366f1"
            strokeWidth={2.5}
            fill="url(#gelirGradient)"
            dot={false}
            activeDot={{ r: 5, fill: '#6366f1', stroke: '#fff', strokeWidth: 2 }}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
