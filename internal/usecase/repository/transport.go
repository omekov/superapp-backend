package repository

import (
	"context"
)

type Mark struct {
	Name string `db:"mark"`
}

type Model struct {
	Name string `db:"model"`
}
type Volume struct {
	Value int `db:"volume"`
}

type Specification struct {
	Year   int `db:"year"`
	Amount int `db:"amount"`
}

func (r Repo) GetMarks(ctx context.Context) ([]Mark, error) {
	marks := make([]Mark, 0)
	rows, err := r.getMarksStmt.Query()
	if err != nil {
		return marks, err
	}

	for rows.Next() {
		mark := Mark{}
		err := rows.Scan(&mark.Name)
		if err != nil {
			return marks, err
		}

		marks = append(marks, mark)
	}

	return marks, nil
}

func (r Repo) GetModels(ctx context.Context, mark string) ([]Model, error) {
	models := make([]Model, 0)
	rows, err := r.getModelsStmt.Query(mark)
	if err != nil {
		return models, err
	}

	for rows.Next() {
		model := Model{}
		err := rows.Scan(&model.Name)
		if err != nil {
			return models, err
		}

		models = append(models, model)
	}
	return models, nil
}

func (r Repo) GetVolumes(ctx context.Context, mark, model string) ([]Volume, error) {
	volumes := make([]Volume, 0)
	rows, err := r.getVolumesStmt.Query(mark, model)
	if err != nil {
		return volumes, err
	}

	for rows.Next() {
		volume := Volume{}
		err := rows.Scan(&volume.Value)
		if err != nil {
			return volumes, err
		}

		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func (r Repo) GetSpecifications(ctx context.Context, mark, model string, volume int) ([]Specification, error) {
	specifications := make([]Specification, 0)
	rows, err := r.getSpecificationsStmt.Query(mark, model, volume)
	if err != nil {
		return specifications, err
	}

	for rows.Next() {
		specification := Specification{}
		err := rows.Scan(&specification.Year, &specification.Amount)
		if err != nil {
			return specifications, err
		}

		specifications = append(specifications, specification)
	}
	return specifications, nil
}

type Delivered struct {
	FromCity string
	ToCity   string
	Amount   int
}

func (r Repo) GetDelivereds(ctx context.Context, country string) ([]Delivered, error) {
	delivereds := make([]Delivered, 0)
	rows, err := r.db.Query("SELECT from_city, to_city, amount FROM delivered FROM country = ?", country)
	if err != nil {
		return delivereds, err
	}

	for rows.Next() {
		delivered := Delivered{}
		err := rows.Scan(&delivered.FromCity, &delivered.ToCity, &delivered.Amount)
		if err != nil {
			return delivereds, err
		}

		delivereds = append(delivereds, delivered)
	}

	return delivereds, nil
}

func (r Repo) GetBrokerAmounts(ctx context.Context, country string) ([]int, error) {
	amounts := make([]int, 0)
	rows, err := r.db.Query("SELECT amount FROM broker_amount FROM country = ?", country)
	if err != nil {
		return amounts, err
	}

	for rows.Next() {
		var amount int
		err := rows.Scan(&amount)
		if err != nil {
			return amounts, err
		}

		amounts = append(amounts, amount)
	}

	return amounts, nil
}
