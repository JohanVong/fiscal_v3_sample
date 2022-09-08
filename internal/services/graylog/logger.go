package graylog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type graylogWriter struct {
	graylogClient *Client
}

/*
// NewConsole create ConsoleWriter returning as LoggerInterface.
func NewConsole() logs.Logger {
	cl, _ := NewClient(beego.AppConfig.String("graylog"))

	return &graylogWriter{
		graylogClient: cl,
	}
}
*/

// Init init console logger.
// jsonConfig like '{"level":LevelTrace}'.
func (c *graylogWriter) Init(jsonConfig string) error {

	return nil
}

// WriteMsg write message in console.
func (c *graylogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if strings.Contains(msg, "panic") {
		m := NewGELFMessage("CrashLog")
		*m.SendDate = when
		if len(msg) > 30000 {
			m.SendMessage = msg[:30000]
		} else {
			m.SendMessage = msg
		}

		payload, _ := json.Marshal(m)

		reader := bytes.NewReader(payload)
		cl := http.DefaultClient
		_, _ = cl.Post("http://192.168.151.110:12222/gelf", "application/json", reader)
	}

	return nil
}

// Destroy implementing method. empty.
func (c *graylogWriter) Destroy() {

}

// Flush implementing method. empty.
func (c *graylogWriter) Flush() {

}
