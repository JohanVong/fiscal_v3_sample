package requests

import (
	"time"

	"github.com/ericlagergren/decimal"
)

/*
 * {
  "Amount": 800
}
*/
type IncomeReq struct {
	IdKkm  uint64      `json:"-"`
	Amount decimal.Big `json:"Amount" valid:"required~Поле Amount не должно быть пустым,twoDecimalPlaces~Поле Price	не должно иметь более двух десятичных знаков,isNonNegative~поле Amount не должно быть отрицательным"`
	Uid    string      `json:"Uid" valid:"required~Заголовок не содержит Uid, uuidv4~Uid не соответствует формату"`
	Date   time.Time   `json:"Date"`
}
