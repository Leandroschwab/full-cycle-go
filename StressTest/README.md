# Ferramenta de Teste de Carga CLI em Go

Esta aplicação permite realizar testes de carga em serviços web, simulando múltiplas requisições simultâneas para avaliar desempenho e estabilidade.

## Instalação

### Compilação local
```bash
go build -o stress-test .
```

### Docker
```bash
docker build -t stress-test .
```

## Como usar

Execute o binário ou o container Docker com os parâmetros obrigatórios:

- `--url`: URL do serviço a ser testado
- `--requests`: Número total de requisições
- `--concurrency`: Número de chamadas simultâneas

### Exemplo com Docker

```bash
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10


docker run stress-test --url=http://httpstat.us/random/200,201,500-504 --requests=100 --concurrency=10

```

## Relatório

Ao final do teste, será apresentado um relatório contendo:
- Tempo total de execução
- Total de requisições realizadas
- Quantidade de respostas HTTP 200
- Distribuição dos demais códigos de status HTTP

# Desafio: 

Objetivo: Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.


O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

Entrada de Parâmetros via CLI:

--url: URL do serviço a ser testado.
--requests: Número total de requests.
--concurrency: Número de chamadas simultâneas.


Execução do Teste:

    Realizar requests HTTP para a URL especificada.
    Distribuir os requests de acordo com o nível de concorrência definido.
    Garantir que o número total de requests seja cumprido.

Geração de Relatório:

    Apresentar um relatório ao final dos testes contendo:
        Tempo total gasto na execução
        Quantidade total de requests realizados.
        Quantidade de requests com status HTTP 200.
        Distribuição de outros códigos de status HTTP (como 404, 500, etc.).

    Execução da aplicação:

    Poderemos utilizar essa aplicação fazendo uma chamada via docker. Ex:
        docker run <sua imagem docker> —url=http://google.com —requests=1000 —concurrency=10
