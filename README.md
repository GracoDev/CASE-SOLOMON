# Mini Analytics Platform - Case Solomon

Plataforma de analytics completa com pipeline de dados, backend e dashboard.

## 📋 Requisitos

- Docker Desktop instalado e rodando
- Git instalado
- Portas disponíveis: 3000, 3001, 5000, 8080, 5432

## 🚀 Como começar

### 1. Clonar o repositório

```bash
git clone <URL_DO_SEU_REPOSITORIO>
cd CASE-SOLOMON
```

### 2. Subir os serviços

```bash
docker compose up --build -d
```

### 3. Verificar se está tudo rodando

```bash
docker compose ps
```

Todos os 7 serviços devem estar com status "Up".

## 🎯 Como usar o sistema

### 1. Acessar o Dashboard

Abra seu navegador e acesse:

```
http://localhost:3001
```

### 2. Fazer Login

**Credenciais padrão:**
- **Usuário:** `admin`
- **Senha:** `admin123`

### 3. Sincronizar dados (IMPORTANTE!)

Após fazer login, você verá o dashboard vazio. Para carregar os dados:

1. Clique no botão **"🔄 Sincronizar Dados"** no topo da página
2. Aguarde a mensagem de sucesso (pode levar alguns segundos)
3. Os dados serão carregados automaticamente

### 4. Explorar o Dashboard

Após sincronizar, você verá:
- **Cards de métricas**: Receitas e pedidos por status (Aprovado, Pendente, Cancelado)
- **Gráfico de séries temporais**: Evolução das métricas ao longo do tempo
- **Filtros**: Data inicial, data final e método de pagamento

## 🌐 URLs dos Serviços

### Frontend (Dashboard Principal)
```
http://localhost:3001
```

### Backend 1 (Auth)
```
http://localhost:5000
http://localhost:5000/health
```

### Backend 2 (Query API)
```
http://localhost:8080
http://localhost:8080/health
```

### Data Source Server
```
http://localhost:3000
http://localhost:3000/health
```

### pgAdmin (Interface gráfica do PostgreSQL)
```
http://localhost:5050
```

**Credenciais pgAdmin:**
- Email: `admin@admin.com`
- Senha: `admin`

### Como visualizar os dados no pgAdmin

1. **Conectar ao servidor PostgreSQL:**
   - Clique com botão direito em "Servers" → "Register" → "Server"
   - Na aba "General":
     - **Name:** `PostgreSQL` (ou qualquer nome)
   - Na aba "Connection":
     - **Host name/address:** `postgres`
     - **Port:** `5432`
     - **Maintenance database:** `analytics_db`
     - **Username:** `postgres`
     - **Password:** `postgres`
   - Clique em "Save"

2. **Visualizar dados brutos (raw_data):**
   - Expanda: `Servers` → `PostgreSQL` → `Databases` → `analytics_db` → `Schemas` → `raw_data` → `Tables` → `orders`
   - Clique com botão direito em `orders` → "View/Edit Data" → "All Rows"
   - Você verá todos os pedidos inseridos pelo pipeline

3. **Visualizar dados agregados (aggregated):**
   - Expanda: `Servers` → `PostgreSQL` → `Databases` → `analytics_db` → `Schemas` → `aggregated` → `Tables` → `daily_metrics`
   - Clique com botão direito em `daily_metrics` → "View/Edit Data" → "All Rows"
   - Você verá as métricas agregadas por dia, status e método de pagamento

**Nota:** Os dados só aparecerão após executar a sincronização no dashboard (botão "🔄 Sincronizar Dados").
