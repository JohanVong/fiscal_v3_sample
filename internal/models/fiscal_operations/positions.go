package fiscal_operations

import (
	"github.com/ericlagergren/decimal"
)

type Position struct {
	Name        string      `json:"Name" valid:"required~Поле Name не должно быть пустым"`
	IdSection   int         `json:"IdSection" valid:"required~Поле IdSection не должно быть пустым"`
	Price       decimal.Big `json:"Price" valid:"required~Поле Price не должно быть пустым, twoDecimalPlaces~Поле Price не должно иметь более двух десятичных знаков,isNonNegative~поле Price должно быть положительным"`
	Markup      decimal.Big `json:"Markup" valid:"twoDecimalPlaces~Поле Markup не должно иметь более двух десятичных знаков,isNonNegative~поле Markup должно быть положительным"`
	Discount    decimal.Big `json:"Discount" valid:"twoDecimalPlaces~Поле Discount не должно иметь более двух десятичных знаков,isNonNegative~поле Discount должно быть положительным"`
	Qty         decimal.Big `json:"Qty" valid:"required~Поле Qty не должно быть пустым,isNonNegative~поле Qty должно быть положительным"`
	Storno      bool        `json:"Storno"`
	ProductCode string      `json:"ProductCode"`
	IdUnit      int         `json:"IdUnit"`
}
