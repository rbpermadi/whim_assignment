package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/rbpermadi/whim_assignment/delivery"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/handler"
	"github.com/rbpermadi/whim_assignment/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type requestCurrencyTestCases struct {
	name           string
	method         string
	endpoint       string
	payload        []byte
	expectedStatus int
}

//NewCurrencyHTTPRequest create new http request for test
func NewCurrencyHTTPRequest(method, path, token string, body []byte) *http.Request {
	request := httptest.NewRequest(method, "http://localhost"+path, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func newCurrencyHandler() (http.Handler, *mocks.CurrencyUsecase) {
	uc := new(mocks.CurrencyUsecase)
	CurrencyHandler := delivery.NewCurrencyHandler(uc)

	h := handler.NewHandler(&CurrencyHandler)
	return h, uc
}

func buildStubCurrencies() []entity.Currency {
	stubbedCurrencies := make([]entity.Currency, 0)
	for i := 0; i < 10; i++ {
		cry := entity.Currency{}
		if err := faker.FakeData(&cry); err != nil {
			panic(err)
		}
		stubbedCurrencies = append(stubbedCurrencies, cry)
	}
	return stubbedCurrencies
}

func buildCurrencyPayload(cry entity.Currency) ([]byte, error) {
	cry.ID = 0
	return json.Marshal(&cry)
}

func TestCurrencyRequest(t *testing.T) {
	handler, uc := newCurrencyHandler()
	stubCurrencies := buildStubCurrencies()
	singleCurrency := stubCurrencies[0]
	examplePayload, err := buildCurrencyPayload(singleCurrency)
	assert.NoError(t, err)

	uc.On("CreateCurrency", mock.Anything, mock.Anything).Return(nil)
	uc.On("GetCurrencies", mock.Anything, mock.Anything).Return(stubCurrencies, int64(len(stubCurrencies)), nil)
	uc.On("GetCurrency", mock.Anything, mock.AnythingOfType("int64")).Return(&singleCurrency, nil)
	uc.On("UpdateCurrency", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(nil)

	testCases := []requestCurrencyTestCases{
		{
			name:           "Create new currency",
			method:         "POST",
			endpoint:       "/v1/currencies",
			payload:        examplePayload,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Get all currencies",
			method:         "GET",
			endpoint:       "/v1/currencies?limit=20&offset=0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get currency",
			method:         "GET",
			endpoint:       fmt.Sprintf("/v1/currencies/%v", singleCurrency.ID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update existing currency",
			method:         "PATCH",
			endpoint:       fmt.Sprintf("/v1/currencies/%v", singleCurrency.ID),
			payload:        examplePayload,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stubRequest := NewCurrencyHTTPRequest(testCase.method, testCase.endpoint, "", testCase.payload)
			ctx := context.TODO()
			stubRequest = stubRequest.WithContext(ctx)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, stubRequest)
			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}

	uc.AssertExpectations(t)
}
