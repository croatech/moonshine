package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ArmorSlot          uint
	Attack             uint                 `json:"attack" sql:"DEFAULT:1"`
	Avatar             *Avatar              `json:"avatar"`
	AvatarID           uint                 `json:"avatar_id"`
	BeltSlot           uint                 `json:"belt_slot"`
	BracersSlot        uint                 `json:"bracers_slot"`
	CloakSlot          uint                 `json:"cloak_slot"`
	CurrentHp          uint                 `json:"current_hp"`
	Defense            uint                 `json:"defense" sql:"DEFAULT:1"`
	Email              string               `gorm:"type:varchar(100);unique;not null"`
	Events             []*Event             `json:"events"`
	Exp                uint                 `json:"exp" sql:"DEFAULT:0"`
	ExpNext            uint                 `json:"exp_next" sql:"DEFAULT:100"`
	Fights             []*Fight             `json:"fights"`
	FishingSkill       uint                 `json:"fishing_skill" sql:"DEFAULT:0"`
	FishingSlot        uint                 `json:"fishing_slot"`
	FootsSlot          uint                 `json:"foots_slot"`
	FreeStats          uint                 `json:"free_stats" sql:"DEFAULT:10"`
	GlovesSlot         uint                 `json:"gloves_slot"`
	Gold               uint                 `json:"gold" sql:"DEFAULT:0"`
	HelmetSlot         uint                 `json:"helmet_slot"`
	Hp                 uint                 `json:"hp" sql:"DEFAULT:20"`
	Level              uint                 `json:"level" sql:"DEFAULT:1"`
	Location           *Location            `json:"location"`
	LocationID         uint                 `json:"location_id"`
	LumberjackingSkill uint                 `json:"lumberjacking_skill" sql:"DEFAULT:0"`
	LumberjackingSlot  uint                 `json:"lumberjacking_slot"`
	MailSlot           uint                 `json:"mail_slot"`
	Messages           []*Message           `json:"messages"`
	Movements          []*Movement          `json:"movements"`
	Name               string               `json:"name"`
	NecklaceSlot       uint                 `json:"necklace_slot"`
	PantsSlot          uint                 `json:"pants_slot"`
	Password           string               `gorm:"type:varchar(255);not null"`
	RingSlot           uint                 `json:"ring_slot"`
	ShieldSlot         uint                 `json:"shield_slot"`
	Stuffs             []*Stuff             `json:"stuffs"`
	Tools              []*UserToolItem      `json:"user_tool_items"`
	Username           string               `json:"username" gorm:"type:varchar(100);unique;not null"`
	WeaponSlot         uint                 `json:"weapon_slot"`
	Equipment          []*UserEquipmentItem `json:"user_equipment_items"`
}
