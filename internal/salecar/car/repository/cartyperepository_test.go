package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/omekov/superapp-backend/internal/salecar/domain"
	"github.com/stretchr/testify/assert"
)

var expectedCarType = domain.CarType{
	ID:   1,
	Name: "Легковой",
}

func TestCarType_CRUD(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	carrepository := instanceReposity(t)

	carType := domain.CarType{Name: "Легковой"}
	t.Run("create", func(t *testing.T) {
		err := carrepository.CarTyper.Create(ctx, &carType)
		if err != nil {
			panic(err)
		}
		// assert.Equal(t, expectedCarType.ID, carType.ID)
	})

	t.Run("get by id", func(t *testing.T) {
		var err error
		carType, err = carrepository.CarTyper.GetByID(ctx, carType.ID)

		if err != nil {
			panic(err)
		}
		// assert.Equal(t, expectedCarType.ID, carType.ID)
	})

	carType.Name = "Легкавая"
	t.Run("update", func(t *testing.T) {
		err := carrepository.CarTyper.Update(ctx, &carType)
		if err == sql.ErrNoRows {
			return
		}
		if err != nil {
			panic(err)
		}
		assert.Equal(t, "Легкавая", carType.Name)
	})

	t.Run("get all", func(t *testing.T) {
		carTypes, err := carrepository.CarTyper.GetAll(ctx)
		if err != nil {
			panic(err)
		}
		if len(carTypes) == 0 {
			assert.Fail(t, "cartypes empty")
		}
	})

	t.Run("delete", func(t *testing.T) {
		err := carrepository.CarTyper.Delete(ctx, carType.ID)
		if err != nil {
			panic(err)
		}
	})

}
