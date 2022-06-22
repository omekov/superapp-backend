package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/omekov/superapp-backend/internal/salecar/domain"
	"github.com/stretchr/testify/assert"
)

var expectedCarModel = domain.CarModel{
	ID:        1,
	Name:      "Копейка",
	CreatedAt: defaultDate,
	UpdatedAt: defaultDate,
	NameRus:   "Копейка",
	CarMarkID: 1,
}

func TestCarModel_CRUD(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	carrepository := instanceReposity(t)

	carType := domain.CarType{Name: "Легковой"}
	err := carrepository.CarTyper.Create(ctx, &carType)
	if err != nil {
		panic(err)
	}

	carMark := domain.CarMark{
		Name:      "Жигули",
		CarTypeID: carType.ID,
		NameRus:   "Жигули",
	}

	err = carrepository.CarMarker.Create(ctx, &carMark)
	if err != nil {
		panic(err)
	}

	carModel := domain.CarModel{
		Name:      "Копейка",
		CarTypeID: carType.ID,
		NameRus:   "Копейка",
		CarMarkID: carMark.ID,
	}

	t.Run("create", func(t *testing.T) {
		err := carrepository.CarModeler.Create(ctx, &carModel)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, expectedCarModel.ID, carModel.ID)
	})

}
