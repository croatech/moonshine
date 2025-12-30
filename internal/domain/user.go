package domain

import "github.com/google/uuid"

type User struct {
	Model
	Attack                uint                 `json:"attack" db:"attack"`
	Avatar                *Avatar              `json:"avatar,omitempty"`
	AvatarID              *uuid.UUID           `json:"avatar_id" db:"avatar_id"`
	CurrentHp             uint                 `json:"current_hp" db:"current_hp"`
	Defense               uint                 `json:"defense" db:"defense"`
	Email                 string               `json:"email" db:"email"`
	Events                []*Event             `json:"events,omitempty"`
	Exp                   uint                 `json:"exp" db:"exp"`
	Fights                []*Fight             `json:"fights,omitempty"`
	FishingSkill          uint                 `json:"fishing_skill" db:"fishing_skill"`
	FishingSlot           uint                 `json:"fishing_slot" db:"fishing_slot"`
	FreeStats             uint                 `json:"free_stats" db:"free_stats"`
	Gold                  uint                 `json:"gold" db:"gold"`
	Hp                    uint                 `json:"hp" db:"hp"`
	Level                 uint                 `json:"level" db:"level"`
	Location              *Location            `json:"location,omitempty"`
	LocationID            uuid.UUID            `json:"location_id" db:"location_id"`
	LumberjackingSkill    uint                 `json:"lumberjacking_skill" db:"lumberjacking_skill"`
	LumberjackingSlot     uint                 `json:"lumberjacking_slot" db:"lumberjacking_slot"`
	Messages              []*Message           `json:"messages,omitempty"`
	Movements             []*Movement          `json:"movements,omitempty"`
	Name                  string               `json:"name" db:"name"`
	Password              string               `json:"-" db:"password"`
	Stuffs                []*Stuff             `json:"stuffs,omitempty"`
	Tools                 []*UserToolItem      `json:"tools,omitempty"`
	Username              string               `json:"username" db:"username"`
	Equipment             []*UserEquipmentItem `json:"equipment,omitempty"`
	ChestEquipmentItemID  *uuid.UUID           `json:"chest_equipment_item_id" db:"chest_equipment_item_id"`
	BeltEquipmentItemID   *uuid.UUID           `json:"belt_equipment_item_id" db:"belt_equipment_item_id"`
	HeadEquipmentItemID   *uuid.UUID           `json:"head_equipment_item_id" db:"head_equipment_item_id"`
	NeckEquipmentItemID   *uuid.UUID           `json:"neck_equipment_item_id" db:"neck_equipment_item_id"`
	WeaponEquipmentItemID *uuid.UUID           `json:"weapon_equipment_item_id" db:"weapon_equipment_item_id"`
	ShieldEquipmentItemID *uuid.UUID           `json:"shield_equipment_item_id" db:"shield_equipment_item_id"`
	LegsEquipmentItemID   *uuid.UUID           `json:"legs_equipment_item_id" db:"legs_equipment_item_id"`
	FeetEquipmentItemID   *uuid.UUID           `json:"feet_equipment_item_id" db:"feet_equipment_item_id"`
	ArmsEquipmentItemID   *uuid.UUID           `json:"arms_equipment_item_id" db:"arms_equipment_item_id"`
	HandsEquipmentItemID  *uuid.UUID           `json:"hands_equipment_item_id" db:"hands_equipment_item_id"`
	Ring1EquipmentItemID  *uuid.UUID           `json:"ring1_equipment_item_id" db:"ring1_equipment_item_id"`
	Ring2EquipmentItemID  *uuid.UUID           `json:"ring2_equipment_item_id" db:"ring2_equipment_item_id"`
	Ring3EquipmentItemID  *uuid.UUID           `json:"ring3_equipment_item_id" db:"ring3_equipment_item_id"`
	Ring4EquipmentItemID  *uuid.UUID           `json:"ring4_equipment_item_id" db:"ring4_equipment_item_id"`
}

var levelMatrix = map[uint]uint{
	1:  100,
	2:  200,
	3:  400,
	4:  800,
	5:  1500,
	6:  3000,
	7:  5000,
	8:  10000,
	9:  15000,
	10: 20000,
}

func (user *User) ReachedNewLevel() bool {
	requiredExp, exists := levelMatrix[user.Level]
	if !exists {
		return false
	}
	return user.Exp >= requiredExp
}

func (user *User) RegenerateHealth(percent float64) uint {
	if user.CurrentHp >= user.Hp {
		return user.Hp
	}

	regeneration := uint(float64(user.Hp) * percent / 100.0)
	
	minRegeneration := uint(5)
	if regeneration < minRegeneration {
		regeneration = minRegeneration
	}

	newHp := user.CurrentHp + regeneration
	
	if newHp > user.Hp {
		return user.Hp
	}

	return newHp
}
