package handlers

import (
	"errors"
	"net/http"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	holderService  *services.AccountHolderService
	accountService *services.AccountService
}

func NewAccountHandler(
	accountHolderService *services.AccountHolderService,
	accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		holderService:  accountHolderService,
		accountService: accountService,
	}
}

type createAccountRequestBody struct {
	Holder  createHolderRequestData  `json:"holder"`
	Account createAccountRequestData `json:"account"`
}

type createHolderRequestData struct {
	UniqueIdentifier string            `json:"uniqueIdentifier"`
	HolderType       domain.HolderType `json:"holderType"`
}

type createAccountRequestData struct {
	Status      domain.AccountStatus   `json:"status"`
	AccountType domain.AccountType     `json:"accountType"`
	Currency    domain.AccountCurrency `json:"currency"`
}

func (h *AccountHandler) Create(c *gin.Context) {
	context := c.Request.Context()

	body, ok := decodeRequestBody[createAccountRequestBody](c)
	if !ok {
		return
	}

	holder := domain.AccountHolder{
		UniqueIdentifier: body.Holder.UniqueIdentifier,
		HolderType:       body.Holder.HolderType,
	}

	createdHolder, err := h.holderService.Create(context, &holder)
	if err != nil {
		writeInternalServerError(c, err)
		return
	}

	account := domain.Account{
		Status:      body.Account.Status,
		AccountType: body.Account.AccountType,
		Currency:    body.Account.Currency,
		HolderId:    createdHolder.Id,
	}

	createdAccount, err := h.accountService.Create(context, &account)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			writeBadRequest(c, err)
			return
		}
		writeInternalServerError(c, err)
		return
	}

	foundAccount, err := h.accountService.GetById(context, createdAccount.Id)
	if err != nil {
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusCreated, foundAccount)
}

func (h *AccountHandler) List(c *gin.Context) {
	accounts, err := h.accountService.GetAll(c.Request.Context())
	if err != nil {
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, accounts)
}

func (h *AccountHandler) GetById(c *gin.Context) {
	account, err := h.accountService.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, account)
}

type updateAccountRequestBody struct {
	Status      *domain.AccountStatus   `json:"status"`
	AccountType *domain.AccountType     `json:"accountType"`
	Currency    *domain.AccountCurrency `json:"currency"`
	HolderId    *string                 `json:"holderId"`
}

func (h *AccountHandler) Update(c *gin.Context) {
	body, ok := decodeRequestBody[updateAccountRequestBody](c)
	if !ok {
		return
	}

	account, err := h.accountService.Update(c.Request.Context(), c.Param("id"), services.AccountPatch{
		Status:      body.Status,
		AccountType: body.AccountType,
		Currency:    body.Currency,
		HolderId:    body.HolderId,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	if err := h.accountService.Deactivate(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
