package requests

import (
	"time"
)

type RefundReq struct {
	IdDocument		int			`json:"IdDocument" valid:"required~Поле IdDocument не должно быть пустым"`
	AFP				uint32		`json:"AFP"`
	ChkLocation   	string      `json:"ChkLocation"`
	GenerateCheck	bool        `json:"GenerateCheck"`
	ReceiptDate   	time.Time	`json:"ReceiptDate"`
} 