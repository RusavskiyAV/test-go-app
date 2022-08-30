package account

import (
	"errors"
	"gorm.io/gorm"
	transaction "test3/src/model/transaction/entity"
	"time"
)

type RepositoryInterface interface {
	Create(a *Account) error
	AddToBalanceByUserId(id uint, money *Money, date time.Time) error
	SubtractFromBalanceByUserId(id uint, money *Money, date time.Time) error
	TransferByUserId(senderId uint, receiverId uint, money *Money, date time.Time) error
}

type Repository struct {
	Db *gorm.DB
}

func (ur *Repository) Create(a *Account) error {
	res := ur.Db.Save(a)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (ur *Repository) AddToBalanceByUserId(id uint, money *Money, date time.Time) error {
	err := ur.Db.Transaction(func(tx *gorm.DB) error {
		res := ur.Db.Exec("UPDATE accounts SET balance = balance + ? WHERE user_id = ?", money.Amount(), id)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ошибка пополнения баланса")
		}

		if err := ur.addTransaction(ur.getAccountIdByUserId(id), money, transaction.ReasonsAdd, nil, date); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ur *Repository) SubtractFromBalanceByUserId(id uint, money *Money, date time.Time) error {
	err := ur.Db.Transaction(func(tx *gorm.DB) error {
		res := ur.Db.Exec("UPDATE accounts SET balance = balance - ? WHERE user_id = ?", money.Amount(), id)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ошибка списания с баланса")
		}

		if err := ur.addTransaction(ur.getAccountIdByUserId(id), money, transaction.ReasonsSubtract, nil, date); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ur *Repository) TransferByUserId(senderId, receiverId uint, money *Money, date time.Time) error {
	err := ur.Db.Transaction(func(tx *gorm.DB) error {
		if err := ur.SubtractFromBalanceByUserId(senderId, money, date); err != nil {
			return err
		}

		if err := ur.AddToBalanceByUserId(receiverId, money, date); err != nil {
			return err
		}

		accountIdSender := ur.getAccountIdByUserId(senderId)
		accountIdReceiver := ur.getAccountIdByUserId(receiverId)

		if err := ur.addTransaction(accountIdSender, money, transaction.ReasonsTransfer, accountIdReceiver, date); err != nil {
			return err
		}

		if err := ur.addTransaction(accountIdReceiver, money, transaction.ReasonsTransfer, accountIdSender, date); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ur *Repository) getAccountIdByUserId(id uint) uint {
	var accountId uint

	ur.Db.Raw("SELECT id FROM accounts WHERE user_id = ?", id).Scan(&accountId)

	return accountId
}

func (ur *Repository) addTransaction(accountId uint, money *Money, reason string, participantAccountId interface{}, date time.Time) error {
	res := ur.Db.Exec("INSERT INTO transactions(account_id, amount, reason, participant_account_id, created_at) VALUES(?, ?, ?, ?, ?)", accountId, money.Amount(), reason, participantAccountId, date)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return errors.New("ошибка записи операции в таблицу транзакций")
	}

	return nil
}
