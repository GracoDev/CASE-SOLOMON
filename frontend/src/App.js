import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import { AuthProvider, useAuth } from './context/AuthContext';
import './App.css';

function PrivateRoute({ children }) { // verifica se o usuário está autenticado
  const { token } = useAuth();
  return token ? children : <Navigate to="/login" />;
}

function App() { // renderiza a aplicação
  return (
    <AuthProvider> // gerencia o estado de autenticação
      <Router>
        <div className="App">
          <Routes>
            <Route path="/login" element={<Login />} /> // rota para a página de login
            <Route
              path="/dashboard"
              element={
                <PrivateRoute> // protege dashboard, redireciona para login se não estiver autenticado
                  <Dashboard />
                </PrivateRoute>
              }
            />
            <Route path="/" element={<Navigate to="/dashboard" />} /> // rota para a página de dashboard
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;




