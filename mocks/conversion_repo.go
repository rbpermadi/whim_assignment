package mocks

import (
	context "context"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	mock "github.com/stretchr/testify/mock"
)

type ConversionRepo struct {
	mock.Mock
}

func (_m *ConversionRepo) CreateConversion(ctx context.Context, cry *entity.Conversion) error {
	ret := _m.Called(ctx, cry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Conversion) error); ok {
		r0 = rf(ctx, cry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *ConversionRepo) DeleteConversion(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *ConversionRepo) GetConversions(ctx context.Context, p *request.ConversionParameter) ([]entity.Conversion, int64, error) {
	ret := _m.Called(ctx, p)

	var r0 []entity.Conversion
	if rf, ok := ret.Get(0).(func(context.Context, *request.ConversionParameter) []entity.Conversion); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Conversion)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, *request.ConversionParameter) int64); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *request.ConversionParameter) error); ok {
		r2 = rf(ctx, p)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetConversion provides a mock function with given fields: ctx, id
func (_m *ConversionRepo) GetConversion(ctx context.Context, id int64) (*entity.Conversion, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.Conversion
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Conversion); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Conversion)
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

// UpdateConversion provides a mock function with given fields: ctx, id, td
func (_m *ConversionRepo) UpdateConversion(ctx context.Context, id int64, td *entity.Conversion) error {
	ret := _m.Called(ctx, id, td)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *entity.Conversion) error); ok {
		r0 = rf(ctx, id, td)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
