package convert_currencies

import (
	"context"
	"fmt"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/repository"
)

// usecase
type ConvertCurrenciesUsecase interface {
	CreateConvertCurrencies(ctx context.Context, cry *entity.ConvertCurrencies) error
}

type Provider struct {
	Repo repository.ConversionRepo
}

//Service book usecase
type Service struct {
	*Provider
}

//NewService create new service
func NewService(prvd *Provider) ConvertCurrenciesUsecase {
	return &Service{prvd}
}

func (s *Service) CreateConvertCurrencies(ctx context.Context, ec *entity.ConvertCurrencies) error {
	params := request.ConversionParameter{
		Limit:          10,
		Offset:         0,
		CurrencyIDFrom: ec.CurrencyIDFrom,
		CurrencyIDTo:   ec.CurrencyIDTo,
	}
	conversions, total, err := s.Repo.GetConversions(ctx, &params)
	if err != nil {
		return err
	}

	if total == 0 {
		return fmt.Errorf("Not Found")
	}

	if conversions[0].CurrencyIDFrom == ec.CurrencyIDFrom && conversions[0].CurrencyIDTo == ec.CurrencyIDTo {
		ec.Result = ec.Amount * conversions[0].Rate
	} else {
		ec.Result = ec.Amount / conversions[0].Rate
	}
	return nil
}
