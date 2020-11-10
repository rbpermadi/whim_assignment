package mocks

import (
	context "context"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	mock "github.com/stretchr/testify/mock"
)

type CurrencyUsecase struct {
	mock.Mock
}

func (_m *CurrencyUsecase) CreateCurrency(ctx context.Context, cry *entity.Currency) error {
	ret := _m.Called(ctx, cry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Currency) error); ok {
		r0 = rf(ctx, cry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *CurrencyUsecase) GetCurrency(ctx context.Context, id int64) (*entity.Currency, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.Currency
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Currency); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(*entity.Currency)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *CurrencyUsecase) GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error) {
	ret := _m.Called(ctx, p)

	var r0 []entity.Currency
	if rf, ok := ret.Get(0).(func(context.Context, *request.CurrencyParameter) []entity.Currency); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Currency)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, *request.CurrencyParameter) int64); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *request.CurrencyParameter) error); ok {
		r2 = rf(ctx, p)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateCampaign provides a mock function with given fields: ctx, id, cmp
func (_m *CurrencyUsecase) UpdateCurrency(ctx context.Context, id int64, cry *entity.Currency) error {
	ret := _m.Called(ctx, id, cry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *entity.Currency) error); ok {
		r0 = rf(ctx, id, cry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
