package conversion_test

import (
	"context"
	"testing"
	"time"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/mocks"
	"github.com/rbpermadi/whim_assignment/usecase/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type writeConversionData struct {
	name    string
	data    entity.Conversion
	IsError bool
}

type readConversionData struct {
	name    string
	id      int64
	IsError bool
}

type mockProvider struct {
	Repo         *mocks.ConversionRepo
	CurrencyRepo *mocks.CurrencyRepo
}

func createService(p *conversion.Provider) conversion.ConversionUsecase {
	return conversion.NewService(p)
}

func provider() mockProvider {
	return mockProvider{
		Repo:         new(mocks.ConversionRepo),
		CurrencyRepo: new(mocks.CurrencyRepo),
	}
}

func getWriteConversionData(data entity.Conversion) []writeConversionData {
	return []writeConversionData{
		// TODO: Add test cases.
		{
			name:    "success",
			data:    data,
			IsError: false,
		},
	}
}

func getReadConversionData(data entity.Conversion) []readConversionData {
	return []readConversionData{
		// TODO: Add test cases.
		{
			name:    "success",
			id:      1,
			IsError: false,
		},
		{
			name:    "not found",
			id:      1,
			IsError: true,
		},
	}
}

func sampleConversion() entity.Conversion {
	now := time.Now()
	return entity.Conversion{
		ID:             1,
		CurrencyIDFrom: 1,
		CurrencyIDTo:   2,
		Rate:           15000,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func TestGetConversions(t *testing.T) {
	ap := provider()
	singleConversion := sampleConversion()
	tests := getReadConversionData(singleConversion)
	tests[1].IsError = false

	ap.Repo.On("GetConversions", mock.Anything, mock.Anything).Return([]entity.Conversion{singleConversion}, int64(1), nil).Times(1)
	//emulate not found occurred
	ap.Repo.On("GetConversions", mock.Anything, mock.Anything).Return(nil, int64(0), nil).Times(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&conversion.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			_, _, err := u.GetConversions(ctx, &request.ConversionParameter{Offset: 10, Limit: 0})
			t.Log(err)
			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}

func TestGetConversion(t *testing.T) {
	ap := provider()
	resultConversion := sampleConversion()
	tests := getReadConversionData(resultConversion)
	ap.Repo.On("GetConversion", mock.Anything, mock.AnythingOfType("int64")).Return(&resultConversion, nil).Times(1)
	//emulate not found occurred
	ap.Repo.On("GetConversion", mock.Anything, mock.AnythingOfType("int64")).Return(nil, response.NotFoundError).Times(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&conversion.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			_, err := u.GetConversion(ctx, tt.id)
			t.Log(err)
			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}

func TestUpdateConversion(t *testing.T) {
	ap := provider()
	resultConversion := sampleConversion()

	tests := getWriteConversionData(resultConversion)

	ap.Repo.On("UpdateConversion", mock.Anything, resultConversion.ID, mock.Anything).Return(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&conversion.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			err := u.UpdateConversion(ctx, resultConversion.ID, &tt.data)

			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}
