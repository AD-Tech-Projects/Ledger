package handlers

import (
	"errors"
	"net/http"

	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
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
	Holder  domain.AccountHolder     `json:"holder"`
	Account createAccountRequestData `json:"account"`
}

type createAccountRequestData struct {
	Status      domain.AccountStatus   `json:"status"`
	AccountType domain.AccountType     `json:"accountType"`
	Currency    domain.AccountCurrency `json:"currency"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	body, ok := decodeRequestBody[createAccountRequestBody](w, r.Body)
	if !ok {
		return
	}

	holder := body.Holder

	createdHolder, err := h.holderService.Create(context, &holder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account := domain.Account{
		Status:      body.Account.Status,
		AccountType: body.Account.AccountType,
		Currency:    body.Account.Currency,
	}
	account.HolderId = createdHolder.Id

	createdAccount, err := h.accountService.Create(context, &account)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	foundAccount, err := h.accountService.GetById(context, createdAccount.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusCreated, foundAccount)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusOK, accounts)
}

func (h *AccountHandler) GetById(w http.ResponseWriter, r *http.Request) {
	account, err := h.accountService.GetById(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusOK, account)
}

type updateAccountRequestBody struct {
	Status      *domain.AccountStatus   `json:"status"`
	AccountType *domain.AccountType     `json:"accountType"`
	Currency    *domain.AccountCurrency `json:"currency"`
	HolderId    *string                 `json:"holderId"`
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	body, ok := decodeRequestBody[updateAccountRequestBody](w, r.Body)
	if !ok {
		return
	}

	account, err := h.accountService.Update(r.Context(), r.PathValue("id"), services.AccountPatch{
		Status:      body.Status,
		AccountType: body.AccountType,
		Currency:    body.Currency,
		HolderId:    body.HolderId,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusOK, account)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.accountService.Deactivate(r.Context(), r.PathValue("id")); err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, services.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
