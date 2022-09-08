package ofd_processor

import (
	"errors"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/ofd"
)

func SendToOfd(idKkm int, doc models.Document, positions []models.Position, command ofd.OfdCommand) (fiscal_number uint64, is_authonomus bool, err error) {
	var (
		kkm models.Kkm
	)

	tx := db.Orm.Begin()
	defer func() {
		tx.RollbackUnlessCommitted()
	}()

	tx.First(&kkm, idKkm)
	header, ticketNumber, ofd_err := ofd.OfdPool.ProcessTicketRequest(command, &kkm, &doc, positions, tx)
	if ofd_err != nil {
		err = ofd_err
		return
	}
	if ticketNumber == nil {
		err = errors.New("Error ticketNumber is nill")
		return
	}

	if ticketNumber.Offline == true {
		tx.Model(&models.Document{}).Where("idDocuments=?", doc.Id).Updates(map[string]interface{}{"AutonomousNumber": ticketNumber.Number})
	} else {
		tx.Model(&models.Document{}).Where("idDocuments=?", doc.Id).Updates(map[string]interface{}{"FiscalNumber": ticketNumber.Number})
	}

	if header != nil {
		tx.Model(&models.Kkm{}).Where("idKKM=?", idKkm).Updates(map[string]interface{}{"idCPCR": header.Id, "tokenCPCR": header.Token, "reqnumCPCR": header.ReqNum})
	}
	tx.Commit()

	fiscal_number = ticketNumber.Number
	is_authonomus = ticketNumber.Offline
	return

}
