
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