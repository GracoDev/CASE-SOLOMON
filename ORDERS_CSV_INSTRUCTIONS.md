# 📄 Instruções sobre o arquivo orders.csv

## ⚠️ IMPORTANTE

O arquivo `orders.csv` **NÃO está salvo no Git** porque:
- É um arquivo de dados grande
- Pode conter informações sensíveis
- É fornecido separadamente no case

## 📥 Como obter o arquivo orders.csv

### Opção 1: Baixar do case original
- O arquivo deve estar disponível no material do case
- Baixe e coloque na **raiz do projeto** (mesmo nível do `docker-compose.yml`)

### Opção 2: Estrutura esperada
Se você precisar criar manualmente, o arquivo deve ter esta estrutura:

```csv
order_id;created_at;status;value;payment_method
20260120-001;2026-01-20T00:10:00Z;approved;199,90;credit_card
20260120-002;2026-01-20T01:20:00Z;approved;349,00;pix
20260120-003;2026-01-20T02:30:00Z;pending;129,90;credit_card
...
```

**Características:**
- Delimitador: `;` (ponto e vírgula)
- Encoding: UTF-8
- Campos obrigatórios:
  - `order_id`
  - `created_at` (formato ISO: `2026-01-20T00:10:00Z`)
  - `status` (valores: `approved`, `pending`, `cancelled`)
  - `value` (formato brasileiro com vírgula: `199,90`)
  - `payment_method` (valores: `credit_card`, `pix`, `boleto`)

## ✅ Verificar se o arquivo está correto

Após colocar o arquivo na raiz, teste:

```bash
# Verificar se o arquivo existe
ls orders.csv

# Testar o Data Source Server
docker compose up -d data-source
curl http://localhost:3000
```

Se retornar uma lista de pedidos em JSON, está funcionando! ✅

