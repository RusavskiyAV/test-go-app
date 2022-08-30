package account

import account "test3/src/model/account/entity"

type Interface interface {
	Handle(c Command) error
}

type Handler struct {
	AccountRepository account.RepositoryInterface
}

func (h *Handler) Handle(c Command) error {
	acc, err := account.New(c.UserId)

	if err != nil {
		return err
	}

	err = h.AccountRepository.Create(acc)

	if err != nil {
		return err
	}

	return nil
}
