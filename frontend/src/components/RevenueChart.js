import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import './RevenueChart.css';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function RevenueChart({ data, viewMode = 'revenue' }) {
  if (!data || data.length === 0) {
    return <div className="chart-placeholder">Nenhum dado disponível para o gráfico</div>;
  }

  // Preparar dados para o gráfico
  const labels = data.map((item) => item.date);
  
  // Determinar quais dados usar baseado no modo de visualização
  const isRevenueMode = viewMode === 'revenue';
  const approvedData = data.map((item) => 
    isRevenueMode ? item.approved_revenue : item.approved_orders
  );
  const pendingData = data.map((item) => 
    isRevenueMode ? item.pending_revenue : item.pending_orders
  );
  const cancelledData = data.map((item) => 
    isRevenueMode ? item.cancelled_revenue : item.cancelled_orders
  );

  const chartData = { // cria o objeto chartData que contém os dados do gráfico
    labels,
    datasets: [
      {
        label: isRevenueMode ? 'Receita Aprovada' : 'Pedidos Aprovados',
        data: approvedData,
        borderColor: '#00D47E',
        backgroundColor: 'rgba(0, 212, 126, 0.1)',
        tension: 0.4,
      },
      {
        label: isRevenueMode ? 'Receita Pendente' : 'Pedidos Pendentes',
        data: pendingData,
        borderColor: '#ffd43b',
        backgroundColor: 'rgba(255, 212, 59, 0.1)',
        tension: 0.4,
        ...(!isRevenueMode && {
          borderDash: [5, 5], // linha tracejada apenas no gráfico de pedidos
          pointStyle: 'cross', // pontos em formato "x" apenas no gráfico de pedidos
          pointRadius: 6, // tamanho dos pontos
          pointHoverRadius: 8, // tamanho ao passar o mouse
        }),
      },
      {
        label: isRevenueMode ? 'Receita Cancelada' : 'Pedidos Cancelados',
        data: cancelledData,
        borderColor: '#ff6b6b',
        backgroundColor: 'rgba(255, 107, 107, 0.1)',
        tension: 0.4,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: isRevenueMode ? 'Evolução Diária de Receita' : 'Evolução Diária de Pedidos',
        font: {
          size: 18,
        },
      },
      tooltip: {
        callbacks: {
          label: function (context) {
            if (isRevenueMode) {
              return `${context.dataset.label}: R$ ${context.parsed.y.toFixed(2)}`;
            } else {
              return `${context.dataset.label}: ${context.parsed.y} pedidos`;
            }
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function (value) {
            if (isRevenueMode) {
              return 'R$ ' + value.toFixed(0);
            } else {
              return value;
            }
          },
        },
      },
    },
  };

  return (
    <div className="revenue-chart">
      <div className="chart-container">
        <Line data={chartData} options={options} /> // exibe o gráfico
      </div>
    </div>
  );
}

export default RevenueChart;

