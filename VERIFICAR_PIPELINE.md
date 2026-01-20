# 🔍 Como Verificar se o Pipeline Funcionou

## 1. Verificar os Logs do Pipeline

```powershell
docker compose logs pipeline
```

**O que procurar:**
- ✅ `Conectado ao PostgreSQL`
- ✅ `Schema e tabela verificados/criados`
- ✅ `X pedidos recebidos do Data Source` (deve ser 220)
- ✅ `X pedidos inseridos com sucesso` (deve ser 220)
- ✅ `Pipeline concluído com sucesso`

**Se houver erro:**
- ❌ Verifique a mensagem de erro nos logs
- ❌ Verifique se o Data Source está rodando: `docker compose ps data-source`
- ❌ Verifique se o PostgreSQL está rodando: `docker compose ps postgres`

---

## 2. Verificar Quantidade de Pedidos no Banco

```powershell
docker exec postgres psql -U postgres -d analytics_db -c "SELECT COUNT(*) as total_pedidos FROM raw_data.orders;"
```

**Resultado esperado:**
```
 total_pedidos 
---------------
           220
```

---

## 3. Ver Alguns Pedidos Inseridos

```powershell
docker exec postgres psql -U postgres -d analytics_db -c "SELECT order_id, created_at, status, value, payment_method FROM raw_data.orders ORDER BY id LIMIT 10;"
```

**Resultado esperado:**
```
   order_id    |      created_at       |  status   | value | payment_method 
---------------+-----------------------+-----------+-------+----------------
 20260120-001  | 2026-01-20 00:10:00  | approved  | 199.9 | credit_card
 20260120-002  | 2026-01-20 01:20:00  | approved  | 349.0 | pix
 20260120-003  | 2026-01-20 02:30:00  | pending   | 129.9 | credit_card
 ...
```

---

## 4. Verificar Estrutura da Tabela

```powershell
docker exec postgres psql -U postgres -d analytics_db -c "\d raw_data.orders"
```

**Mostra:**
- Colunas da tabela
- Tipos de dados
- Constraints (chaves, unique, etc.)

---

## 5. Estatísticas dos Dados Inseridos

```powershell
# Por status
docker exec postgres psql -U postgres -d analytics_db -c "SELECT status, COUNT(*) as quantidade FROM raw_data.orders GROUP BY status;"

# Por método de pagamento
docker exec postgres psql -U postgres -d analytics_db -c "SELECT payment_method, COUNT(*) as quantidade FROM raw_data.orders GROUP BY payment_method;"

# Valor total
docker exec postgres psql -U postgres -d analytics_db -c "SELECT SUM(value) as valor_total FROM raw_data.orders;"
```

---

## 6. Verificar Schema raw_data

```powershell
docker exec postgres psql -U postgres -d analytics_db -c "\dn raw_data"
```

**Deve mostrar:**
```
  Name   | Owner  
---------+--------
 raw_data | postgres
```

---

## 7. Ver Todas as Tabelas do Schema raw_data

```powershell
docker exec postgres psql -U postgres -d analytics_db -c "\dt raw_data.*"
```

**Deve mostrar:**
```
         List of relations
 Schema   | Name  | Type  |  Owner   
----------+-------+-------+----------
 raw_data | orders | table | postgres
```

---

## ✅ Checklist de Sucesso

- [ ] Pipeline executou sem erros nos logs
- [ ] Tabela `raw_data.orders` existe
- [ ] 220 pedidos foram inseridos (ou quantidade esperada)
- [ ] Dados estão corretos (order_id, status, value, payment_method)
- [ ] Schema `raw_data` existe

---

## 🐛 Troubleshooting

### Erro: "relation raw_data.orders does not exist"
- O pipeline não executou ou falhou antes de criar a tabela
- Verifique os logs: `docker compose logs pipeline`

### Erro: "0 pedidos inseridos"
- Verifique se o Data Source está respondendo: `curl http://localhost:3000`
- Verifique os logs do pipeline para erros de inserção

### Erro de conexão com PostgreSQL
- Verifique se o PostgreSQL está rodando: `docker compose ps postgres`
- Verifique a connection string no docker-compose.yml

