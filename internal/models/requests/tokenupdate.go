package requests

type TokenUpdateRequest struct {
	Token         uint32 `json:"token"`
	IdCPCR        uint32 `json:"id_cpcr"`
	RNM           string `json:"rnm"`
	OFDType       int    `json:"ofd_type" valid:"isNonNegative~поле ofd_type не должно быть отрицательным"`
	PingAndNotify bool   `json:"ping_and_notify"`
}
