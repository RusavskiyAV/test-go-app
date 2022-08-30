package account

import (
	transaction "test3/src/model/transaction/entity"
	"time"
)

const defaultBalance int64 = 0

type Account struct {
	ID           uint
	Balance      *Money
	UserID       uint
	CreatedAt    time.Time
	Transactions []transaction.Transaction
}

type APIAccount struct {
	Balance *Money
}

func New(userId uint) (*Account, error) {
	money, err := NewMoney(defaultBalance)

	if err != nil {
		return &Account{}, err
	}

	return &Account{UserID: userId, Balance: money}, nil
}
