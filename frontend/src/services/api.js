import axios from 'axios';

const BACKEND1_URL = process.env.REACT_APP_BACKEND1_URL || 'http://localhost:5000';
const BACKEND2_URL = process.env.REACT_APP_BACKEND2_URL || 'http://localhost:8080';

// API do Backend 1 (Auth & Trigger)
export const backend1API = { // API do Backend 1 (Auth & Trigger)
  login: async (username, password) => { // faz uma requisição POST para o endpoint /login do Backend 1
    const response = await axios.post(`${BACKEND1_URL}/login`, { // envia o username e password
      username, // envia o username
      password, // envia o password
    });
    return response.data; // retorna a resposta da requisição
  },

  sync: async (token) => { // faz uma requisição POST para o endpoint /sync do Backend 1
    const response = await axios.post(`${BACKEND1_URL}/sync`, {}, { // envia o token
      headers: { // envia o token
        Authorization: `Bearer ${token}`, // envia o token
      },
    });
    return response.data; // retorna a resposta da requisição
  },
};

// API do Backend 2 (Query API)
// Nota: Backend 2 não valida JWT, mas enviamos o token conforme requisito
export const backend2API = { // API do Backend 2 (Query API)
  getMetrics: async (token, filters = {}) => { // faz uma requisição GET para o endpoint /api/metrics do Backend 2
    const params = new URLSearchParams(); // cria um objeto URLSearchParams para os filtros
    if (filters.startDate) params.append('start_date', filters.startDate); // adiciona o filtro de data inicial à query
    if (filters.endDate) params.append('end_date', filters.endDate); // adiciona o filtro de data final à query
    if (filters.paymentMethod) params.append('payment_method', filters.paymentMethod); // adiciona o filtro de método de pagamento à query

    const headers = {}; // cria um objeto para os headers
    if (token) {
      headers.Authorization = `Bearer ${token}`; // envia o token       
    }

    const response = await axios.get(
      `${BACKEND2_URL}/api/metrics?${params.toString()}`, // envia o endpoint /api/metrics do Backend 2
      { headers } // envia os headers
    );
    return response.data; // retorna a resposta da requisição
  },

  getTimeSeries: async (token, filters = {}) => { // faz uma requisição GET para o endpoint /api/metrics/time-series do Backend 2
    const params = new URLSearchParams(); // cria um objeto URLSearchParams para os filtros
    if (filters.startDate) params.append('start_date', filters.startDate); // adiciona o filtro de data inicial à query
    if (filters.endDate) params.append('end_date', filters.endDate); // adiciona o filtro de data final à query
    if (filters.paymentMethod) params.append('payment_method', filters.paymentMethod); // adiciona o filtro de método de pagamento à query

    const headers = {}; // cria um objeto para os headers
    if (token) {
      headers.Authorization = `Bearer ${token}`; // envia o token
    }

    const response = await axios.get(
      `${BACKEND2_URL}/api/metrics/time-series?${params.toString()}`, // envia o endpoint /api/metrics/time-series do Backend 2
      { headers } // envia os headers
    );
    return response.data; // retorna a resposta da requisição
  },
};

