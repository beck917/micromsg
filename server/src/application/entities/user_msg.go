package entities

import ()

type UserMsg struct {
	Id          int      `json:"id" xorm:"not null pk INT(11)"`
	SendUid     int      `json:"send_uid" xorm:"not null index INT(11)"`
	RecvUid     int      `json:"recv_uid" xorm:"not null index INT(11)"`
	Msg         string   `json:"msg" xorm:"not null VARCHAR(1023)"`
	DeletedTime int      `json:"deleted_time" xorm:"not null INT(11)"`
	CreatedTime JsonTime `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdatedTime JsonTime `json:"updated_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}
