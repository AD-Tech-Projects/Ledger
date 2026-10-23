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
	Holder  domain.AccountHolder `json:"holder"`
	Account domain.Account       `json:"account"`
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
		http.Error(w, "failed to create account holder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	account := body.Account
	account.HolderId = createdHolder.Id

	createdAccount, err := h.accountService.Create(context, &account)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	foundAccount, err := h.accountService.GetById(context, createdAccount.Id)
	if err != nil {
		http.Error(w, "failed to retrieve created account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, http.StatusCreated, foundAccount)
}
