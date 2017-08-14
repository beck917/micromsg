package entities

import ()

type User struct {
	Id          int      `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name        string   `json:"name" xorm:"not null unique VARCHAR(30)"`
	Password    string   `json:"password" xorm:"not null VARCHAR(32)"`
	CreatedTime JsonTime `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdatedTime JsonTime `json:"updated_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}
