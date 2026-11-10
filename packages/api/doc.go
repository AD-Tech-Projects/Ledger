// Package api defines the HTTP routes of the Ledger, the request/response logging middleware and the server startup.
//
// Responses are the JSON of the domain entities. Since they have no json tags, fields are serialized in
// PascalCase (ID, FreeBalance, OwnerID...). Error responses are plain text with the error message.
//
// # Routes
//
//	Method  Route              Description
//	GET     /health            Health check
//	POST    /wallet            Creates a wallet, creating its owner when it does not exist yet
//	GET     /wallets           Lists wallets. Filters: id, status, walletType, assetId, ownerId
//	GET     /wallet/:id        Gets a wallet with its asset, owner and external account
//	DELETE  /wallet/:id        Deactivates the wallet (INACTIVE), preserving its history
//	POST    /transaction       Creates and processes a transaction
//	GET     /transactions      Lists transactions. Filters: id, status, transactionType, sourceWalletId, destinationWalletId
//	GET     /transaction/:id   Gets a transaction with its wallets and their owners
//
// # Creating a wallet
//
//	curl -X POST http://localhost:3000/wallet \
//	  -H 'Content-Type: application/json' \
//	  -d '{
//	    "owner":  { "uniqueIdentifier": "12345678909", "ownerType": "INDIVIDUAL" },
//	    "wallet": { "status": "ACTIVE", "walletType": "INVESTMENT", "assetCode": "BRL" }
//	  }'
//
// Request fields:
//
//   - ownerType: INDIVIDUAL or COMPANY.
//   - walletType: INVESTMENT or CRYPTO. CHECKING and SAVINGS wallets require an external account (AccountID),
//     which cannot be informed through the API yet.
//   - assetCode must reference an asset already registered in the database.
//   - Every wallet starts with its three balances (free, reserved and hold) set to zero.
//
// # Creating a transaction
//
//	curl -X POST http://localhost:3000/transaction \
//	  -H 'Content-Type: application/json' \
//	  -d '{
//	    "amount": 10000,
//	    "transactionType": "TRANSFER",
//	    "sourceWalletId": "<source-wallet-uuid>",
//	    "destinationWalletId": "<destination-wallet-uuid>",
//	    "description": "Rent payment",
//	    "referenceId": "<external-operation-id>"
//	  }'
//
// Request fields:
//
//   - amount: integer in the smallest unit of the asset (10000 = BRL 100.00). For REVERSAL it may be omitted
//     (the original amount is used) or must be equal to the original amount.
//   - transactionType: one of the types in the table below.
//   - sourceWalletId / destinationWalletId: source and destination wallets, according to the table below.
//   - reversalOriginId: required only for REVERSAL; the transaction to be reversed.
//   - description / referenceId: optional description and identifier of the external operation that
//     originated the transaction.
//
// Wallets accepted by each transaction type and its balance effect. Informed wallets must be active, distinct
// and hold the same asset, but they may have different wallet types and owners:
//
//	Type              Source                                  Destination                    Balance effect
//	DEPOSIT           not accepted (external origin)          required                       destination free += amount
//	WITHDRAWAL        required                                not accepted (external target) source free -= amount
//	TRANSFER          required                                required                       source free -> destination free, through hold
//	YIELD             optional (wallet that generated it,     required                       destination free += amount
//	                  not debited)
//	SALE              optional (buyer)                        required                       transfer with source; destination free += amount without it
//	FEE               required                                optional (revenue wallet)      transfer with destination; source free -= amount without it
//	PURCHASE          required                                optional (seller)              transfer with destination; source free -= amount without it
//	RESERVE/RELEASE   required                                not accepted                   free -> reserved / reserved -> free
//	HOLD/UNHOLD       required                                not accepted                   free -> hold / hold -> free
//	REVERSAL          not accepted (from the original)        not accepted (from original)   undoes the original effect
//
// # Reversing a transaction
//
// A REVERSAL is accepted only when the original transaction is COMPLETED, is not a reversal itself and was not
// reversed before. The reversal swaps the source and destination wallets of the original transaction; balance
// state changes (RESERVE, RELEASE, HOLD, UNHOLD) keep the same wallet, and the reversal of a YIELD leaves out the
// wallet that generated it, since that wallet was never debited.
//
//	curl -X POST http://localhost:3000/transaction \
//	  -H 'Content-Type: application/json' \
//	  -d '{ "transactionType": "REVERSAL", "reversalOriginId": "<transaction-uuid>", "description": "Wrong recipient" }'
//
// # Transaction errors
//
//	Status  When
//	400     Invalid input or wallets for the type, inactive wallet, different assets, reversal not allowed
//	404     Wallet not found
//	422     Balance rule violated, such as insufficient balance; the transaction is recorded as FAILED
//
// Inactive wallets and wallets with different assets also leave the transaction recorded as FAILED.
// The wallets involved are locked with SELECT ... FOR UPDATE, always in the same order, and the original
// transaction of a reversal is locked as well, so concurrent operations neither leave inconsistent balances
// nor reverse the same transaction twice.
package api
