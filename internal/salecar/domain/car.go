package domain

import (
	"context"
	"time"
)

// CarType ...
type CarType struct {
	ID   uint
	Name string
}

// CarMark ...
type CarMark struct {
	ID        uint
	Name      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	CarTypeID uint      `db:"car_type_id"`
	NameRus   string    `db:"name_rus"`
	IsPopular bool
}

// CarModel ...
type CarModel struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CarTypeID uint
	CarMarkID uint
	NameRus   string
	IsPopular bool
}

// CarGeneration ...
type CarGeneration struct {
	ID         uint
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CarTypeID  uint
	CarModelID uint
	BeginYear  time.Time
	EndYear    time.Time
}

// CarSerie ...
type CarSerie struct {
	ID              uint
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CarTypeID       uint
	CarModelID      uint
	CarGenerationID uint
}

// CarModification ...
type CarModification struct {
	ID         uint
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CarTypeID  uint
	CarSerieID uint
	CarModelID uint
}

// CarOption ...
type CarOption struct {
	ID          uint
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CarTypeID   uint
	CarOptionID uint
}

// CarOptionValue ...
type CarOptionValue struct {
	ID             uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CarTypeID      uint
	IsBase         bool
	CarOptionID    uint
	CarEquipmentID uint
}

// CarEquipment ...
type CarEquipment struct {
	ID                uint
	Name              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CarModificationID uint
	PriceMin          float64
	CarTypeID         uint
	Year              time.Time
}

// CarCharacteristic ...
type CarCharacteristic struct {
	ID                  uint
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CarCharacteristicID uint
	CarTypeID           uint
}

// CarCharacteristicValue ...
type CarCharacteristicValue struct {
	ID                  uint
	Value               string
	Unit                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CarCharacteristicID uint
	CarModificationID   uint
	CarTypeID           uint
}

// CarTyper ...
type CarTyper interface {
	Create(ctx context.Context, carType *CarType) error
	GetByID(ctx context.Context, ID uint) (CarType, error)
	GetAll(ctx context.Context) ([]CarType, error)
	Update(ctx context.Context, carType *CarType) error
	Delete(ctx context.Context, ID uint) error
}

// CarMarker ...
type CarMarker interface {
	Create(ctx context.Context, carMark *CarMark) error
	GetByID(ctx context.Context, ID uint) (CarMark, error)
	GetAll(ctx context.Context) ([]CarMark, error)
	Update(ctx context.Context, carMark *CarMark) error
	Delete(ctx context.Context, ID uint) error
}

// CarModeler ...
type CarModeler interface {
	Create(ctx context.Context, carModel *CarModel) error
	GetByID(ctx context.Context, ID uint) (CarModel, error)
}
