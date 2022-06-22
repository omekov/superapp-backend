package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/internal/salecar/domain"
)

type carMarkRepository struct {
	db *sqlx.DB
}

func newCarMarkRepository(db *sqlx.DB) *carMarkRepository {
	return &carMarkRepository{
		db: db,
	}
}

func (r *carMarkRepository) Create(ctx context.Context, carMark *domain.CarMark) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO car_mark (
			name,
			name_rus,
			car_type_id
		)
		VALUES ($1,$2,$3)
		RETURNING id
	`, carMark.Name, carMark.NameRus, carMark.CarTypeID).Scan(&carMark.ID)
}

func (r *carMarkRepository) GetByID(ctx context.Context, ID uint) (domain.CarMark, error) {
	var carMark domain.CarMark
	err := r.db.QueryRowxContext(ctx, `SELECT 
			id, 
			name, 
			name_rus, 
			created_at,
			updated_at,
			car_type_id 
		FROM car_mark WHERE id=$1`, ID).StructScan(
		&carMark,
	)
	return carMark, err
}

func (r *carMarkRepository) GetAll(ctx context.Context) ([]domain.CarMark, error) {
	carMarks := make([]domain.CarMark, 0)
	rows, err := r.db.QueryxContext(ctx, `SELECT 
			id, 
			name, 
			name_rus, 
			created_at, 
			updated_at, 
			car_type_id 
		FROM car_mark`)
	if err != nil {
		return carMarks, err
	}
	defer rows.Close()

	for rows.Next() {
		carMark := domain.CarMark{}
		if err := rows.StructScan(&carMark); err != nil {
			return carMarks, err
		}

		carMarks = append(carMarks, carMark)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carMarks, nil
}

func (r *carMarkRepository) Update(ctx context.Context, carMark *domain.CarMark) error {
	carMark.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx,
		`UPDATE car_mark SET 
			name = $2,
			name_rus = $3, 
			updated_at = $4, 
			car_type_id = $5
		WHERE id = $1`,
		carMark.ID,
		carMark.Name,
		carMark.NameRus,
		carMark.UpdatedAt,
		carMark.CarTypeID,
	)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *carMarkRepository) Delete(ctx context.Context, ID uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM car_mark WHERE id = $1", ID)
	return err
}
