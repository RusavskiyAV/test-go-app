package account

import "testing"

func TestNewMoney(t *testing.T) {
	_, err := NewMoney(-1)

	if err == nil {
		t.Error("ожидалась ошибка")
	}
}

func TestNewMoneyFromFloatWithCurrency(t *testing.T) {
	_, err := NewMoney(-1)

	if err == nil {
		t.Error("ожидалась ошибка")
	}
}
