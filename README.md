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

A escolha deste tema foi desenvolvido com indicação do professor, gostamos da ideia sujerida,
além de não ser um tema muito fora do que já foi visto durante o decorrer do curso de Ciência da Computação.

## Arquitetura do Sistema

O sistema segue uma arquitetura baseada em microsserviços com comunicação via mensageria, conforme o modelo especificado:

• **S1:** Serviço gerador de mensagens (dados de usuários, cotações e transações);

• **S2:** Serviços de processamento e armazenamento;

• **S3:** Serviço de log e validação de mensagens;

• **Mensageria:** Sistema de filas para comunicação entre serviços;

• **Bancos de Dados:** PostgreSQL (RDB), MongoDB (DB1) e Cassandra (DB2).

## Justificativa para Escolha dos Bancos de Dados

### PostgreSQL (RDB)

**Dados armazenados**: Informações dos usuários (cadastro, dados pessoais, credenciais).

**Justificativa:**

• O PostgreSQL é um banco de dados relacional robusto, ideal para armazenar dados estruturados que exigem integridade referencial;

• Os dados de usuários possuem um esquema bem definido e relacionamentos claros;

• Oferece suporte a transações ACID, essenciais para operações financeiras;

• Permite implementar restrições de segurança e validações no nível do banco;

• Facilita consultas complexas envolvendo dados de usuários;

• Containerização com Docker facilita a configuração e manutenção do ambiente;

• Excelente suporte para backups e recuperação de dados;

• Alta confiabilidade e maturidade para dados críticos do negócio.

### MongoDB (DB1)

**Dados armazenados**: Dados da moeda (cotações, variações, histórico de preços).

**Justificativa:**

• O MongoDB é um banco de dados NoSQL orientado a documentos;

• Ideal para armazenar dados de cotação que podem variar em estrutura ao longo do tempo;

• Permite consultas rápidas e eficientes para recuperar histórico de preços;

• Oferece boa performance para operações de leitura frequentes (consultas de cotação);

• Facilita o armazenamento de dados semi-estruturados como informações de mercado;

• Escalabilidade horizontal para lidar com grandes volumes de dados históricos.

### Cassandra (DB2)

**Dados armazenados:** Registros de transações de compra e venda.

**Justificativa:**

• O Cassandra é um banco de dados NoSQL orientado a colunas, projetado para alta disponibilidade e escalabilidade;

• Excelente para operações de escrita intensiva, como o registro de transações de compra e venda;

• Arquitetura distribuída que permite processamento de grande volume de transações;

• Modelo de dados otimizado para consultas por timestamp, ideal para histórico de transações;

• Alta disponibilidade sem ponto único de falha, essencial para um sistema financeiro;

• Excelente desempenho para escritas sequenciais, como logs de transações em ordem cronológica.

## Implementação do Serviço S2

O serviço S2 será implementado como um conjunto de microserviços responsáveis por:

1. **Processador de Dados de Usuário:**

▪ Recebe mensagens relacionadas a usuários da fila de mensageria;

▪ Valida e processa os dados;

▪ Armazena ou recupera informações no PostgreSQL;

▪ Retorna resultados via mensageria;

2. **Processador de Cotações:**

▪ Consome mensagens com dados de cotação de moedas;

▪ Processa, normaliza e enriquece os dados;

▪ Armazena no MongoDB para consultas rápidas;

▪ Publica atualizações de cotação para outros serviços.

3. **Processador de Transações:**

▪ Recebe solicitações de compra/venda;

▪ Valida a operação contra dados do usuário (PostgreSQL);

▪ Verifica cotações atuais (MongoDB);

▪ Registra a transação no Cassandra;

▪ Retorna confirmação ou erro via mensageria.

Cada componente do S2 será implementado como um serviço independente,
permitindo escalabilidade individual conforme a demanda.
A comunicação entre os componentes será realizada exclusivamente via sistema de mensageria,
garantindo baixo acoplamento e alta resiliência.

O serviço S2 implementará padrões como Circuit Breaker para lidar com falhas nos bancos de dados
e Retry Pattern para garantir a consistência eventual das operações.

## Fluxo de Dados

1. S1 gera mensagens com dados fictícios (usuários, cotações, solicitações de transação).

2. As mensagens são enviadas para o sistema de mensageria.

3. S2 consome as mensagens e realiza o processamento adequado.

4. S2 armazena ou recupera dados dos bancos apropriados.

5. S2 retorna resultados via mensageria.

6. S3 registra todas as mensagens enviadas e recebidas para auditoria e validação.

Este fluxo garante que cada tipo de dado seja tratado pelo banco mais adequado às suas características,
otimizando o desempenho e a confiabilidade do sistema.

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

#### Variáveis de Ambiente

Configure as seguintes variáveis no arquivo `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
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

## Executando os Serviços

O sistema foi simplificado para rodar com apenas dois comandos principais, um para cada serviço:

### 1. Iniciando o Serviço Gerador (S1)

Primeiro, configure as variáveis de ambiente no arquivo `.env` do S1:

```env
# Servidor HTTP
PORT=8080

# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/

# JWT (para autenticação)
JWT_SECRET=seu_jwt_secret
```

Então execute:

```bash
cd services/s1-generator
go run cmd/main.go
```

Este comando inicia o serviço S1 que:

- Expõe endpoints HTTP para registro e autenticação de usuários
- Gera cotações automáticas de moedas
- Envia mensagens para o RabbitMQ

### 2. Iniciando o Serviço Processador (S2)

Configure as variáveis de ambiente no arquivo `.env` do S2:

```env
# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=sua_senha
POSTGRES_DB=seu_banco

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=financial

# Cassandra
CASSANDRA_HOSTS=localhost
CASSANDRA_KEYSPACE=seu_keyspace
CASSANDRA_USERNAME=seu_usuario
CASSANDRA_PASSWORD=sua_senha
```

Então execute:

```bash
cd services/s2-processor
go run cmd/main.go
```

Este comando inicia o serviço S2 unificado que:

- Processa mensagens de usuários (PostgreSQL)
- Processa cotações (MongoDB)
- Processa transações (Cassandra)
- Gerencia todos os consumidores de forma centralizada

### Pré-requisitos

Antes de executar os serviços, certifique-se de que:

1. O RabbitMQ está em execução
2. Os bancos de dados estão disponíveis:
   - PostgreSQL
   - MongoDB
   - Cassandra
3. As variáveis de ambiente estão configuradas em cada serviço
4. As migrações do banco de dados foram aplicadas

### Observações

- Mantenha os arquivos `.env` seguros e nunca os compartilhe ou commite no repositório
- Para desenvolvimento local, você pode copiar o arquivo `.env.example` de cada serviço para criar seu próprio `.env`
