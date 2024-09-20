package usecase

import (
	"context"
	"time"

	"github.com/omekov/dubaicarkzv2/internal/usecase/repository"
)

type Mark struct {
	Name string
}

type Model struct {
	Name string
}
type Volume struct {
	Value int
}

type Specification struct {
	Year   int
	Amount int
}

type Assessment struct {
	AmountKZT               int
	USD                     int
	Delivereds              []repository.Delivered
	SBKTS                   int
	CustomsCollectionAmount int
	CustomsDutyAmount       int
	VATAmount               int
	ButtonSOSAmount         []int
	BrokerAmouts            []int
	FirstRegistrationAmount int
	UtilAmount              int
}

type Delivered struct {
	FromCity string
	ToCity   string
	Amount   int
}

func (u UseCase) GetMarks(ctx context.Context) ([]Mark, error) {
	marks := make([]Mark, 0)
	marksData, err := u.repo.GetMarks(ctx)
	if err != nil {
		return nil, err
	}

	for _, mark := range marksData {
		marks = append(marks, Mark{
			Name: mark.Name,
		})
	}
	return marks, nil
}

func (u UseCase) GetModels(ctx context.Context, mark string) ([]Model, error) {
	models := make([]Model, 0)
	modelsData, err := u.repo.GetModels(ctx, mark)
	if err != nil {
		return nil, err
	}

	for _, model := range modelsData {
		models = append(models, Model{
			Name: model.Name,
		})
	}
	return models, nil
}
func (u UseCase) GetVolumes(ctx context.Context, mark, model string) ([]Volume, error) {
	volumes := make([]Volume, 0)
	volumesData, err := u.repo.GetVolumes(ctx, mark, model)
	if err != nil {
		return nil, err
	}

	for _, v := range volumesData {
		volumes = append(volumes, Volume{
			Value: v.Value,
		})
	}
	return volumes, nil
}

func (u UseCase) GetSpecifications(ctx context.Context, mark, model string, volume int) ([]Specification, error) {
	specifications := make([]Specification, 0)
	specificationsData, err := u.repo.GetSpecifications(ctx, mark, model, volume)
	if err != nil {
		return nil, err
	}

	for _, s := range specificationsData {
		specifications = append(specifications, Specification{
			Year:   s.Year,
			Amount: s.Amount,
		})
	}
	return specifications, nil
}

func (u UseCase) AssessmentAuto(ctx context.Context, amount, volume, year int) (Assessment, error) {

	currency, err := u.external.GetCurrency(ctx)
	if err != nil {
		return Assessment{}, err
	}

	delivereds, err := u.repo.GetDelivereds(ctx, "kz")
	if err != nil {
		return Assessment{}, err
	}

	brokerAmouts, err := u.repo.GetBrokerAmounts(ctx, "kz")
	if err != nil {
		return Assessment{}, err
	}

	var amountKZT = int(currency.Rates.KZT) * amount
	customsDutyAmount := ((amountKZT) * 15) / 100
	customsCollectionAmount := u.mrp * 6

	return Assessment{
		AmountKZT:               amountKZT,
		USD:                     int(currency.Rates.KZT),
		Delivereds:              delivereds,
		SBKTS:                   0,
		CustomsDutyAmount:       customsDutyAmount,
		CustomsCollectionAmount: customsCollectionAmount,
		BrokerAmouts:            brokerAmouts,
		VATAmount:               ((amountKZT + customsDutyAmount + customsCollectionAmount) * 12) / 100,
		FirstRegistrationAmount: u.calcFirstRegistration(year),
		UtilAmount:              u.calcUtilAmount(float64(volume)),
	}, nil
}

func (u UseCase) calcUtilAmount(volume float64) int {
	var mrpRate float64 = float64(u.mrp) * 50
	if volume <= 1000 {
		totalAmount := mrpRate * 1.5
		return int(totalAmount)
	} else if volume >= 1001 || volume <= 2000 {
		totalAmount := mrpRate * 3.5
		return int(totalAmount)
	} else if volume >= 2001 || volume <= 3000 {
		totalAmount := mrpRate * 5
		return int(totalAmount)
	} else if volume >= 300 {
		totalAmount := mrpRate * 11.5
		return int(totalAmount)
	}
	return 0
}
func (u UseCase) calcFirstRegistration(year int) int {
	currentYear := time.Now().Year()
	age := currentYear - year
	if age >= 1 {
		totalAmount := float64(u.mrp) * 0.25
		return int(totalAmount)
	} else if age >= 2 && age <= 3 {
		totalAmount := float64(u.mrp) * 50
		return int(totalAmount)
	} else if age > 3 {
		totalAmount := float64(u.mrp) * 500
		return int(totalAmount)
	}
	return 0
}
