package user

import "time"

type User struct {
	ID        uint
	firstName string
	lastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
