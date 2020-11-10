package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/usecase/currency"
)

type CurrencyHandler struct {
	uc currency.CurrencyUsecase
}

func NewCurrencyHandler(usecase currency.CurrencyUsecase) CurrencyHandler {
	return CurrencyHandler{uc: usecase}
}

func (ch *CurrencyHandler) Register(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Passed router cannot be nil or empty")
	}

	r.GET("/v1/currencies", ch.GetCurrencies)
	r.GET("/v1/currencies/:id", ch.GetCurrency)
	r.POST("/v1/currencies", ch.CreateCurrency)
	r.PATCH("/v1/currencies/:id", ch.UpdateCurrency)

	return nil
}

func (ch *CurrencyHandler) GetCurrencies(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	helper := request.NewQueryHelper(r)

	params := request.CurrencyParameter{
		Limit:  helper.GetInt("limit", 10),
		Offset: helper.GetInt("offset", 0),
		Query:  helper.GetString("query", ""),
	}

	context := r.Context()
	currencies, total, err := ch.uc.GetCurrencies(context, &params)

	if err != nil {
		fmt.Println(err)
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	if len(currencies) <= 0 {
		m := response.MetaInfo{HTTPStatus: http.StatusNoContent}
		response.Write(w, response.BuildSuccess(currencies, m), http.StatusOK)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
		Offset:     params.Offset,
		Limit:      params.Limit,
		Total:      total,
	}
	response.Write(w, response.BuildSuccess(currencies, meta), http.StatusOK)
	return
}

func (ch *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	currencyID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		response.Write(w, err, http.StatusInternalServerError)
		return
	}

	context := r.Context()

	currency, err := ch.uc.GetCurrency(context, currencyID)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
	}
	response.Write(w, response.BuildSuccess(currency, meta), http.StatusOK)
	return
}

func (ch *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var currency entity.Currency
	if err := decoder.Decode(&currency); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	context := r.Context()
	if err := ch.uc.CreateCurrency(context, &currency); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusCreated,
	}
	response.Write(w, response.BuildSuccess(currency, meta), http.StatusCreated)
	return
}

func (ch *CurrencyHandler) UpdateCurrency(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	currencyID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		response.Write(w, err, http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var currency entity.Currency
	if err := decoder.Decode(&currency); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	context := r.Context()
	if err := ch.uc.UpdateCurrency(context, currencyID, &currency); err != nil {
		fmt.Println(err)
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	curr, err := ch.uc.GetCurrency(context, currencyID)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
	}
	response.Write(w, response.BuildSuccess(curr, meta), http.StatusOK)
	return
}
