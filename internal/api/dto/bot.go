package dto

import (
	"time"

	"moonshine/internal/domain"
)

type Bot struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Attack    int       `json:"attack"`
	Defense   int       `json:"defense"`
	Hp        int       `json:"hp"`
	Level     int       `json:"level"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
}

func BotFromDomain(bot *domain.Bot) *Bot {
	if bot == nil {
		return nil
	}

	return &Bot{
		ID:        bot.ID.String(),
		Name:      bot.Name,
		Slug:      bot.Slug,
		Attack:    int(bot.Attack),
		Defense:   int(bot.Defense),
		Hp:        int(bot.Hp),
		Level:     int(bot.Level),
		Avatar:    bot.Avatar,
		CreatedAt: bot.CreatedAt,
	}
}

func BotsFromDomain(bots []*domain.Bot) []*Bot {
	result := make([]*Bot, len(bots))
	for i, bot := range bots {
		result[i] = BotFromDomain(bot)
	}
	return result
}
