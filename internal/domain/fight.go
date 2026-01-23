package domain

import "github.com/google/uuid"

type FightStatus string

const (
	FightStatusInProgress FightStatus = "IN_PROGRESS"
	FightStatusFinished   FightStatus = "FINISHED"
)

type Fight struct {
	Model
	UserID        uuid.UUID  `db:"user_id"`
	BotID         uuid.UUID  `db:"bot_id"`
	Status        FightStatus `db:"status"`
	DroppedGold   uint       `db:"dropped_gold"`
	DroppedItemID *uuid.UUID `db:"dropped_item_id"`
}
