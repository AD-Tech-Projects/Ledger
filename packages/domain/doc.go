// Package domain contains the entities and enumerations of the Ledger,
// a financial system built around a double-entry bookkeeping ledger.
//
// The ledger is meant to be the financial source of truth. Wallets are the operational view
// that gives owners access to the resources recorded in it. Business operations
// ([Transaction]) are recorded in the books as balanced debit/credit lines
// ([LedgerTransaction] + [LedgerEntry]) classified by a chart of accounts ([LedgerAccount]).
//
// # Relationships
//
// In "A x:y B", each B is linked to x records of A and each A is linked to y records of B;
// "0..1" marks an optional link.
//
//	Owner             1:N    Wallet
//	Owner             1:N    Account
//	Account           1:N    AccountIdentifier
//	Account           0..1:N Wallet             (a wallet points to at most one external account)
//	Asset             1:N    Wallet
//	Asset             1:N    AssetPrice         (as the priced asset, AssetID)
//	Asset             1:N    AssetPrice         (as the quote asset, QuoteAssetID)
//	Asset             1:N    Transaction
//	Asset             1:N    LedgerAccount
//	Asset             1:N    LedgerEntry
//	Wallet            0..1:N Transaction        (as source, SourceWalletID)
//	Wallet            0..1:N Transaction        (as destination, DestinationWalletID)
//	Transaction       0..1:N Transaction        (reversals, ReversalOriginID)
//	Transaction       1:0..1 LedgerTransaction  (at most one posting per transaction)
//	LedgerTransaction 0..1:N LedgerTransaction  (reversals, ReversalOriginID)
//	LedgerTransaction 1:N    LedgerEntry
//	LedgerAccount     1:N    LedgerEntry
//
// # Three different "accounts"
//
// The domain has three entities that could be called "account". They must not be confused:
//
//   - [Wallet]: the owner's position (balances) in one asset inside the platform.
//     Belongs to an owner and stores an operational projection of the balance.
//   - [Account]: the owner's external account (bank account or blockchain address).
//     Belongs to an owner, may serve many wallets and is used to move funds into and out of
//     the platform; stores no balance.
//   - [LedgerAccount]: an account of the accounting chart of accounts.
//     Belongs to the system; its balance is computed from ledger entries.
//
// # Owner
//
// [Owner] is the individual or company that owns wallets and external accounts
// (also known as account holder or customer; it was previously called AccountHolder).
// It is identified by a unique document (CPF/CNPJ in Brazil, SSN/EIN in the United States),
// which is sensitive personal data: it must be normalized, unique, and never exposed in public
// audit records. An owner is not an authentication identity; login credentials are a separate concern.
//
// Its CountryCode ([CountryCode]) is the jurisdiction of the document. The unique identifier follows
// the format COUNTRY:TYPE:DOCUMENT (e.g. "BR:INDIVIDUAL:12345678909"), where COUNTRY is the [CountryCode],
// TYPE is the [OwnerType] and DOCUMENT is the document with everything but letters and digits removed.
// The whole identifier is converted to uppercase.
//
// # Wallet
//
// [Wallet] is an owner's financial position in a single asset inside the platform
// (e.g. "João's BRL wallet"). It was previously called Account and was renamed to avoid
// confusion with [LedgerAccount]. An owner that holds several assets has several wallets.
//
// A wallet has three mutually exclusive balance states, all stored as integers in the smallest
// unit of the asset (the total balance is their sum):
//
//   - free: available to be moved;
//   - reserved: committed to an operation in progress (e.g. a scheduled transaction);
//   - hold: temporarily blocked by a restriction or analysis (e.g. a dispute). It is also the
//     intermediate state of a transfer (source free → source hold → destination hold →
//     destination free), which prevents partial transfers.
//
// These balances are a projection used for fast reads and validations. They must always be
// reconcilable with, and reconstructable from, the ledger entries.
//
// [WalletType] describes what the wallet is used for (checking, savings, investment, crypto),
// while [AssetType] describes what it stores. Checking and savings wallets must be linked to an
// external [Account]. Deleting a wallet is modeled as
// deactivation, so its history is preserved.
//
// # Account and AccountIdentifier
//
// [Account] is the owner's financial account outside the platform: a bank account
// (institution, branch, number, routing code, IBAN, BIC/SWIFT) or a crypto account
// (network, address, memo). It is the origin of deposits and the destination of withdrawals.
// It belongs to an [Owner] and may be linked to many of the owner's wallets; an owner may have
// several accounts (they can be looked up by owner and currency).
//
// [AccountIdentifier] is an alias that addresses an [Account] on a specific payment rail
// without its full bank details, such as a PIX key (CPF, CNPJ, phone, email or random key).
// An account may have many identifiers. Its CountryCode ([CountryCode]) is the jurisdiction of the
// account, which defines the available payment rails; [CountryCode] is shared with [Owner].
//
// # Asset and AssetPrice
//
// [Asset] is anything that can be held and moved by the platform: a fiat currency, a crypto asset
// or an unbacked asset issued by the platform (e.g. points). It replaced the old fixed currency list
// (USD, EUR, BRL). Its decimal places define the precision of the amounts, which are always stored
// as int64 in the smallest unit (never as floats). Amounts of different assets are never added
// together or compensated against each other.
//
// [AssetPrice] is a historical quote of an asset expressed in another (quote) asset,
// e.g. 1 BTC = 350,000.00 BRL at a given moment. It is used to compute the market value of wallets
// without changing their balances. An appreciation is not a transaction and generates no ledger
// entries, while a yield ([TRANSACTION_YIELD]) increases the quantity of the asset held.
//
// # Transaction
//
// [Transaction] is a business operation requested to the platform: transfer, deposit, withdrawal,
// reversal, yield, fee, purchase, sale, or a change of balance state (reserve, release, hold, unhold).
// It describes what happened (type, amount, asset, wallets, status). How it is recorded in the books is
// the responsibility of the ledger: each posted transaction originates one [LedgerTransaction].
//
// Source and destination wallets are optional because not every operation has both sides.
// Completed transactions are never edited or deleted; mistakes are corrected by a new reversal
// transaction that references the original one.
//
// # LedgerAccount
//
// [LedgerAccount] is an account of the chart of accounts. It defines how financial movements are
// classified in the books, through its accounting nature ([LedgerAccountType]: asset, liability,
// equity, revenue, expense), its normal balance direction and its asset.
//
// It is not a user wallet: it is a persistent accounting category, created when the chart of accounts
// is configured (or when a new accounting need appears) and reused by many operations. A single
// ledger account such as "Users available balance" receives entries related to thousands of wallets.
// It stores no balance; its balance is computed from its entries.
//
// The normal balance direction is the side that increases the account balance: DEBIT for asset and
// expense accounts, CREDIT for liability, equity and revenue accounts. It is used to interpret
// entries (balance = debits − credits, or credits − debits), not to restrict them.
//
// Example chart of accounts:
//
//	Code  Name                      Type       Normal balance
//	1001  Platform reserves         ASSET      DEBIT
//	2001  Users available balance   LIABILITY  CREDIT
//	2002  Users reserved balance    LIABILITY  CREDIT
//	2003  Users held balance        LIABILITY  CREDIT
//	4001  Fee revenue               REVENUE    CREDIT
//	5001  Yield expenses            EXPENSE    DEBIT
//
// # LedgerTransaction
//
// [LedgerTransaction] is the accounting record (posting) of a [Transaction]. It groups the entries
// generated by that operation, which must be balanced per asset (sum of debits = sum of credits) and
// persisted atomically. Its lifecycle ([LedgerTransactionStatus]) is independent of the business
// transaction status. Once posted it is immutable; corrections are made by another ledger transaction
// that reverses it, reproducing the original entries with inverted directions.
//
// # LedgerEntry
//
// [LedgerEntry] is a single immutable debit or credit line of the general ledger. This is where the
// amounts live; ledger accounts only classify them. Each entry has a positive amount and a direction
// (a single column instead of separate debit/credit columns), and its asset must match the asset of
// its ledger account. The sequence gives a deterministic order to the entries of the same ledger
// transaction (ordering, auditing, replay, detection of missing or duplicated entries); it is not
// a global chronological order.
//
// Example: a transfer of BRL 100.00 between two wallets.
//
//	Seq  LedgerAccount             Direction  Amount
//	1    Users available balance   DEBIT      10000
//	2    Users available balance   CREDIT     10000
package domain
