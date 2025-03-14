# Projeto Tópicos Avançados de Bancos de Dados

### Integrantes:

#### - João Paulo Paggi Zuanon Dias  RA: 22.222.058-4 
#### - Leandro de Brito Alencar  RA: 22.222.034-5
#### - Mateus Rocha              RA: 22.222.002-2

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

• Facilita consultas complexas envolvendo dados de usuários

**Usando Supabase**
• O Supabase oferece uma interface simplificada para gerenciar o PostgreSQL, facilitando a administração do banco;
• Possui autenticação e autorização integradas, garantindo segurança no acesso aos dados;
• Oferece suporte a WebSockets para atualizações em tempo real, útil para operações financeiras dinâmicas;
• Conta com APIs automáticas para manipulação dos dados, acelerando o desenvolvimento;
• Inclui armazenamento de arquivos e funções serverless, expandindo as possibilidades do projeto.

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
