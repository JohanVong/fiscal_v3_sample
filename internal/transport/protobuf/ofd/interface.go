package ofd

import (
	"time"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

type HandlerOFD interface {
	SendOfflineModeInformation(start time.Time, args ...interface{}) (interface{}, error)
	SendSystemMessage(args ...interface{}) (interface{}, error)
	GetKKMInfo(args ...interface{}) (interface{}, error)
	SendOperationSale(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error)
	SendOperationSaleReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error)
	SendOperationPurchase(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error)
	SendOperationPurchaseReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error)
	SendTicketRollback(args ...interface{}) (interface{}, error)
	CloseShift(document *models.Document, args ...interface{}) (interface{}, error)
	RequestZReport(document *models.Document, args ...interface{}) (interface{}, error)
	RequestXReport(document *models.Document, args ...interface{}) (interface{}, error)
	DepositMoney(document *models.Document, args ...interface{}) (interface{}, error)
	WithdrawMoney(document *models.Document, args ...interface{}) (interface{}, error)
	SetTimeout(timeout time.Duration)
}
