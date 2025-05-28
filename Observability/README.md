# Desafio Observability

## Sobre o Projeto
Sistema distribuído em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin). O sistema implementa OpenTelemetry (OTEL) e Zipkin para observabilidade completa.

## Arquitetura

O sistema é composto por dois microserviços:

- **Serviço A (inputvalidator)**: Valida o formato do CEP e encaminha para o Serviço B
- **Serviço B (orchestrator)**: Busca informações do CEP e temperaturas, retornando os dados formatados

### Componentes de Observabilidade
- **OTEL Collector**: Coleta telemetria de ambos os serviços e exporta para outros sistemas
- **Jaeger/Zipkin**: Armazenam e visualizam traces distribuídos
- **Prometheus**: Coleta e armazena métricas

## Requisitos

- Docker e Docker Compose
- Conexão com internet (para acessar APIs externas)

## Executando o Projeto

1. Clone o repositório
2. Navegue até a pasta do projeto
3. Inicie todos os serviços usando Docker Compose:

```bash
docker compose up -d
```

## Testando a Aplicação

### Endpoints Disponíveis

- **Serviço A**: `http://localhost:28080/temperature`
  - Aceita requisições POST com um CEP para validação

### Exemplos de Uso

#### Teste com CEP válido:

```bash
curl -X POST http://localhost:28080/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "22021001"}'
```

Resposta esperada (HTTP 200):
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

#### Teste com CEP inválido (formato incorreto):

```bash
curl -X POST http://localhost:28080/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```

Resposta esperada (HTTP 422):
```json
{
  "error": "invalid zipcode"
}
```

#### Teste com CEP inexistente:

```bash
curl -X POST http://localhost:28080/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "99999999"}'
```

Resposta esperada (HTTP 404):
```json
{
  "error": "can not find zipcode"
}
```

## Acessando as Ferramentas de Observabilidade

### Jaeger (Tracing)
- **URL**: http://localhost:16686
- Use a interface para explorar os traces gerados pelos serviços
- Filtre por serviço "serviceA" ou "serviceB"

### Zipkin (Tracing)
- **URL**: http://localhost:9411/zipkin/
- Pesquise traces por serviço, tags ou intervalo de tempo

### Prometheus (Métricas)
- **URL**: http://localhost:9090
- Acesse métricas de sistema e personalizadas

## Desafio Observability
Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

 

Requisitos - Serviço A (responsável pelo input):

    O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
    O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
        Caso seja válido, será encaminhado para o Serviço B via HTTP
        Caso não seja válido, deve retornar:
            Código HTTP: 422
            Mensagem: invalid zipcode

Requisitos - Serviço B (responsável pela orquestração):

    O sistema deve receber um CEP válido de 8 digitos
    O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
    O sistema deve responder adequadamente nos seguintes cenários:
        Em caso de sucesso:
            Código HTTP: 200
            Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
        Em caso de falha, caso o CEP não seja válido (com formato correto):
            Código HTTP: 422
            Mensagem: invalid zipcode
        ​​​Em caso de falha, caso o CEP não seja encontrado:
            Código HTTP: 404
            Mensagem: can not find zipcode

Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

    Implementar tracing distribuído entre Serviço A - Serviço B
    Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

Dicas:

    Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
    Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
    Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
    Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
        Sendo F = Fahrenheit
        Sendo C = Celsius
        Sendo K = Kelvin
    Para dúvidas da implementação do OTEL: https://opentelemetry.io/docs/languages/go/getting-started/
    Para implementação de spans: https://opentelemetry.io/docs/languages/go/instrumentation/#creating-spans
    Você precisará utilizar um serviço de collector do OTEL https://opentelemetry.io/docs/collector/quick-start/
    Para mais informações sobre Zipkin: https://zipkin.io/

Entrega:

    O código-fonte completo da implementação.
    Documentação explicando como rodar o projeto em ambiente dev.
    Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
