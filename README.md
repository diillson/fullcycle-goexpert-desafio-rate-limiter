# Rate Limiter

## Descrição
Rate limiter implementado em Go que limita requisições por IP e token de acesso.

## Configuração
Configure através de variáveis de ambiente ou arquivo .env:

```env
REDIS_URL=redis://redis:6379/0
IP_MAX_REQUESTS=5           # Máximo de requisições/segundo por IP
IP_BLOCK_DURATION=5         # Tempo de bloqueio (minutos) para IP
TOKEN_MAX_REQUESTS=10       # Máximo de requisições/segundo por token
TOKEN_BLOCK_DURATION=5      # Tempo de bloqueio (minutos) para token
```

## Instalação e Execução

Com Docker:
```bash
docker-compose up -d
```

## Local:
```bash
go mod download
go run main.go
```

## Testes

# Testes unitários
```bash
go test ./... -short
```

# Todos os testes (necessita Redis)
```bash 
go test ./...
```

## Uso

O rate limiter pode ser usado de duas formas:

    Por IP (automático)
    Por Token (via header API_KEY)

Exemplo de requisição com token:

curl -H "API_KEY: seu-token" http://localhost:8080/api

# Uso do Makefile

## Compilar a aplicação
make build

## Executar a aplicação localmente
make run

## Executar todos os testes
make test

## Executar testes sem Redis (modo short)
make test-short

## Iniciar a aplicação com Docker
make docker-up

## Parar a aplicação Docker
make docker-down

## Executar testes no ambiente Docker
make docker-test
