import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { backend1API, backend2API } from '../services/api';
import MetricsCards from './MetricsCards';
import RevenueChart from './RevenueChart';
import Filters from './Filters';
import SyncButton from './SyncButton';
import './Dashboard.css';

function Dashboard() { // componente Dashboard que exibe os dados e filtros
  const { token, logout } = useAuth(); // obtém o token e a função logout do contexto de autenticação
  const navigate = useNavigate();
  const [metrics, setMetrics] = useState(null); 
  const [timeSeries, setTimeSeries] = useState([]);
  const [loading, setLoading] = useState(true); 
  const [error, setError] = useState(''); 
  const [filters, setFilters] = useState({ 
    startDate: '',
    endDate: '',
    paymentMethod: '',
  });
  const [tempFilters, setTempFilters] = useState({ // estado para os filtros temporários
    startDate: '',
    endDate: '',
    paymentMethod: '',
  });
  const [chartViewMode, setChartViewMode] = useState('revenue'); // estado inicial do modo de visualização do gráfico ('revenue' ou 'orders') nesse caso 'revenue'

  const loadData = async (filtersToUse = filters) => { // função para carregar os dados
    setLoading(true);
    setError('');
    try {
      const [metricsData, timeSeriesData] = await Promise.all([ // faz uma requisição GET para o endpoint /api/metrics e /api/metrics/time-series do Backend 2
        backend2API.getMetrics(token, filtersToUse),
        backend2API.getTimeSeries(token, filtersToUse),
      ]);
      setMetrics(metricsData);
      // Backend 2 retorna { filters, data } para time-series
      setTimeSeries(timeSeriesData.data || timeSeriesData);
    } catch (err) {
      setError(
        err.response?.data?.error || 'Erro ao carregar dados. Verifique sua conexão.'
      );
    } finally {
      setLoading(false);
    }
  };

  // Carregar dados apenas uma vez ao montar o componente
  useEffect(() => {
    loadData();
  }, [token]);

  const handleApplyFilters = () => { // função para aplicar os filtros
    setFilters({ ...tempFilters });
    loadData(tempFilters);
  };

  const handleClearFilters = () => { // função para limpar os filtros
    const clearedFilters = { startDate: '', endDate: '', paymentMethod: '' };
    setTempFilters(clearedFilters);
    setFilters(clearedFilters);
    loadData(clearedFilters);
  };

  const handleLogout = () => { // função para fazer logout
    logout();
    navigate('/login');
  };

  const handleSync = async () => { // função para sincronizar os dados. acessa o endpoint POST /sync do Backend 1 e inicia a sincronização
    try {
      await backend1API.sync(token); 
      alert('Sincronização iniciada com sucesso!');
      // Recarregar dados após sincronização (usando filtros atuais)
      setTimeout(() => {
        loadData(filters);
      }, 2000); // espera 2 segundos após a sincronização para carregar os dados novamente
    } catch (err) {
      alert(
        err.response?.data?.error || 'Erro ao sincronizar dados. Tente novamente.'
      );
    }
  };

  return ( // retorna o componente Dashboard que exibe os dados e filtros 
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Analytics Dashboard</h1>
        <div className="header-actions">
          <SyncButton onSync={handleSync} />
          <button onClick={handleLogout} className="logout-button">
            Sair
          </button>
        </div>
      </header>

      <div className="dashboard-content">
        <Filters 
          filters={tempFilters} 
          onFiltersChange={setTempFilters}
          onApplyFilters={handleApplyFilters}
          onClearFilters={handleClearFilters}
        />

        {error && <div className="error-banner">{error}</div>}

        {loading ? (
          <div className="loading">Carregando dados...</div>
        ) : (
          <>
            {metrics && <MetricsCards metrics={metrics} />} // passa os dados das métricas para o componente MetricsCards
            {timeSeries && timeSeries.length > 0 && (
              <div className="chart-section">
                <div className="chart-controls">
                  <label className="chart-view-toggle">
                    <span>Visualizar por:</span>
                    <select 
                      value={chartViewMode} 
                      onChange={(e) => setChartViewMode(e.target.value)}
                      className="chart-mode-select"
                    >
                      <option value="revenue">Receita</option>
                      <option value="orders">Número de Pedidos</option>
                    </select>
                  </label>
                </div>
                <RevenueChart data={timeSeries} viewMode={chartViewMode} /> // passa os dados da série temporal para o componente RevenueChart
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default Dashboard;

