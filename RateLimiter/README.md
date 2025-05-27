# Rate Limiter

## Visão Geral

O Rate Limiter é projetado para controlar o fluxo de requisições para um serviço web com base no endereço IP ou token de acesso. Ele previne o abuso limitando quantas requisições podem ser feitas em um determinado período de tempo.

## Arquitetura

- **Config**: Carrega configurações a partir de variáveis de ambiente ou arquivo .env
- **Storage**: Interface para persistência com implementação para Redis
- **Limiter**: Lógica principal para decisões de limitação de taxa
- **Middleware**: Middleware HTTP para integração com servidores web

## Configuração

O Rate Limiter pode ser configurado usando variáveis de ambiente ou editando o arquivo .env:

| Variável | Descrição | Padrão |
|----------|-------------|---------|
| IP_RATE_LIMIT | Máximo de requisições permitidas por IP | 10 |
| TOKEN_RATE_LIMIT | Máximo de requisições permitidas por token | 100 |
| BLOCK_DURATION | Duração do bloqueio em segundos | 300 |
| REDIS_URL | URL para conexão com Redis | redis:6379 |

## Como Funciona

1. **Processamento de Requisições**:
   - Quando uma requisição chega, o middleware extrai o IP do cliente e o token de API opcional
   - O limitador verifica se o identificador (IP ou token) está atualmente bloqueado
   - Se não estiver bloqueado, incrementa um contador para esse identificador
   - Se o contador exceder o limite configurado, o identificador é bloqueado pela duração configurada

2. **Regras de Precedência**:
   - Limites baseados em token têm precedência sobre limites baseados em IP
   - Se uma requisição tiver um cabeçalho API_KEY válido, o limite de token é usado
   - Caso contrário, o limite de IP é usado

3. **Comportamento de Bloqueio**:
   - Uma vez que um identificador é bloqueado, todas as requisições desse IP ou usando esse token receberão um erro 429
   - O bloqueio expirará após a duração de bloqueio configurada

## Testes

O limitador de taxa foi testado sob várias condições:

1. **Testes Unitários**: Verificam se a lógica principal funciona corretamente
2. **Testes de Concorrência**: Garantem que o limitador funcione sob alta carga
3. **Testes de Duração**: Confirmam que os períodos de bloqueio funcionam conforme esperado

### Scripts Shell para Testes

Incluímos vários scripts shell no diretório `scripts` para ajudar a testar o limitador de taxa:

1. **test_ip_limit.sh**: Testa a limitação baseada em IP enviando múltiplas requisições e verificando quando elas começam a ser bloqueadas
2. **test_token_limit.sh**: Testa a limitação baseada em token para verificar se os limites de token são aplicados corretamente
3. **test_concurrent_limits.sh**: Simula tráfego alto enviando requisições concorrentes de múltiplos clientes
4. **test_blocking_duration.sh**: Verifica se IPs/tokens bloqueados permanecem bloqueados pela duração configurada

Para executar estes scripts:

```bash
# Execute os testes individualmente
./scripts/test_ip_limit.sh
./scripts/test_token_limit.sh
./scripts/test_concurrent_limits.sh
./scripts/test_blocking_duration.sh
```

Estes scripts fornecem testes práticos do limitador de taxa sob diferentes condições e são úteis para validar a configuração.

## Implantação com Docker

O projeto inclui configurações Docker e docker-compose para fácil implantação:

```bash
# Iniciar o serviço
docker-compose up -d

# Testar com curl
curl -H "API_KEY: token-teste" http://localhost:8080/
```

# Desafio 
Objetivo: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

Descrição: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

    Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
    Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
        API_KEY: <TOKEN>
    As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

Requisitos:

    O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
    O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
    O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
    As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
    Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
    O sistema deve responder adequadamente quando o limite é excedido:
        Código HTTP: 429
        Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
    Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
    Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
    A lógica do limiter deve estar separada do middleware.

Exemplos:

    Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
    Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
    Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

Dicas:

    Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

Entrega:

    O código-fonte completo da implementação.
    Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
    Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
    Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
    O servidor web deve responder na porta 8080.
