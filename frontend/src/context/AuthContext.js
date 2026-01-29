import React, { createContext, useState, useContext } from 'react';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => {
    // Recuperar token do localStorage ao inicializar
    return localStorage.getItem('jwt_token') || null;
  });

  const login = (newToken) => { // função para login
    setToken(newToken); // define o token
    localStorage.setItem('jwt_token', newToken); // salva o token no localStorage
  };

  const logout = () => { // função para logout
    setToken(null); // define o token como null
    localStorage.removeItem('jwt_token'); // remove o token do localStorage
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}> 
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth deve ser usado dentro de AuthProvider');
  }
  return context;
}






