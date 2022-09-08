package requests

type ReceiptRequest struct {
	Date         string `form:"datetime" valid:"required~Поле datetime не должно быть пустым, isEscapedDateTimeNoSeconds~Формат параметра datefrom неверный: ошибка в процентном кодировании или в формате даты (YYYY-MM-DDThh:mm)"`
	ZNM          int    `form:"znm" valid:"required~Поле znm не должно быть пустым"`
	FiscalNumber int64  `form:"fn" valid:"required~Поле fn не должно быть пустым"`
	TotalValue   string `form:"total" valid:"required~Поле total не должно быть пустым, twoDecimalPlaces~Поле total не должно иметь более двух десятичных знаков"`
}

type ReceiptResponse struct {
	Receipt string `json:"Receipt"`
}
