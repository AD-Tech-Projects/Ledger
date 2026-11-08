package handlers

import (
	"errors"
	"net/http"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	"github.com/AD-Tech-Projects/Ledger/packages/domain"
	"github.com/AD-Tech-Projects/Ledger/packages/services"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	ownerService  *services.OwnerService
	assetService  *services.AssetService
	walletService *services.WalletService
}

func NewWalletHandler(
	walletOwnerService *services.OwnerService,
	walletAssetService *services.AssetService,
	walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		ownerService:  walletOwnerService,
		assetService:  walletAssetService,
		walletService: walletService,
	}
}

type createOwnerRequestData struct {
	UniqueIdentifier string           `json:"uniqueIdentifier"`
	OwnerType        domain.OwnerType `json:"ownerType"`
}

type createWalletRequestData struct {
	Status     domain.WalletStatus `json:"status"`
	WalletType domain.WalletType   `json:"walletType"`
	AssetCode  domain.AssetCode    `json:"assetCode"`
}

type createWalletRequestBody struct {
	Owner  createOwnerRequestData  `json:"owner"`
	Wallet createWalletRequestData `json:"wallet"`
}

func (h *WalletHandler) Create(c *gin.Context) {
	log := logger.WithContext("WalletHandler.create")
	context := c.Request.Context()

	body, ok := decodeRequestBody[createWalletRequestBody](c)
	if !ok {
		return
	}

	asset, err := h.assetService.FindByCode(context, string(body.Wallet.AssetCode))
	if err != nil {
		log.Error(err)
		writeInternalServerError(c, err)
		return
	}

	owner := domain.Owner{
		UniqueIdentifier: body.Owner.UniqueIdentifier,
		Type:             body.Owner.OwnerType,
	}

	existingOrCreatedOwner, err := h.ownerService.FindByUniqueIdentifier(context, owner.UniqueIdentifier)
	if err != nil {
		if errors.Is(err, services.ErrOwnerNotFound) {
			existingOrCreatedOwner, err = h.ownerService.Create(context, &owner)
			if err != nil {
				log.Error(err)
				writeInternalServerError(c, err)
				return
			}
		} else {
			log.Error(err)
			writeInternalServerError(c, err)
			return
		}
	}

	wallet := domain.Wallet{
		Status:  body.Wallet.Status,
		Type:    body.Wallet.WalletType,
		AssetID: asset.ID,
		OwnerID: existingOrCreatedOwner.ID,
	}

	createdWallet, err := h.walletService.Create(context, &wallet)
	if err != nil {
		log.Error(err)
		if errors.Is(err, services.ErrInvalidWallet) {
			writeBadRequest(c, err)
			return
		}
		writeInternalServerError(c, err)
		return
	}

	foundWallet, err := h.walletService.GetById(context, createdWallet.ID)
	if err != nil {
		log.Error(err)
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusCreated, foundWallet)
}

func (h *WalletHandler) List(c *gin.Context) {
	log := logger.WithContext("WalletHandler.list")
	filter := &domain.Wallet{
		ID:      c.Query("id"),
		Status:  domain.WalletStatus(c.Query("status")),
		Type:    domain.WalletType(c.Query("walletType")),
		AssetID: c.Query("assetId"),
		OwnerID: c.Query("ownerId"),
	}
	wallets, err := h.walletService.GetAll(c.Request.Context(), filter)
	if err != nil {
		log.Error(err)
		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, wallets)
}

func (h *WalletHandler) GetById(c *gin.Context) {
	log := logger.WithContext("WalletHandler.getById")
	wallet, err := h.walletService.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		log.Error(err)
		if errors.Is(err, services.ErrInvalidWallet) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrWalletNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	encodeJSONResponse(c, http.StatusOK, wallet)
}

func (h *WalletHandler) Deactivate(c *gin.Context) {
	log := logger.WithContext("WalletHandler.deactivate")
	if err := h.walletService.Deactivate(c.Request.Context(), c.Param("id")); err != nil {
		log.Error(err)
		if errors.Is(err, services.ErrInvalidWallet) {
			writeBadRequest(c, err)
			return
		}
		if errors.Is(err, services.ErrWalletNotFound) {
			writeNotFound(c, err)
			return
		}

		writeInternalServerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
