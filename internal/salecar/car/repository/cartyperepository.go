package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/internal/salecar/domain"
)

type carTypeRepository struct {
	db *sqlx.DB
}

func newCarTypeRepository(db *sqlx.DB) *carTypeRepository {
	return &carTypeRepository{
		db: db,
	}
}

func (r *carTypeRepository) Create(ctx context.Context, carType *domain.CarType) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO car_type (name) VALUES ($1) RETURNING id", carType.Name).Scan(&carType.ID)
}

func (r *carTypeRepository) GetByID(ctx context.Context, ID uint) (domain.CarType, error) {
	var carType domain.CarType
	err := r.db.QueryRowxContext(ctx, "SELECT id, name FROM car_type WHERE id=$1", ID).StructScan(&carType)
	return carType, err
}

func (r *carTypeRepository) GetAll(ctx context.Context) ([]domain.CarType, error) {
	carTypes := make([]domain.CarType, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM car_type")
	if err != nil {
		return carTypes, err
	}
	defer rows.Close()

	for rows.Next() {
		carType := domain.CarType{}
		if err := rows.Scan(
			&carType.ID,
			&carType.Name,
		); err != nil {
			return carTypes, err
		}

		carTypes = append(carTypes, carType)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carTypes, nil
}

func (r *carTypeRepository) Update(ctx context.Context, carType *domain.CarType) error {
	result, err := r.db.ExecContext(ctx, "UPDATE car_type SET name = $2 WHERE id = $1", carType.ID, carType.Name)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *carTypeRepository) Delete(ctx context.Context, ID uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM car_type WHERE id = $1", ID)
	return err
}
