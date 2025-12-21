package domain

import "github.com/google/uuid"

type User struct {
	Model
	Attack                uint                 `json:"attack"`
	Avatar                *Avatar              `json:"avatar,omitempty"`
	AvatarID              *uuid.UUID           `json:"avatar_id"`
	CurrentHp             uint                 `json:"current_hp"`
	Defense               uint                 `json:"defense"`
	Email                 string               `json:"email"`
	Events                []*Event             `json:"events,omitempty"`
	Exp                   uint                 `json:"exp"`
	Fights                []*Fight             `json:"fights,omitempty"`
	FishingSkill          uint                 `json:"fishing_skill"`
	FishingSlot           uint                 `json:"fishing_slot"`
	FreeStats             uint                 `json:"free_stats"`
	Gold                  uint                 `json:"gold"`
	Hp                    uint                 `json:"hp"`
	Level                 uint                 `json:"level"`
	Location              *Location            `json:"location,omitempty"`
	LocationID            uuid.UUID            `json:"location_id"`
	LumberjackingSkill    uint                 `json:"lumberjacking_skill"`
	LumberjackingSlot     uint                 `json:"lumberjacking_slot"`
	Messages              []*Message           `json:"messages,omitempty"`
	Movements             []*Movement          `json:"movements,omitempty"`
	Name                  string               `json:"name"`
	Password              string               `json:"-"`
	Stuffs                []*Stuff             `json:"stuffs,omitempty"`
	Tools                 []*UserToolItem      `json:"tools,omitempty"`
	Username              string               `json:"username"`
	Equipment             []*UserEquipmentItem `json:"equipment,omitempty"`
	ChestEquipmentItemID  *uuid.UUID           `json:"chest_equipment_item_id"`
	BeltEquipmentItemID   *uuid.UUID           `json:"belt_equipment_item_id"`
	HeadEquipmentItemID   *uuid.UUID           `json:"head_equipment_item_id"`
	NeckEquipmentItemID   *uuid.UUID           `json:"neck_equipment_item_id"`
	WeaponEquipmentItemID *uuid.UUID           `json:"weapon_equipment_item_id"`
	ShieldEquipmentItemID *uuid.UUID           `json:"shield_equipment_item_id"`
	LegsEquipmentItemID   *uuid.UUID           `json:"legs_equipment_item_id"`
	FeetEquipmentItemID   *uuid.UUID           `json:"feet_equipment_item_id"`
	ArmsEquipmentItemID   *uuid.UUID           `json:"arms_equipment_item_id"`
	HandsEquipmentItemID  *uuid.UUID           `json:"hands_equipment_item_id"`
	Ring1EquipmentItemID  *uuid.UUID           `json:"ring1_equipment_item_id"`
	Ring2EquipmentItemID  *uuid.UUID           `json:"ring2_equipment_item_id"`
	Ring3EquipmentItemID  *uuid.UUID           `json:"ring3_equipment_item_id"`
	Ring4EquipmentItemID  *uuid.UUID           `json:"ring4_equipment_item_id"`
}
