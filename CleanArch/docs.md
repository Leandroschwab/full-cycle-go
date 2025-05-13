## Environment Setup

To run MySQL and RabbitMQ on another machine with Docker, update the environment variables accordingly:  
```CleanArch/cmd/ordersystem/.env```


## Update .env file

   Precisei atualizar o arquivo .env trocando o endereço do mysql e rabbitmq para o ip da minha máquina docker.


## important Notes

1. Install Wire:
   ```bash
   go install github.com/google/wire/cmd/wire@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
    Generate dependency injection code:
   ```bash
   wire
   ```
3. Install Evans (gRPC client):
   ```bash
   go install github.com/ktr0731/evans@latest
   ```
4. Install Protocol Buffers and gRPC tools:
   ```bash
   sudo apt install -y protobuf-compiler
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

5. Install gqlgen:
   ```bash
   go get github.com/99designs/gqlgen
   go run github.com/99designs/gqlgen generate
   ```

## Running the Application

Run the application with the following command:
```bash
go run main.go wire_gen.go
```

## Usage

### gRPC

1. Start Evans in REPL mode:
   ```bash
   evans -r repl
   ```
2. Call the `CreateOrder` or `ListOrders` methods:
   ```bash
   package pb
   service orderService
   call ListOrders
   call CreateOrder or call ListOrders
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
3. Example query:
   ```graphql
   query ListOrders {
       listOrders {
           orders {
               id
               Price
               Tax
               FinalPrice
           }
       }
   }
   ```

### REST API

1. Use the provided HTTP file to test the REST API:
   ```bash
   CleanArch/api/create_order.http
   ```
   ```bash
   CleanArch/api/list_orders.http
   ```
2. Open the file in an HTTP client (e.g., VS Code REST Client or Postman) and execute the requests.

## RabbitMQ

- URL: `<docker Machine ip>:15672`
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


## Progress Track

- Created `CleanArch/internal/usecase/list_orders.go`
- Added `FindAll()` to `CleanArch/internal/infra/database/order_repository.go`
- Added `FindAll()` to `CleanArch/internal/entity/interface`
- UseCase created
- WebServer `GET /order` working
- gRPC `ListOrders` working
- GraphQL `ListOrders` working

## Additional Tools

1. Install gqlgen:
   ```bash
   go get github.com/99designs/gqlgen
   go run github.com/99designs/gqlgen generate
   ```