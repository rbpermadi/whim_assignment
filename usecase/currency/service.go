package currency

import (
	"context"
	"time"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/repository"
)

// usecase
type CurrencyUsecase interface {
	CreateCurrency(ctx context.Context, cry *entity.Currency) error
	UpdateCurrency(ctx context.Context, id int64, cry *entity.Currency) error
	GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error)
	GetCurrency(ctx context.Context, id int64) (*entity.Currency, error)
}

type Provider struct {
	Repo repository.CurrencyRepo
}

//Service book usecase
type Service struct {
	*Provider
}

//NewService create new service
func NewService(prvd *Provider) CurrencyUsecase {
	return &Service{prvd}
}

func (s *Service) CreateCurrency(ctx context.Context, ec *entity.Currency) error {
	ec.CreatedAt = time.Now()
	ec.UpdatedAt = time.Now()

	err := s.Repo.CreateCurrency(ctx, ec)
	return err
}

func (s *Service) UpdateCurrency(ctx context.Context, id int64, ec *entity.Currency) error {
	ec.UpdatedAt = time.Now()

	err := s.Repo.UpdateCurrency(ctx, id, ec)
	return err
}

func (s *Service) GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error) {
	currencies, length, err := s.Repo.GetCurrencies(ctx, p)

	return currencies, length, err
}

func (s *Service) GetCurrency(ctx context.Context, id int64) (*entity.Currency, error) {
	return s.Repo.GetCurrency(ctx, id)
}
