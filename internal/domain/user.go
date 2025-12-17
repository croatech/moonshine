package domain

import "github.com/google/uuid"

type User struct {
	Model
	ArmorSlot          uint                 `json:"armor_slot"`
	Attack             uint                 `json:"attack" gorm:"default:1"`
	Avatar             *Avatar              `json:"avatar,omitempty"`
	AvatarID           *uuid.UUID           `json:"avatar_id" gorm:"type:uuid"`
	BeltSlot           uint                 `json:"belt_slot"`
	BracersSlot        uint                 `json:"bracers_slot"`
	CloakSlot          uint                 `json:"cloak_slot"`
	CurrentHp          uint                 `json:"current_hp"`
	Defense            uint                 `json:"defense" gorm:"default:1"`
	Email              string               `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Events             []*Event             `json:"events,omitempty"`
	Exp                uint                 `json:"exp" gorm:"default:0"`
	ExpNext            uint                 `json:"exp_next" gorm:"default:100"`
	Fights             []*Fight             `json:"fights,omitempty"`
	FishingSkill       uint                 `json:"fishing_skill" gorm:"default:0"`
	FishingSlot        uint                 `json:"fishing_slot"`
	FootsSlot          uint                 `json:"foots_slot"`
	FreeStats          uint                 `json:"free_stats" gorm:"default:10"`
	GlovesSlot         uint                 `json:"gloves_slot"`
	Gold               uint                 `json:"gold" gorm:"default:0"`
	HelmetSlot         uint                 `json:"helmet_slot"`
	Hp                 uint                 `json:"hp" gorm:"default:20"`
	Level              uint                 `json:"level" gorm:"default:1"`
	Location           *Location            `json:"location,omitempty"`
	LocationID         uuid.UUID            `json:"location_id" gorm:"type:uuid;not null"`
	LumberjackingSkill uint                 `json:"lumberjacking_skill" gorm:"default:0"`
	LumberjackingSlot  uint                 `json:"lumberjacking_slot"`
	MailSlot           uint                 `json:"mail_slot"`
	Messages           []*Message           `json:"messages,omitempty"`
	Movements          []*Movement          `json:"movements,omitempty"`
	Name               string               `json:"name"`
	NecklaceSlot       uint                 `json:"necklace_slot"`
	PantsSlot          uint                 `json:"pants_slot"`
	Password           string               `json:"-" gorm:"type:varchar(255);not null"`
	RingSlot           uint                 `json:"ring_slot"`
	ShieldSlot         uint                 `json:"shield_slot"`
	Stuffs             []*Stuff             `json:"stuffs,omitempty"`
	Tools              []*UserToolItem      `json:"tools,omitempty"`
	Username           string               `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	WeaponSlot         uint                 `json:"weapon_slot"`
	Equipment          []*UserEquipmentItem `json:"equipment,omitempty"`
}
