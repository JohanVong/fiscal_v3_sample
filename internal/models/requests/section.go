package requests

type SectionsRequest struct {
	KKM string `json:"kkm" valid:"required~Поле kkm не должно быть пустым"`
}
