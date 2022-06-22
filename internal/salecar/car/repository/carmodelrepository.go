package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/internal/salecar/domain"
)

type carModelRepository struct {
	db *sqlx.DB
}

func newCarModelRepository(db *sqlx.DB) *carModelRepository {
	return &carModelRepository{
		db: db,
	}
}

func (r *carModelRepository) Create(ctx context.Context, carModel *domain.CarModel) error {
	return r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO car_model (
			name,
			name_rus,
			car_type_id,
			car_mark_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id
		`,
		carModel.Name,
		carModel.NameRus,
		carModel.CarTypeID,
		carModel.CarMarkID,
	).Scan(&carModel.ID)
}

func (r *carModelRepository) GetByID(ctx context.Context, ID uint) (domain.CarModel, error) {
	var carModel domain.CarModel
	err := r.db.QueryRowxContext(
		ctx,
		``,
		ID,
	).Scan(
		&carModel.ID,
	)
	return carModel, err
}
