package models

import (
	"strconv"
)

type ResponseV2 struct {
	Status int    `json:"Status"`
	Msg    string `json:"Message,omitempty"`
}

func (self *ResponseV2) Error() string {
	return strconv.Itoa(self.Status) + " : " + self.Msg
}

func (self *ResponseV2) GetStatus() int {
	return self.Status
}

func (self *ResponseV2) GetMsg() string {
	return self.Msg
}
