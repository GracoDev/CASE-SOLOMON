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
    <AuthProvider>  
      <Router>
        <div className="App">
          <Routes>
            <Route path="/login" element={<Login />} /> 
            <Route
              path="/dashboard"
              element={
                <PrivateRoute> 
                  <Dashboard />
                </PrivateRoute>
              }
            />
            <Route path="/" element={<Navigate to="/dashboard" />} /> 
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;




