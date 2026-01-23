package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Model
	UpdatedAt             time.Time  `db:"updated_at"`
	Attack                uint16     `db:"attack"`
	AvatarID              *uuid.UUID `db:"avatar_id"`
	CurrentHp             uint16     `db:"current_hp"`
	Defense               uint16     `db:"defense"`
	Email                 string     `db:"email"`
	Exp                   uint       `db:"exp"`
	FreeStats             uint8      `db:"free_stats"`
	Gold                  uint       `db:"gold"`
	Hp                    uint16     `db:"hp"`
	Level                 uint8      `db:"level"`
	LocationID            uuid.UUID  `db:"location_id"`
	Name                  string     `db:"name"`
	Password              string     `db:"password"`
	Username              string     `db:"username"`
	ChestEquipmentItemID  *uuid.UUID `db:"chest_equipment_item_id"`
	BeltEquipmentItemID   *uuid.UUID `db:"belt_equipment_item_id"`
	HeadEquipmentItemID   *uuid.UUID `db:"head_equipment_item_id"`
	NeckEquipmentItemID   *uuid.UUID `db:"neck_equipment_item_id"`
	WeaponEquipmentItemID *uuid.UUID `db:"weapon_equipment_item_id"`
	ShieldEquipmentItemID *uuid.UUID `db:"shield_equipment_item_id"`
	LegsEquipmentItemID   *uuid.UUID `db:"legs_equipment_item_id"`
	FeetEquipmentItemID   *uuid.UUID `db:"feet_equipment_item_id"`
	ArmsEquipmentItemID   *uuid.UUID `db:"arms_equipment_item_id"`
	HandsEquipmentItemID  *uuid.UUID `db:"hands_equipment_item_id"`
	Ring1EquipmentItemID  *uuid.UUID `db:"ring1_equipment_item_id"`
	Ring2EquipmentItemID  *uuid.UUID `db:"ring2_equipment_item_id"`
	Ring3EquipmentItemID  *uuid.UUID `db:"ring3_equipment_item_id"`
	Ring4EquipmentItemID  *uuid.UUID `db:"ring4_equipment_item_id"`
	Avatar                string     `db:"avatar"`
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
	requiredExp, exists := levelMatrix[uint(user.Level)]
	if !exists {
		return false
	}
	return user.Exp >= requiredExp
}

func (user *User) RegenerateHealth(percent float64) uint16 {
	if user.CurrentHp >= user.Hp {
		return user.Hp
	}

	regeneration := uint16(float64(user.Hp) * percent / 100.0)

	minRegeneration := uint16(5)
	if regeneration < minRegeneration {
		regeneration = minRegeneration
	}

	newHp := user.CurrentHp + regeneration

	if newHp > user.Hp {
		return user.Hp
	}

	return newHp
}
