package handlers

import (
	"log"
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

// GetEquipmentItems returns equipment items filtered by category
// GET /api/equipment_items?category=slug
func (h *EquipmentItemHandler) GetEquipmentItems(c echo.Context) error {
	category := c.QueryParam("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "category parameter is required"})
	}

	log.Printf("[EquipmentItemHandler] Fetching items for category: %s", category)
	items, err := h.equipmentItemService.GetByCategorySlug(c.Request().Context(), category)
	if err != nil {
		log.Printf("[EquipmentItemHandler] Error fetching items for category %s: %+v", category, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	log.Printf("[EquipmentItemHandler] Found %d items for category: %s", len(items), category)
	return c.JSON(http.StatusOK, dto.EquipmentItemsFromDomain(items))
}

// BuyEquipmentItem handles equipment item purchase
// POST /api/equipment_items/:slug/buyЩ
func (h *EquipmentItemHandler) BuyEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		log.Printf("[EquipmentItemHandler] Bad Request: empty item slug")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "item slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[EquipmentItemHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err = h.equipmentItemBuyService.BuyEquipmentItem(c.Request().Context(), userID, itemSlug)
	if err != nil {
		switch err {
		case services.ErrEquipmentItemNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: equipment item with slug %s not found", itemSlug)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "equipment item not found"})
		case services.ErrInsufficientGold:
			log.Printf("[EquipmentItemHandler] Bad Request: insufficient gold for user %s to buy item %s", userID, itemSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient gold"})
		case services.ErrAlreadyOwned:
			log.Printf("[EquipmentItemHandler] Bad Request: user %s already owns item %s", userID, itemSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "item already owned"})
		case repository.ErrUserNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Printf("[EquipmentItemHandler] Internal Server Error: failed to buy equipment item: %+v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item purchased successfully"})
}

// TakeOnEquipmentItem handles equipping an item from inventory
// POST /api/equipment_items/:slug/take_on
func (h *EquipmentItemHandler) TakeOnEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		log.Printf("[EquipmentItemHandler] Bad Request: empty item slug")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "item slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[EquipmentItemHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Get item by slug to get its ID
	equipmentItemRepo := repository.NewEquipmentItemRepository(h.db)
	item, err := equipmentItemRepo.FindBySlug(itemSlug)
	if err != nil {
		log.Printf("[EquipmentItemHandler] Not Found: equipment item with slug %s not found", itemSlug)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "equipment item not found"})
	}

	err = h.equipmentItemTakeOnService.TakeOnEquipmentItem(c.Request().Context(), userID, item.ID)
	if err != nil {
		switch err {
		case services.ErrEquipmentItemNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: equipment item with slug %s not found", itemSlug)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "equipment item not found"})
		case services.ErrItemNotInInventory:
			log.Printf("[EquipmentItemHandler] Bad Request: item %s not in inventory for user %s", itemSlug, userID)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "item not in inventory"})
		case services.ErrInsufficientLevel:
			log.Printf("[EquipmentItemHandler] Bad Request: insufficient level for user %s to equip item %s", userID, itemSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient level"})
		case services.ErrInvalidEquipmentType:
			log.Printf("[EquipmentItemHandler] Bad Request: invalid equipment type for item %s", itemSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid equipment type"})
		case repository.ErrUserNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Printf("[EquipmentItemHandler] Internal Server Error: failed to equip item: %+v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item equipped successfully"})
}

// TakeOffEquipmentItem handles removing an item from equipment slot
// POST /api/equipment_items/take_off/:slot
func (h *EquipmentItemHandler) TakeOffEquipmentItem(c echo.Context) error {
	slotName := c.Param("slot")
	if slotName == "" {
		log.Printf("[EquipmentItemHandler] Bad Request: empty slot name")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slot name is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[EquipmentItemHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err = h.equipmentItemTakeOffService.TakeOffEquipmentItem(c.Request().Context(), userID, slotName)
	if err != nil {
		switch err {
		case services.ErrNoItemEquipped:
			log.Printf("[EquipmentItemHandler] Bad Request: no item in slot %s for user %s", slotName, userID)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "no item equipped in this slot"})
		case services.ErrInvalidEquipmentType:
			log.Printf("[EquipmentItemHandler] Bad Request: invalid slot name %s", slotName)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid slot name"})
		case repository.ErrUserNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Printf("[EquipmentItemHandler] Internal Server Error: failed to take off item: %+v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item removed successfully"})
}

// SellEquipmentItem handles selling an item from inventory
// POST /api/equipment_items/:slug/sell
func (h *EquipmentItemHandler) SellEquipmentItem(c echo.Context) error {
	itemSlug := c.Param("slug")
	if itemSlug == "" {
		log.Printf("[EquipmentItemHandler] Bad Request: empty item slug")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "item slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[EquipmentItemHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err = h.equipmentItemSellService.SellEquipmentItem(c.Request().Context(), userID, itemSlug)
	if err != nil {
		switch err {
		case services.ErrItemNotOwned:
			log.Printf("[EquipmentItemHandler] Bad Request: user %s does not own item %s", userID, itemSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "item not owned"})
		case services.ErrEquipmentItemNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: equipment item %s not found", itemSlug)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "equipment item not found"})
		case repository.ErrUserNotFound:
			log.Printf("[EquipmentItemHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Printf("[EquipmentItemHandler] Internal Server Error: failed to sell item: %+v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "item sold successfully"})
}
