package requests

import "time"

type OperRep struct {
	KKM      string    `json:"-" valid:"required~Не указан IdKkm"`
	FromDate time.Time `json:"FromDate" valid:"required~Поле FromDate не должно быть пустым"`
	ToDate   time.Time `json:"ToDate" valid:"required~Поле ToDate не должно быть пустым"`
}
