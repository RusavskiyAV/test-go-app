package account

import (
	account "test3/src/model/account/entity"
	accountUseCaseCreate "test3/src/model/account/useCase/create"
	"time"
)

type Interface interface {
	Handle(c Command, createHandler accountUseCaseCreate.Interface) error
}

type Handler struct {
	AccountRepository account.RepositoryInterface
}

func (h *Handler) Handle(c Command, createHandler accountUseCaseCreate.Interface) error {
	_ = createHandler.Handle(accountUseCaseCreate.Command{UserId: c.UserId})
	money, err := account.NewMoney(c.Amount)

	if err != nil {
		return err
	}

	err = h.AccountRepository.AddToBalanceByUserId(c.UserId, money, time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}
