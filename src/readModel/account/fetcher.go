package account

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	account "test3/src/model/account/entity"
	transaction "test3/src/model/transaction/entity"
)

type Interface interface {
	Get(id uint) (*account.APIAccount, error)
	GetTransactions(req *http.Request, userId uint) *[]transaction.Transaction
}

type Fetcher struct {
	Db *gorm.DB
}

func (f *Fetcher) Get(id uint) (*account.APIAccount, error) {
	acc := account.APIAccount{}
	res := f.Db.Model(&account.Account{}).Where("user_id = ?", id).First(&acc)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &acc, errors.New("у пользователя отсутствует баланс")
	}

	return &acc, nil
}

func (f *Fetcher) GetTransactions(req *http.Request, userId uint) *[]transaction.Transaction {
	var res []transaction.Transaction

	orders := map[string]struct{}{
		"id":         {},
		"amount":     {},
		"created_at": {},
	}

	f.Db.Scopes(f.paginate(req, orders)).Joins("join accounts on accounts.id = transactions.account_id and accounts.user_id = ?", userId).Find(&res)

	return &res
}

func (f *Fetcher) paginate(req *http.Request, orders map[string]struct{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		const defaultOrder = "id"

		q := req.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))

		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		orderString := ""

		order := q.Get("order")

		if _, inMap := orders[order]; inMap {
			orderString = order
		} else {
			orderString = defaultOrder
		}

		isDesc, _ := strconv.ParseBool(q.Get("desc"))

		if isDesc {
			orderString += " desc"
		} else {
			orderString += " asc"
		}

		return db.Offset(offset).Limit(pageSize).Order(orderString)
	}
}
