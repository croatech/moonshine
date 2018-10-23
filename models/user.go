package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username           string `gorm:"type:varchar(100);unique;not null"`
	Email              string `gorm:"type:varchar(100);unique;not null"`
	Password           string `gorm:"type:varchar(255);not null"`
	LocationID         int    `gorm:"not null"`
	HelmetSlot         int
	ArmorSlot          int
	MailSlot           int
	GlovesSlot         int
	FootsSlot          int
	BracersSlot        int
	BeltSlot           int
	WeaponSlot         int
	ShieldSlot         int
	CloakSlot          int
	PantsSlot          int
	RingSlot           int
	NecklaceSlot       int
	Gold               int `gorm:"not null; default: 1500"`
	Attack             int `gorm:"not null; default: 1"`
	Defense            int `gorm:"not null; default: 1"`
	Hp                 int `gorm:"not null; default: 20"`
	Level              int `gorm:"not null; default: 1"`
	Exp                int `gorm:"not null; default: 0"`
	ExpNext            int `gorm:"not null; default: 100"`
	FreeStats          int `gorm:"not null; default: 10"`
	LumberjackingSkill int `gorm:"not null; default: 0"`
	FishingSkill       int `gorm:"not null; default: 0"`
	CurrentHp          int `gorm:"not null; default: 20"`
	AvatarID           int
	Equipment          []int `json:"-" sql:"-"`
	Tools              []int `json:"-" sql:"-"`
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
