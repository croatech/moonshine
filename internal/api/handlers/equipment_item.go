package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/middleware"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type EquipmentItemHandler struct {
	db                          *sqlx.DB
	equipmentItemService        *services.EquipmentItemService
	equipmentItemBuyService     *services.EquipmentItemBuyService
	equipmentItemSellService    *services.EquipmentItemSellService
	equipmentItemTakeOnService  *services.EquipmentItemTakeOnService
	equipmentItemTakeOffService *services.EquipmentItemTakeOffService
}

func NewEquipmentItemHandler(db *sqlx.DB) *EquipmentItemHandler {
	equipmentItemRepo := repository.NewEquipmentItemRepository(db)
	equipmentItemService := services.NewEquipmentItemService(equipmentItemRepo)

	userEquipmentItemRepo := repository.NewUserEquipmentItemRepository(db)
	userRepo := repository.NewUserRepository(db)
	equipmentItemBuyService := services.NewEquipmentItemBuyService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)
	equipmentItemSellService := services.NewEquipmentItemSellService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)
	equipmentItemTakeOnService := services.NewEquipmentItemTakeOnService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)
	equipmentItemTakeOffService := services.NewEquipmentItemTakeOffService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)

	return &EquipmentItemHandler{
		db:                          db,
		equipmentItemService:        equipmentItemService,
		equipmentItemBuyService:     equipmentItemBuyService,
		equipmentItemSellService:    equipmentItemSellService,
		equipmentItemTakeOnService:  equipmentItemTakeOnService,
		equipmentItemTakeOffService: equipmentItemTakeOffService,
	}
}

func (h *EquipmentItemHandler) GetEquipmentItems(c echo.Context) error {
	category := c.QueryParam("category")
	if category == "" {
		return ErrBadRequest(c, "category parameter is required")
	}

	items, err := h.equipmentItemService.GetByCategorySlug(c.Request().Context(), category)
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, dto.EquipmentItemsFromDomain(items))
}

func (h *EquipmentItemHandler) BuyEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		return ErrBadRequest(c, "item slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	err = h.equipmentItemBuyService.BuyEquipmentItem(c.Request().Context(), userID, itemSlug)
	if err != nil {
		switch err {
		case services.ErrEquipmentItemNotFound:
			return ErrNotFound(c, "equipment item not found")
		case services.ErrInsufficientGold:
			return ErrBadRequest(c, "insufficient gold")
		case repository.ErrUserNotFound:
			return ErrNotFound(c, "user not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item purchased successfully"})
}

func (h *EquipmentItemHandler) TakeOnEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		return ErrBadRequest(c, "item slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	equipmentItemRepo := repository.NewEquipmentItemRepository(h.db)
	item, err := equipmentItemRepo.FindBySlug(itemSlug)
	if err != nil {
		return ErrNotFound(c, "equipment item not found")
	}

	err = h.equipmentItemTakeOnService.TakeOnEquipmentItem(c.Request().Context(), userID, item.ID)
	if err != nil {
		switch err {
		case services.ErrEquipmentItemNotFound:
			return ErrNotFound(c, "equipment item not found")
		case services.ErrItemNotInInventory:
			return ErrBadRequest(c, "item not in inventory")
		case services.ErrInsufficientLevel:
			return ErrBadRequest(c, "insufficient level")
		case services.ErrInvalidEquipmentType:
			return ErrBadRequest(c, "invalid equipment type")
		case repository.ErrUserNotFound:
			return ErrNotFound(c, "user not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item equipped successfully"})
}

func (h *EquipmentItemHandler) TakeOffEquipmentItem(c echo.Context) error {
	slotName := c.Param("slot")
	if slotName == "" {
		return ErrBadRequest(c, "slot name is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	err = h.equipmentItemTakeOffService.TakeOffEquipmentItem(c.Request().Context(), userID, slotName)
	if err != nil {
		switch err {
		case services.ErrNoItemEquipped:
			return ErrBadRequest(c, "no item equipped in this slot")
		case services.ErrInvalidEquipmentType:
			return ErrBadRequest(c, "invalid slot name")
		case repository.ErrUserNotFound:
			return ErrNotFound(c, "user not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item removed successfully"})
}

func (h *EquipmentItemHandler) SellEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		return ErrBadRequest(c, "item slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	err = h.equipmentItemSellService.SellEquipmentItem(c.Request().Context(), userID, itemSlug)
	if err != nil {
		switch err {
		case services.ErrItemNotOwned:
			return ErrBadRequest(c, "item not owned")
		case services.ErrEquipmentItemNotFound:
			return ErrNotFound(c, "equipment item not found")
		case repository.ErrUserNotFound:
			return ErrNotFound(c, "user not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item sold successfully"})
}
