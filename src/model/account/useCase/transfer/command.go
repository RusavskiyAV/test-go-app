package account

type Command struct {
	UserId     uint  `json:"-" validate:"required,gte=0"`
	Amount     int64 `validate:"required,gte=0"`
	ReceiverId uint  `validate:"required,gte=0,necsfield=UserId"`
}
