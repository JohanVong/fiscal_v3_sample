package logbuilder

import (
	"fmt"
	"strings"
	"time"
)

type LogBuilder interface {
	Duration() time.Duration
	String() string
	WriteEntry(string, int)
	WriteSingleString(string)
	Send()
	SetUID(string)
	SetKKM(int)
	SetZNM(int64)
	SetShift(int)
	SetUserID(int)
}

type StringLogBuilder struct {
	builder             strings.Builder
	firstMessageTime    *time.Time
	previousMessageTime *time.Time
}

func (l *StringLogBuilder) SetUserID(int) {

}

func (l *StringLogBuilder) SetUID(string) {

}

func (l *StringLogBuilder) SetKKM(int) {

}

func (l *StringLogBuilder) SetZNM(int64) {

}

func (l *StringLogBuilder) SetShift(int) {

}

func (l *StringLogBuilder) Send() {

}

func (l *StringLogBuilder) Duration() time.Duration {
	if l.firstMessageTime == nil {
		return 0
	}
	return time.Now().Sub(*l.firstMessageTime)
}

func (l *StringLogBuilder) String() string {
	l.WriteSingleString(fmt.Sprintf("TOTAL TIME: %s", time.Now().Sub(*l.firstMessageTime).String()))
	return l.builder.String()
}

func (l *StringLogBuilder) WriteEntry(message string, code int) {
	eventTime := time.Now()
	l.builder.WriteString(fmt.Sprintf("%d. %s\n", code, message))
	l.builder.WriteString(eventTime.String() + "\n")
	if l.previousMessageTime == nil {
		l.previousMessageTime = new(time.Time)
		l.firstMessageTime = new(time.Time)
		*l.previousMessageTime = time.Now()
		*l.firstMessageTime = time.Now()
		l.builder.WriteString("\n")
	} else {
		eventDuration := eventTime.Sub(*l.previousMessageTime)
		l.builder.WriteString("Time since previous event: " + eventDuration.String())
		if eventDuration > time.Second*6 {
			l.builder.WriteString("<---------------------!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n")
		} else {
			l.builder.WriteString("\n\n")
		}
	}
	l.previousMessageTime = &eventTime
}

func (l *StringLogBuilder) WriteSingleString(message string) {
	l.builder.WriteString(message + "\n")
}
