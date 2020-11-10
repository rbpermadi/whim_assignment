package conversion

import (
	"context"
	"fmt"
	"time"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/repository"
)

// usecase
type ConversionUsecase interface {
	CreateConversion(ctx context.Context, cry *entity.Conversion) error
	UpdateConversion(ctx context.Context, id int64, cry *entity.Conversion) error
	GetConversions(ctx context.Context, p *request.ConversionParameter) ([]entity.Conversion, int64, error)
	GetConversion(ctx context.Context, id int64) (*entity.Conversion, error)
}

type Provider struct {
	Repo         repository.ConversionRepo
	CurrencyRepo repository.CurrencyRepo
}

//Service book usecase
type Service struct {
	*Provider
}

//NewService create new service
func NewService(prvd *Provider) ConversionUsecase {
	return &Service{prvd}
}

func (s *Service) CreateConversion(ctx context.Context, ec *entity.Conversion) error {
	_, err := s.CurrencyRepo.GetCurrency(ctx, ec.CurrencyIDFrom)
	if err != nil {
		return fmt.Errorf("Bad Request")
	}

	_, err = s.CurrencyRepo.GetCurrency(ctx, ec.CurrencyIDTo)
	if err != nil {
		return fmt.Errorf("Bad Request")
	}

	params := request.ConversionParameter{
		Limit:          10,
		Offset:         0,
		CurrencyIDFrom: ec.CurrencyIDFrom,
		CurrencyIDTo:   ec.CurrencyIDTo,
	}
	_, total, err := s.Repo.GetConversions(ctx, &params)
	if err == nil && total > 0 {
		return fmt.Errorf("Duplicate entry")
	}

	ec.CreatedAt = time.Now()
	ec.UpdatedAt = time.Now()

	err = s.Repo.CreateConversion(ctx, ec)
	return err
}

func (s *Service) UpdateConversion(ctx context.Context, id int64, ec *entity.Conversion) error {
	ec.UpdatedAt = time.Now()

	err := s.Repo.UpdateConversion(ctx, id, ec)

	return err
}

func (s *Service) GetConversions(ctx context.Context, p *request.ConversionParameter) ([]entity.Conversion, int64, error) {
	conversions, length, err := s.Repo.GetConversions(ctx, p)

	return conversions, length, err
}

func (s *Service) GetConversion(ctx context.Context, id int64) (*entity.Conversion, error) {
	return s.Repo.GetConversion(ctx, id)
}
