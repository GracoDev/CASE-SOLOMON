import React from 'react';
import './MetricsCards.css';

function MetricsCards({ metrics }) {
  if (!metrics) return null;

  const { financial_metrics, operational_metrics } = metrics;

  const formatCurrency = (value) => { // função para formatar o valor em reais
    return new Intl.NumberFormat('pt-BR', { // cria um objeto Intl.NumberFormat para formatar o valor em reais
      style: 'currency', // estilo de moeda
      currency: 'BRL', // moeda
    }).format(value); // formata o valor em reais
  };

  return ( // retorna o componente MetricsCards que exibe as métricas financeiras e operacionais
    <div className="metrics-cards">
      <div className="metrics-section">
        <h2>Métricas Financeiras</h2>
        <div className="cards-grid">
          <div className="metric-card revenue-approved">
            <h3>Receita Aprovada</h3>
            <p className="metric-value">{formatCurrency(financial_metrics.approved_revenue)}</p>
          </div>
          <div className="metric-card revenue-pending">
            <h3>Receita Pendente</h3>
            <p className="metric-value">{formatCurrency(financial_metrics.pending_revenue)}</p>
          </div>
          <div className="metric-card revenue-cancelled">
            <h3>Receita Cancelada</h3>
            <p className="metric-value">{formatCurrency(financial_metrics.cancelled_revenue)}</p>
          </div>
        </div>
      </div>

      <div className="metrics-section">
        <h2>Métricas Operacionais</h2>
        <div className="cards-grid">
          <div className="metric-card orders-approved">
            <h3>Pedidos Aprovados</h3>
            <p className="metric-value">{operational_metrics.approved_orders}</p>
          </div>
          <div className="metric-card orders-pending">
            <h3>Pedidos Pendentes</h3>
            <p className="metric-value">{operational_metrics.pending_orders}</p>
          </div>
          <div className="metric-card orders-cancelled">
            <h3>Pedidos Cancelados</h3>
            <p className="metric-value">{operational_metrics.cancelled_orders}</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default MetricsCards;




