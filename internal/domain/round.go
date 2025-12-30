package domain

import "github.com/google/uuid"

type Round struct {
	Model
	FightID            uuid.UUID `json:"fight_id"`
	Fight              *Fight `json:"fight,omitempty"`
	PlayerDamage       uint   `json:"player_damage"`
	BotDamage          uint   `json:"bot_damage"`
	Status             uint   `json:"status"`
	PlayerHp           uint   `json:"player_hp"`
	BotHp              uint   `json:"bot_hp"`
	PlayerAttackPoint  string `json:"player_attack_point"`
	PlayerDefensePoint string `json:"player_defense_point"`
	BotAttackPoint     string `json:"bot_attack_point"`
	BotDefensePoint    string `json:"bot_defense_point"`
}
