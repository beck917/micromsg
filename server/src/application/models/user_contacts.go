package models

import (
	"application/entities"
	"application/libraries/helpers"
)

type UserContacts struct {
	helpers.Model
	UserContactsEntity *entities.UserContacts
}

func NewUserContacts() *UserContacts {
	userContacts := new(UserContacts)
	userContacts.UserContactsEntity = &entities.UserContacts{}
	userContacts.Model.Construct("user_contacts")
	return userContacts
}

func (this *UserContacts) GetContactsListByUid(uid int) (contactsList []*entities.UserContacts, err error) {
	userContacts := entities.UserContacts{}
	rows, err := this.Model.DB.XORM.Where("uid = ?", uid).Rows(&userContacts)

	if err != nil {
		return
	}

	if rows != nil {
		for rows.Next() {
			userContactsRow := entities.UserContacts{}
			err = rows.Scan(&userContactsRow)
			if err != nil {
				break
			}
			contactsList = append(contactsList, &userContactsRow)
		}
	}
	return
}

func (this *UserContacts) Insert(user *entities.UserContacts) (int, error) {
	affected, err := this.Model.DB.XORM.Insert(user)

	return int(affected), err
}
