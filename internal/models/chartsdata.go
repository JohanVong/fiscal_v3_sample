package models

import "time"

type ChartsDataResp struct {
	AvgCheck     Decimal     `json:"AvgCheck"`
	Sells        opersInfo   `json:"Sells"`
	Refunds      opersInfo   `json:"Refunds"`
	OldSells     int64       `json:"OldSells"`
	OldRefunds   int64       `json:"OldRefunds"`
	OldAvgCheck  int64       `json:"OldAvgCheck"`
	TopPositions interface{} `json:"TopPositions"`
	ChartsData   chartsData  `json:"ChartsData"`
	Time         time.Time   `json:"Time"`
}

type opersInfo struct {
	Cash          Decimal
	NonCash       Decimal
	Total         Decimal
	TotalReceipts int `json:"TotalReceipts"`
}

type chartsData struct {
	LineChart lineChart `json:"LineChart"`
	PieChart  pieChart  `json:"PieChart"`
}

type lineChart struct {
	Operations []shrunkOperLineChart `json:"Operations"`
}

type pieChart struct {
	Cash          Decimal
	NonCash       Decimal
	TotalReceipts int `json:"TotalReceipts"`
}

type shrunkOperLineChart struct {
	Time          time.Time `json:"Time"`
	SellCash      Decimal   `json:"SellCash"`
	SellNonCash   Decimal   `json:"SellNonCash"`
	SellTotal     Decimal   `json:"SellTotal"`
	RefundCash    Decimal   `json:"RefundCash"`
	RefundNonCash Decimal   `json:"RefundNonCash"`
	RefundTotal   Decimal   `json:"RefundTotal"`
}

func InitShrunkOperLineChart(doc *Document) *shrunkOperLineChart {
	var pbc Decimal // paid by cash

	if doc.Cash.Cmp(&doc.Value.Big) >= 0 {
		pbc.Big = doc.Value.Big
	} else {
		pbc.Big = doc.Cash.Big
	}

	solc := &shrunkOperLineChart{}

	solc.Time = doc.DateDocument

	if doc.IdTypedocument == 1 {
		solc.SellCash = pbc
		solc.SellNonCash = doc.NonCash
		solc.SellTotal = doc.Value
	}

	if doc.IdTypedocument == 5 {
		solc.RefundCash = pbc
		solc.RefundNonCash = doc.NonCash
		solc.RefundTotal = doc.Value
	}

	return solc
}
