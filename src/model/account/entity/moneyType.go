package account

import (
	"database/sql/driver"
	"errors"
	"fmt"
	goMoney "github.com/Rhymond/go-money"
)

type Money struct {
	*goMoney.Money
}

func (m *Money) Scan(value interface{}) error {
	intValue, ok := value.(int64)

	if !ok {
		return errors.New(fmt.Sprint("failed to get int64 value:", value))
	}

	money, err := NewMoney(intValue)

	if err != nil {
		return err
	}

	*m = *money

	return nil
}

func (m Money) Value() (driver.Value, error) {
	return m.Amount(), nil
}

func (m *Money) GormDataType() string {
	return "int64"
}

func NewMoney(s int64) (*Money, error) {
	if s < 0 {
		return &Money{}, errors.New("значение не может быть отрицательным")
	}

	return &Money{goMoney.New(s, goMoney.RUB)}, nil
}

func NewMoneyFromFloatWithCurrency(s float64, c string) (*Money, error) {
	if s < 0 {
		return &Money{}, errors.New("значение не может быть отрицательным")
	}

	return &Money{goMoney.NewFromFloat(s, c)}, nil
}
