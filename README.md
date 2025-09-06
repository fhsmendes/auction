# Auction

**Desafio fullcycle - Leilão (Auction)**

## Objetivo
Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: clique para acessar o repositório.

## Descrição
Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

## Você deverá desenvolver:

- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

## Dicas:

- Concentre-se na no arquivo `internal/infra/database/auction/create_auction.go`, você deverá implementar a solução nesse arquivo;
- Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
- Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
- Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;

## Entrega:

- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

# Como rodar o projeto

## Pré-requisitos

- Go 1.24 ou superior
- Docker
- Docker Compose

## Instalação

### 1. Clonar o repositório

```bash
git clone https://github.com/fhsmendes/auction.git
cd auction
```

### 2. Atualizar dependências

```bash
go mod tidy
```

### 3. Configurar variáveis de ambiente

Configure as seguintes variáveis no arquivo `cmd/auction/.env`:

```env
# Configurações de batch para lances
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4

# Intervalo para fechamento automático dos leilões
AUCTION_INTERVAL=20s

# Configuração do MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

## Executando a aplicação

### 1. Rodar o Docker Compose

```bash
docker-compose up -d
```

### 2. Verificar se a aplicação está funcionando

```bash
# Teste básico - deve retornar lista vazia se não houver leilões
curl -X GET "http://localhost:8080/auction?status=0"
```

A aplicação estará disponível em: `http://localhost:8080`

**Importante:** O leilão será automaticamente fechado após o tempo definido em `AUCTION_INTERVAL` (padrão: 20 segundos).

# Endpoints da API

A aplicação disponibiliza os seguintes endpoints para interação com o sistema de leilões:

## Leilões (Auctions)

### 1. Criar Leilão
**POST** `/auction`

Cria um novo leilão no sistema.

```bash
curl -X POST http://localhost:8080/auction \
    -H "Content-Type: application/json" \
    -d '{
        "product_name": "iPhone 15 Pro",
        "category": "Eletrônicos",
        "description": "iPhone 15 Pro em perfeito estado, sem riscos, com todos os acessórios originais",
        "condition": 1
    }'
```

**Condições disponíveis:**
- `1` = Novo
- `2` = Usado
- `3` = Recondicionado

### 2. Buscar Leilões
**GET** `/auction`

Busca leilões com filtros opcionais.

```bash
# Buscar todos os leilões
curl -X GET http://localhost:8080/auction

# Buscar leilões por status (0=Active, 1=Completed)
curl -X GET "http://localhost:8080/auction?status=0"

# Buscar leilões por categoria
curl -X GET "http://localhost:8080/auction?category=Eletrônicos"

# Buscar leilões por nome do produto
curl -X GET "http://localhost:8080/auction?productName=iPhone"

# Buscar com múltiplos filtros
curl -X GET "http://localhost:8080/auction?status=0&category=Eletrônicos&productName=iPhone"
```

### 3. Buscar Leilão por ID
**GET** `/auction/:auctionId`

Busca um leilão específico pelo seu ID.

```bash
curl -X GET http://localhost:8080/auction/41cc6b10-fad1-4523-9bf5-1729a2cdfb53
```

### 4. Buscar Lance Vencedor
**GET** `/auction/winner/:auctionId`

Busca informações do lance vencedor de um leilão.

```bash
curl -X GET http://localhost:8080/auction/winner/41cc6b10-fad1-4523-9bf5-1729a2cdfb53
```

## Lances (Bids)

### 1. Criar Lance
**POST** `/bid`

Cria um novo lance em um leilão.

```bash
curl -X POST http://localhost:8080/bid \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "b2f9198e-11c3-4cfa-b8bd-a7af65814cc0",
        "auction_id": "41cc6b10-fad1-4523-9bf5-1729a2cdfb53",
        "amount": 1500.00
    }'
```

### 2. Buscar Lances por Leilão
**GET** `/bid/:auctionId`

Busca todos os lances de um leilão específico.

```bash
curl -X GET http://localhost:8080/bid/41cc6b10-fad1-4523-9bf5-1729a2cdfb53
```

## Usuários (Users)

### 1. Buscar Usuário por ID
**GET** `/user/:userId`

Busca informações de um usuário pelo seu ID.

```bash
curl -X GET http://localhost:8080/user/b2f9198e-11c3-4cfa-b8bd-a7af65814cc0
```

# Exemplos de uso

## Fluxo completo: Criar um leilão e fazer lances

### 1. Criar um leilão

```bash
curl -X POST http://localhost:8080/auction \
    -H "Content-Type: application/json" \
    -d '{
        "product_name": "MacBook Pro M3",
        "category": "Informática",
        "description": "MacBook Pro 14 polegadas com chip M3, 16GB RAM, 512GB SSD",
        "condition": 1
    }'
```

### 2. Buscar leilões ativos

```bash
curl -X GET "http://localhost:8080/auction?status=0"
```

### 3. Fazer lances

```bash
# Primeiro lance
curl -X POST http://localhost:8080/bid \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "user-uuid-aqui",
        "auction_id": "auction-uuid-aqui",
        "amount": 8000.00
    }'

# Segundo lance maior
curl -X POST http://localhost:8080/bid \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "outro-user-uuid",
        "auction_id": "auction-uuid-aqui",
        "amount": 8500.00
    }'
```

### 4. Verificar lances do leilão

```bash
curl -X GET http://localhost:8080/bid/auction-uuid-aqui
```

### 5. Aguardar fechamento automático

Aguarde o tempo definido em `AUCTION_INTERVAL` para o fechamento automático do leilão.

### 6. Verificar lance vencedor

```bash
curl -X GET http://localhost:8080/auction/winner/auction-uuid-aqui
```

## Respostas da API

### Exemplo de resposta - Criar leilão
**Status: 201 Created**

### Exemplo de resposta - Buscar leilão
```json
{
    "id": "41cc6b10-fad1-4523-9bf5-1729a2cdfb53",
    "product_name": "iPhone 15 Pro",
    "category": "Eletrônicos",
    "description": "iPhone 15 Pro em perfeito estado",
    "condition": 1,
    "status": 0,
    "timestamp": "2025-09-06T10:30:00Z"
}
```

### Exemplo de resposta - Lance vencedor
```json
{
    "auction": {
        "id": "41cc6b10-fad1-4523-9bf5-1729a2cdfb53",
        "product_name": "iPhone 15 Pro",
        "category": "Eletrônicos",
        "description": "iPhone 15 Pro em perfeito estado",
        "condition": 1,
        "status": 1,
        "timestamp": "2025-09-06T10:30:00Z"
    },
    "bid": {
        "id": "bid-uuid",
        "user_id": "user-uuid",
        "auction_id": "41cc6b10-fad1-4523-9bf5-1729a2cdfb53",
        "amount": 1500.00,
        "timestamp": "2025-09-06T10:35:00Z"
    }
}
```
