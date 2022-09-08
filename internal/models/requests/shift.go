package requests

import "time"

type CreateShiftRequest struct {
	KKM  uint64    `json:"IdKKM" valid:"required~Поле IdKKM не должно быть пустым"`
	Date time.Time `json:"Date"`
}

type CloseShiftRequest struct {
	KKM  uint64    `json:"IdKKM" valid:"required~Поле IdKKM не должно быть пустым"`
	Date time.Time `json:"Date"`
}

type CreateShiftResponse struct {
	Shift struct {
		Id            int
		IdUser        int
		IdKkm         int
		IdStatusShift int
		DateOpen      time.Time
	}
}

type EditShiftRequest struct {
	Status string `json:"status" valid:"required~Поле status не должно быть пустым,in(closed)~Значение поля status может быть только 'closed'"`
}

type GetShiftsFilterRequest struct {
	Status string `json:"Status" valid:"in(open|closed)~Значение поля Status должно быть либо open, либо closed"`
}
