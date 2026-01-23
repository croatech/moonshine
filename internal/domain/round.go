package domain

import "github.com/google/uuid"

type BodyPart string

const (
	BodyPartHead  BodyPart = "HEAD"
	BodyPartNeck  BodyPart = "NECK"
	BodyPartChest BodyPart = "CHEST"
	BodyPartBelt  BodyPart = "BELT"
	BodyPartLegs  BodyPart = "LEGS"
	BodyPartHands BodyPart = "HANDS"
)

type Round struct {
	Model
	FightID            uuid.UUID   `db:"fight_id"`
	PlayerDamage       uint        `db:"player_damage"`
	BotDamage          uint        `db:"bot_damage"`
	Status             FightStatus `db:"status"`
	PlayerHp           uint        `db:"player_hp"`
	BotHp              uint        `db:"bot_hp"`
	PlayerAttackPoint  *BodyPart   `db:"player_attack_point"`
	PlayerDefensePoint *BodyPart   `db:"player_defense_point"`
	BotAttackPoint     *BodyPart   `db:"bot_attack_point"`
	BotDefensePoint    *BodyPart   `db:"bot_defense_point"`
}
