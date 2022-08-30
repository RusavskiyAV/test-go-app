package account

type Command struct {
	UserId uint `json:"-" validate:"required,gte=0"`
}
