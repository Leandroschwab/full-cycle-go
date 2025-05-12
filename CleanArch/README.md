
## Environment Setup

To run MySQL and RabbitMQ on another machine with Docker, update the environment variables accordingly. ```CleanArch/cmd/ordersystem/.env ```
## Installation

1. Install Wire:
   ```bash
   go install github.com/google/wire/cmd/wire@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
2. Generate dependency injection code:
   ```bash
   wire
   ```

3. Install Evans (gRPC client):
   ```bash
   go install github.com/ktr0731/evans@latest
   ```

## Usage

### gRPC

1. Start Evans in REPL mode:
   ```bash
   evans -r repl
   ```
2. Call the `CreateOrder` method:
   ```bash
   call CreateOrder
   ```

### GraphQL

1. Forward port `8080` and open it in your browser.
2. Example mutation:
   ```graphql
   mutation createOrder { 
       createOrder(input: {id: "ccc", Price: 100, Tax: 2.0}) { 
           id 
           Price 
           Tax 
           FinalPrice
       }
   }
   ```

### REST API

1. Use the provided HTTP file to test the REST API:
   ```bash
   CleanArch/api/create_order.http
   ```
2. Open the file in an HTTP client (e.g., VS Code REST Client or Postman) and execute the requests.

## Running the Application

Run the application with the following command:
```bash
go run main.go wire_gen.go
```

## RabbitMQ

- URL: `172.20.20.15:15672`
- User: `guest`
- Password: `guest`

### Desafio

Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.



Progress track
Created CleanArch/internal/usecase/list_orders.go
added "FindAll()" CleanArch/internal/infra/database/order_repository.go
added "FindAll()" CleanArch/internal/entity/interface.

UseCase created.

WebServer GET /order Working


go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

export PATH=$PATH:$(go env GOPATH)/bin
protoc --go_out=. --go-grpc_out=. protofiles/order.proto
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto


wget https://github.com/protocolbuffers/protobuf/releases/download/v21.9/protoc-21.9-linux-x86_64.zip

unzip protoc-21.9-linux-x86_64.zip -d protoc-21.9
sudo mv protoc-21.9/bin/protoc /usr/local/bin/
sudo mv protoc-21.9/include/* /usr/local/include/
protoc --version

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0