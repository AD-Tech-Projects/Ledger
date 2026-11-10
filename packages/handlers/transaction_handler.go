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
	SourceWalletID      *string                `json:"sourceWalletId"`
	DestinationWalletID *string                `json:"destinationWalletId"`
	ReversalOriginID    *string                `json:"reversalOriginId"`
	ReferenceID         *string                `json:"referenceId"`
	Description         string                 `json:"description"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.create")
	context := c.Request.Context()
	log.Info("creating transaction")

	transactionPayload, ok := decodeRequestBody[createTransactionRequestBody](c)
	if !ok {
		log.Warn("transaction creation rejected: invalid request body")
		return
	}

	createdTransaction, err := h.transactionService.Create(context, services.CreateTransactionInput{
		Amount:              transactionPayload.Amount,
		Type:                transactionPayload.TransactionType,
		SourceWalletID:      transactionPayload.SourceWalletID,
		DestinationWalletID: transactionPayload.DestinationWalletID,
		ReversalOriginID:    transactionPayload.ReversalOriginID,
		ReferenceID:         transactionPayload.ReferenceID,
		Description:         transactionPayload.Description,
	})
	if err != nil {
		log.Error(err)
		switch {
		case errors.Is(err, services.ErrInvalidTransaction):
			writeBadRequest(c, err)
		case errors.Is(err, services.ErrWalletNotFound):
			writeNotFound(c, err)
		case errors.Is(err, services.ErrTransactionFailed):
			writeUnprocessableEntity(c, err)
		default:
			writeInternalServerError(c, err)
		}
		return
	}

	log.Info("transaction created successfully")
	encodeJSONResponse(c, http.StatusCreated, createdTransaction)
}

func (h *TransactionHandler) List(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.list")
	log.Info("listing transactions")
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

	log.Info("transactions listed successfully")
	encodeJSONResponse(c, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetById(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.getById")
	log.Info("getting transaction by id")
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

	log.Info("transaction retrieved successfully")
	encodeJSONResponse(c, http.StatusOK, transaction)
}
