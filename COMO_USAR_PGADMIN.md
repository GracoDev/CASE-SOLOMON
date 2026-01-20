# 🗄️ Como Usar o pgAdmin para Ver os Dados

## 📥 Opção 1: pgAdmin via Docker (Recomendado)

### 1. Subir o pgAdmin
```powershell
docker compose up -d pgadmin
```

### 2. Acessar no navegador
```
http://localhost:5050
```

### 3. Login
- **Email:** `admin@admin.com`
- **Senha:** `admin`

### 4. Conectar ao PostgreSQL

1. Clique com botão direito em **"Servers"** → **"Register"** → **"Server"**

2. Na aba **"General"**:
   - **Name:** `PostgreSQL Local`

3. Na aba **"Connection"**:
   - **Host name/address:** `postgres` (nome do container)
   - **Port:** `5432`
   - **Maintenance database:** `analytics_db`
   - **Username:** `postgres`
   - **Password:** `postgres`
   - ✅ Marque **"Save password"**

4. Clique em **"Save"**

### 5. Navegar pelos dados

#### Ver dados brutos (raw_data.orders)

1. Expanda: **Servers** → **PostgreSQL Local** → **Databases** → **analytics_db** → **Schemas** → **raw_data** → **Tables** → **orders**

2. Clique com botão direito em **"orders"** → **"View/Edit Data"** → **"All Rows"**

3. Você verá todos os 220 pedidos em uma tabela visual! 🎉

#### Ver dados agregados (aggregated.daily_metrics)

1. Expanda: **Servers** → **PostgreSQL Local** → **Databases** → **analytics_db** → **Schemas** → **aggregated** → **Tables** → **daily_metrics**

2. Clique com botão direito em **"daily_metrics"** → **"View/Edit Data"** → **"All Rows"**

3. Você verá os dados agregados por:
   - **date** (data)
   - **status** (approved/pending/cancelled)
   - **payment_method** (credit_card/pix/billet)
   - **total_orders** (quantidade de pedidos)
   - **total_value** (soma dos valores)

4. **Dica**: Você pode ordenar clicando no cabeçalho da coluna (date, status, payment_method, etc.)

#### Executar queries personalizadas no pgAdmin

1. Clique com botão direito em **"daily_metrics"** → **"View/Edit Data"** → **"Filtered Rows..."**

2. Ou use o **Query Tool**:
   - Clique com botão direito em **"daily_metrics"** → **"Query Tool"**
   - Digite sua query, por exemplo:
     ```sql
     SELECT * FROM aggregated.daily_metrics 
     ORDER BY date DESC, total_value DESC;
     ```
   - Clique em **"Execute"** (F5)

---

## 💡 Dica

Se você já tem o Docker rodando, use o **pgAdmin via Docker** - é a opção mais rápida e não precisa instalar nada!

```powershell
docker compose up -d pgadmin
# Depois acesse: http://localhost:5050
```

