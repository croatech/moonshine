package handlers

import (
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
	locationRepo := repository.NewLocationRepository(db)
	userService := services.NewUserService(userRepo, avatarRepo, locationRepo)

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := services.NewInventoryService(inventoryRepo)

	return &UserHandler{
		db:               db,
		userService:      userService,
		inventoryService: inventoryService,
	}
}

func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	user, avatar, location, inFight, err := h.userService.GetCurrentUserWithRelations(c.Request().Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ErrNotFound(c, "user not found")
		}
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, dto.UserFromDomain(user, avatar, location, inFight))
}

func (h *UserHandler) GetUserInventory(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	items, err := h.inventoryService.GetUserInventory(c.Request().Context(), userID)
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, dto.EquipmentItemsFromDomain(items))
}

func (h *UserHandler) GetUserEquippedItems(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound(c, "user not found")
	}

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

	equipmentItemRepo := repository.NewEquipmentItemRepository(h.db)
	equippedItems := make(map[string]*dto.EquipmentItem)

	for slot, itemIDPtr := range itemIDs {
		if itemIDPtr != nil {
			itemID, err := uuid.Parse(*itemIDPtr)
			if err != nil {
				continue
			}
			item, err := equipmentItemRepo.FindByID(itemID)
			if err != nil {
				continue
			}
			equippedItems[slot] = dto.EquipmentItemFromDomain(item)
		}
	}

	return c.JSON(http.StatusOK, equippedItems)
}

func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return ErrBadRequest(c, "invalid request")
	}

	var avatarID *uuid.UUID
	if req.AvatarID != nil {
		parsedID, err := uuid.Parse(*req.AvatarID)
		if err != nil {
			return ErrBadRequest(c, "invalid avatar ID")
		}
		avatarID = &parsedID
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), userID, avatarID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ErrNotFound(c, "user not found")
		}
		if err == repository.ErrAvatarNotFound {
			return ErrNotFound(c, "avatar not found")
		}
		return ErrInternalServerError(c)
	}

	user, avatar, location, inFight, err := h.userService.GetCurrentUserWithRelations(c.Request().Context(), userID)
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, dto.UserFromDomain(user, avatar, location, inFight))
}
