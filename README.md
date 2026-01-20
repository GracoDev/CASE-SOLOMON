# Mini Analytics Platform - Case Solomon

Plataforma de analytics completa com pipeline de dados, backend e dashboard.

## 📋 Requisitos

- Docker Desktop instalado e rodando
- Git instalado
- Portas disponíveis: 3000, 3001, 5000, 8080, 5432

## 🚀 Como começar em um novo PC

### 1. Clonar o repositório

```bash
git clone <URL_DO_SEU_REPOSITORIO>
cd CASE-SOLOMON
```

### 2. Arquivos que precisam ser criados/baixados

#### 📄 `orders.csv` (OBRIGATÓRIO)
Este arquivo **NÃO** está no Git (por ser grande). Você precisa:

- **Opção 1**: Baixar do repositório original do case
- **Opção 2**: Criar manualmente na raiz do projeto com a estrutura:
  ```
  order_id;created_at;status;value;payment_method
  20260120-001;2026-01-20T00:10:00Z;approved;199,90;credit_card
  ...
  ```

**IMPORTANTE**: O arquivo `orders.csv` deve estar na raiz do projeto para o Docker montar corretamente.

### 3. Subir os serviços

```bash
docker compose up --build -d
```

### 4. Verificar se está tudo rodando

```bash
docker compose ps
```

Todos os 7 serviços devem estar com status "Up":
- ✅ postgres
- ✅ backend1-auth
- ✅ backend2-api
- ✅ data-source
- ✅ pipeline
- ✅ transformer
- ✅ frontend

## 🧪 Testar os serviços

### Data Source Server
```bash
# No navegador ou PowerShell:
http://localhost:3000          # Retorna todos os pedidos
http://localhost:3000/health  # Health check
```

### Outros serviços
```bash
http://localhost:5000   # Backend 1 (Auth)
http://localhost:8080     # Backend 2 (Query API)
http://localhost:3001    # Frontend
```

## 📁 Estrutura do Projeto

```
CASE-SOLOMON/
├── docker-compose.yml       # Orquestração de todos os serviços
├── orders.csv              # ⚠️ NÃO está no Git - precisa baixar/criar
├── README.md               # Este arquivo
│
├── backend1-auth/          # Backend 1 - Autenticação e Trigger (Flask/Python)
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app.py
│
├── backend2-api/           # Backend 2 - Query API (Go)
│   ├── Dockerfile
│   └── main.go
│
├── data-source/            # Servidor de Origem de Dados (Python)
│   ├── Dockerfile
│   ├── requirements.txt
│   └── server.py          # ✅ Implementado
│
├── pipeline/               # Pipeline de Dados (Go)
│   ├── Dockerfile
│   └── main.go
│
├── transformer/            # Serviço de Transformação (Python)
│   ├── Dockerfile
│   ├── requirements.txt
│   └── transform.py
│
├── frontend/               # Dashboard React
│   ├── Dockerfile
│   ├── package.json
│   └── server.js
│
└── postgres/               # Scripts de inicialização do banco
    └── init.sql
```

## 📝 Arquivos ignorados pelo Git (.gitignore)

Estes arquivos **NÃO** são salvos no Git (e não precisam ser):

- `node_modules/` - Dependências do Node.js (instaladas automaticamente)
- `__pycache__/`, `*.pyc` - Cache do Python (gerado automaticamente)
- `postgres_data/` - Dados do PostgreSQL (volume do Docker, recriado automaticamente)
- `*.ps1` - Scripts temporários do PowerShell
- `*.log` - Arquivos de log

**Todos esses arquivos são gerados automaticamente quando você roda o projeto!**

## 🔧 Comandos úteis

### Ver logs de um serviço
```bash
docker compose logs data-source
docker compose logs -f data-source  # Seguir logs em tempo real
```

### Parar todos os serviços
```bash
docker compose down
```

### Reconstruir um serviço específico
```bash
docker compose build data-source
docker compose up -d data-source
```

### Ver status de todos os containers
```bash
docker compose ps
```

## ✅ Status da Implementação

- [x] **Estrutura de pastas** - Criada
- [x] **docker-compose.yml** - Configurado com todos os serviços
- [x] **Data Source Server** - ✅ Implementado e testado
- [ ] **Pipeline** - Pendente
- [ ] **Transformer** - Pendente
- [ ] **Backend 1 (Auth)** - Pendente
- [ ] **Backend 2 (Query API)** - Pendente
- [ ] **Frontend** - Pendente
- [ ] **PostgreSQL** - Schemas criados, tabelas pendentes

## 🐛 Troubleshooting

### Erro: "orders.csv não encontrado"
- Certifique-se de que o arquivo `orders.csv` está na raiz do projeto
- Verifique se o arquivo tem conteúdo (não está vazio)

### Erro: "Port already in use"
- Pare outros serviços que usam as portas 3000, 3001, 5000, 8080, 5432
- Ou altere as portas no `docker-compose.yml`

### Container não inicia
```bash
docker compose logs <nome-do-servico>
```

## 📚 Próximos Passos

1. Implementar o **Pipeline** (Go) - Consumir dados do Data Source
2. Implementar o **Transformer** (Python) - Agregar dados
3. Implementar o **Backend 1** (Flask) - Autenticação
4. Implementar o **Backend 2** (Go) - Query API
5. Implementar o **Frontend** (React) - Dashboard

## 📞 Contato

Para dúvidas sobre o case, consulte a documentação original.

