package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/middleware"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type UserHandler struct {
	db               *sqlx.DB
	userService      *services.UserService
	inventoryService *services.InventoryService
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewAvatarRepository(db)
	userService := services.NewUserService(userRepo, avatarRepo)

	inventoryService := services.NewInventoryService(db)

	return &UserHandler{
		db:               db,
		userService:      userService,
		inventoryService: inventoryService,
	}
}

// GetCurrentUser returns the currently authenticated user
// GET /api/user/me
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user, err := h.userService.GetCurrentUser(c.Request().Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, dto.UserFromDomain(user))
}

// GetUserInventory returns all equipment items in the user's inventory
// GET /api/users/me/inventory
func (h *UserHandler) GetUserInventory(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[UserHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	log.Printf("[UserHandler] Fetching inventory for user ID: %s", userID)
	items, err := h.inventoryService.GetUserInventory(c.Request().Context(), userID)
	if err != nil {
		log.Printf("[UserHandler] Error fetching inventory: %+v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	log.Printf("[UserHandler] Found %d items in inventory for user: %s", len(items), userID)
	return c.JSON(http.StatusOK, dto.EquipmentItemsFromDomain(items))
}

// GetUserEquippedItems returns all equipment items currently equipped by the user
// GET /api/users/me/equipped
func (h *UserHandler) GetUserEquippedItems(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[UserHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	log.Printf("[UserHandler] Fetching equipped items for user ID: %s", userID)
	
	// Get user to get equipment IDs
	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[UserHandler] User not found: %+v", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	// Collect all equipped item IDs
	itemIDs := make(map[string]*string)
	if user.ChestEquipmentItemID != nil {
		id := user.ChestEquipmentItemID.String()
		itemIDs["chest"] = &id
	}
	if user.BeltEquipmentItemID != nil {
		id := user.BeltEquipmentItemID.String()
		itemIDs["belt"] = &id
	}
	if user.HeadEquipmentItemID != nil {
		id := user.HeadEquipmentItemID.String()
		itemIDs["head"] = &id
	}
	if user.NeckEquipmentItemID != nil {
		id := user.NeckEquipmentItemID.String()
		itemIDs["neck"] = &id
	}
	if user.WeaponEquipmentItemID != nil {
		id := user.WeaponEquipmentItemID.String()
		itemIDs["weapon"] = &id
	}
	if user.ShieldEquipmentItemID != nil {
		id := user.ShieldEquipmentItemID.String()
		itemIDs["shield"] = &id
	}
	if user.LegsEquipmentItemID != nil {
		id := user.LegsEquipmentItemID.String()
		itemIDs["legs"] = &id
	}
	if user.FeetEquipmentItemID != nil {
		id := user.FeetEquipmentItemID.String()
		itemIDs["feet"] = &id
	}
	if user.ArmsEquipmentItemID != nil {
		id := user.ArmsEquipmentItemID.String()
		itemIDs["arms"] = &id
	}
	if user.HandsEquipmentItemID != nil {
		id := user.HandsEquipmentItemID.String()
		itemIDs["hands"] = &id
	}
	if user.Ring1EquipmentItemID != nil {
		id := user.Ring1EquipmentItemID.String()
		itemIDs["ring1"] = &id
	}
	if user.Ring2EquipmentItemID != nil {
		id := user.Ring2EquipmentItemID.String()
		itemIDs["ring2"] = &id
	}
	if user.Ring3EquipmentItemID != nil {
		id := user.Ring3EquipmentItemID.String()
		itemIDs["ring3"] = &id
	}
	if user.Ring4EquipmentItemID != nil {
		id := user.Ring4EquipmentItemID.String()
		itemIDs["ring4"] = &id
	}

	// Fetch actual items
	equipmentItemRepo := repository.NewEquipmentItemRepository(h.db)
	equippedItems := make(map[string]*dto.EquipmentItem)
	
	for slot, itemIDPtr := range itemIDs {
		if itemIDPtr != nil {
			itemID, err := uuid.Parse(*itemIDPtr)
			if err != nil {
				log.Printf("[UserHandler] Invalid item ID for slot %s: %v", slot, err)
				continue
			}
			item, err := equipmentItemRepo.FindByID(itemID)
			if err != nil {
				log.Printf("[UserHandler] Item not found for slot %s: %v", slot, err)
				continue
			}
			equippedItems[slot] = dto.EquipmentItemFromDomain(item)
		}
	}

	log.Printf("[UserHandler] Found %d equipped items for user: %s", len(equippedItems), userID)
	return c.JSON(http.StatusOK, equippedItems)
}

// UpdateCurrentUser updates the currently authenticated user
// PUT /api/user/me
func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[UserHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[UserHandler] Bad Request: failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	var avatarID *uuid.UUID
	if req.AvatarID != nil {
		parsedID, err := uuid.Parse(*req.AvatarID)
		if err != nil {
			log.Printf("[UserHandler] Bad Request: invalid avatar ID: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid avatar ID"})
		}
		avatarID = &parsedID
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), userID, avatarID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			log.Printf("[UserHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		if err == repository.ErrAvatarNotFound {
			log.Printf("[UserHandler] Not Found: avatar not found")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "avatar not found"})
		}
		log.Printf("[UserHandler] Internal Server Error: failed to update user: %+v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	log.Printf("[UserHandler] Successfully updated user %s", userID)
	return c.JSON(http.StatusOK, dto.UserFromDomain(user))
}
