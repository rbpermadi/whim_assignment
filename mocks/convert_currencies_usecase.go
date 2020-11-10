package mocks

import (
	context "context"

	"github.com/rbpermadi/whim_assignment/entity"
	mock "github.com/stretchr/testify/mock"
)

type ConvertCurrenciesUsecase struct {
	mock.Mock
}

func (_m *ConvertCurrenciesUsecase) CreateConvertCurrencies(ctx context.Context, cc *entity.ConvertCurrencies) error {
	ret := _m.Called(ctx, cc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.ConvertCurrencies) error); ok {
		r0 = rf(ctx, cc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
