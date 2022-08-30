package user

import "gorm.io/gorm"

type RepositoryInterface interface {
	Find(id uint) *User
}

type Repository struct {
	Db *gorm.DB
}

func (ur *Repository) Find(id uint) *User {
	user := User{ID: id}
	ur.Db.First(&user)

	return &user
}
