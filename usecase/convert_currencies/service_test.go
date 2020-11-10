package convert_currencies_test

import (
	"context"
	"testing"
	"time"

	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/mocks"
	"github.com/rbpermadi/whim_assignment/usecase/convert_currencies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type writeConvertCurrenciesData struct {
	name    string
	data    entity.ConvertCurrencies
	IsError bool
}

type mockProvider struct {
	Repo *mocks.ConversionRepo
}

func createService(p *convert_currencies.Provider) convert_currencies.ConvertCurrenciesUsecase {
	return convert_currencies.NewService(p)
}

func provider() mockProvider {
	return mockProvider{
		Repo: new(mocks.ConversionRepo),
	}
}

func getWriteConvertCurrenciesData(data entity.ConvertCurrencies) []writeConvertCurrenciesData {
	return []writeConvertCurrenciesData{
		// TODO: Add test cases.
		// TODO: Add test cases.
		{
			name:    "success",
			IsError: false,
		},
		{
			name:    "not found",
			IsError: true,
		},
	}
}

func sampleConvertCurrencies() entity.ConvertCurrencies {
	return entity.ConvertCurrencies{
		CurrencyIDFrom: 1,
		CurrencyIDTo:   2,
		Amount:         580,
		Result:         20,
	}
}

func sampleConversion() entity.Conversion {
	now := time.Now()
	return entity.Conversion{
		ID:             1,
		CurrencyIDFrom: 1,
		CurrencyIDTo:   2,
		Rate:           29,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func TestCreateConvertCurrencies(t *testing.T) {
	ap := provider()
	resultConvertCurrencies := sampleConvertCurrencies()
	singleConversion := sampleConversion()

	tests := getWriteConvertCurrenciesData(resultConvertCurrencies)

	ap.Repo.On("GetConversions", mock.Anything, mock.Anything).Return([]entity.Conversion{singleConversion}, int64(1), nil).Times(1)
	//emulate not found occurred
	ap.Repo.On("GetConversions", mock.Anything, mock.Anything).Return(nil, int64(0), nil).Times(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&convert_currencies.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			err := u.CreateConvertCurrencies(ctx, &tt.data)

			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}
