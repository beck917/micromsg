package models

import (
	"application/entities"
	"application/libraries/constant"
	"application/libraries/helpers"
	"strconv"
	"testing"
	"time"
)

func Test_Insert(t *testing.T) {
	userModel := NewUser()
	userContactsModel := NewUserContacts()

	for i := 1; i <= 13; i++ {
		id := strconv.Itoa(i)
		userModel.UserEntity.Id = 0
		userModel.UserEntity.Name = "user" + id
		userModel.UserEntity.Password = helpers.Sha1Str("123456" + constant.PASSWORD_SALT)
		userModel.UserEntity.CreatedTime = entities.JsonTime(time.Now())
		userModel.UserEntity.UpdatedTime = entities.JsonTime(time.Now())
		userModel.Insert(userModel.UserEntity)

		for j := 1; j <= 13; j++ {
			if i != j {
				userContactsModel.UserContactsEntity.Id = 0
				userContactsModel.UserContactsEntity.Uid = i
				userContactsModel.UserContactsEntity.Cid = j
				userContactsModel.UserContactsEntity.Cname = "user" + strconv.Itoa(j)
				userContactsModel.UserContactsEntity.CreatedTime = entities.JsonTime(time.Now())
				userContactsModel.UserContactsEntity.UpdatedTime = entities.JsonTime(time.Now())
				userContactsModel.Insert(userContactsModel.UserContactsEntity)
			}
		}
	}
}
