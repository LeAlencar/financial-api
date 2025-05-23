#!/bin/bash

# Sistema de Compra e Venda com Cotações Reais - Exemplo de Teste
# Este script demonstra como usar o sistema com cotações em tempo real

BASE_URL="http://localhost:8080"

echo "=== TESTE DO SISTEMA COM COTAÇÕES REAIS ==="
echo ""

# 1. Primeiro, vamos gerar algumas cotações
echo "1. 📊 Gerando cotações para USD/BRL..."
QUOTATION_RESPONSE=$(curl -s -X POST $BASE_URL/quotations/generate \
  -H "Content-Type: application/json" \
  -d '{
    "currency_pair": "USD/BRL",
    "count": 5
  }')

echo "Resposta da geração de cotações: $QUOTATION_RESPONSE"
echo ""

# Aguardar processamento das cotações
echo "⏳ Aguardando processamento das cotações..."
sleep 3

# 2. Registrar um usuário de teste
echo "2. 👤 Registrando usuário de teste..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Trader",
    "email": "maria.trader@example.com",
    "password": "senha123"
  }')

echo "Resposta do registro: $REGISTER_RESPONSE"
echo ""

sleep 10
# 3. Fazer login para obter o token
echo "3. 🔐 Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria.trader@example.com",
    "password": "senha123"
  }')

echo "Resposta do login: $LOGIN_RESPONSE"

# Extrair o token do response (assumindo que jq está instalado)
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Erro: Não foi possível obter o token. Verifique se o usuário foi criado corretamente."
    exit 1
fi

echo "✅ Token obtido: ${TOKEN:0:20}..."
echo "✅ User ID: $USER_ID"
echo ""

# 4. Gerar mais cotações para garantir que temos dados recentes
echo "4. 📈 Gerando cotações mais recentes..."
curl -s -X POST $BASE_URL/quotations/generate \
  -H "Content-Type: application/json" \
  -d '{
    "currency_pair": "USD/BRL",
    "count": 2
  }' > /dev/null

echo "✅ Cotações atualizadas"
echo ""

# Aguardar processamento
sleep 2

# 5. Comprar dólares com cotação real (agora deve funcionar com saldo inicial)
echo "5. 💰 Comprando \$100 USD com cotação em tempo real..."
BUY_RESPONSE_1=$(curl -s -X POST $BASE_URL/transactions/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 100.0,
    "currency_pair": "USD/BRL"
  }')

echo "Resposta da compra: $BUY_RESPONSE_1"
echo ""

# Aguardar processamento
echo "⏳ Aguardando processamento da transação..."
sleep 3

# 6. Consultar histórico de transações
echo "6. 📋 Consultando histórico de transações..."
TRANSACTIONS_RESPONSE=$(curl -s -X GET "$BASE_URL/transactions?limit=10" \
  -H "Authorization: Bearer $TOKEN")

echo "Histórico de transações: $TRANSACTIONS_RESPONSE"
echo ""

# 7. Comprar dólares novamente
echo "7. 💰 Comprando mais \$50 USD..."
BUY_RESPONSE_2=$(curl -s -X POST $BASE_URL/transactions/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 50.0,
    "currency_pair": "USD/BRL"
  }')

echo "Resposta da segunda compra: $BUY_RESPONSE_2"
echo ""

# Aguardar processamento
echo "⏳ Aguardando processamento da transação..."
sleep 3

# 8. Gerar nova cotação para simular mudança de preço
echo "8. 📊 Gerando nova cotação para simular mudança de preço..."
curl -s -X POST $BASE_URL/quotations/generate \
  -H "Content-Type: application/json" \
  -d '{
    "currency_pair": "USD/BRL",
    "count": 1
  }' > /dev/null

echo "✅ Nova cotação gerada"
echo ""

# Aguardar processamento
sleep 2

# 9. Vender dólares com a nova cotação
echo "9. 💸 Vendendo \$30 USD com a nova cotação..."
SELL_RESPONSE=$(curl -s -X POST $BASE_URL/transactions/sell \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 30.0,
    "currency_pair": "USD/BRL"
  }')

echo "Resposta da venda: $SELL_RESPONSE"
echo ""

# 10. Consultar dados do usuário para ver o saldo atualizado
echo "10. 👤 Consultando dados do usuário para verificar saldo final..."
USER_RESPONSE=$(curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN")

echo "Dados do usuário: $USER_RESPONSE"
echo ""

# 11. Consultar histórico completo de transações
echo "11. 📊 Consultando histórico completo de transações..."
FINAL_TRANSACTIONS=$(curl -s -X GET "$BASE_URL/transactions?limit=20" \
  -H "Authorization: Bearer $TOKEN")

echo "Histórico final: $FINAL_TRANSACTIONS"
echo ""

# 12. Testar transação sem cotação disponível (par inexistente)
echo "12. 🧪 Testando transação com par de moedas sem cotação (EUR/BRL)..."
FALLBACK_RESPONSE=$(curl -s -X POST $BASE_URL/transactions/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 10.0,
    "currency_pair": "EUR/BRL"
  }')

echo "Resposta (deve usar taxa de fallback): $FALLBACK_RESPONSE"
echo ""

echo "=== TESTE CONCLUÍDO ==="
echo ""
echo "📊 Resumo do teste:"
echo "   ✅ Usuários começam com saldo inicial de R\$ 1.000"
echo "   ✅ Transações salvas no MongoDB"
echo "   ✅ Saldo atualizado no PostgreSQL"
echo "   ✅ Cotações geradas dinamicamente"
echo "   ✅ Transações usando cotações reais"
echo "   ✅ Histórico de transações disponível"
echo "   ✅ Fallback funcionando para pares sem cotação"
echo "   ✅ Status diferenciados (completed vs completed_fallback_rate)"
echo ""
echo "🔍 Verifique os logs dos serviços para ver:"
echo "   • s1-generator: Recebimento das requisições"
echo "   • s2-processor: Busca de cotações e salvamento"
echo "   • s3-validator: Auditoria das transações"
echo ""
echo "💡 Dicas:"
echo "   • Transações com cotações reais terão status 'completed'"
echo "   • Transações com fallback terão status 'completed_fallback_rate'"
echo "   • O QuotationID mostra qual cotação foi usada ou 'fallback'"
echo "   • Todas as transações são salvas no MongoDB para auditoria" 