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

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusCreated, createdTransaction)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.transactionService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	transaction, err := h.transactionService.GetById(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, services.ErrInvalidTransaction) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrTransactionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusOK, transaction)
}
