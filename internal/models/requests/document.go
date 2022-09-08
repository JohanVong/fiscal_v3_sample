package requests

import (
	"time"

	"github.com/ericlagergren/decimal"
)

type ZreqList struct {
	Phone    string `form:"phone"`
	IdUser   int    `form:"iduser"`
	DateFrom string `form:"datefrom" valid:"isEscapedDateTime~Формат параметра datefrom неверный: ошибка в процентном кодировании или в формате даты (YYYY-MM-DDThh:mm:ss)"`
	DateTo   string `form:"dateto" valid:"isEscapedDateTime~Формат параметра dateto неверный: ошибка в процентном кодировании или в формате даты (YYYY-MM-DDThh:mm:ss)"`
}

type ZreqOpersByShift struct {
	IdShift int `form:"idshift"`
}

type RequestDateRange struct {
	DateFrom string `form:"datefrom" valid:"isEscapedDateTime~Формат параметра datefrom неверный: ошибка в процентном кодировании или в формате даты (YYYY-MM-DDThh:mm:ss)"`
	DateTo   string `form:"dateto" valid:"isEscapedDateTime~Формат параметра dateto неверный: ошибка в процентном кодировании или в формате даты (YYYY-MM-DDThh:mm:ss)"`
}

type OfflineDocument struct {
	IdKkm            int         `json:"IdKkm"`
	PhoneLogin       string      `json:"PhoneLogin"`
	IdDomain         int         `json:"IdDomain"`
	IdTypeDocument   int         `json:"IdTypeDocument"`
	Positions        []Position  `json:"Positions"`
	Value            decimal.Big `json:"Value"`
	Cash             decimal.Big `json:"Cash"`
	NonCash          decimal.Big `json:"NonCash"`
	Change           decimal.Big `json:"Change"`
	AutonomousNumber uint32      `json:"AutonomousNumber"`
	CustomerIIN      string      `json:"CustomerIIN"`
	ReceiptDate      time.Time   `json:"ReceiptDate"`
}
