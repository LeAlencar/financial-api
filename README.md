# Projeto Tópicos Avançados de Bancos de Dados

### Integrantes:

#### - João Paulo Paggi Zuanon Dias RA: 22.222.058-4

#### - Leandro de Brito Alencar RA: 22.222.034-5

#### - Mateus Rocha RA: 22.222.002-2

## Descrição do Projeto

Este projeto implementa um sistema de compra e venda de moedas (Dólar/BRL),
utilizando uma arquitetura distribuída e múltiplos bancos de dados para otimizar o
armazenamento de diferentes tipos de informações.
O sistema permite aos usuários realizar operações de câmbio,
acompanhar cotações em tempo real e manter um histórico de transações.

O objetivo principal é demonstrar como diferentes tipos de bancos de dados podem ser
utilizados em conjunto para atender às necessidades específicas de cada tipo de dado,
seguindo o princípio de "escolher o banco certo para o dado certo".

## Arquitetura do Sistema

O sistema segue uma arquitetura baseada em microsserviços com comunicação via mensageria:

• **S1 (s1-generator):** Serviço gerador com API HTTP para usuários, transações e cotações;

• **S2 (s2-processor):** Serviço de processamento via RabbitMQ consumers;

• **S3 (s3-validator):** Serviço de log e validação de mensagens utilizando Cassandra;

• **Mensageria:** RabbitMQ para comunicação entre serviços;

• **Bancos de Dados:** PostgreSQL (usuários), MongoDB (transações e cotações).

## Funcionalidades Implementadas

### 🔐 Autenticação e Usuários

- Registro de novos usuários com saldo inicial de R$ 1.000,00
- Sistema de login com JWT
- CRUD completo de usuários
- Middleware de autenticação para rotas protegidas

### 💱 Sistema de Câmbio

- Compra e venda de moedas (USD/BRL)
- Geração automática de cotações
- Histórico completo de transações
- Validação de saldo antes das operações

### 📊 Banco de Dados Distribuído

- **PostgreSQL:** Armazena dados dos usuários e seus saldos
- **MongoDB:** Armazena transações de câmbio e cotações
- **RabbitMQ:** Sistema de mensageria para comunicação assíncrona

## API Endpoints

### 🏥 Health Check

```http
GET /health
```

### 👤 Usuários

#### Registro de Usuário

```http
POST /users/register
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@email.com",
  "password": "123456"
}
```

#### Login

```http
POST /users/login
Content-Type: application/json

{
  "email": "joao@email.com",
  "password": "123456"
}
```

#### Buscar Usuário (Autenticado)

```http
GET /users/{id}
Authorization: Bearer {jwt_token}
```

#### Atualizar Usuário (Autenticado)

```http
PATCH /users/{id}
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "name": "João Santos",
  "email": "joao.santos@email.com"
}
```

#### Deletar Usuário (Autenticado)

```http
DELETE /users/{id}
Authorization: Bearer {jwt_token}
```

### 💰 Transações

#### Histórico de Transações (Autenticado)

```http
GET /transactions?limit=50
Authorization: Bearer {jwt_token}
```

#### Comprar Moeda (Autenticado)

```http
POST /transactions/buy
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "amount": 100.0,
  "currency_pair": "USD/BRL"
}
```

#### Vender Moeda (Autenticado)

```http
POST /transactions/sell
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "amount": 50.0,
  "currency_pair": "USD/BRL"
}
```

### 📈 Cotações

#### Gerar Cotações

```http
POST /quotations/generate
```

## Exemplo de Uso Completo

### 1. Registrar um novo usuário

```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Silva",
    "email": "maria@email.com",
    "password": "123456"
  }'
```

### 2. Fazer login

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@email.com",
    "password": "123456"
  }'
```

### 3. Comprar USD (usar o token do login)

```bash
curl -X POST http://localhost:8080/transactions/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_JWT_TOKEN_AQUI" \
  -d '{
    "amount": 100.0,
    "currency_pair": "USD/BRL"
  }'
```

### 4. Ver histórico de transações

```bash
curl -X GET http://localhost:8080/transactions \
  -H "Authorization: Bearer SEU_JWT_TOKEN_AQUI"
```

## Justificativa para Escolha dos Bancos de Dados

### PostgreSQL (RDB)

**Dados armazenados**: Informações dos usuários (cadastro, dados pessoais, credenciais, saldo).

**Justificativa:**

• O PostgreSQL é um banco de dados relacional robusto, ideal para armazenar dados estruturados que exigem integridade referencial;

• Os dados de usuários possuem um esquema bem definido e relacionamentos claros;

• Oferece suporte a transações ACID, essenciais para operações financeiras;

• Permite implementar restrições de segurança e validações no nível do banco;

• Facilita consultas complexas envolvendo dados de usuários;

• Excelente suporte para backups e recuperação de dados;

• Alta confiabilidade e maturidade para dados críticos do negócio.

### MongoDB (DB1)

**Dados armazenados**: Transações de câmbio, cotações e histórico de preços.

**Justificativa:**

• O MongoDB é um banco de dados NoSQL orientado a documentos;

• Ideal para armazenar dados de transações que podem variar em estrutura ao longo do tempo;

• Permite consultas rápidas e eficientes para recuperar histórico de transações;

• Oferece boa performance para operações de leitura frequentes (consultas de transações);

• Facilita o armazenamento de dados semi-estruturados como informações de cotações;

• Escalabilidade horizontal para lidar com grandes volumes de dados históricos.

### Cassandra (DB2)

**Dados armazenados**: Logs de validação, auditoria de transações e histórico de eventos.

**Justificativa:**

• O Cassandra é um banco de dados NoSQL distribuído projetado para alta disponibilidade e escalabilidade;

• Ideal para armazenar grandes volumes de dados de log e auditoria;

• Excelente performance para operações de escrita intensiva (logging);

• Oferece replicação automática e tolerância a falhas;

• Permite consultas eficientes por timestamp para análise temporal;

• Estrutura otimizada para dados de séries temporais e eventos;

• Escalabilidade linear para lidar com crescimento de dados de auditoria.

## Configuração e Execução

### Pré-requisitos

1. **Go 1.21+**
2. **Docker e Docker Compose**
3. **PostgreSQL** (via Docker)
4. **MongoDB** (via Docker)
5. **RabbitMQ** (via Docker)

### Variáveis de Ambiente

#### S1-Generator (.env)

```env
# Servidor HTTP
PORT=8080

# JWT
JWT_SECRET_KEY=seu_jwt_secret_super_seguro

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=financial

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=financial

# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/
```

#### S2-Processor (.env)

```env
# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=financial

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=financial
```

### Executando o Sistema

#### 1. Iniciar os Bancos de Dados

```bash
docker-compose up -d
```

#### 2. Configurar o Schema do Cassandra

```bash
docker exec -i $(docker ps -q --filter "ancestor=cassandra:latest") cqlsh < setup_cassandra.cql
```

#### 3. Executar as Migrações (se necessário)

```bash
cd services/s2-processor/internal/infra/database/migrations
tern migrate
```

#### 4. Iniciar o S2-Processor (Consumidores)

```bash
cd services/s2-processor
go run cmd/main.go
```

#### 5. Iniciar o S3-Validator (Logs e Auditoria)

```bash
cd services/s3-validator
go run cmd/main.go
```

#### 6. Iniciar o S1-Generator (API HTTP)

```bash
cd services/s1-generator
go run cmd/main.go
```

O sistema estará disponível em `http://localhost:8080`

## Fluxo de Dados

1. **Usuário se registra** → S1 envia evento via RabbitMQ → S2 processa e salva no PostgreSQL
2. **Usuário faz login** → S1 autentica direto no PostgreSQL e retorna JWT
3. **Usuário compra/vende moeda** → S1 envia evento via RabbitMQ → S2 processa e salva no MongoDB
4. **Usuário consulta transações** → S1 busca diretamente no MongoDB
5. **Sistema gera cotações** → S1 envia evento via RabbitMQ → S2 processa e salva no MongoDB

## Gerenciamento do Banco de Dados

### Configuração do PostgreSQL com Docker

O projeto utiliza PostgreSQL containerizado via Docker para maior portabilidade e facilidade de configuração.

#### Iniciando o Banco de Dados

1. **Iniciar os Serviços:**

```bash
docker-compose up -d
```

2. **Verificar Status:**

```bash
docker-compose ps
```

3. **Acessar o PostgreSQL:**

```bash
docker exec -it postgres psql -U postgres
```

4. **Parar os Serviços:**

```bash
docker-compose down
```

### Migrações com Tern

O projeto utiliza [Tern](https://github.com/jackc/tern) para gerenciamento de migrações do PostgreSQL.

#### Instalação do Tern

```bash
go install github.com/jackc/tern@latest
```

#### Configuração

O arquivo de configuração `tern.conf` está localizado em `services/s2-processor/internal/infra/database/migrations/` e contém as configurações de conexão com o banco de dados:

```conf
[database]
host = localhost
port = 5432
database = your_database
user = your_user
password = your_password
```

#### Comandos Básicos

1. **Criar Nova Migração:**

```bash
cd services/s2-processor/internal/infra/database/migrations
tern new nome_da_migracao
```

Isso criará um novo arquivo de migração com o formato: `YYYYMMDDHHMMSS_nome_da_migracao.sql`

2. **Executar Migrações Pendentes:**

```bash
cd services/s2-processor/internal/infra/database/migrations
tern migrate
```

3. **Verificar Status das Migrações:**

```bash
tern status
```

4. **Reverter Última Migração:**

```bash
tern migrate --destination -1
```

### Geração de Código com SQLC

O projeto utiliza [SQLC](https://sqlc.dev/) para gerar código Go type-safe a partir das queries SQL.

#### Instalação do SQLC

```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

#### Estrutura dos Arquivos SQLC

- `services/s2-processor/internal/infra/database/queries/`: Contém os arquivos .sql com as queries
- `services/s2-processor/internal/infra/database/sqlc.yaml`: Arquivo de configuração do SQLC

#### Comandos Básicos

1. **Gerar Código:**

```bash
cd services/s2-processor
sqlc generate
```

2. **Verificar SQL:**

```bash
sqlc vet
```

#### Criando Novas Queries

1. Adicione suas queries no diretório `queries/` com a sintaxe SQLC:

```sql
-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
```

2. Execute `sqlc generate` para gerar o código Go correspondente

### Boas Práticas

1. **Migrações:**

   - Sempre teste as migrações em um ambiente de desenvolvimento
   - Mantenha as migrações idempotentes
   - Inclua sempre instruções de rollback (down migrations)
   - Use nomes descritivos para os arquivos de migração

2. **SQLC:**

   - Mantenha as queries simples e focadas
   - Use comentários para documentar queries complexas
   - Verifique os tipos gerados no Go
   - Utilize os tipos corretos para cada coluna

3. **Versionamento:**
   - Mantenha as migrações no controle de versão
   - Documente alterações significativas no esquema
   - Nunca modifique migrações já aplicadas em produção

### Configuração do Cassandra

O projeto utiliza Cassandra para armazenar logs de validação e auditoria das transações no serviço s3-validator.

#### Iniciando o Cassandra

O Cassandra é iniciado automaticamente junto com os outros serviços:

```bash
docker-compose up -d
```

#### Configurando o Schema

Após iniciar o Cassandra, é necessário criar o keyspace e as tabelas necessárias:

1. **Executar o Script de Setup:**

```bash
docker exec -i $(docker ps -q --filter "ancestor=cassandra:latest") cqlsh < setup_cassandra.cql
```

2. **Verificar se o Keyspace foi Criado:**

```bash
docker exec -it $(docker ps -q --filter "ancestor=cassandra:latest") cqlsh -e "DESCRIBE KEYSPACES;"
```

3. **Verificar as Tabelas:**

```bash
docker exec -it $(docker ps -q --filter "ancestor=cassandra:latest") cqlsh -e "USE financial; DESCRIBE TABLES;"
```

#### Estrutura do Schema

O arquivo `setup_cassandra.cql` cria:

- **Keyspace `financial`**: Com replicação simples para desenvolvimento
- **Tabela `transactions`**: Para armazenar cópias das transações para auditoria
- **Tabela `validation_logs`**: Para logs de validação e eventos do sistema
- **Índices**: Para otimizar consultas por `user_id`, `status` e `transaction_id`

#### Comandos Úteis do Cassandra

1. **Acessar o CQL Shell:**

```bash
docker exec -it $(docker ps -q --filter "ancestor=cassandra:latest") cqlsh
```

2. **Consultar Dados:**

```cql
USE financial;
SELECT * FROM transactions LIMIT 10;
SELECT * FROM validation_logs LIMIT 10;
```

3. **Limpar Dados (Desenvolvimento):**

```cql
USE financial;
TRUNCATE transactions;
TRUNCATE validation_logs;
```

## Observações de Segurança

- **Nunca compartilhe** arquivos `.env` ou os commite no repositório
- Use **senhas fortes** para JWT_SECRET_KEY em produção
- Configure **CORS** adequadamente para produção
- Implemente **rate limiting** nas APIs públicas
- Use **HTTPS** em produção
