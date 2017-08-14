package models

import (
	"application/entities"
	"application/libraries/helpers"
)

type User struct {
	helpers.Model
	UserEntity *entities.User
}

func NewUser() *User {
	user := new(User)
	user.UserEntity = &entities.User{}
	user.Model.Construct("user")
	return user
}

func (this *User) GetUserByName(username string) (has bool, err error) {
	has, err = this.Model.DB.XORM.Where("name = ?", username).Get(this.UserEntity)

	return has, err
}

func (this *User) GetUserById(id int) (has bool, err error) {
	has, err = this.Model.DB.XORM.Id(id).Get(this.UserEntity)

	return has, err
}

func (this *User) Insert(user *entities.User) (int, error) {
	affected, err := this.Model.DB.XORM.Insert(user)
	if err != nil {
		panic(err)
	}
	return int(affected), err
}
