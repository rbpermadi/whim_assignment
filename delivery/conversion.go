package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/usecase/conversion"
)

type ConversionHandler struct {
	uc conversion.ConversionUsecase
}

func NewConversionHandler(usecase conversion.ConversionUsecase) ConversionHandler {
	return ConversionHandler{uc: usecase}
}

func (ch *ConversionHandler) Register(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Passed router cannot be nil or empty")
	}

	r.GET("/v1/conversions", ch.GetConversions)
	r.GET("/v1/conversions/:id", ch.GetConversion)
	r.POST("/v1/conversions", ch.CreateConversion)
	r.PATCH("/v1/conversions/:id", ch.UpdateConversion)

	return nil
}

func (ch *ConversionHandler) GetConversions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	helper := request.NewQueryHelper(r)

	params := request.ConversionParameter{
		Limit:          helper.GetInt("limit", 10),
		Offset:         helper.GetInt("offset", 0),
		CurrencyIDFrom: helper.GetInt64("currency_id_from", 0),
		CurrencyIDTo:   helper.GetInt64("currency_id_to", 0),
	}

	context := r.Context()
	conversions, total, err := ch.uc.GetConversions(context, &params)

	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	if len(conversions) <= 0 {
		m := response.MetaInfo{HTTPStatus: http.StatusNoContent}
		response.Write(w, response.BuildSuccess(conversions, m), http.StatusOK)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
		Offset:     params.Offset,
		Limit:      params.Limit,
		Total:      total,
	}
	response.Write(w, response.BuildSuccess(conversions, meta), http.StatusOK)
	return
}

func (ch *ConversionHandler) GetConversion(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conversionID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	context := r.Context()

	conversion, err := ch.uc.GetConversion(context, conversionID)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusOK,
	}
	response.Write(w, response.BuildSuccess(conversion, meta), http.StatusOK)
	return
}

func (ch *ConversionHandler) CreateConversion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var conversion entity.Conversion
	if err := decoder.Decode(&conversion); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	context := r.Context()
	if err := ch.uc.CreateConversion(context, &conversion); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	meta := response.MetaInfo{
		HTTPStatus: http.StatusCreated,
	}
	response.Write(w, response.BuildSuccess(conversion, meta), http.StatusCreated)
	return
}

func (ch *ConversionHandler) UpdateConversion(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conversionID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var conversion entity.Conversion
	if err := decoder.Decode(&conversion); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	context := r.Context()
	if err := ch.uc.UpdateConversion(context, conversionID, &conversion); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	curr, err := ch.uc.GetConversion(context, conversionID)
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
