package models

import (
	"application/entities"
	"application/libraries/helpers"
)

type User struct {
	helpers.Model
	UserEntity *entities.UserMsg
}

func NewUser() *User {
	user := new(UserMsg)
	user.UserEntity = &entities.UserMsg{}
	user.Model.Construct("user_msg")
	return user
}

func (this *UserMsg) Insert(user *entities.UserMsg) (int, error) {
	affected, err := this.Model.DB.XORM.Insert(user)

	return int(affected), err
}
