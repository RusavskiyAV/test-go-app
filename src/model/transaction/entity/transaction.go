package transaction

import (
	"database/sql"
	"time"
)

const (
	ReasonsAdd      = "Add"
	ReasonsSubtract = "Subtract"
	ReasonsTransfer = "Transfer"
)

type Transaction struct {
	ID                   uint
	AccountId            uint
	Amount               int64
	Reason               string
	ParticipantAccountId sql.NullInt64
	CreatedAt            time.Time
}
