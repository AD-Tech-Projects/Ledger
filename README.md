# Ledger

Sistema financeiro em Go baseado em um **ledger (livro-razão) de partidas dobradas**, com carteiras (wallets) multiativo (moedas fiduciárias, criptoativos e ativos sem lastro), operações financeiras (transferências, depósitos, saques, rendimentos, tarifas, reservas e retenções) e foco em consistência, auditabilidade e rastreabilidade.

> A finalidade do projeto, os conceitos de domínio e o roadmap estão em [`CONTEXT.md`](./CONTEXT.md).

## Stack

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.27.2 |
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Banco de dados | PostgreSQL |
| ORM | [GORM](https://gorm.io) |
| Migrações | [Atlas](https://atlasgo.io) + `atlas-provider-gorm` |
| Logs | [`Go-Libs/logger`](https://github.com/AD-Tech-Projects/Go-Libs) |
| Cache | Redis (provisionado no `docker-compose.yml`, ainda não utilizado) |

## Pré-requisitos

- [Go](https://go.dev/dl/) **1.27.2** ou superior
- [Docker](https://docs.docker.com/get-docker/) e Docker Compose
- [Atlas CLI](https://atlasgo.io/getting-started#installation): `curl -sSf https://atlasgo.sh | sh`
- `make`

## Como rodar

### 1. Variáveis de ambiente

Copie o arquivo de exemplo e ajuste se necessário:

```bash
cp .env.local .env
```

O `Makefile` e o `docker-compose.yml` leem o `.env`. Se o arquivo não existir, a aplicação usa as variáveis de ambiente do sistema.

### 2. Subir a infraestrutura

```bash
docker compose up -d database   # PostgreSQL
# docker compose up -d          # PostgreSQL + Redis
```

### 3. Instalar as dependências

```bash
make tidy   # go mod tidy
```

A biblioteca compartilhada `github.com/AD-Tech-Projects/Go-Libs` é baixada como um módulo Go comum. Para desenvolver a lib localmente, clone-a em `modules/Go-Libs` e descomente a diretiva `replace` no `go.mod`.

### 4. Criar o schema do banco (migrações)

O schema é gerado pelo Atlas a partir dos modelos GORM listados em `loader/main.go`, e as migrações versionadas em `migrations/` estão sincronizadas com esses modelos.

```bash
make migrate-apply   # aplica as migrações no banco do .env
```

O Atlas exige um banco dedicado ao Ledger. Ele recusa aplicar migrações em um banco que já tenha tabelas de outro projeto.

Quando os modelos mudarem:

```bash
make migrate-create name=<descricao>   # atlas migrate diff --env gorm (usa um container postgres:15 temporário; o Docker precisa estar rodando)
make migrate-hash                      # recalcula migrations/atlas.sum após editar uma migração manualmente
```

### 5. Cadastrar os ativos

A criação de carteiras busca o ativo pelo código (`BRL`, `USD`, `BTC`...), mas ainda não existe seed nem endpoint de ativos. Segue um exemplo de inserção manual no banco:

```sql
INSERT INTO assets (id, created_at, updated_at, name, type, status, code, decimal_places) VALUES
(gen_random_uuid(), now(), now(), 'Brazilian Real', 'FIAT',     'ACTIVE', 'BRL',    2),
(gen_random_uuid(), now(), now(), 'US Dollar',      'FIAT',     'ACTIVE', 'USD',    2),
(gen_random_uuid(), now(), now(), 'Euro',           'FIAT',     'ACTIVE', 'EUR',    2),
(gen_random_uuid(), now(), now(), 'Bitcoin',        'CRYPTO',   'ACTIVE', 'BTC',    8),
(gen_random_uuid(), now(), now(), 'Ethereum',       'CRYPTO',   'ACTIVE', 'ETH',    18),
(gen_random_uuid(), now(), now(), 'Tether',         'CRYPTO',   'ACTIVE', 'USDT',   6),
(gen_random_uuid(), now(), now(), 'Points',         'UNBACKED', 'ACTIVE', 'POINTS', 0);
```

### 6. Executar a aplicação

```bash
make run            # go run ./main.go
# ou
make build && ./out/app
```

A API sobe em `http://localhost:${APP_PORT}` (padrão `3000`):

```bash
curl http://localhost:3000/health   # OK
```

Para depurar no VS Code, use a configuração **Go: Debug** em `.vscode/launch.json`.

## Testes

```bash
make test            # go test -v ./...
make test-coverage   # gera coverage.out e abre o relatório HTML
```

Os testes atuais são unitários e não usam banco. Eles cobrem:

- as operações de saldo da `Wallet`, inclusive os erros;
- o efeito de cada tipo de transação, inclusive as que envolvem uma carteira só ou carteiras de tipos diferentes;
- a reversão de cada tipo, que precisa restaurar os saldos iniciais;
- a preservação dos saldos quando falta saldo;
- a validação da entrada e das carteiras de cada tipo de transação.

## API

A documentação das rotas fica em [`packages/api/doc.go`](./packages/api/doc.go), no formato godoc. Ela traz exemplos de chamadas, os campos aceitos, as carteiras e o efeito de cada tipo de transação, as reversões e os códigos de erro.

```bash
go doc ./packages/api
```

## Estrutura do projeto

```text
.
├── main.go                         # Carrega o .env, conecta o banco, monta handlers e sobe o servidor
├── loader/main.go                  # Exporta o schema dos modelos GORM para o Atlas
├── migrations/                     # Migrações SQL versionadas (Atlas)
├── packages/
│   ├── api/                        # Rotas Gin, middleware de log HTTP e recovery (documentadas em api/doc.go)
│   ├── handlers/                   # Entrada/saída HTTP, validação do payload e mapeamento de erros
│   ├── services/                   # Casos de uso, validações, operações de saldo e controle transacional
│   ├── domain/                     # Entidades e enums (documentados em domain/doc.go)
│   └── infra/database/
│       ├── connection.go           # Conexão GORM/PostgreSQL
│       ├── models/                 # Modelos de persistência e conversão domínio <-> banco
│       └── repositories/           # Repositório genérico (CRUD, filtros, joins, transações)
├── docker-compose.yml              # PostgreSQL e Redis
├── atlas.hcl                       # Configuração do Atlas
└── Makefile
```

## Problemas conhecidos

- **Sem seed de ativos**: o cadastro é manual (ver passo 5).
- **Variáveis `DB_*` exportadas no shell têm prioridade sobre o `.env`**: isso vale para o `docker compose` e também para a aplicação (`godotenv` não sobrescreve variáveis já definidas). Só o `make migrate-*` lê o `.env` diretamente. Confira `echo $DB_NAME` antes de subir o banco ou rodar o app.
