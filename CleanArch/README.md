# CleanArch 

### Desafio

Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando `docker-compose up` tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

---

### Passos para Execução

1. **Clone o repositório**:
   ```bash
   git clone https://github.com/Leandroschwab/full-cycle-go.git
   cd CleanArch
   ```

2. **Edite o arquivo `.env`**:
   - Configure as variáveis de ambiente necessárias no arquivo `.env`.

3. **Execute o Docker Compose**:
   ```bash
   docker-compose up -d
   ```

---

### Testando a Aplicação

#### GRPC

1. **Inicie o Evans em modo REPL**:
   ```bash
   evans -r repl --host localhost --port 50051
   ```

2. **Chame os métodos `CreateOrder` ou `ListOrders`**:
   ```bash
   package pb
   service orderService
   call CreateOrder
   call ListOrders
   ```

---

#### GraphQL

1. **Acesse o playground do GraphQL**:
   - Abra o navegador e acesse  pela porta definida no `.env` (padrão `8080`):
    ```bash
   `http://localhost:8080`.
    ```

2. **Exemplo de Mutation**:
   ```graphql
   mutation createOrder { 
       createOrder(input: {id: "ccc", price: 100, tax: 2.0}) { 
           id 
           price 
           tax 
           finalPrice
       }
   }
   ```

3. **Exemplo de Query**:
   ```graphql
   query {
       listOrders {
           orders {
               id
               price
               tax
               finalPrice
           }
       }
   }
   ```

---

#### REST API

1. **Use o arquivo HTTP para testar a API REST**:
   - Arquivo para criar pedidos:
     ```bash
     CleanArch/api/create_order.http
     ```
   - Arquivo para listar pedidos:
     ```bash
     CleanArch/api/list_orders.http
     ```

2. **Execute as requisições**:
   - Use um cliente HTTP como o VS Code REST Client ou o Postman para executar as requisições. Não esqueça de alterar a URL base para o servidor se necessário.

---

### Notas Adicionais

- Certifique-se de que as portas configuradas no `.env` estão disponíveis no seu sistema.
- Para reiniciar os serviços, use:
  ```bash
  docker-compose down && docker-compose up -d
  ```
- Para limpar o banco de dados, exclua os volumes do Docker:
  ```bash
  docker-compose down -v
  ```