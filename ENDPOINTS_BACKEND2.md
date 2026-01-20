# 📋 Endpoints do Backend 2 API

## Base URL
```
http://localhost:8080
```

---

## 1. Health Check
Verifica se a API está funcionando.

```
http://localhost:8080/health
```

**Resposta:**
```json
{"status":"healthy"}
```

---

## 2. Métricas Agregadas por Período

### Sem filtros (todos os dados)
```
http://localhost:8080/api/metrics
```

### Com período específico
```
http://localhost:8080/api/metrics?start_date=2026-01-20&end_date=2026-02-10
```

### Apenas cartão de crédito
```
http://localhost:8080/api/metrics?payment_method=credit_card
```

### Período + método de pagamento
```
http://localhost:8080/api/metrics?start_date=2026-01-20&end_date=2026-02-10&payment_method=credit_card
```

**Resposta exemplo:**
```json
{
  "filters": {
    "start_date": "2026-01-20",
    "end_date": "2026-02-10",
    "payment_method": "credit_card"
  },
  "financial_metrics": {
    "approved_revenue": 33527.2,
    "pending_revenue": 6375.6,
    "cancelled_revenue": 9495.6
  },
  "operational_metrics": {
    "approved_orders": 132,
    "pending_orders": 44,
    "cancelled_orders": 44
  }
}
```

**Campos da resposta:**
- `approved_revenue`: Receita de vendas aprovadas
- `pending_revenue`: Receita de vendas pendentes
- `cancelled_revenue`: Receita de vendas canceladas
- `approved_orders`: Quantidade de pedidos aprovados
- `pending_orders`: Quantidade de pedidos pendentes
- `cancelled_orders`: Quantidade de pedidos cancelados

---

## 3. Séries Temporais (para gráficos)

### Sem filtros (todos os dados)
```
http://localhost:8080/api/metrics/time-series
```

### Com período específico
```
http://localhost:8080/api/metrics/time-series?start_date=2026-01-20&end_date=2026-02-10
```

### Apenas cartão de crédito
```
http://localhost:8080/api/metrics/time-series?payment_method=credit_card
```

### Período + método de pagamento
```
http://localhost:8080/api/metrics/time-series?start_date=2026-01-20&end_date=2026-02-10&payment_method=credit_card
```

**Resposta exemplo:**
```json
{
  "filters": {
    "start_date": "2026-01-20",
    "end_date": "2026-02-10",
    "payment_method": ""
  },
  "data": [
    {
      "date": "2026-01-20",
      "approved_revenue": 1457.6,
      "pending_revenue": 289.8,
      "cancelled_revenue": 439.8,
      "approved_orders": 6,
      "pending_orders": 1,
      "cancelled_orders": 2
    },
    {
      "date": "2026-01-21",
      "approved_revenue": 1607.6,
      "pending_revenue": 259.8,
      "cancelled_revenue": 359.8,
      "approved_orders": 6,
      "pending_orders": 1,
      "cancelled_orders": 2
    }
  ]
}
```

**Campos da resposta:**
- `date`: Data do ponto da série temporal
- `approved_revenue`: Receita aprovada no dia
- `pending_revenue`: Receita pendente no dia
- `cancelled_revenue`: Receita cancelada no dia
- `approved_orders`: Pedidos aprovados no dia
- `pending_orders`: Pedidos pendentes no dia
- `cancelled_orders`: Pedidos cancelados no dia

---

## 📌 Informações Importantes

- **Período dos dados:** 2026-01-20 até 2026-02-10
- **Formato de data:** `YYYY-MM-DD` (ex: `2026-01-20`)
- **Métodos de pagamento:** `credit_card`, `billet`, `pix`
- **Filtros:** Todos são opcionais

---

## 🔍 Exemplos Práticos

### Ver todas as métricas:
```
http://localhost:8080/api/metrics
```

### Ver apenas vendas aprovadas de janeiro:
```
http://localhost:8080/api/metrics?start_date=2026-01-01&end_date=2026-01-31
```

### Ver séries temporais de cartão de crédito:
```
http://localhost:8080/api/metrics/time-series?payment_method=credit_card
```

### Ver séries temporais de vendas canceladas (período específico):
```
http://localhost:8080/api/metrics/time-series?start_date=2026-01-25&end_date=2026-02-05
```
