package currency_test

import (
	"context"
	"testing"
	"time"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/mocks"
	"github.com/rbpermadi/whim_assignment/usecase/currency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type writeCurrencyData struct {
	name    string
	data    entity.Currency
	IsError bool
}

type readCurrencyData struct {
	name    string
	id      int64
	IsError bool
}

type mockProvider struct {
	Repo *mocks.CurrencyRepo
}

func createService(p *currency.Provider) currency.CurrencyUsecase {
	return currency.NewService(p)
}

func provider() mockProvider {
	return mockProvider{
		Repo: new(mocks.CurrencyRepo),
	}
}

func getWriteCurrencyData(data entity.Currency) []writeCurrencyData {
	return []writeCurrencyData{
		// TODO: Add test cases.
		{
			name:    "success",
			data:    data,
			IsError: false,
		},
	}
}

func getReadCurrencyData(data entity.Currency) []readCurrencyData {
	return []readCurrencyData{
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

func sampleCurrency() entity.Currency {
	now := time.Now()
	return entity.Currency{
		ID:        1,
		Name:      "USD",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestGetCurrencies(t *testing.T) {
	ap := provider()
	singleCurrency := sampleCurrency()
	tests := getReadCurrencyData(singleCurrency)
	tests[1].IsError = false

	ap.Repo.On("GetCurrencies", mock.Anything, mock.Anything).Return([]entity.Currency{singleCurrency}, int64(1), nil).Times(1)
	//emulate not found occurred
	ap.Repo.On("GetCurrencies", mock.Anything, mock.Anything).Return(nil, int64(0), nil).Times(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&currency.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			_, _, err := u.GetCurrencies(ctx, &request.CurrencyParameter{Offset: 10, Limit: 0})
			t.Log(err)
			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}

func TestGetCurrency(t *testing.T) {
	ap := provider()
	resultCurrency := sampleCurrency()
	tests := getReadCurrencyData(resultCurrency)
	ap.Repo.On("GetCurrency", mock.Anything, mock.AnythingOfType("int64")).Return(&resultCurrency, nil).Times(1)
	//emulate not found occurred
	ap.Repo.On("GetCurrency", mock.Anything, mock.AnythingOfType("int64")).Return(nil, response.NotFoundError).Times(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&currency.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			_, err := u.GetCurrency(ctx, tt.id)
			t.Log(err)
			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}

func TestUpdateCurrency(t *testing.T) {
	ap := provider()
	resultCurrency := sampleCurrency()

	tests := getWriteCurrencyData(resultCurrency)

	ap.Repo.On("UpdateCurrency", mock.Anything, resultCurrency.ID, mock.Anything).Return(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createService(&currency.Provider{Repo: ap.Repo})
			ctx := context.TODO()

			err := u.UpdateCurrency(ctx, resultCurrency.ID, &tt.data)

			if !assert.Equal(t, err != nil, tt.IsError) {
				t.Error("Something wrong")
			}
		})
	}
	ap.Repo.AssertExpectations(t)
}
