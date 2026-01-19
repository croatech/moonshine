package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type BotService struct {
	locationRepo *repository.LocationRepository
	botRepo      *repository.BotRepository
	userRepo     *repository.UserRepository
	fightRepo    *repository.FightRepository
	roundRepo    *repository.RoundRepository
	db           *sqlx.DB
}

func NewBotService(db *sqlx.DB) *BotService {
	return &BotService{
		locationRepo: repository.NewLocationRepository(db),
		botRepo:      repository.NewBotRepository(db),
		userRepo:     repository.NewUserRepository(db),
		fightRepo:    repository.NewFightRepository(db),
		roundRepo:    repository.NewRoundRepository(db),
		db:           db,
	}
}

func (s *BotService) GetBotsByLocationSlug(locationSlug string) ([]*domain.Bot, error) {
	if locationSlug == "" {
		return nil, errors.New("location slug is required")
	}

	location, err := s.locationRepo.FindBySlug(locationSlug)
	if err != nil {
		return nil, err
	}

	bots, err := s.botRepo.FindBotsByLocationID(location.ID)
	if err != nil {
		return nil, err
	}

	return bots, nil
}

type AttackResult struct {
	User *domain.User
	Bot  *domain.Bot
}

func (s *BotService) Attack(ctx context.Context, botSlug string, userID uuid.UUID) (*AttackResult, error) {
	if botSlug == "" {
		return nil, errors.New("bot slug is required")
	}

	bot, err := s.botRepo.FindBySlug(botSlug)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if exists, err := s.locationRepo.HasBot(user.LocationID, bot.ID); err != nil || !exists {
		return nil, errors.New("bot is not in the same location as user")
	}

	inFight, err := s.userRepo.InFight(userID)
	if err != nil {
		return nil, err
	}
	if inFight {
		fight, err := s.fightRepo.FindActiveByUserID(userID)
		if err != nil {
			return nil, err
		}

		currentBot, err := s.botRepo.FindByID(fight.BotID)
		if err != nil {
			return nil, err
		}

		return &AttackResult{
			User: user,
			Bot:  currentBot,
		}, nil
	}

	fightID, err := s.fightRepo.Create(&domain.Fight{
		UserID: user.ID,
		BotID:  bot.ID,
	})
	if err != nil {
		return nil, err
	}

	err = s.roundRepo.Create(fightID)
	if err != nil {
		return nil, err
	}

	return &AttackResult{
		User: user,
		Bot:  bot,
	}, nil
}
