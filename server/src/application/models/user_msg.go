package models

import (
	"application/entities"
	"application/libraries/helpers"
	"time"
)

type UserMsg struct {
	helpers.Model
	UserMsgEntity *entities.UserMsg
}

func NewUserMsg() *UserMsg {
	user := new(UserMsg)
	user.UserMsgEntity = &entities.UserMsg{}
	user.Model.Construct("user_msg")
	return user
}

func (this *UserMsg) GetMsgListByUid(uid int, open_id int, page int, pageSize int) (msgList []*entities.UserMsg, err error) {
	userMsg := entities.UserMsg{}
	rows, err := this.Model.DB.XORM.Where("((send_uid = ? and recv_uid = ? ) or (recv_uid = ? and send_uid = ?)) and deleted_time = 0", uid, open_id, uid, open_id).Desc("id").Limit(pageSize, page*pageSize).Rows(&userMsg)

	if err != nil {
		return
	}

	if rows != nil {
		for rows.Next() {
			userMsgRow := entities.UserMsg{}
			err = rows.Scan(&userMsgRow)
			if err != nil {
				break
			}
			msgList = append(msgList, &userMsgRow)
		}
	}
	return
}

func (this *UserMsg) Insert(user *entities.UserMsg) (int, error) {
	affected, err := this.Model.DB.XORM.Insert(user)

	return int(affected), err
}

func (this *UserMsg) DeleteMsgById(id int) (int, error) {
	userMsg := &entities.UserMsg{}
	userMsg.DeletedTime = int(time.Now().Unix())
	affected, err := this.Model.DB.XORM.Id(id).Update(userMsg)

	return int(affected), err
}
