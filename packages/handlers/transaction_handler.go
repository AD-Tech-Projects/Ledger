package handlers

import (
	"errors"
	"net/http"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	account, ok := decodeRequestBody[domain.Transaction](w, r.Body)
	if !ok {
		return
	}

	createdTransaction, err := h.transactionService.Create(context, account)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTransaction) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusCreated, createdTransaction)
}
