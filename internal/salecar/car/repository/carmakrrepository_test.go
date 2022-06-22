package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/omekov/superapp-backend/internal/salecar/domain"
	"github.com/stretchr/testify/assert"
)

var defaultDate = time.Date(time.Now().Year(), time.January, time.Now().Day(), 0, 0, 0, 0, time.FixedZone("Almaty", 6))
var expectedCarMark = domain.CarMark{
	ID:        1,
	Name:      "Жигули",
	CreatedAt: defaultDate,
	UpdatedAt: defaultDate,
	NameRus:   "Жигули",
}

func TestCarMark_CRUD(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	carrepository := instanceReposity(t)

	carType := domain.CarType{Name: "Легковой"}
	err := carrepository.CarTyper.Create(ctx, &carType)
	if err != nil {
		panic(err)
	}

	expectedCarMark.CarTypeID = carType.ID

	carMark := domain.CarMark{
		Name:      "Жигули",
		CarTypeID: carType.ID,
		NameRus:   "Жигули",
	}

	t.Run("create", func(t *testing.T) {
		err := carrepository.CarMarker.Create(ctx, &carMark)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, expectedCarMark.ID, carMark.ID)
	})

	t.Run("getByID", func(t *testing.T) {
		var err error
		carMark, err = carrepository.CarMarker.GetByID(ctx, carMark.ID)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, expectedCarMark.Name, carMark.Name)
		assert.Equal(t, expectedCarMark.NameRus, carMark.NameRus)
		assert.Equal(t, expectedCarMark.CarTypeID, carMark.CarTypeID)
	})

	t.Run("get all", func(t *testing.T) {
		carMarks, err := carrepository.CarMarker.GetAll(ctx)
		if err != nil {
			panic(err)
		}
		if len(carMarks) == 0 {
			assert.Fail(t, "carMarks empty")
		}
	})

	expectedCarMark.Name = "Москивич"
	expectedCarMark.NameRus = "Москивич"
	expectedCarMark.UpdatedAt = carMark.UpdatedAt
	t.Run("update", func(t *testing.T) {
		carMark.Name = "Москивич"
		carMark.NameRus = "Москивич"
		err := carrepository.CarMarker.Update(ctx, &carMark)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, expectedCarMark.Name, carMark.Name)
		assert.Equal(t, expectedCarMark.NameRus, carMark.NameRus)
		assert.NotEqual(t, expectedCarMark.UpdatedAt, carMark.UpdatedAt)
	})

	t.Run("delete", func(t *testing.T) {
		err := carrepository.CarMarker.Delete(ctx, carMark.ID)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, expectedCarMark.ID, carMark.ID)
	})
}
