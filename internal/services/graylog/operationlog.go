package graylog

import (
	"encoding/json"
	"time"
)

const (
	AuthenticationLoginStart = iota + 1
	AuthenticationLoginReadRequestBody
	AuthenticationLoginReadBodyIntoStruct
	AuthenticationLoginRetrievingUserFromDB
	AuthenticationLoginComparingPasswordHash
	AuthenticationLoginGeneratingToken
	AuthenticationLoginPH5
	AuthenticationLoginPH6
	AuthenticationLoginPH7
	AuthenticationTokenStart
	AuthenticationTokenParse
	AuthenticationTokenRetrievingUserFromDB
	AuthenticationTokenRetrievingUserPermissionsFromDB
	AuthenticationTokenUnauthorized
	AuthenticationTokenPH5
	AuthenticationTokenPH6
	AuthenticationTokenPH7
	AuthorizationGeneralStart
	AuthorizationKKMStart
	AuthorizationShiftStart
	AuthorizationRetrievingUserGroup
	AuthorizationRetrievingHTTPMethod
	AuthorizationRetrievingRequestURI
	AuthorizationForbidden
	AuthorizationRetrievingShiftFromDB
	AuthorizationRetrievingGroupUserArrayFromDb
	AuthorizationRetrievingTypePermArrayFromDB
	AuthorizationRetrievingTypeObjectArrayFromDb
	AuthorizationRetrievingPermissionsForUser
	AuthorizationCacheKKM
	AuthorizationCheckingRedisForKKM
	AuthorizationPH2
	AuthorizationPH3
	ControllerSaleStart
	ControllerSaleReadintRequestBody
	ControllerSaleRetrievingUIDFromRedis
	ControllerSaleRegeneratingUID
	ControllerSaleReadingBodyIntoStruct
	ControllerSaleValidatingMessageFields
	ControllerSaleRetrievingLastDocument
	ControllerSaleCheckingCacheIfOperationInProgress
	ControllerSaleFetchingReeiptDuplicate
	ControllerSaleWritingUIDToRedis
	ControllerSaleFetchingKKMFromContext
	ControllerSaleReturnFromProcessing
	ControllerSaleSendingOperationTimeToMetrics
	ControllerSaleDeletingUIDFromRedis
	ControllerSalePH14
	ControllerSalePH15
	ControllerSalePH16
	ControllerSalePH17
	ControllerSalePH18
	ControllerSalePH19
	ControllerSalePH20
	ControllerSalePH21
	ControllerSalePH22
	ControllerSalePH23
	ControllerSalePH24
	ControllerSalePH25
	ControllerSalePH26
	ControllerSalePH27
	ControllerSalePH28
	ControllerSalePH29
	ControllerSalePH30
	ControllerSalePH31
	ControllerSalePH32
	ControllerSalePH33
	ControllerSalePH34
	ControllerSalePH35
	ControllerSalePH36
	ControllerSalePH37
	ControllerSalePH38
	ControllerSalePH39
	ControllerSalePH40
	ControllerSalePH41
	ControllerSalePH42
	ControllerSalePH43
	ControllerSalePH44
	ControllerSalePH45
	ControllerSalePH46
	ControllerSalePH47
	ControllerSalePH48
	ControllerSalePH49
	ControllerSalePH50
)

type OperationLog struct {
	firstMessageTime    *time.Time
	Uid                 string         `json:"uid"`
	ZNM                 int64          `json:"znm"`
	KKM                 int            `json:"kkm"`
	Shift               int            `json:"shift"`
	UserID              int            `json:"user_id"`
	Stages              []*GELFMessage `json:"Stages"`
	previousMessageTime *time.Time
}

func (l *OperationLog) SetUserID(userID int) {
	l.UserID = userID
}

func (l *OperationLog) SetUID(uid string) {
	l.Uid = uid
}

func (l *OperationLog) SetKKM(kkm int) {
	l.KKM = kkm
}

func (l *OperationLog) SetZNM(znm int64) {
	l.ZNM = znm
}

func (l *OperationLog) SetShift(shift int) {
	l.Shift = shift
}

func (l *OperationLog) Duration() time.Duration {
	if l.firstMessageTime != nil {
		return time.Now().Sub(*l.firstMessageTime)
	}
	return 0
}

func (l *OperationLog) String() string {
	res, _ := json.Marshal(l)
	return string(res)
}

func (l *OperationLog) WriteEntry(message string, code int) {

	s := NewGELFMessage("OperationLog")
	if l.previousMessageTime == nil {
		l.previousMessageTime = new(time.Time)
		l.firstMessageTime = new(time.Time)
		*l.previousMessageTime = time.Now()
		*l.firstMessageTime = time.Now()
	}
	eventTime := time.Now()
	s.Duration = eventTime.Sub(*l.previousMessageTime).Seconds()
	s.OperationStageCode = code
	s.SendMessage = message
	s.SendDate = new(time.Time)
	*s.SendDate = time.Now()

	l.Stages = append(l.Stages, s)
	l.previousMessageTime = &eventTime
}

func (l *OperationLog) WriteSingleString(string) {

}

func (l *OperationLog) Send() {
	duration := l.Duration()
	go func() {
		for _, s := range l.Stages {
			s.TotalDuration = duration.Seconds()
			s.UID = l.Uid
			s.ZNM = l.ZNM
			s.Kkm = l.KKM
			s.Shift = l.Shift
			s.UserID = l.UserID
			_, _ = SendMessage(s)
		}
	}()
}
