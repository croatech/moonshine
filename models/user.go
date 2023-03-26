package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username           string               `gorm:"type:varchar(100);unique;not null"`
	Email              string               `gorm:"type:varchar(100);unique;not null"`
	Password           string               `gorm:"type:varchar(255);not null"`
	Name               string               `json:"name"`
	LocationID         uint                 `json:"location_id"`
	Location           *Location            `json:"location"`
	HelmetSlot         uint                 `json:"helmet_slot"`
	ArmorSlot          uint                 `json:"armor_slot"`
	MailSlot           uint                 `json:"mail_slot"`
	GlovesSlot         uint                 `json:"gloves_slot"`
	FootsSlot          uint                 `json:"foots_slot"`
	BracersSlot        uint                 `json:"bracers_slot"`
	BeltSlot           uint                 `json:"belt_slot"`
	WeaponSlot         uint                 `json:"weapon_slot"`
	ShieldSlot         uint                 `json:"shield_slot"`
	RingSlot           uint                 `json:"ring_slot"`
	NecklaceSlot       uint                 `json:"necklace_slot"`
	Gold               uint                 `json:"gold" sql:"DEFAULT:0"`
	Attack             uint                 `json:"attack" sql:"DEFAULT:1"`
	Defense            uint                 `json:"defense" sql:"DEFAULT:1"`
	Hp                 uint                 `json:"hp" sql:"DEFAULT:20"`
	Level              uint                 `json:"level" sql:"DEFAULT:1"`
	Exp                uint                 `json:"exp" sql:"DEFAULT:0"`
	ExpNext            uint                 `json:"exp_next" sql:"DEFAULT:100"`
	FreeStats          uint                 `json:"free_stats" sql:"DEFAULT:10"`
	LumberjackingSkill uint                 `json:"lumberjacking_skill" sql:"DEFAULT:0"`
	FishingSkill       uint                 `json:"fishing_skill" sql:"DEFAULT:0"`
	LumberjackingSlot  uint                 `json:"lumberjacking_slot"`
	FishingSlot        uint                 `json:"fishing_slot"`
	CloakSlot          uint                 `json:"cloak_slot"`
	PantsSlot          uint                 `json:"pants_slot"`
	CurrentHp          uint                 `json:"current_hp"`
	AvatarID           uint                 `json:"avatar_id"`
	Avatar             *Avatar              `json:"avatar"`
	Messages           []*Message           `json:"messages"`
	Stuffs             []*Stuff             `json:"stuffs"`
	Fights             []*Fight             `json:"fights"`
	Movements          []*Movement          `json:"movements"`
	Events             []*Event             `json:"events"`
	Equipment          []*UserEquipmentItem `json:"user_equipment_items"`
	Tools              []*UserToolItem      `json:"user_tool_items"`
}
