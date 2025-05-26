# Full Cycle Auction GoExpert

## Objetivo

Este projeto implementa um sistema de leilão com fechamento automático após um tempo definido via variável de ambiente. O leilão é criado, permanece aberto por um período configurável e é fechado automaticamente por uma goroutine.

## Funcionalidades

- Criação de leilão com tempo de duração configurável (`AUCTION_INTERVAL`)
- Fechamento automático do leilão após o tempo definido
- API REST para gerenciamento de leilões, lances e usuários
- Teste automatizado para validar o fechamento automático do leilão

## Como rodar o projeto em ambiente de desenvolvimento

### 1. Pré-requisitos

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [MongoDB](https://www.mongodb.com/) (ou utilize o serviço via Docker)

### 2. Configuração das variáveis de ambiente

Crie um arquivo `.env` em `cmd/auction/.env` com o seguinte conteúdo:

```
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=1m

MONGO_INITDB_ROOT_USERNAME: admin
MONGO_INITDB_ROOT_PASSWORD: admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

- `AUCTION_INTERVAL`: define o tempo de duração do leilão (exemplo: `30s`, `1m`, `5m`).

### 3. Subindo o ambiente com Docker Compose

No diretório raiz do projeto, execute:

```sh
docker-compose up --build
```

Isso irá subir o MongoDB e a aplicação Go.

### 4. Utilizando a API

Veja exemplos de requisições no arquivo [`api/http.api`](api/http.api) ou utilize o VSCode REST Client.

### 5. Rodando os testes

Para rodar os testes automatizados (incluindo o teste de fechamento automático do leilão):

```sh
go test ./internal/infra/database/auction
```

Certifique-se de que o MongoDB está rodando localmente ou ajuste a URI de conexão conforme necessário.

## Observações

- O fechamento automático do leilão é realizado por uma goroutine iniciada na criação do leilão.
- O tempo de expiração é definido pela variável de ambiente `AUCTION_INTERVAL`.
- O status do leilão é atualizado para `Completed` automaticamente após o tempo definido.

## Objetivo:

Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: clique para acessar o repositório.

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

Você deverá desenvolver:

    Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
    Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
    Um teste para validar se o fechamento está acontecendo de forma automatizada;


Dicas:

    Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
    Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
    Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
    Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;
     

Entrega:

    O código-fonte completo da implementação.
    Documentação explicando como rodar o projeto em ambiente dev.
    Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.
