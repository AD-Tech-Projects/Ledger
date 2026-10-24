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
	accountService     *services.AccountService
}

func NewTransactionHandler(
	transactionService *services.TransactionService,
	accountService *services.AccountService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		accountService:     accountService,
	}
}

type createTransactionRequestBody struct {
	Amount               int64                    `json:"amount"`
	Status               domain.TransactionStatus `json:"status"`
	TransactionType      domain.TransactionType   `json:"transactionType"`
	SourceAccountId      string                   `json:"sourceAccountId"`
	DestinationAccountId string                   `json:"destinationAccountId"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	log := logger.WithContext("TransactionHandler.create")
	context := c.Request.Context()

	transactionPayload, ok := decodeRequestBody[createTransactionRequestBody](c)
	if !ok {
		return
	}

	sourceAccount, err := getAccountById(c, log, h.accountService, transactionPayload.SourceAccountId)
	if err != nil {
		return
	}
	destinationAccount, err := getAccountById(c, log, h.accountService, transactionPayload.DestinationAccountId)
	if err != nil {
		return
	}

	transaction := &domain.Transaction{
		Amount:             transactionPayload.Amount,
		Status:             transactionPayload.Status,
		TransactionType:    transactionPayload.TransactionType,
		SourceAccount:      sourceAccount,
		DestinationAccount: destinationAccount,
	}

	createdTransaction, err := h.transactionService.Create(context, transaction)
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
		Id:              c.Query("id"),
		Status:          domain.TransactionStatus(c.Query("status")),
		TransactionType: domain.TransactionType(c.Query("transactionType")),
	}
	if sourceAccountID := c.Query("sourceAccountId"); sourceAccountID != "" {
		filter.SourceAccount = &domain.Account{Id: sourceAccountID}
	}
	if destinationAccountID := c.Query("destinationAccountId"); destinationAccountID != "" {
		filter.DestinationAccount = &domain.Account{Id: destinationAccountID}
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

func getAccountById(c *gin.Context, log logger.ContextLogger, accountService *services.AccountService, accountId string) (*domain.Account, error) {
	account, err := accountService.GetById(c.Request.Context(), accountId)
	if err != nil {
		log.Error(err)
		if errors.Is(err, services.ErrInvalidAccount) {
			writeBadRequest(c, err)
			return nil, err
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			writeNotFound(c, err)
			return nil, err
		}

		writeInternalServerError(c, err)
		return nil, err
	}

	return account, nil
}
