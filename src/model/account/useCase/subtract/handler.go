package account

import (
	account "test3/src/model/account/entity"
	"time"
)

type Interface interface {
	Handle(c Command) error
}

type Handler struct {
	AccountRepository account.RepositoryInterface
}

func (h *Handler) Handle(c Command) error {
	money, err := account.NewMoney(c.Amount)

	if err != nil {
		return err
	}

	err = h.AccountRepository.SubtractFromBalanceByUserId(c.UserId, money, time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}
