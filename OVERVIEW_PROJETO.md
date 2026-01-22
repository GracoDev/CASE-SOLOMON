# 📊 Overview Completo do Projeto - Mini Analytics Platform

## 🏗️ Arquitetura Geral

```
┌─────────────┐
│   Frontend  │ (React - Porta 3001)
│  (React)    │
└──────┬──────┘
       │
       ├─────────────────┬─────────────────┐
       │                 │                 │
       ▼                 ▼                 ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Backend 1  │  │  Backend 2  │  │   pgAdmin   │
│  (Auth)     │  │  (Query API)│  │  (Visual)  │
│  Porta 5000 │  │  Porta 8080 │  │  Porta 5050 │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                 │                 │
       │                 │                 │
       ▼                 │                 │
┌─────────────┐         │                 │
│   Pipeline  │         │                 │
│  (Go)       │         │                 │
│  Porta 8081 │         │                 │
└──────┬──────┘         │                 │
       │                 │                 │
       ├─────────┐       │                 │
       │         │       │                 │
       ▼         ▼       │                 │
┌─────────────┐ ┌─────────────┐           │
│ Data Source │ │ Transformer │           │
│  (Python)   │ │  (Python)   │           │
│  Porta 3000 │ │  Porta 8082 │           │
└─────────────┘ └──────┬───────┘           │
                       │                   │
                       └───────────┬───────┘
                                   │
                                   ▼
                          ┌─────────────┐
                          │  PostgreSQL │
                          │  Porta 5432 │
                          └─────────────┘
```

---

## 📁 Estrutura de Arquivos e Conexões

### 1. **Frontend** (`frontend/`)
**Linguagem:** React/JavaScript  
**Porta:** 3001

#### Arquivos Principais:
- `src/App.js` → Roteamento principal (Login/Dashboard)
- `src/components/Login.js` → Tela de login
- `src/components/Dashboard.js` → Dashboard com métricas
- `src/services/api.js` → Comunicação com backends

#### Conexões:
```
Frontend → POST http://localhost:5000/login (Backend 1)
Frontend → POST http://localhost:5000/sync (Backend 1) + JWT
Frontend → GET http://localhost:8080/api/metrics (Backend 2) + JWT
Frontend → GET http://localhost:8080/api/metrics/time-series (Backend 2) + JWT
```

---

### 2. **Backend 1 - Auth & Trigger** (`backend1-auth/`)
**Linguagem:** Python/Flask  
**Porta:** 5000

#### Arquivos Principais:
- `app.py` → Autenticação JWT e trigger do pipeline
- `requirements.txt` → Dependências (flask, flask-cors, pyjwt, requests)
- `Dockerfile` → Build da imagem Docker

#### Endpoints:
- `POST /login` → Autenticação (retorna JWT)
- `POST /sync` → Dispara pipeline (requer JWT)

#### Conexões:
```
Backend 1 → POST http://pipeline:8080/trigger (Pipeline)
```

#### Variáveis de Ambiente:
- `JWT_SECRET` → Chave para assinar tokens
- `PIPELINE_URL` → URL do pipeline

---

### 3. **Backend 2 - Query API** (`backend2-api/`)
**Linguagem:** Go  
**Porta:** 8080

#### Arquivos Principais:
- `main.go` → API de consulta de métricas
- `Dockerfile` → Build da imagem Docker

#### Endpoints:
- `GET /api/metrics` → Métricas agregadas (requer JWT)
- `GET /api/metrics/time-series` → Séries temporais (requer JWT)
- `GET /health` → Health check

#### Conexões:
```
Backend 2 → PostgreSQL (SELECT em aggregated.daily_metrics)
```

#### Variáveis de Ambiente:
- `DATABASE_URL` → String de conexão PostgreSQL
- `JWT_SECRET` → Validação de tokens

---

### 4. **Data Source Server** (`data-source/`)
**Linguagem:** Python/Flask  
**Porta:** 3000

#### Arquivos Principais:
- `server.py` → Servidor que expõe dados do CSV
- `requirements.txt` → Dependências (flask)
- `Dockerfile` → Build da imagem Docker

#### Endpoints:
- `GET /` → Retorna todos os pedidos do CSV
- `GET /health` → Health check

#### Conexões:
```
Data Source → Lê orders.csv (volume montado)
Pipeline → GET http://data-source:3000/ (busca dados)
```

#### Volumes:
- `./orders.csv:/app/orders.csv` → Arquivo CSV montado no container

---

### 5. **Pipeline** (`pipeline/`)
**Linguagem:** Go  
**Porta:** 8081 (externa) / 8080 (interna)

#### Arquivos Principais:
- `main.go` → Pipeline de ingestão de dados
- `Dockerfile` → Build da imagem Docker

#### Endpoints:
- `POST /trigger` → Dispara ingestão de dados
- `GET /health` → Health check

#### Funções Principais:
- `triggerHandler()` → Handler HTTP que recebe POST /trigger
- `runPipeline()` → Executa o pipeline completo
- `fetchOrders()` → Busca dados do Data Source (GET)
- `insertOrders()` → Insere dados no PostgreSQL (SQL INSERT)
- `callTransformer()` → Chama transformer via HTTP

#### Conexões:
```
Pipeline → GET http://data-source:3000/ (busca dados)
Pipeline → PostgreSQL (INSERT em raw_data.orders)
Pipeline → POST http://transformer:8080/transform (chama transformer)
```

#### Variáveis de Ambiente:
- `DATA_SOURCE_URL` → URL do Data Source
- `DATABASE_URL` → String de conexão PostgreSQL
- `TRANSFORMER_URL` → URL do Transformer

---

### 6. **Transformer** (`transformer/`)
**Linguagem:** Python/Flask  
**Porta:** 8082 (externa) / 8080 (interna)

#### Arquivos Principais:
- `transform.py` → Agregação de dados
- `requirements.txt` → Dependências (psycopg2-binary, flask, flask-cors)
- `Dockerfile` → Build da imagem Docker

#### Endpoints:
- `POST /transform` → Executa agregação de dados
- `GET /health` → Health check

#### Funções Principais:
- `aggregate_data()` → Agrega dados de raw_data.orders
- `insert_aggregated_data()` → Insere em aggregated.daily_metrics

#### Conexões:
```
Transformer → PostgreSQL (SELECT de raw_data.orders)
Transformer → PostgreSQL (INSERT em aggregated.daily_metrics)
Pipeline → POST http://transformer:8080/transform (chama transformer)
```

#### Variáveis de Ambiente:
- `DATABASE_URL` → String de conexão PostgreSQL

---

### 7. **PostgreSQL** (`postgres/`)
**Banco de Dados:** PostgreSQL 15  
**Porta:** 5432

#### Arquivos Principais:
- `init.sql` → Script de inicialização (cria schemas)

#### Schemas:
- `raw_data` → Dados brutos (criado por init.sql)
- `aggregated` → Dados agregados (criado por init.sql)

#### Tabelas:
- `raw_data.orders` → Criada pelo Pipeline
  - Campos: id, order_id, created_at, status, value, payment_method, created_at_pipeline
- `aggregated.daily_metrics` → Criada pelo Transformer
  - Campos: id, date, status, payment_method, total_orders, total_value, created_at

#### Conexões:
```
Pipeline → PostgreSQL (INSERT em raw_data.orders)
Transformer → PostgreSQL (SELECT de raw_data.orders, INSERT em aggregated.daily_metrics)
Backend 2 → PostgreSQL (SELECT de aggregated.daily_metrics)
```

---

### 8. **Docker Compose** (`docker-compose.yml`)
**Arquivo:** Orquestração de todos os serviços

#### Configurações:
- **Rede:** `analytics-network` (bridge)
- **Volumes:** 
  - `postgres_data` → Dados persistentes do PostgreSQL
  - `./orders.csv` → CSV montado no Data Source
- **Dependências:** Define ordem de inicialização

---

## 🔄 Fluxo de Dados Completo

### Fluxo 1: Login e Autenticação
```
1. Frontend → POST /login (Backend 1)
   Envia: {username: "admin", password: "admin123"}
   
2. Backend 1 → Valida credenciais
   Retorna: {token: "eyJhbGci..."}
   
3. Frontend → Salva token no localStorage
```

### Fluxo 2: Sincronização de Dados
```
1. Frontend → POST /sync (Backend 1) + JWT
   
2. Backend 1 → Valida JWT → POST /trigger (Pipeline)
   
3. Pipeline → GET / (Data Source)
   Data Source → Lê orders.csv → Retorna JSON
   
4. Pipeline → INSERT em raw_data.orders (PostgreSQL)
   
5. Pipeline → POST /transform (Transformer)
   
6. Transformer → SELECT de raw_data.orders (PostgreSQL)
   Transformer → Agrega dados
   Transformer → INSERT em aggregated.daily_metrics (PostgreSQL)
   
7. Pipeline → Retorna resultado para Backend 1
   Backend 1 → Retorna resultado para Frontend
```

### Fluxo 3: Consulta de Métricas
```
1. Frontend → GET /api/metrics (Backend 2) + JWT
   
2. Backend 2 → Valida JWT
   
3. Backend 2 → SELECT de aggregated.daily_metrics (PostgreSQL)
   
4. Backend 2 → Retorna JSON com métricas
   
5. Frontend → Exibe métricas no dashboard
```

### Fluxo 4: Séries Temporais
```
1. Frontend → GET /api/metrics/time-series (Backend 2) + JWT
   
2. Backend 2 → Valida JWT
   
3. Backend 2 → SELECT agrupado por data de aggregated.daily_metrics
   
4. Backend 2 → Retorna JSON com séries temporais
   
5. Frontend → Exibe gráfico
```

---

## 🔐 Autenticação e Segurança

### JWT (JSON Web Token)
- **Geração:** Backend 1 (`/login`)
- **Validação:** Backend 1 e Backend 2
- **Secret:** `JWT_SECRET` (mesma chave em ambos)
- **Expiração:** 24 horas

### Fluxo de Autenticação:
```
Login → Token JWT → Header Authorization: Bearer <token>
```

---

## 📊 Banco de Dados - Estrutura

### Schema: `raw_data`
```sql
raw_data.orders
├── id (SERIAL PRIMARY KEY)
├── order_id (VARCHAR UNIQUE)
├── created_at (TIMESTAMP)
├── status (VARCHAR)
├── value (NUMERIC)
├── payment_method (VARCHAR)
└── created_at_pipeline (TIMESTAMP)
```

### Schema: `aggregated`
```sql
aggregated.daily_metrics
├── id (SERIAL PRIMARY KEY)
├── date (DATE)
├── status (VARCHAR)
├── payment_method (VARCHAR)
├── total_orders (INTEGER)
├── total_value (NUMERIC)
├── created_at (TIMESTAMP)
└── UNIQUE(date, status, payment_method)
```

---

## 🌐 Rede Docker

### Rede: `analytics-network`
Todos os serviços estão na mesma rede Docker, permitindo comunicação interna usando nomes dos containers:
- `http://postgres:5432`
- `http://data-source:3000`
- `http://pipeline:8080`
- `http://transformer:8080`
- `http://backend1-auth:5000`
- `http://backend2-api:8080`

---

## 📦 Dependências entre Serviços

```
PostgreSQL (base)
    ↑
    ├── Backend 2 (depende de postgres)
    ├── Pipeline (depende de postgres, data-source, transformer)
    ├── Transformer (depende de postgres)
    └── pgAdmin (depende de postgres)

Data Source (independente)
    ↑
    └── Pipeline (depende de data-source)

Pipeline
    ↑
    └── Backend 1 (depende de pipeline)

Backend 1 + Backend 2
    ↑
    └── Frontend (depende de ambos)
```

---

## 🚀 Como Tudo Se Conecta

1. **Docker Compose** orquestra todos os serviços
2. **Rede Docker** permite comunicação interna
3. **Variáveis de Ambiente** configuram URLs e credenciais
4. **Volumes** montam arquivos (CSV, init.sql)
5. **Health Checks** garantem ordem de inicialização
6. **JWT** autentica requisições entre Frontend e Backends
7. **HTTP APIs** conectam serviços (Pipeline ↔ Data Source, Pipeline ↔ Transformer)
8. **PostgreSQL** centraliza todos os dados

---

## 📝 Resumo por Serviço

| Serviço | Porta | Linguagem | Função Principal | Conecta Com |
|---------|-------|-----------|------------------|-------------|
| Frontend | 3001 | React | Interface do usuário | Backend 1, Backend 2 |
| Backend 1 | 5000 | Python | Auth + Trigger | Pipeline |
| Backend 2 | 8080 | Go | Query API | PostgreSQL |
| Data Source | 3000 | Python | Serve CSV | Pipeline (via GET) |
| Pipeline | 8081 | Go | Ingestão | Data Source, PostgreSQL, Transformer |
| Transformer | 8082 | Python | Agregação | PostgreSQL |
| PostgreSQL | 5432 | SQL | Banco de dados | Pipeline, Transformer, Backend 2 |
| pgAdmin | 5050 | Web | Visualização DB | PostgreSQL |

. Acessar o projeto
Frontend: http://localhost:3001
Backend 1: http://localhost:5000
Backend 2: http://localhost:8080
Data Source: http://localhost:3000
Pipeline: http://localhost:8081
Transformer: http://localhost:8082
pgAdmin: http://localhost:5050

---

## 🎯 Pontos-Chave

1. **Separação de Responsabilidades:** Cada serviço tem uma função específica
2. **Microserviços:** Serviços independentes e comunicam via HTTP
3. **Pipeline de Dados:** Data Source → Pipeline → Transformer → Backend 2
4. **Autenticação Centralizada:** Backend 1 gerencia JWT
5. **Dados em Camadas:** raw_data (bruto) → aggregated (processado)
6. **Docker:** Tudo containerizado e orquestrado

---

Este é o overview completo do projeto! Cada arquivo tem seu papel específico e se conecta através de APIs HTTP, banco de dados ou rede Docker.

