package handlers

import (
	"errors"
	"net/http"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(
	transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

type createTransactionRequestBody struct {
	Amount              int64                  `json:"amount"`
	TransactionType     domain.TransactionType `json:"transactionType"`
	SourceWalletID      string                 `json:"sourceWalletId"`
	DestinationWalletID string                 `json:"destinationWalletId"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.create")
	context := c.Request.Context()

	transactionPayload, ok := decodeRequestBody[createTransactionRequestBody](c)
	if !ok {
		return
	}

	createdTransaction, err := h.transactionService.Create(context,
		transactionPayload.Amount,
		transactionPayload.TransactionType,
		transactionPayload.SourceWalletID,
		transactionPayload.DestinationWalletID)
	if err != nil {
		log.Error(err)
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
	log := logger.WithContext("TransactionHandler.list")
	filter := &domain.Transaction{
		ID:     c.Query("id"),
		Status: domain.TransactionStatus(c.Query("status")),
		Type:   domain.TransactionType(c.Query("transactionType")),
	}
	if sourceWalletID := c.Query("sourceWalletId"); sourceWalletID != "" {
		filter.SourceWallet = &domain.Wallet{ID: sourceWalletID}
	}
	if destinationWalletID := c.Query("destinationWalletId"); destinationWalletID != "" {
		filter.DestinationWallet = &domain.Wallet{ID: destinationWalletID}
	}
	transactions, err := h.transactionService.GetAll(c.Request.Context(), filter)
	if err != nil {
		log.Error(err)
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetById(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.getById")
	transaction, err := h.transactionService.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		log.Error(err)
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
