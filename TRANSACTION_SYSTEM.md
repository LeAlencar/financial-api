# Sistema de Compra e Venda de Moedas

Este documento descreve o sistema de compra e venda de moedas (dólares) implementado seguindo a arquitetura de microserviços com RabbitMQ.

## Arquitetura

O sistema segue o mesmo padrão dos outros recursos:

1. **s1-generator**: Recebe as requisições de compra/venda e envia eventos para RabbitMQ
2. **s2-processor**: Processa os eventos de transação, valida saldo, **salva transações no MongoDB** e atualiza o usuário no PostgreSQL
3. **s3-validator**: Recebe as transações processadas para logging e auditoria

## Persistência de Dados

### PostgreSQL (Usuários)

- **Tabela**: `users`
- **Saldo inicial**: R$ 1.000,00 para novos usuários
- **Atualização**: Saldo é atualizado após cada transação bem-sucedida

### MongoDB (Transações)

- **Collection**: `currency_transactions`
- **Todas as transações**: Sucessos, falhas e fallbacks são salvos
- **Rastreabilidade completa**: Histórico completo de todas as operações

## APIs Disponíveis

### Consultar Histórico de Transações

```bash
GET /transactions?limit=50
Authorization: Bearer <token>
```

### Comprar Dólares

```bash
POST /transactions/buy
Authorization: Bearer <token>
{
  "amount": 100.0,
  "currency_pair": "USD/BRL"
}
```

### Vender Dólares

```bash
POST /transactions/sell
Authorization: Bearer <token>
{
  "amount": 50.0,
  "currency_pair": "USD/BRL"
}
```

## Fluxo de Funcionamento

### 1. Compra de Dólares (BUY)

```
POST /transactions/buy
Headers: Authorization: Bearer <token>
Body: {
  "amount": 100.0,
  "currency_pair": "USD/BRL",
  "quotation_id": "optional"
}
```

**Fluxo:**

1. s1-generator recebe a requisição
2. Valida o token JWT e extrai o user_id
3. Cria um `TransactionEvent` com action "BUY"
4. Envia para a fila "transactions" no RabbitMQ
5. s2-processor consome o evento
6. Verifica o saldo do usuário
7. Se o saldo for suficiente:
   - Cria a transação com status "completed"
   - Debita o valor do saldo do usuário
   - Envia para s3-validator
8. Se o saldo for insuficiente:
   - Cria a transação com status "failed: insufficient_balance"
   - Envia para s3-validator

### 2. Venda de Dólares (SELL)

```
POST /transactions/sell
Headers: Authorization: Bearer <token>
Body: {
  "amount": 50.0,
  "currency_pair": "USD/BRL",
  "quotation_id": "optional"
}
```

**Fluxo:**

1. s1-generator recebe a requisição
2. Valida o token JWT e extrai o user_id
3. Cria um `TransactionEvent` com action "SELL"
4. Envia para a fila "transactions" no RabbitMQ
5. s2-processor consome o evento
6. Cria a transação com status "completed"
7. Credita o valor no saldo do usuário
8. Envia para s3-validator

## Estrutura dos Eventos

### TransactionEvent

```json
{
  "action": "BUY" | "SELL",
  "data": {
    "user_id": "123",
    "currency_pair": "USD/BRL",
    "amount": 100.0,
    "quotation_id": "optional",
    "timestamp": "2023-12-01T10:00:00Z"
  }
}
```

### Transaction (modelo final)

```json
{
  "id": "TXN_1701421200000000000",
  "user_id": "123",
  "type": "BUY" | "SELL",
  "currency_pair": "USD/BRL",
  "amount": 100.0,
  "exchange_rate": 5.5,
  "total_value": 550.0,
  "status": "completed" | "failed: reason",
  "timestamp": "2023-12-01T10:00:00Z",
  "quotation_id": "optional"
}
```

## Taxas de Câmbio

O sistema agora busca **cotações em tempo real** do serviço de quotations:

- **Busca automática**: Para cada transação, o sistema consulta a cotação mais recente do par de moedas
- **Fallback inteligente**: Se não houver cotação disponível, usa taxas de segurança:
  - **Compra**: 1 USD = 5.5 BRL
  - **Venda**: 1 USD = 5.3 BRL (spread de 0.2)
- **Status diferenciados**:
  - `"completed"`: Transação com cotação real
  - `"completed_fallback_rate"`: Transação com taxa de fallback

### Como funciona a busca de cotações:

1. Sistema consulta `quotations` collection no MongoDB
2. Busca a cotação mais recente (`created_at` DESC) para o `currency_pair`
3. Usa `BuyPrice` para compras e `SellPrice` para vendas
4. Se não encontrar cotação, aplica taxas de fallback com status especial

> **Integração completa**: O sistema agora está totalmente integrado com o serviço de cotações, garantindo preços sempre atualizados!

## Exemplo de Uso

### 1. Registrar um usuário

```bash
POST /users/register
{
  "name": "João Silva",
  "email": "joao@example.com",
  "password": "senha123"
}
```

### 2. Fazer login

```bash
POST /users/login
{
  "email": "joao@example.com",
  "password": "senha123"
}
```

Resposta: `{"token": "eyJ...", "user": {...}}`

### 3. Comprar dólares

```bash
POST /transactions/buy
Authorization: Bearer eyJ...
{
  "amount": 100.0,
  "currency_pair": "USD/BRL"
}
```

### 4. Vender dólares

```bash
POST /transactions/sell
Authorization: Bearer eyJ...
{
  "amount": 50.0,
  "currency_pair": "USD/BRL"
}
```

## Segurança

- Todas as rotas de transação requerem autenticação JWT
- O user_id é extraído do token, não pode ser falsificado
- Validação de saldo para compras
- Todas as transações são auditadas no s3-validator

## Filas RabbitMQ

- **transactions**: Eventos de compra/venda do s1-generator para s2-processor
- **transactions-validator**: Transações processadas do s2-processor para s3-validator

## Melhorias Futuras

1. **Integração com cotações reais**: Usar o serviço de quotations para taxas dinâmicas
2. **Validação de posse**: Para vendas, verificar se o usuário possui a quantidade de dólares
3. **Limites de transação**: Implementar limites diários/mensais
4. **Histórico de transações**: Endpoint para consultar transações do usuário
5. **Notificações**: Notificar usuários sobre transações concluídas
6. **Rollback**: Mecanismo de rollback para transações com falha
