import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { backend1API } from '../services/api';
import './Login.css';

function Login() { // componente Login que exibe o formulário de login
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState(''); 
  const [error, setError] = useState(''); 
  const [loading, setLoading] = useState(false); 
  const { login } = useAuth(); 
  const navigate = useNavigate(); 

  const handleSubmit = async (e) => { // função para enviar o formulário de login
    e.preventDefault();
    setError('');
    setLoading(true);

    try { // tenta fazer o login
      const response = await backend1API.login(username, password); // Acessa o endpoint POST /login do Backend 1 e recebe o token
      login(response.token); // salva o token
      navigate('/dashboard'); // redireciona para a página de dashboard
    } catch (err) {
      setError(
        err.response?.data?.error || 'Erro ao fazer login. Verifique suas credenciais.'
      );
    } finally {
      setLoading(false);
    }
  };

  return ( // retorna o componente Login que exibe o formulário de login
    <div className="login-container">
      <div className="login-box">
        <h1>Analytics Dashboard</h1>
        <h2>Login</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">Usuário:</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              placeholder="Digite seu usuário"
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">Senha:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Digite sua senha"
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <button type="submit" disabled={loading} className="login-button">
            {loading ? 'Entrando...' : 'Entrar'}
          </button>
        </form>
        <div className="login-hint">
          <p>Credenciais padrão:</p>
          <p>Usuário: <strong>admin</strong></p>
          <p>Senha: <strong>admin123</strong></p>
        </div>
      </div>
    </div>
  );
}

export default Login;




