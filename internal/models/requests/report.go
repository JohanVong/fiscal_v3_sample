package requests

import "time"

type XReportRequest struct {
	Date time.Time `json:"Date"`
}
