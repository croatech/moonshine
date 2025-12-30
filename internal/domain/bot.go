package domain

type Bot struct {
	Model
	Name         string              `json:"name"`
	Attack       uint                `json:"attack"`
	Defense      uint                `json:"defense"`
	Hp           uint                `json:"hp"`
	Level        uint                `json:"level"`
	Avatar       string              `json:"avatar"`
	LocationBots []*LocationBot      `json:"location_bots,omitempty"`
	Fights       []*Fight            `json:"fights,omitempty"`
	Equipment    []*BotEquipmentItem `json:"equipment,omitempty"`
}
