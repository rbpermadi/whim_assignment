package mocks

import (
	context "context"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"

	mock "github.com/stretchr/testify/mock"
)

type CurrencyRepo struct {
	mock.Mock
}

func (_m *CurrencyRepo) CreateCurrency(ctx context.Context, cry *entity.Currency) error {
	ret := _m.Called(ctx, cry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Currency) error); ok {
		r0 = rf(ctx, cry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *CurrencyRepo) DeleteCurrency(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *CurrencyRepo) GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error) {
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

// GetCurrency provides a mock function with given fields: ctx, id
func (_m *CurrencyRepo) GetCurrency(ctx context.Context, id int64) (*entity.Currency, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.Currency
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Currency); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Currency)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCurrency provides a mock function with given fields: ctx, id, td
func (_m *CurrencyRepo) UpdateCurrency(ctx context.Context, id int64, td *entity.Currency) error {
	ret := _m.Called(ctx, id, td)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *entity.Currency) error); ok {
		r0 = rf(ctx, id, td)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
