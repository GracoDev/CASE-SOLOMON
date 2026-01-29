import React, { useState } from 'react';
import './SyncButton.css';

function SyncButton({ onSync }) {
  const [loading, setLoading] = useState(false);

  const handleClick = async () => { // função para sincronizar os dados, chamada quando o botão é clicado
    setLoading(true); // inicia o loading
    try {
      await onSync(); // faz uma requisição POST para o endpoint /sync do Backend 1
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      onClick={handleClick}
      disabled={loading}
      className="sync-button" // mesmo do arquivo CSS
    >
      {loading ? 'Sincronizando...' : '🔄 Sincronizar Dados'}
    </button>
  );
}

export default SyncButton;






