package calculations

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ericlagergren/decimal"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

func GetDocumentChecksum(d models.Document) (string, error) {

	var hasher struct {
		Shift, User, Kkm int
		Value            decimal.Big
		Number           string
		Date             time.Time
		Host             string
	}
	var (
		data []byte
		err  error
	)
	//-----------------------------------
	hasher.Shift = d.IdShift
	hasher.User = d.IdUser
	hasher.Kkm = d.IdKkm
	hasher.Value = d.Value.Big
	hasher.Number = d.NumberDoc
	hasher.Date = d.DateDocument
	hasher.Host = os.Getenv("CONTAINER_HOST_ADDRESS")
	//-------------------------------------------
	data, err = json.Marshal(hasher)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:]), nil
}
