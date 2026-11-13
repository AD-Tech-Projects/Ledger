# Plano de Contas Inicial

> Contas contábeis (`LedgerAccount`). Para o conceito de conta contábil, veja a seção 4.8 do [Contexto do Projeto](./project-overview.md).

## Convenções

- **Uma conta por ativo.** Cada `LedgerAccount` pertence a um único ativo (`AssetID`), e valores de ativos diferentes nunca são somados. Por isso cada conta base é criada uma vez para cada ativo em que faz sentido.
- **Código.** Como `Code` é único no banco, o código de cada conta é `BASE-ATIVO`, por exemplo `2001-BRL` e `2001-BTC`. A base identifica a finalidade da conta, e o sufixo identifica o ativo.
- **Nome.** É o nome da conta base seguido do ativo entre parênteses, por exemplo `Saldo disponível dos usuários (BRL)`.
- **Status.** Todas as contas são criadas como `ACTIVE`.
- **Saldo.** As contas não guardam saldo. O saldo é calculado pelos lançamentos (`LedgerEntry`) e interpretado pela direção normal: **débitos − créditos** quando a direção é `DEBIT`, e **créditos − débitos** quando é `CREDIT`.
- **Carteira individual.** As contas guardam o total de todos os titulares. A carteira de cada lançamento será identificada pelo `WalletID` planejado em `LedgerEntry`, sem contas separadas por carteira ou por `WalletType`.

## Faixas de código

A base do código tem quatro dígitos e segue a hierarquia classe → grupo → conta:

| Posição | Significado | Valores |
|---|---|---|
| 1º dígito | Classe: o tipo da conta (`LedgerAccountType`) | `1` `ASSET`, `2` `LIABILITY`, `3` `EQUITY`, `4` `REVENUE`, `5` `EXPENSE` |
| 2º dígito | Grupo dentro do tipo | `0` a `9` (ex.: `20xx` obrigações com os titulares, `31xx` conversão entre ativos) |
| 3º e 4º dígitos | Número sequencial da conta dentro do grupo | `01` a `99` |

Exemplos:

- `2003`: `2` passivo, `0` obrigações com os titulares, `03` terceira conta do grupo (saldo retido).
- `3101`: `3` patrimônio líquido, `1` conversão entre ativos, `01` primeira conta do grupo.
- `5002`: `5` despesa, `0` despesas com os titulares, `02` segunda conta do grupo (programa de pontos).

Cada tipo comporta até 10 grupos, e cada grupo até 99 contas. O sufixo `-ATIVO` não faz parte da base: ele só distingue as cópias da mesma conta por ativo.

Grupos em uso e reservados (as faixas marcadas como reservadas ainda não têm contas):

| Faixa | Tipo | Direção normal | Grupo |
|---|---|---|---|
| `10xx` | `ASSET` | `DEBIT` | Recursos custodiados pela plataforma (reservas em bancos e em carteiras cripto) |
| `11xx` | `ASSET` | `DEBIT` | Reservada: valores em liquidação (depósitos e saques pendentes de confirmação externa) |
| `20xx` | `LIABILITY` | `CREDIT` | Obrigações com os titulares (saldos das carteiras, por estado) |
| `21xx` | `LIABILITY` | `CREDIT` | Reservada: outras obrigações (ex.: valores a repassar a parceiros) |
| `30xx` | `EQUITY` | `CREDIT` | Reservada: capital próprio da plataforma |
| `31xx` | `EQUITY` | `CREDIT` | Contas de conversão entre ativos (negociação) |
| `40xx` | `REVENUE` | `CREDIT` | Receitas com serviços (tarifas) |
| `41xx` | `REVENUE` | `CREDIT` | Reservada: receitas com negociação (ex.: spread de conversão) |
| `50xx` | `EXPENSE` | `DEBIT` | Despesas com os titulares (rendimentos e programa de pontos) |

## Contas

Ativos com lastro: `BRL`, `USD`, `EUR`, `BTC`, `ETH` e `USDT`. `POINTS` é a moeda interna sem lastro (1 ponto = R$ 0,50).

| Código | Conta | Tipo | Direção normal | Ativos | Descrição e utilização |
|---|---|---|---|---|---|
| `1001` | Reservas da plataforma | `ASSET` | `DEBIT` | BRL, USD, EUR, BTC, ETH, USDT | Dinheiro e criptoativos que a plataforma mantém em bancos e carteiras de custódia em nome dos titulares. É a contrapartida de depósitos (débito) e saques (crédito). Só muda quando recursos entram ou saem da plataforma. `POINTS` não tem lastro e não tem reserva. |
| `2001` | Saldo disponível dos usuários | `LIABILITY` | `CREDIT` | Todos | Obrigação da plataforma com os titulares pelo saldo livre das carteiras (`FreeBalance`). Recebe crédito em depósitos, rendimentos, entradas de transferências e compras, e débito em saques, tarifas, saídas de transferências e quando o saldo passa para reservado ou retido. |
| `2002` | Saldo reservado dos usuários | `LIABILITY` | `CREDIT` | Todos | Obrigação pelo saldo comprometido com operações em andamento (`ReservedBalance`). Movimentada por `RESERVE` (crédito, a partir de `2001`) e `RELEASE` (débito, de volta para `2001`). |
| `2003` | Saldo retido dos usuários | `LIABILITY` | `CREDIT` | Todos | Obrigação pelo saldo bloqueado (`HoldBalance`). Movimentada por `HOLD` e `UNHOLD`, e usada como etapa intermediária das transferências (livre da origem → retido da origem → retido do destino → livre do destino). |
| `3101` | Conversão de ativos | `EQUITY` | `CREDIT` | Todos | Conta-ponte para compras e vendas entre ativos diferentes. Em cada ativo, a operação fecha débitos = créditos sem somar valores de ativos distintos: o ativo entregue sai de `2001-X` para `3101-X`, e o ativo recebido sai de `3101-Y` para `2001-Y`. O saldo pode ficar devedor em um ativo e credor em outro; a diferença entre eles, avaliada pelas cotações, é o resultado das conversões. |
| `4001` | Receita com tarifas | `REVENUE` | `CREDIT` | BRL, USD, EUR, BTC, ETH, USDT | Tarifas cobradas dos titulares. Contrapartida de `FEE` sem carteira de destino: débito em `2001-X` e crédito em `4001-X`. |
| `5001` | Despesas com rendimentos | `EXPENSE` | `DEBIT` | BRL, USD, EUR, BTC, ETH, USDT | Rendimentos pagos aos titulares (`YIELD`): débito em `5001-X` e crédito em `2001-X`. Quando o rendimento vier de um parceiro (ex.: um banco que remunera a caixinha), a contrapartida poderá ser uma conta a receber em vez de uma despesa. |
| `5002` | Despesas com programa de pontos | `EXPENSE` | `DEBIT` | POINTS | Emissão de pontos para os titulares, como custo do programa de fidelidade: débito em `5002-POINTS` e crédito em `2001-POINTS`. |

Total: 47 contas.

| Conta | Quantidade |
|---|---|
| `1001`, `4001` e `5001` | 6 cada |
| `2001`, `2002`, `2003` e `3101` | 7 cada |
| `5002` | 1 |

## Contabilização das operações

`X` é o ativo da transação.

| Operação | Débito | Crédito |
|---|---|---|
| Depósito (`DEPOSIT`) | `1001-X` | `2001-X` |
| Saque (`WITHDRAWAL`) | `2001-X` | `1001-X` |
| Transferência (`TRANSFER`), por etapa | `2001-X` (origem), depois `2003-X` (origem → destino) | `2003-X`, depois `2001-X` (destino) |
| Reserva / liberação (`RESERVE` / `RELEASE`) | `2001-X` / `2002-X` | `2002-X` / `2001-X` |
| Retenção / liberação (`HOLD` / `UNHOLD`) | `2001-X` / `2003-X` | `2003-X` / `2001-X` |
| Rendimento (`YIELD`) | `5001-X` | `2001-X` |
| Tarifa sem destino (`FEE`) | `2001-X` | `4001-X` |
| Compra de US$ 1,00 por R$ 5,51 | `2001-BRL` 551 e `3101-USD` 100 | `3101-BRL` 551 e `2001-USD` 100 |
| Emissão de pontos | `5002-POINTS` | `2001-POINTS` |
| Reversão (`REVERSAL`) | Os lançamentos da original com as direções invertidas | |

## Fora do plano inicial

- **Valores em liquidação (`11xx`).** Serão necessários quando depósitos e saques ficarem `PENDING` até a confirmação externa (PIX, confirmações na blockchain).
- **Receita com spread de conversão (`41xx`).** Será necessária quando as conversões cobrarem uma margem sobre a cotação.
- **Capital próprio (`30xx`).** Será necessário quando a plataforma registrar o próprio capital ou operar com saldo próprio.
