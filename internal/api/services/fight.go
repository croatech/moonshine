package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type FightService struct {
	fightRepo *repository.FightRepository
	botRepo   *repository.BotRepository
	userRepo  *repository.UserRepository
	db        *sqlx.DB
}

func NewFightService(db *sqlx.DB) *FightService {
	return &FightService{
		fightRepo: repository.NewFightRepository(db),
		botRepo:   repository.NewBotRepository(db),
		userRepo:  repository.NewUserRepository(db),
		db:        db,
	}
}

type GetCurrentFightResult struct {
	User *domain.User
	Bot  *domain.Bot
}

var ErrNoActiveFight = errors.New("no active fight")
var ErrUserNotFound = errors.New("user not found")
var ErrBotNotFound = errors.New("bot not found")

func (s *FightService) GetCurrentFight(ctx context.Context, userID uuid.UUID) (*GetCurrentFightResult, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	fight, err := s.fightRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, ErrNoActiveFight
	}

	bot, err := s.botRepo.FindByID(fight.BotID)
	if err != nil {
		return nil, ErrBotNotFound
	}

	return &GetCurrentFightResult{
		User: user,
		Bot:  bot,
	}, nil
}