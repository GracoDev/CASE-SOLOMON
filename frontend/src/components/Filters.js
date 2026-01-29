import React from 'react';
import './Filters.css';

function Filters({ filters, onFiltersChange, onApplyFilters, onClearFilters }) { // componente Filters que exibe os filtros
  const handleChange = (field, value) => { // função para mudar o valor de um filtro
    onFiltersChange({
      ...filters,
      [field]: value,
    });
  };

  const handleClearFilters = () => { // função para limpar os filtros
    if (onClearFilters) {
      onClearFilters();
    } else {
      const clearedFilters = { startDate: '', endDate: '', paymentMethod: '' };
      onFiltersChange(clearedFilters);
    }
  };

  return ( // retorna o componente Filters que exibe os filtros
    <div className="filters">
      <h2>Filtros</h2>
      <div className="filters-grid">
        <div className="filter-group">
          <label htmlFor="start-date">Data Inicial:</label>
          <input
            type="date"
            id="start-date"
            value={filters.startDate}
            onChange={(e) => handleChange('startDate', e.target.value)}
          />
        </div>
        <div className="filter-group">
          <label htmlFor="end-date">Data Final:</label>
          <input
            type="date"
            id="end-date"
            value={filters.endDate}
            onChange={(e) => handleChange('endDate', e.target.value)}
          />
        </div>
        <div className="filter-group">
          <label htmlFor="payment-method">Método de Pagamento:</label>
          <select
            id="payment-method"
            value={filters.paymentMethod}
            onChange={(e) => handleChange('paymentMethod', e.target.value)}
          >
            <option value="">Todos</option>
            <option value="credit_card">Cartão de Crédito</option>
            <option value="boleto">Boleto</option>
            <option value="pix">PIX</option>
          </select>
        </div>
        <div className="filter-group">
          <button
            type="button"
            onClick={onApplyFilters}
            className="apply-filters-button"
          >
            Aplicar Filtros
          </button>
        </div>
        <div className="filter-group">
          <button
            type="button"
            onClick={handleClearFilters}
            className="clear-filters-button"
          >
            Limpar Filtros
          </button>
        </div>
      </div>
    </div>
  );
}

export default Filters;

