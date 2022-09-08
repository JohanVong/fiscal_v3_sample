package models

import (
	"sync"
	"time"
)

//var ShiftCache = make(map[int]Shift, 100)
type ShiftCacheST struct {
	mx sync.Mutex
	m  map[int]Shift
}

func (c *ShiftCacheST) Load(key int) (Shift, bool) {
	c.mx.Lock()
	val, ok := c.m[key]
	c.mx.Unlock()
	return val, ok
}

func (c *ShiftCacheST) Store(key int, value Shift) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *ShiftCacheST) Delete(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, key)
}

var ShiftCache = ShiftCacheST{m: make(map[int]Shift, 100)}

type AvrgOper struct {
	Count int
	Time  int64
}

var AvOp AvrgOper
var Och = make(chan int64, 150)

var KKMOFDStatusStorage map[int]*KKMOFDStatus
var OfdStatusMutex sync.Mutex

type KKMOFDStatus struct {
	IDKKM     int       `json:"kkm"`
	ZNM       int       `json:"znm"`
	StatusOFD string    `json:"status"`
	Code      uint32    `json:"code"`
	Time      time.Time `json:"time"`
}
