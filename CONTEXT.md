# Contexto do projeto — Ledger em Go

> Fontes: conversas do projeto do ChatGPT "Ledger em Go" (*Documentação para ledger em Go* e *Documentação para ledger em Go (1)*) e o código do repositório no commit `9ef1469`.
> Para instruções de execução, veja o [`README.md`](./README.md).

Documentação inicial do projeto: [Ledger em Go](https://chatgpt.com/g/g-p-6aefd4b3f97c8191bc88b938025d71e8-ledger-em-go/project).

## 1. Finalidade

Construir um **sistema financeiro em Go** cujo núcleo é um **ledger (livro-razão) de partidas dobradas**. O ledger é a fonte de verdade de todas as movimentações financeiras.

Requisitos da visão original:

- Ter um ledger financeiro **consistente, rastreável e auditável**.
- Oferecer ao usuário **uma ou mais carteiras (wallets)** que respeitem o ledger.
- Trabalhar com **diferentes ativos**: moedas fiduciárias, criptoativos e ativos sem lastro.
- Guardar o **histórico de transações** e listá-lo respeitando a **privacidade** de cada carteira.
- Autenticar cada usuário por **login (CPF/CNPJ) e senha**.
- Usar **criptografia e blockchain**, ser **descentralizado** e **auditável por qualquer indivíduo**, sem expor dados pessoais.
- **Processar grandes volumes** de dados.

### Princípio orientador

> O ledger é a fonte de verdade financeira. A wallet é a interface e a identidade operacional que dá acesso aos recursos registrados nele. A blockchain pode fortalecer a verificação distribuída, mas não substitui as regras contábeis nem resolve a privacidade sozinha.

### Ordem de desenvolvimento

1. **Domínio financeiro e regras do ledger**: partidas dobradas, ativos, reversões, invariantes. ← *fase atual*
2. **Persistência e transações**: constraints, atomicidade, consistência dos saldos.
3. **Wallets e operações financeiras**: depósitos, transferências, saques, consulta de saldo.
4. **Identidade e autorização**: login CPF/CNPJ, senhas seguras, permissões por carteira.
5. **Auditoria e integridade criptográfica**: histórico imutável, hashes, assinaturas e, depois, blockchain.
6. **Escalabilidade e descentralização**: replicação, concorrência, processamento assíncrono, particionamento, testes de carga.

Segurança, testes e observabilidade acompanham todas as fases. O que fica para depois é a complexidade da blockchain, não a integridade dos registros.

## 2. Estado atual (resumo)

O projeto tem uma API HTTP em camadas (API → Handlers → Services → Domain → Repositories → PostgreSQL/GORM) que permite:

- criar, listar, consultar e desativar **carteiras**, criando o **titular** automaticamente quando ele não existe;
- criar, listar e consultar **transações** de 12 tipos, que alteram os saldos das carteiras dentro de uma transação de banco atômica.

As entidades contábeis (`LedgerAccount`, `LedgerTransaction`, `LedgerEntry`), de ativos (`Asset`, `AssetPrice`) e de contas externas (`Account`, `AccountIdentifier`) já estão **modeladas e mapeadas no banco**. Ainda **não existem serviços que as usem**: as transações alteram os saldos das carteiras, mas ainda **não geram lançamentos contábeis**.

## 3. Arquitetura

| Camada | Pacote | Responsabilidade |
|---|---|---|
| API | `packages/api` | Rotas Gin, log estruturado de requisição/resposta com `requestId`, recovery de panics |
| Handlers | `packages/handlers` | Decodificar o payload, chamar os serviços, mapear erros para status HTTP |
| Serviços | `packages/services` | Casos de uso, validações, geração de UUIDs, controle transacional |
| Domínio | `packages/domain` | Entidades e enums |
| Persistência | `packages/infra/database/models` | Modelos GORM e conversão `ToDomain`/`FromDomain` |
| Repositórios | `packages/infra/database/repositories` | `GormRepository[M, E]` genérico: `Create`, `GetByID`, `FindOne`, `FindAll`, `Update`, `Delete`, `Transaction`, `Join` |
| Migrações | `loader/`, `migrations/`, `atlas.hcl` | Schema gerado pelo Atlas a partir dos modelos GORM |

Decisões de modelagem:

- **O modelo de domínio é separado do modelo de persistência.** No domínio, os relacionamentos usam um ID (`OwnerID string`) e um ponteiro opcional (`Owner *Owner`). O ponteiro só é preenchido quando a relação é carregada (join/preload).
- **IDs opcionais são ponteiros** (`*string`) quando a relação pode não existir, como a origem de um depósito ou o destino de um saque. Os campos opcionais não devem existir só por conveniência: a obrigatoriedade segue as regras de cada `TransactionType`.
- **Todos os IDs são UUID** e todos os modelos têm `CreatedAt`, `UpdatedAt` e `DeletedAt` (exclusão lógica).

## 4. Entidades do domínio (`packages/domain`)

> A documentação de referência das entidades, em inglês, fica em [`packages/domain/doc.go`](./packages/domain/doc.go) (`go doc ./packages/domain`).

```text
Owner       1──N    Wallet
Owner       1──N    Account 1──N AccountIdentifier
Account     0..1──N Wallet (uma carteira pode apontar para uma conta externa)
Asset       1──N    Wallet, AssetPrice (como ativo base ou de cotação)
Wallet      0..1──N Transaction (como origem ou destino)
Transaction 1──1    LedgerTransaction 1──N LedgerEntry N──1 LedgerAccount
```

Há três "contas" diferentes, e elas não devem ser confundidas:

| Entidade | Representa | Pertence a | Guarda saldo? |
|---|---|---|---|
| `Wallet` | Posição do titular em um ativo **dentro da plataforma** | Um titular | Sim (projeção operacional) |
| `Account` | Conta **externa** do titular (banco ou endereço cripto) | Um titular (pode servir várias carteiras) | Não |
| `LedgerAccount` | Conta do **plano de contas contábil** | O sistema | Não (calculado pelos lançamentos) |

### 4.1 `Owner`: titular

A pessoa física (`INDIVIDUAL`) ou jurídica (`COMPANY`) dona das carteiras e contas externas. Antes se chamava `AccountHolder` mas foi renomeada para `Owner` por não ter vínculo direto a uma conta e sim a uma carteira.

- `UniqueIdentifier`: documento único (CPF/CNPJ no Brasil, SSN/EIN nos EUA), com índice único no banco. É dado pessoal sensível: deve ser normalizado e nunca aparecer em registros públicos de auditoria. Está planejado guardá-lo criptografado, com um HMAC determinístico para verificar a unicidade.
- O `Owner` não é uma identidade de autenticação. Login e senha serão uma responsabilidade separada.

### 4.2 `Wallet`: carteira

A posição financeira de um titular em **um único ativo** dentro da plataforma, por exemplo "carteira BRL do João". Antes se chamava `Account` e foi renomeada para não ser confundida com `LedgerAccount`.

- **Um ativo por carteira.** Quem tem BRL e BTC tem duas carteiras. Isso facilita validar, contabilizar e conciliar. Uma restrição `UNIQUE (owner_id, asset_id)` opcional está em avaliação.
- **Saldos**, em inteiros na menor unidade do ativo. Os três estados são mutuamente exclusivos, então o total é a soma deles:

  | Saldo | Significado |
  |---|---|
  | `FreeBalance` | Disponível para movimentação |
  | `ReservedBalance` | Comprometido com uma operação em andamento (ex.: agendamento, obrigação futura) |
  | `HoldBalance` | Bloqueado temporariamente por uma restrição ou análise (ex.: disputa). Também é o estado intermediário das transferências |

- Os saldos são uma **projeção**, não a fonte de verdade. Eles precisam poder ser reconstruídos a partir dos lançamentos do ledger e conciliados com eles.
- `Type` (`WalletType`) define **para que** a carteira serve: `CHECKING` (pagamentos), `SAVINGS` (guardar recursos), `INVESTMENT` (investimentos), `CRYPTO` (custódia de cripto). Já `Asset.Type` define **o que** ela guarda. `CHECKING` e `SAVINGS` exigem uma conta externa (validado em `services/wallet_validate.go`).
- `Status`: `ACTIVE` | `INACTIVE`. "Excluir" uma carteira é desativá-la, o que preserva o histórico.
- `Version`: reservado para optimistic locking (ainda não aplicado).
- `AccountID`: conta externa opcional, usada como origem de depósitos e destino de saques.

### 4.3 `Account`: conta externa

A conta financeira do titular **fora da plataforma**. É a origem e o destino das entradas e saídas de recursos.

- Pertence a um `Owner` (`OwnerID`) e pode estar ligada a várias carteiras dele (`wallets.account_id`). Há um índice por `owner_id` + `currency` para consultar as contas de um titular por moeda.

- `Kind`: `BANK` (banco ou instituição de pagamento) ou `CRYPTO` (endereço em uma blockchain). Define quais campos são preenchidos.
- Campos bancários: `Currency`, `InstitutionName`, `BranchCode` (agência), `AccountNumber`, `RoutingCode` (código de compensação do país), `IBAN`, `BIC/SWIFT`.
- Campos cripto: `CryptoNetwork`, `CryptoAddress`, `CryptoMemo` (memo/tag exigido por algumas redes).
- `CountryCode`: `US`, `BR`, `EU` ou `UNKNOWN` (ex.: contas cripto).

### 4.4 `AccountIdentifier`: chave de pagamento

Um apelido que endereça uma `Account` em um trilho de pagamento sem os dados bancários completos, como uma **chave PIX**.

- `Network`: `BR_PIX`, `EU_SEPA`, `US_WIRE`, `UNKNOWN`.
- `Type`: `PERSONAL_CODE` (CPF/SSN), `CORPORATE_CODE` (CNPJ/EIN), `PHONE`, `EMAIL`, `RANDOM`.

### 4.5 `Asset`: ativo

Generaliza a antiga lista fixa de moedas (`USD`, `EUR`, `BRL`). Carteiras, transações, contas contábeis e lançamentos referenciam um ativo.

- `Type`: `FIAT` (BRL, USD, EUR), `CRYPTO` (BTC, ETH, USDT), `UNBACKED` (ex.: `POINTS`).
- `DecimalPlaces`: a precisão do ativo. Valores são guardados como `valor × 10^DecimalPlaces` em `int64`, **nunca em float**. Exemplos: R$ 100,00 = `10000`; 0,00000001 BTC = `1`.
- Valores de ativos diferentes **nunca são somados nem compensados entre si**.

### 4.6 `AssetPrice`: cotação

Histórico de preços de um ativo (`Asset`) expresso em outro (`QuoteAsset`), por exemplo 1 BTC = R$ 350.000,00 em uma data. Serve para calcular o **valor de mercado** das carteiras sem alterar seus saldos.

- **Valorização não é rendimento.** Se 1 BTC passa de R$ 300 mil para R$ 350 mil, a quantidade do titular continua igual e não há transação. Um rendimento (`YIELD`) aumenta a **quantidade** de ativos e gera lançamentos.

### 4.7 `Transaction`: operação de negócio

Uma operação financeira pedida à plataforma. Ela descreve **o que** aconteceu: tipo, valor, ativo, carteiras e status. **Como** isso entra na contabilidade é responsabilidade do ledger (`LedgerTransaction` + `LedgerEntry`).

- Tipos (`TransactionType`):
  - Movimentações: `DEPOSIT`, `WITHDRAWAL`, `TRANSFER`, `REVERSAL`
  - Receitas, custos e negociação: `YIELD`, `FEE`, `PURCHASE`, `SALE`
  - Mudança de estado do saldo: `RESERVE`, `RELEASE`, `HOLD`, `UNHOLD`
- Status (`TransactionStatus`): `PENDING` → `COMPLETED` | `FAILED`. Transações rejeitadas por regra de negócio também são gravadas, como `FAILED`.
- `SourceWalletID` e `DestinationWalletID` são opcionais, porque nem toda operação tem os dois lados.
- `ReferenceID`: identifica a operação ou evento externo que originou a transação.
- `ReversalOriginID`: a transação original, quando esta é uma reversão.
- Uma transação concluída **nunca é editada nem excluída**. Erros são corrigidos com uma reversão.

### 4.8 `LedgerAccount`: conta contábil

Uma conta do **plano de contas**. Define como cada movimentação é classificada na contabilidade.

- É uma categoria contábil **persistente**, criada na configuração do plano de contas (ou quando surge uma nova necessidade contábil), e **não a cada operação**. Uma mesma conta, como "Saldo disponível dos usuários", recebe lançamentos de milhares de carteiras.
- **Não guarda saldo.** O saldo é calculado a partir dos `LedgerEntry` dela.
- `Type`: `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`.
- `NormalBalanceDirection`: o lado que **aumenta** o saldo da conta. `ASSET` e `EXPENSE` aumentam por `DEBIT`; `LIABILITY`, `EQUITY` e `REVENUE` aumentam por `CREDIT`. O campo serve para **interpretar** o saldo, não para restringir lançamentos: uma conta de ativo também recebe créditos. Ele deve ser compatível com `Type`.
- `Code` (único) e `Name`. Por exemplo, `2001 - Saldo disponível dos usuários`.
- `AssetID`: contas com a mesma finalidade podem ser duplicadas por ativo.
- `Status`: `ACTIVE` (aceita lançamentos) | `INACTIVE`. Não muda quando operações são contabilizadas.

Exemplo de plano de contas:

| Código | Conta | Tipo | Saldo normal |
|---|---|---|---|
| 1001 | Reservas da plataforma | ASSET | DEBIT |
| 2001 | Saldo disponível dos usuários | LIABILITY | CREDIT |
| 2002 | Saldo reservado dos usuários | LIABILITY | CREDIT |
| 2003 | Saldo retido dos usuários | LIABILITY | CREDIT |
| 4001 | Receita com tarifas | REVENUE | CREDIT |
| 5001 | Despesas com rendimentos | EXPENSE | DEBIT |

### 4.9 `LedgerTransaction`: contabilização

O registro contábil de uma `Transaction`. Agrupa os lançamentos gerados por aquela operação.

- `TransactionID` é único: uma contabilização principal por operação.
- `Status`: `PENDING` → `POSTED` | `FAILED`. Esse ciclo é independente do status da `Transaction`.
- `PostedAt`: o momento da contabilização.
- `ReversalOriginID`: a contabilização revertida, quando esta é uma reversão. A reversão reproduz os lançamentos originais com as direções invertidas, e o efeito líquido é zero.
- Depois de `POSTED`, é **imutável**.

### 4.10 `LedgerEntry`: lançamento

Uma linha **imutável** de débito ou crédito no livro-razão. É aqui que os **valores** ficam.

- `Amount` é sempre positivo. O lado vem de `Direction` (`DEBIT` | `CREDIT`), numa coluna única em vez de colunas separadas de débito e crédito.
- `AssetID` deve ser o mesmo ativo da conta contábil.
- `Sequence` dá uma ordem determinística aos lançamentos de uma mesma `LedgerTransaction`, a partir de 1. Serve para ordenação, auditoria, replay e detecção de lançamentos ausentes ou duplicados. **Não** é a ordem cronológica global do ledger.
- Está planejado acrescentar uma referência à carteira afetada (`WalletID`). Assim, uma única conta contábil guarda o saldo agregado de todos os usuários, e o lançamento identifica a carteira individual.

## 5. Conceitos implementados

- **Saldos em inteiros na menor unidade do ativo**, com três estados: livre, reservado e retido.
- **Ativos genéricos** (fiat, cripto, sem lastro) com precisão configurável. As carteiras são criadas pelo código do ativo.
- **Uma carteira por ativo**, ligada a um titular. O titular é criado automaticamente quando não existe, e seu documento é único.
- **Novas carteiras começam com saldos zerados.** "Excluir" é desativar (`INACTIVE`), e o soft delete preserva o histórico.
- **Validações de transação**: tipo suportado, valor positivo, carteiras aceitas pelo tipo, carteiras distintas, ativas e do mesmo ativo, e saldo suficiente. As carteiras podem ser de tipos e titulares diferentes.
- **Transferência em duas etapas via retenção**: origem livre → origem retida → destino retido → destino livre. Isso impede que uma transferência fique pela metade.
- **Carteiras e efeito de cada tipo de transação** (`transactionWalletRules` e `applyTransaction`):

  | Tipo | Origem | Destino | Efeito |
  |---|---|---|---|
  | `DEPOSIT` | não aceita | obrigatória | destino: livre += valor |
  | `WITHDRAWAL` | obrigatória | não aceita | origem: livre −= valor |
  | `TRANSFER` | obrigatória | obrigatória | origem → destino (saldo livre) |
  | `YIELD` | opcional (carteira que gerou o rendimento; não é debitada) | obrigatória | destino: livre += valor |
  | `SALE` | opcional (comprador) | obrigatória | transferência ou só crédito no destino |
  | `FEE`, `PURCHASE` | obrigatória | opcional (receita ou vendedor) | transferência ou só débito na origem |
  | `RESERVE` / `RELEASE` | obrigatória | não aceita | livre ↔ reservado |
  | `HOLD` / `UNHOLD` | obrigatória | não aceita | livre ↔ retido |

- **Reversão de transações** (`REVERSAL`): exige `reversalOriginId` e uma original `COMPLETED` que não seja reversão nem tenha sido revertida. O valor é o mesmo da original. A reversão troca origem e destino (nas mudanças de estado do saldo, mantém a carteira; na de `YIELD`, deixa de fora a carteira geradora, que nunca foi debitada) e aplica o efeito inverso. A ligação com a original fica em `ReversalOriginID`.
- **Concorrência**: as carteiras são travadas com `SELECT ... FOR UPDATE`, sempre na mesma ordem de ID, e a transação original de uma reversão também é travada. Isso evita saldos inconsistentes, deadlocks e reversões duplicadas. Foi testado com transferências concorrentes contra o PostgreSQL.
- **Atomicidade**: travamento das carteiras, validação, atualização dos saldos e gravação da transação acontecem numa única transação de banco. Uma transação rejeitada depois de carregar as carteiras é gravada como `FAILED` nessa mesma transação de banco, sem alterar os saldos.
- **Transação completa pela API**: `AssetID` vem das carteiras, e `Description`, `ReferenceID` e `ReversalOriginID` são informados na requisição. Os erros viram `400` (inválida), `404` (carteira inexistente) e `422` (saldo insuficiente).
- **Entidades contábeis, de ativos, de cotações e de contas externas** modeladas no domínio e mapeadas no banco.
- **Repositório genérico** com filtros por exemplo, joins e transações (`WithDB`).
- **Migrações com Atlas** geradas a partir dos modelos GORM e sincronizadas com eles.
- **Constraints no banco**: `CHECK` de valor positivo (`transactions.amount`, `ledger_entries.amount`), de saldos não negativos (`wallets`) e de direção válida (`ledger_entries.direction`), além de `UNIQUE (ledger_transaction_id, sequence)`.
- **Observabilidade**: log estruturado de cada requisição e resposta (parâmetros, corpo, IP, `requestId`, status, duração) e logs com contexto por handler.
- **Configuração por variáveis de ambiente**, Docker Compose (PostgreSQL e Redis) e Makefile.
- **Testes unitários** das operações de saldo da carteira (inclusive erros), do efeito e da reversão de cada tipo de transação, da preservação dos saldos quando falta saldo e da validação da entrada e das carteiras de cada tipo.

## 6. Conceitos planejados (não implementados)

### 6.1 Partidas dobradas

Toda movimentação registra **pelo menos dois lançamentos**, e em cada transação **Σ débitos = Σ créditos, por ativo**.

| Categoria | Aumenta com | Diminui com |
|---|---|---|
| Ativo, Despesa | Débito | Crédito |
| Passivo, Patrimônio líquido, Receita | Crédito | Débito |

O saldo de uma conta contábil é **Débitos − Créditos** quando o saldo normal é `DEBIT`, e **Créditos − Débitos** quando é `CREDIT`.

Regras: equilíbrio, atomicidade, imutabilidade, reversão (em vez de edição), separação por ativo e rastreabilidade (todo lançamento pertence a uma transação identificável).

Contabilização prevista para cada operação:

| Operação | Débito | Crédito |
|---|---|---|
| Depósito de R$ 100 | 1001 Reservas da plataforma | 2001 Disponível (carteira do João) |
| Saque de R$ 30 | 2001 Disponível (carteira do João) | 1001 Reservas da plataforma |
| Transferência João → Maria | 2001 Disponível (João) | 2001 Disponível (Maria) |
| Reserva | 2001 Disponível (João) | 2002 Reservado (João) |
| Confirmação da reserva | 2002 Reservado (João) | 2001 Disponível (Maria) |
| Rendimento de R$ 50 | 5001 Despesas com rendimentos | 2001 Disponível (João) |
| Tarifa | 2001 Disponível (João) | 4001 Receita com tarifas |

A contrapartida real de um rendimento depende da origem econômica dele. Pode ser uma obrigação, uma conta de liquidação etc., não necessariamente uma despesa. A recomendação é ter **contas contábeis separadas por estado de saldo** (disponível, reservado e retido), para facilitar a conciliação.

### 6.2 Reversões contábeis

- Só transações contabilizadas podem ser revertidas, e **uma única vez**.
- A reversão cria uma nova `Transaction` (`REVERSAL`) e uma nova `LedgerTransaction`, com os lançamentos invertidos. As duas ficam ligadas por `ReversalOriginID`. A original fica intacta.
- Campos que podem ser acrescentados: `ReversedBy` e `ReversalReason`.
- **Reversão não é cancelamento.** Uma transação pendente é cancelada liberando o valor reservado (status `CANCELLED`). Uma transação contabilizada precisa de reversão.

### 6.3 Idempotência

Uma entidade `IdempotencyKey` (`Key`, `OwnerID`, `Operation`, `RequestHash`, `TransactionID`, `Status` `PROCESSING|COMPLETED|FAILED`, `Response`, `ExpiresAt`) com `UNIQUE(owner_id, operation, key)`.

- Se a chave já existe e está concluída, a resposta anterior é devolvida.
- Se a chave já existe com um `RequestHash` diferente, a requisição é rejeitada.
- A chave e a transação financeira são gravadas de forma atômica.

O plano é começar pela criação de transações.

### 6.4 Concorrência

- O bloqueio pessimista já está implementado (ver seção 5).
- Usar optimistic locking com `Wallet.Version` como complemento.
- Usar `SERIALIZABLE` só onde necessário, com política de retry. As constraints `CHECK` e `UNIQUE` básicas já existem (ver seção 5).
- Testar transferências simultâneas.

### 6.5 Conciliação e reconstrução de saldos

- **Equilíbrio**: validar Σ débitos − Σ créditos = 0 por ativo em cada contabilização.
- **Conciliação** (`ReconciliationService`): comparar os saldos materializados (livre, reservado, retido) com os saldos calculados pelo ledger, gerar relatórios de divergência e **não** corrigir nada automaticamente.
- **Reconstrução** (`BalanceReconstructionService`): recalcular o saldo de uma conta ou carteira a partir dos lançamentos, ordenados e interpretados pelo saldo normal, inclusive em uma data passada.
- **Snapshots** (`BalanceSnapshot`: `LedgerAccountID`, ativo, saldos de débito e crédito, `LastEntrySequence`) para acelerar reconstruções. Snapshots e saldos materializados são **projeções** e nunca substituem os lançamentos.

### 6.6 Auditoria imutável

- Uma entidade `AuditLog` (`ActorID`, `EventType`, `EntityType`, `EntityID`, `Payload`, `PreviousHash`, `EventHash`, `CreatedAt`), com eventos como `WALLET_CREATED`, `TRANSACTION_COMPLETED`, `TRANSACTION_REVERSED` e `LEDGER_ENTRY_POSTED`.
- **Hash chaining**: Hₙ = SHA256(Hₙ₋₁ + eventoₙ), calculado sobre uma representação canônica do evento.
- Proibir `UPDATE`/`DELETE` em lançamentos contabilizados e em eventos de auditoria, com permissões específicas no banco.
- Publicar **checkpoints** dos hashes, assinados ou ancorados em blockchain, para detectar a reescrita da cadeia inteira.
- Nunca incluir CPF, CNPJ, senhas ou chaves nos eventos públicos.

### 6.7 Rendimentos e valorização

- Suportar os dois tipos, **em registros separados**:
  - **Rendimento por quantidade**: transação `YIELD` com lançamentos.
  - **Valorização de mercado**: calculada com `AssetPrice`, sem lançamentos.

### 6.8 Fases posteriores

- **Identidade**: login por CPF/CNPJ e senha com hash seguro (OWASP), sessões, autorização por carteira e privacidade do histórico. O documento deve ser guardado criptografado, com HMAC para garantir a unicidade.
- **Criptografia e blockchain**: assinaturas (Ed25519), consenso (referências: Hyperledger Fabric, CometBFT, Ethereum), um ledger auditável publicamente sem dados pessoais e integração com redes cripto.
- **Escala**: Kafka ou processamento assíncrono, particionamento, replicação, SQL distribuído (CockroachDB), base analítica (ClickHouse) e testes de carga.
- **Integração bancária**: depósitos e saques via `Account` e `AccountIdentifier` (PIX, SEPA, wire).

### 6.9 Referências de estudo

TigerBeetle, Formance, Modern Treasury (ledger), PostgreSQL (isolamento e locking), GORM (transações e locking), Hyperledger Fabric, CometBFT, Ethereum, pacote `crypto` do Go, OWASP Cheat Sheets, OpenID Connect e Apache Kafka.

## 7. Divergências entre o código e o modelo planejado

- As transações **alteram os saldos diretamente** e não geram `LedgerTransaction` nem `LedgerEntry`.
- `DEPOSIT` e `WITHDRAWAL` alteram só a carteira interna. A contrapartida externa (conta bancária ou contábil, como as reservas da plataforma) ainda não é registrada.
- `REVERSAL` desfaz o efeito nos saldos, mas ainda não gera lançamentos contábeis invertidos.
- Carteiras `CHECKING` e `SAVINGS` não podem ser criadas pela API, porque ela não recebe `AccountID`.

## 8. Checklist de funcionalidades e conceitos

### Fundamentos do domínio

- [x] Representação monetária em inteiros na menor unidade do ativo (sem float)
- [x] Entidade genérica de ativos (`Asset`) com tipo (fiat, cripto, sem lastro) e casas decimais
- [x] Carteira de um único ativo, ligada a um titular
- [x] Saldos separados em livre, reservado e retido
- [x] Distinção entre `Wallet`, `Account` (conta externa) e `LedgerAccount` (conta contábil)
- [x] Regras de partidas dobradas documentadas
- [x] Entidades contábeis modeladas (`LedgerAccount`, `LedgerTransaction`, `LedgerEntry`)
- [x] Entidade de cotações modelada (`AssetPrice`)
- [x] Entidades de conta externa e chaves de pagamento modeladas (`Account`, `AccountIdentifier`)
- [ ] Seed de ativos suportados
- [ ] Plano de contas inicial cadastrado (seed de `LedgerAccount`)
- [ ] Validar a compatibilidade entre `NormalBalanceDirection` e `LedgerAccount.Type`
- [ ] Regras de negócio específicas por `WalletType`

### Persistência

- [x] PostgreSQL com GORM
- [x] Separação entre modelos de domínio e de persistência
- [x] Repositório genérico (CRUD, filtros, joins, transações)
- [x] Geração de migrações com Atlas a partir dos modelos GORM
- [x] Transações atômicas no banco
- [x] UUIDs como identificadores
- [x] Documento do titular único (índice único)
- [x] Exclusão lógica (soft delete) e desativação de carteiras
- [x] Migração versionada sincronizada com os modelos atuais
- [x] `UNIQUE (ledger_transaction_id, sequence)` em `ledger_entries`
- [x] Constraints `CHECK` (valor > 0, direção válida, saldos ≥ 0)
- [ ] `UNIQUE (owner_id, asset_id)` em `wallets` (se for regra de negócio)

### Carteiras e titulares

- [x] Criar carteira, criando o titular quando ele não existe
- [x] Listar carteiras com filtros
- [x] Consultar carteira com ativo, titular e conta externa
- [x] Desativar carteira preservando o histórico
- [x] Novas carteiras com saldos zerados
- [ ] Vincular conta externa (`Account`) pela API, habilitando `CHECKING` e `SAVINGS`
- [ ] CRUD de contas externas e chaves de pagamento
- [ ] Endpoints de ativos e cotações

### Transações

- [x] Criar transação com validações (valor positivo, carteiras distintas, mesmo ativo, ativas, saldo suficiente)
- [x] Transferência em duas etapas via retenção, sem conclusão parcial
- [x] Aplicação dos 12 tipos de transação aos saldos
- [x] Gravação de transações rejeitadas como `FAILED`
- [x] Listar e consultar transações com filtros
- [x] Preencher `AssetID` da transação a partir das carteiras
- [x] Informar `Description` e `ReferenceID` pela API
- [x] Usar `ReservedBalance` em `RESERVE`/`RELEASE`
- [x] Depósitos e saques com uma carteira só (sem exigir duas carteiras internas)
- [x] Transações com carteira opcional (`YIELD`, `SALE`, `FEE`, `PURCHASE`) e entre carteiras de tipos diferentes
- [x] Reversão de transações concluídas, ligada à original e sem reversão dupla
- [x] Respostas `400`/`404`/`422` conforme o tipo de erro da transação
- [ ] Registrar a contrapartida externa de depósitos e saques (conta externa ou contábil)
- [ ] Cancelamento de transações pendentes (`CANCELLED`)

### Núcleo contábil (Marco 1)

- [ ] Gerar `LedgerTransaction` e `LedgerEntry` para cada transação
- [ ] Validar o equilíbrio débito = crédito por ativo
- [ ] Lançamentos imutáveis (bloquear `UPDATE`/`DELETE` depois de `POSTED`)
- [ ] Referência à carteira afetada em `LedgerEntry`
- [ ] Contas contábeis por estado de saldo (disponível, reservado, retido)
- [ ] Reversões contábeis ligadas à original, sem reversão dupla
- [ ] Idempotência das requisições (`IdempotencyKey`)
- [x] Bloqueio pessimista (`SELECT ... FOR UPDATE`) em ordem determinística
- [ ] Optimistic locking com `Wallet.Version`
- [ ] Reconstrução de saldos pelo histórico de lançamentos
- [ ] Conciliação entre saldos materializados e o ledger
- [ ] Snapshots de saldo (`BalanceSnapshot`)

### Rendimentos e valorização

- [x] Tipo de transação `YIELD`
- [x] Modelagem de cotações (`AssetPrice`)
- [ ] Contabilização de rendimentos (despesa/obrigação × passivo do usuário)
- [ ] Cálculo do valor de mercado das carteiras a partir das cotações

### Auditoria e observabilidade

- [x] Logs estruturados de requisições e respostas HTTP com `requestId`
- [x] Logs com contexto por handler
- [ ] Trilha de auditoria imutável (`AuditLog`)
- [ ] Hash chaining dos eventos de auditoria
- [ ] Checkpoints assinados ou ancorados em blockchain

### Testes

- [x] Testes unitários das operações de saldo da carteira
- [x] Testes da aplicação de todos os tipos de transação
- [x] Testes de preservação dos saldos em caso de falha (saldo insuficiente e transações inválidas)
- [ ] Testes do fluxo completo de `TransactionService.Create` (status `COMPLETED`/`FAILED` e rollback no banco)
- [ ] Testes automatizados de concorrência (hoje validados só manualmente, contra um PostgreSQL local)
- [ ] Testes de integração com banco real
- [ ] Testes dos handlers e da API

### Infraestrutura

- [x] Configuração por variáveis de ambiente
- [x] Docker Compose (PostgreSQL e Redis)
- [x] Makefile (run, build, test, migrações)
- [x] Health check
- [ ] Uso do Redis (cache)
- [ ] Lint (`go vet`, `golangci-lint`) e CI

### Identidade e segurança (fase 4)

- [ ] Login por CPF/CNPJ e senha com hash seguro
- [ ] Autorização por carteira e privacidade do histórico
- [ ] Documento do titular criptografado, com HMAC para unicidade

### Descentralização e escala (fases 5 e 6)

- [ ] Assinaturas digitais das operações
- [ ] Blockchain e consenso distribuído
- [ ] Ledger auditável publicamente sem exposição de dados pessoais
- [ ] Processamento assíncrono (Kafka), particionamento e replicação
- [ ] Testes de carga para grandes volumes
