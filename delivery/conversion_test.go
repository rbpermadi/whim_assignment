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

type requestConversionTestCase struct {
	name           string
	method         string
	endpoint       string
	payload        []byte
	expectedStatus int
}

//NewConversionHTTPRequest create new http request for test
func NewConversionHTTPRequest(method, path, token string, body []byte) *http.Request {
	request := httptest.NewRequest(method, "http://localhost"+path, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func newConversionHandler() (http.Handler, *mocks.ConversionUsecase) {
	uc := new(mocks.ConversionUsecase)
	ConversionHandler := delivery.NewConversionHandler(uc)

	h := handler.NewHandler(&ConversionHandler)
	return h, uc
}

func buildStubConversions() []entity.Conversion {
	stubbedConversions := make([]entity.Conversion, 0)
	for i := 0; i < 10; i++ {
		cry := entity.Conversion{}
		if err := faker.FakeData(&cry); err != nil {
			panic(err)
		}
		stubbedConversions = append(stubbedConversions, cry)
	}
	return stubbedConversions
}

func buildConversionPayload(cry entity.Conversion) ([]byte, error) {
	cry.ID = 0
	return json.Marshal(&cry)
}

func TestConversionRequest(t *testing.T) {
	handler, uc := newConversionHandler()
	stubConversions := buildStubConversions()
	singleConversion := stubConversions[0]
	examplePayload, err := buildConversionPayload(singleConversion)
	assert.NoError(t, err)

	uc.On("CreateConversion", mock.Anything, mock.Anything).Return(nil)
	uc.On("GetConversions", mock.Anything, mock.Anything).Return(stubConversions, int64(len(stubConversions)), nil)
	uc.On("GetConversion", mock.Anything, mock.AnythingOfType("int64")).Return(&singleConversion, nil)
	uc.On("UpdateConversion", mock.Anything, mock.AnythingOfType("int64"), mock.Anything).Return(nil)

	testCases := []requestConversionTestCase{
		{
			name:           "Create new conversion",
			method:         "POST",
			endpoint:       "/v1/conversions",
			payload:        examplePayload,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Get all conversions",
			method:         "GET",
			endpoint:       "/v1/conversions?limit=20&offset=0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get conversion",
			method:         "GET",
			endpoint:       fmt.Sprintf("/v1/conversions/%v", singleConversion.ID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update existing conversion",
			method:         "PATCH",
			endpoint:       fmt.Sprintf("/v1/conversions/%v", singleConversion.ID),
			payload:        examplePayload,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stubRequest := NewConversionHTTPRequest(testCase.method, testCase.endpoint, "", testCase.payload)
			ctx := context.TODO()
			stubRequest = stubRequest.WithContext(ctx)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, stubRequest)
			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}

	uc.AssertExpectations(t)
}
