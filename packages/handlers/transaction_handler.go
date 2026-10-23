package handlers

import (
	"errors"
	"net/http"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	context := c.Request.Context()

	account, ok := decodeRequestBody[domain.Transaction](c)
	if !ok {
		return
	}

	createdTransaction, err := h.transactionService.Create(context, account)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTransaction) {
			writeBadRequest(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusCreated, createdTransaction)
}

func (h *TransactionHandler) List(c *gin.Context) {
	transactions, err := h.transactionService.GetAll(c.Request.Context())
	if err != nil {
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetById(c *gin.Context) {
	transaction, err := h.transactionService.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrInvalidTransaction) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrTransactionNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, transaction)
}
