package requests

import (
	"time"

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

type SellReq struct {
	IdKkm         uint64      `json:"-"`
	IdDomain      int         `json:"IdDomain" valid:"required~Поле IdDomain не должно быть пустым"`
	Cash          decimal.Big `json:"Cash" valid:"twoDecimalPlaces~Поле Cash не должно иметь более двух десятичных знаков,isNonNegative~поле Cash должно быть положительным"`
	NonCash       decimal.Big `json:"NonCash" valid:"twoDecimalPlaces~Поле NonCash не должно иметь более двух десятичных знаков,isNonNegative~поле NonCash должно быть положительным"`
	Mobile        decimal.Big `json:"Mobile" valid:"twoDecimalPlaces~Поле Mobile не должно иметь более двух десятичных знаков,isNonNegative~поле Mobile должно быть положительным"`
	Positions     []Position  `json:"Positions" valid:"required~Поле Positions не должно быть пустым"`
	Total         decimal.Big `json:"Total" valid:"required~Поле Total не должно быть пустым,twoDecimalPlaces~Поле Total не должно иметь более двух десятичных знаков,isNonNegative~поле Total должно быть положительным"`
	Uid           string      `json:"Uid" valid:"required~Заголовок не содержит Uid, uuidv4~Uid не соответствует формату"`
	AFP           uint32      `json:"AFP"`
	ChkLocation   string      `json:"ChkLocation"`
	GenerateCheck bool        `json:"GenerateCheck"`
	ReceiptDate   time.Time   `json:"ReceiptDate"`
	CustomerIIN   string      `json:"CustomerIIN"`
}
