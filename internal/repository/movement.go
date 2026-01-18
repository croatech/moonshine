package repository

import (
	"moonshine/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MovementRepository struct {
	db *sqlx.DB
}

func NewMovementRepository(db *sqlx.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

func (r *MovementRepository) FindActiveByUserID(userID uuid.UUID) (*domain.Movement, error) {
	var movement domain.Movement
	query := `
		SELECT id, user_id, status, created_at, updated_at, deleted_at
		FROM movements
		WHERE user_id = $1 AND status = 'active' AND deleted_at IS NULL
		LIMIT 1
	`
	err := r.db.Get(&movement, query, userID)
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (r *MovementRepository) Create(movement *domain.Movement) error {
	query := `
		INSERT INTO movements (user_id, status)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowx(query, movement.UserID, movement.Status).
		Scan(&movement.ID, &movement.CreatedAt, &movement.UpdatedAt)
}

func (r *MovementRepository) UpdateStatus(movementID uuid.UUID, status domain.MovementStatus) error {
	query := `
		UPDATE movements
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(query, status, movementID)
	return err
}

func (r *MovementRepository) CreateMovementCell(cell *domain.MovementCell) error {
	query := `
		INSERT INTO movements_cells (movement_id, from_cell_id, to_cell_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowx(query, cell.MovementID, cell.FromCellID, cell.ToCellID).
		Scan(&cell.ID, &cell.CreatedAt)
}






