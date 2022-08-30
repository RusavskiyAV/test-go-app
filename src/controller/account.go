package account

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	account "test3/src/model/account/entity"
	accountAddHandler "test3/src/model/account/useCase/add"
	accountCreateHandler "test3/src/model/account/useCase/create"
	accountSubtractHandler "test3/src/model/account/useCase/subtract"
	accountTransferHandler "test3/src/model/account/useCase/transfer"
	accountReadModel "test3/src/readModel/account"
	currencyConvertor "test3/src/service"
)

type Controller struct {
	CreateHandler     accountCreateHandler.Interface
	AddHandler        accountAddHandler.Interface
	SubtractHandler   accountSubtractHandler.Interface
	TransferHandler   accountTransferHandler.Interface
	Fetcher           accountReadModel.Interface
	Validate          *validator.Validate
	CurrencyConvertor currencyConvertor.Interface
}

func (c Controller) Add(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		jsonBadRequest(res, map[string]string{"id": "недопустимое значение"})

		return
	}

	command := accountAddHandler.Command{}
	decoder := json.NewDecoder(req.Body)
	command.UserId = uint(id)

	if err = decoder.Decode(&command); err != nil {
		jsonBadRequest(res, err)

		return
	}

	if err = c.Validate.Struct(command); err != nil {
		jsonValidationError(res, err)

		return
	}

	err = c.AddHandler.Handle(command, c.CreateHandler)

	if err != nil {
		jsonBadRequest(res, err)

		return
	}

	jsonOk(res)
}

func (c Controller) Subtract(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		jsonBadRequest(res, map[string]string{"id": "недопустимое значение"})

		return
	}

	command := accountSubtractHandler.Command{}
	decoder := json.NewDecoder(req.Body)
	command.UserId = uint(id)

	if err = decoder.Decode(&command); err != nil {
		jsonBadRequest(res, err)

		return
	}

	if err = c.Validate.Struct(command); err != nil {
		jsonValidationError(res, err)

		return
	}

	err = c.SubtractHandler.Handle(command)

	if err != nil {
		jsonBadRequest(res, err)

		return
	}

	jsonOk(res)
}

func (c Controller) Transfer(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		jsonBadRequest(res, map[string]string{"id": "недопустимое значение"})

		return
	}

	command := accountTransferHandler.Command{}
	decoder := json.NewDecoder(req.Body)
	command.UserId = uint(id)

	if err = decoder.Decode(&command); err != nil {
		jsonBadRequest(res, err)

		return
	}

	if err = c.Validate.Struct(command); err != nil {
		jsonValidationError(res, err)

		return
	}

	err = c.TransferHandler.Handle(command, c.CreateHandler)

	if err != nil {
		jsonBadRequest(res, err)

		return
	}

	jsonOk(res)
}

func (c Controller) GetBalance(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		jsonBadRequest(res, map[string]string{"id": "недопустимое значение"})

		return
	}

	acc, err := c.Fetcher.Get(uint(id))

	if err != nil {
		jsonBadRequest(res, err)

		return
	}

	if currency := req.URL.Query().Get("currency"); currency != "" {
		money, err := c.CurrencyConvertor.Convert(acc.Balance, currency)

		if err != nil {
			jsonBadRequest(res, err)

			return
		}

		jsonResponse(res, account.APIAccount{Balance: money})

		return
	}

	jsonResponse(res, acc)
}

func (c Controller) GetTransactions(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)

	if err != nil {
		jsonBadRequest(res, map[string]string{"id": "недопустимое значение"})

		return
	}

	jsonResponse(res, c.Fetcher.GetTransactions(req, uint(id)))
}

func jsonError(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
}

func jsonOk(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func jsonResponse(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)

	if err != nil {
		jsonError(res)
	}

	if _, err = res.Write(jsonData); err != nil {
		panic(err)
	}
}

func jsonBadRequest(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadRequest)
	jsonData, err := json.Marshal(data)

	if err != nil {
		jsonError(res)
	}

	if _, err = res.Write(jsonData); err != nil {
		panic(err)
	}
}

func jsonValidationError(res http.ResponseWriter, errs interface{}) {
	data := make(map[string]string)

	for _, err := range errs.(validator.ValidationErrors) {
		data[err.Field()] = err.Error()
	}

	jsonBadRequest(res, data)
}
