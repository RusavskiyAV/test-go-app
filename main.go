package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	accountControl "test3/src/controller"
	accountModel "test3/src/model/account/entity"
	accountUseCaseAdd "test3/src/model/account/useCase/add"
	accountUseCaseCreate "test3/src/model/account/useCase/create"
	accountUseCaseSubtract "test3/src/model/account/useCase/subtract"
	accountUseCaseTransfer "test3/src/model/account/useCase/transfer"
	accountReadModel "test3/src/readModel/account"
	currencyConvertor "test3/src/service"
)

func init() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=localhost user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_DATABASE"),
			viper.GetInt("DB_PORT"),
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	accountRepository := accountModel.Repository{Db: db}
	accountController := accountControl.Controller{
		CreateHandler:     &accountUseCaseCreate.Handler{AccountRepository: &accountRepository},
		AddHandler:        &accountUseCaseAdd.Handler{AccountRepository: &accountRepository},
		SubtractHandler:   &accountUseCaseSubtract.Handler{AccountRepository: &accountRepository},
		TransferHandler:   &accountUseCaseTransfer.Handler{AccountRepository: &accountRepository},
		Fetcher:           &accountReadModel.Fetcher{Db: db},
		Validate:          validator.New(),
		CurrencyConvertor: &currencyConvertor.CurrencyConverter{APIKey: viper.GetString("API_LAYER_KEY")},
	}

	router := httprouter.New()

	router.POST("/users/:id/balance/add", accountController.Add)
	router.POST("/users/:id/balance/subtract", accountController.Subtract)
	router.POST("/users/:id/balance/transfer", accountController.Transfer)
	router.GET("/users/:id/balance", accountController.GetBalance)
	router.GET("/users/:id/balance/transactions", accountController.GetTransactions)

	if err = http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
