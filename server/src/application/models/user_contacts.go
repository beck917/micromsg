package models

import (
	"application/entities"
	"application/libraries/helpers"
	"database/sql"
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

func (this *UserContacts) UpdateUnread(uid int, cid int, step int) (sql.Result, error) {
	var sql string
	if step == 0 {
		sql = "update `user_contacts` set unread = ? where uid=? and cid = ?"
	} else {
		sql = "update `user_contacts` set unread = unread + ? where uid=? and cid = ?"
	}
	res, err := this.Model.DB.XORM.Exec(sql, step, uid, cid)

	return res, err
}

func (this *UserContacts) Update(id int, user *entities.UserContacts) (int, error) {
	affected, err := this.Model.DB.XORM.Id(id).Update(user)

	return int(affected), err
}
