package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/usecase/convert_currencies"
)

type ConvertCurrenciesHandler struct {
	uc convert_currencies.ConvertCurrenciesUsecase
}

func NewConvertCurrenciesHandler(usecase convert_currencies.ConvertCurrenciesUsecase) ConvertCurrenciesHandler {
	return ConvertCurrenciesHandler{uc: usecase}
}

func (ch *ConvertCurrenciesHandler) Register(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Passed router cannot be nil or empty")
	}

	r.POST("/v1/convert-currencies", ch.CreateConvertCurrencies)

	return nil
}

func (ch *ConvertCurrenciesHandler) CreateConvertCurrencies(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var convert entity.ConvertCurrencies
	if err := decoder.Decode(&convert); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	context := r.Context()
	if err := ch.uc.CreateConvertCurrencies(context, &convert); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
	}
	response.Write(w, response.BuildSuccess(convert, meta), http.StatusOK)
	return
}
