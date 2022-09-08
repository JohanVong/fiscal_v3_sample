package graylog

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JohanVong/fiscal_v3_sample/internal/services/logbuilder"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/systemconfig"
)

var GraylogClient *Client

const MaxMessageSize = 30000

type Client struct {
	client *http.Client
	ip     *net.IP
	port   string
}

func NewClient(address string) (c *Client, err error) {
	//defer recoverpanic.RecoverPanic()
	c = new(Client)
	ip, port, _, err := ParseIPPort(address)
	if err != nil {
		return nil, err
	}

	if port == "" {
		return nil, errors.New("port should not be empty")
	}

	c.client = http.DefaultClient
	c.ip = &ip
	c.port = port

	return c, nil
}

type GELFMessage struct {
	Host                      string                `json:"host"`
	ShortMessage              string                `json:"short_message"`
	SendMessage               string                `json:"sendMessage,omitempty"`
	SendDate                  *time.Time            `json:"sendDate,omitempty"`
	ReceiveMessage            string                `json:"receiveMessage,omitempty"`
	ReceiveDate               *time.Time            `json:"receiveDate,omitempty"`
	Facility                  string                `json:"facility,omitempty"`
	Address                   string                `json:"address,omitempty"`
	Endpoint                  string                `json:"endpoint,omitempty"`
	Method                    string                `json:"method,omitempty"`
	UserID                    int                   `json:"userId,omitempty"`
	Shift                     int                   `json:"shift,omitempty"`
	Kkm                       int                   `json:"kkm,omitempty"`
	IDOFD                     uint32                `json:"idOfd,omitempty"`
	ErrorMessageOFD           string                `json:"errorMessage,omitempty"`
	ErrorCodeOFD              uint32                `json:"errorCode,omitempty"`
	ZNM                       int64                 `json:"znm,omitempty"`
	TimeToBlock               string                `json:"timeToBlock,omitempty"`
	Environment               string                `json:"environment,omitempty"`
	TokenReqnum               string                `json:"tokenReqnum,omitempty"`
	Duration                  float64               `json:"duration,omitempty"`
	CheckLocation             string                `json:"check_location,omitempty"`
	DocNumber                 string                `json:"doc_number,omitempty"`
	UID                       string                `json:"uid,omitempty"`
	DetailedLog               string                `json:"detailed_log,omitempty"`
	OperationStageCode        int                   `json:"operation_stage_code,omitempty"`
	OperationStageDescription string                `json:"operation_stage_desc,omitempty"`
	TotalDuration             float64               `json:"total_duration,omitempty"`
	LogEvent                  string                `json:"log_event,omitempty"`
	logBuilder                logbuilder.LogBuilder `json:"-"`
}

func (m *GELFMessage) InitLogBuilder(builder logbuilder.LogBuilder) {
	v, ok := systemconfig.SystemConfig.GetOption("operation_log_enabled")
	if ok == true {
		enabled, err := strconv.ParseBool(v)
		if err == nil {
			if enabled == true {
				m.logBuilder = builder
			}
		}
	}
}

func (m *GELFMessage) SetUserID(id int) {
	m.UserID = id
	if m.logBuilder != nil {
		m.logBuilder.SetUserID(id)
	}
}

func (m *GELFMessage) SetUID(uid string) {
	m.UID = uid
	if m.logBuilder != nil {
		m.logBuilder.SetUID(uid)
	}
}

func (m *GELFMessage) SetKKM(kkm int) {
	m.Kkm = kkm
	if m.logBuilder != nil {
		m.logBuilder.SetKKM(kkm)
	}
}

func (m *GELFMessage) SetZNM(znm int64) {
	m.ZNM = znm
	if m.logBuilder != nil {
		m.logBuilder.SetZNM(znm)
	}
}

func (m *GELFMessage) SetShift(shift int) {
	m.Shift = shift
	if m.logBuilder != nil {
		m.logBuilder.SetShift(shift)
	}
}

func (m *GELFMessage) WriteEntry(message string, code int) {
	if m.logBuilder != nil {
		m.logBuilder.WriteEntry(message, code)
	}
}

func (m *GELFMessage) WriteSingleString(message string) {
	if m.logBuilder != nil {
		m.logBuilder.WriteSingleString(message)
	}
}

func NewGELFMessage(logEvent string) (m *GELFMessage) {
	//defer recoverpanic.RecoverPanic()
	m = new(GELFMessage)
	m.SendDate = new(time.Time)
	m.ReceiveDate = new(time.Time)
	host := os.Getenv("HOSTNAME")
	if host == "" {
		m.Host = "fiscalv2_local"
	} else {
		m.Host = host
	}
	m.Facility = "fiscalv2-api-" + strings.ToLower(os.Getenv("STAGE"))
	m.LogEvent = logEvent
	m.ShortMessage = "fiscalv2"
	return
}

func ParseIPPort(s string) (ip net.IP, port, space string, err error) {
	//defer recoverpanic.RecoverPanic()
	ip = net.ParseIP(s)
	if ip == nil {
		var host string
		host, port, err = net.SplitHostPort(s)
		if err != nil {
			return
		}
		if port != "" {
			if _, err = strconv.ParseUint(port, 10, 16); err != nil {
				return
			}
		}
		ip = net.ParseIP(host)
	}
	if ip == nil {
		err = errors.New("invalid address format")
	} else {
		space = "IPv6"
		if ip4 := ip.To4(); ip4 != nil {
			space = "IPv4"
			ip = ip4
		}
	}
	return
}

/*
func GetOperlog(c *context.Context) (l *GELFMessage) {
	v, ok := c.Input.Data()["oper_log"]
	if ok == true {
		l = v.(*GELFMessage)
	} else {
		l = NewGELFMessage("OperationLogFull")
		l.logBuilder = new(OperationLog)
		c.Input.SetData("oper_log", l)
	}

	return l
}
*/

func SendMessage(m *GELFMessage) (resp *http.Response, err error) {
	//defer recoverpanic.RecoverPanic()
	if m.logBuilder != nil {
		m.DetailedLog = m.logBuilder.String()
		m.Duration = m.logBuilder.Duration().Seconds()
		m.logBuilder.Send()
	}
	payload, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(payload)

	//return GraylogClient.client.Post("http://"+GraylogClient.ip.String()+":"+GraylogClient.port+"/gelf", "application/json", reader)
	resp, err = GraylogClient.client.Post("http://"+GraylogClient.ip.String()+":"+GraylogClient.port+"/gelf", "application/json", reader)
	return
}
