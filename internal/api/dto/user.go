package dto

import (
	"time"

	"moonshine/internal/domain"
)

type User struct {
	ID                    string    `json:"id"`
	Username              string    `json:"username"`
	Email                 string    `json:"email"`
	Hp                    int       `json:"hp"`
	CurrentHp             int       `json:"currentHp"`
	Attack                int       `json:"attack"`
	Defense               int       `json:"defense"`
	Level                 int       `json:"level"`
	Gold                  int       `json:"gold"`
	Exp                   int       `json:"exp"`
	FreeStats             int       `json:"freeStats"`
	CreatedAt             time.Time `json:"createdAt"`
	Avatar                *Avatar   `json:"avatar,omitempty"`
	ChestEquipmentItemID  *string   `json:"chestEquipmentItemId,omitempty"`
	BeltEquipmentItemID   *string   `json:"beltEquipmentItemId,omitempty"`
	HeadEquipmentItemID   *string   `json:"headEquipmentItemId,omitempty"`
	NeckEquipmentItemID   *string   `json:"neckEquipmentItemId,omitempty"`
	WeaponEquipmentItemID *string   `json:"weaponEquipmentItemId,omitempty"`
	ShieldEquipmentItemID *string   `json:"shieldEquipmentItemId,omitempty"`
	LegsEquipmentItemID   *string   `json:"legsEquipmentItemId,omitempty"`
	FeetEquipmentItemID   *string   `json:"feetEquipmentItemId,omitempty"`
	ArmsEquipmentItemID   *string   `json:"armsEquipmentItemId,omitempty"`
	HandsEquipmentItemID  *string   `json:"handsEquipmentItemId,omitempty"`
	Ring1EquipmentItemID  *string   `json:"ring1EquipmentItemId,omitempty"`
	Ring2EquipmentItemID  *string   `json:"ring2EquipmentItemId,omitempty"`
	Ring3EquipmentItemID  *string   `json:"ring3EquipmentItemId,omitempty"`
	Ring4EquipmentItemID  *string   `json:"ring4EquipmentItemId,omitempty"`
	LocationSlug          *string   `json:"locationSlug,omitempty"`
	InFight               *bool     `json:"inFight,omitempty"`
}

type Avatar struct {
	ID      string `json:"id"`
	Image   string `json:"image"`
	Private bool   `json:"private"`
}

func UserFromDomain(user *domain.User, avatar *domain.Avatar, location *domain.Location, inFight *bool) *User {
	if user == nil {
		return nil
	}

	result := &User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Hp:        int(user.Hp),
		CurrentHp: int(user.CurrentHp),
		Attack:    int(user.Attack),
		Defense:   int(user.Defense),
		Level:     int(user.Level),
		Gold:      int(user.Gold),
		Exp:       int(user.Exp),
		FreeStats: int(user.FreeStats),
		CreatedAt: user.CreatedAt,
		InFight:   inFight,
	}

	if avatar != nil {
		result.Avatar = &Avatar{
			ID:      avatar.ID.String(),
			Image:   avatar.Image,
			Private: avatar.Private,
		}
	}

	if user.ChestEquipmentItemID != nil {
		id := user.ChestEquipmentItemID.String()
		result.ChestEquipmentItemID = &id
	}
	if user.BeltEquipmentItemID != nil {
		id := user.BeltEquipmentItemID.String()
		result.BeltEquipmentItemID = &id
	}
	if user.HeadEquipmentItemID != nil {
		id := user.HeadEquipmentItemID.String()
		result.HeadEquipmentItemID = &id
	}
	if user.NeckEquipmentItemID != nil {
		id := user.NeckEquipmentItemID.String()
		result.NeckEquipmentItemID = &id
	}
	if user.WeaponEquipmentItemID != nil {
		id := user.WeaponEquipmentItemID.String()
		result.WeaponEquipmentItemID = &id
	}
	if user.ShieldEquipmentItemID != nil {
		id := user.ShieldEquipmentItemID.String()
		result.ShieldEquipmentItemID = &id
	}
	if user.LegsEquipmentItemID != nil {
		id := user.LegsEquipmentItemID.String()
		result.LegsEquipmentItemID = &id
	}
	if user.FeetEquipmentItemID != nil {
		id := user.FeetEquipmentItemID.String()
		result.FeetEquipmentItemID = &id
	}
	if user.ArmsEquipmentItemID != nil {
		id := user.ArmsEquipmentItemID.String()
		result.ArmsEquipmentItemID = &id
	}
	if user.HandsEquipmentItemID != nil {
		id := user.HandsEquipmentItemID.String()
		result.HandsEquipmentItemID = &id
	}
	if user.Ring1EquipmentItemID != nil {
		id := user.Ring1EquipmentItemID.String()
		result.Ring1EquipmentItemID = &id
	}
	if user.Ring2EquipmentItemID != nil {
		id := user.Ring2EquipmentItemID.String()
		result.Ring2EquipmentItemID = &id
	}
	if user.Ring3EquipmentItemID != nil {
		id := user.Ring3EquipmentItemID.String()
		result.Ring3EquipmentItemID = &id
	}
	if user.Ring4EquipmentItemID != nil {
		id := user.Ring4EquipmentItemID.String()
		result.Ring4EquipmentItemID = &id
	}

	if location != nil && location.Slug != "" {
		result.LocationSlug = &location.Slug
	}

	return result
}

type UpdateUserRequest struct {
	AvatarID *string `json:"avatarId,omitempty"`
}

func AvatarFromDomain(avatar *domain.Avatar) *Avatar {
	if avatar == nil {
		return nil
	}
	return &Avatar{
		ID:      avatar.ID.String(),
		Image:   avatar.Image,
		Private: avatar.Private,
	}
}

func AvatarsFromDomain(avatars []*domain.Avatar) []*Avatar {
	result := make([]*Avatar, len(avatars))
	for i, avatar := range avatars {
		result[i] = AvatarFromDomain(avatar)
	}
	return result
}
