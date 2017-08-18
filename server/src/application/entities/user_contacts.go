package entities

type UserContacts struct {
	Id          int      `json:"id" xorm:"not null pk autoincr INT(11)"`
	Uid         int      `json:"uid" xorm:"not null index INT(11)"`
	Cid         int      `json:"cid" xorm:"not null INT(11)"`
	Cname       string   `json:"cname" xorm:"not null VARCHAR(30)"`
	Unread      int      `json:"unread" xorm:"not null SMALLINT(11)"`
	CreatedTime JsonTime `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdatedTime JsonTime `json:"updated_time" xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}
