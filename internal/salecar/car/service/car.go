package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/omekov/superapp-backend/internal/salecar/delivery/http/dto"
	"github.com/omekov/superapp-backend/internal/salecar/domain"
)

// создание справочника например
// выбирает лекговой автомобиль есть ID
// если нет создается новый ID
func (s *Service) CreateCar(ctx context.Context, car dto.CarRequestBody) error {

	carType, err := s.carrepository.CarTyper.GetByID(ctx, car.CarTypeID)
	if errors.Is(err, sql.ErrNoRows) {
		carType = domain.CarType{
			Name: car.CarTypeName,
		}
		s.carrepository.CarTyper.Create(ctx, &carType)
	} else if err != nil {
		return fmt.Errorf("CarTyper.GetByID %s", err)
	}

	carMark, err := s.carrepository.CarMarker.GetByID(ctx, car.CarMarkID)
	if errors.Is(err, sql.ErrNoRows) {
		carMark = domain.CarMark{
			Name:      car.CarMarkName,
			NameRus:   car.CarMarkNameRus,
			CarTypeID: carType.ID,
		}
		s.carrepository.CarMarker.Create(ctx, &carMark)
	} else if err != nil {
		return fmt.Errorf("CarMarker.GetByID %s", err)
	}

	carModel, err := s.carrepository.CarModeler.GetByID(ctx, car.CarModelID)
	if errors.Is(err, sql.ErrNoRows) {
		carModel = domain.CarModel{
			Name:      car.CarMarkName,
			NameRus:   car.CarMarkNameRus,
			CarTypeID: carType.ID,
			CarMarkID: carMark.ID,
		}
		s.carrepository.CarModeler.Create(ctx, &carModel)
	} else if err != nil {
		return fmt.Errorf("CarModeler.GetByID %s", err)
	}
	return nil
}

func (s *Service) GetCarTypes(ctx context.Context) ([]dto.CarType, error) {
	var carTypeDTO []dto.CarType
	carType, err := s.carrepository.CarTyper.GetAll(ctx)
	if err != nil {
		return carTypeDTO, err
	}
	carTypeDTO = dto.ToDTO(carType).([]dto.CarType)

	return carTypeDTO, nil
}

// {
// 	"carType": [{
// 		"id": 1,
// 		"name": "Легковой",
// 		"carMark": [{
// 			"id":1,
// 			"name": "Aston Martin"
// 		}]
// 	}],
// }
