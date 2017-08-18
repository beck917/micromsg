package controllers

import (
	"application/entities"
	"application/libraries/helpers"
	"application/libraries/logger"
	"application/models"
	"encoding/json"
	"time"

	"github.com/beck917/pillX/pillx"
)

type SendJson struct {
	SendId int    `json:"send_id"`
	RecvId int    `json:"recv_id"`
	Msg    string `json:"msg"`
}

type PushJson struct {
	SendId  int      `json:"send_id"`
	RecvId  int      `json:"recv_id"`
	Msg     string   `json:"msg"`
	Contact *Contact `json:"contact"`
}

//发送消息
func SendHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &SendJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		retjson, _ := retJson("send", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "send").Error(jsonErr)
		return
	}

	uid := helpers.GlobalClientIdBindUid[req.Header.ClientId]
	if uid == 0 {
		//记录错误
		retjson, _ := retJson("send", "用户没有登录", 90002, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "send").Error("")
		return
	}

	//判断A是否是B的联系人
	userContactsModel := models.NewUserContacts()
	has, _ := userContactsModel.GetContactByUidCid(jsonData.RecvId, uid)
	userModel := models.NewUser()
	if has == false {
		userModel.GetUserById(uid)
		//添加联系人
		userContactsEntity := userContactsModel.UserContactsEntity
		userContactsEntity.Id = 0
		userContactsEntity.Uid = jsonData.RecvId
		userContactsEntity.Cid = uid
		userContactsEntity.Cname = userModel.UserEntity.Name
		userContactsEntity.Unread = 0
		userContactsEntity.CreatedTime = entities.JsonTime(time.Now())
		userContactsEntity.UpdatedTime = entities.JsonTime(time.Now())
		userContactsModel.Insert(userContactsEntity)
	}

	if helpers.GlobalUidBindClientId[jsonData.RecvId] != nil {
		jsonPushData := &PushJson{}
		jsonPushData.SendId = jsonData.SendId
		jsonPushData.RecvId = jsonData.RecvId
		jsonPushData.Msg = jsonData.Msg

		if has == false {
			contact := &Contact{}
			contact.Cid = uid
			contact.Cname = userModel.UserEntity.Name
			contact.Unread = 1
			jsonPushData.Contact = contact
		}

		retjson, _ := retJson("pushmsg", "推送信息", 1, jsonPushData)
		returnMsg(helpers.GlobalUidBindClientId[jsonData.RecvId].ClientId, retjson)
	}

	userMsgModel := models.NewUserMsg()
	userMsgModel.UserMsgEntity.Id = 0
	userMsgModel.UserMsgEntity.Msg = jsonData.Msg
	userMsgModel.UserMsgEntity.SendUid = uid
	userMsgModel.UserMsgEntity.RecvUid = jsonData.RecvId
	userMsgModel.UserMsgEntity.CreatedTime = entities.JsonTime(time.Now())
	userMsgModel.UserMsgEntity.UpdatedTime = entities.JsonTime(time.Now())
	userMsgModel.UserMsgEntity.DeletedTime = 0
	userMsgModel.Insert(userMsgModel.UserMsgEntity)

	//记录条数
	userContactsModel.UpdateUnread(jsonData.RecvId, uid, 1)

	retjson, _ := retJson("send", "发送成功", 1, nil)
	returnMsg(req.Header.ClientId, retjson)
}

type OpenJson struct {
	OpenId int `json:"open_id"`
}

type RetOpenData struct {
	MsgList []*entities.UserMsg `json:"msg_list"`
}

//打开消息界面
func OpenHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &OpenJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		retjson, _ := retJson("open", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "open").Error(jsonErr)
		return
	}

	uid := helpers.GlobalClientIdBindUid[req.Header.ClientId]
	if uid == 0 {
		//记录错误
		retjson, _ := retJson("open", "用户没有登录", 90002, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "open").Error("")
		return
	}

	//清除未读记录
	userContactsModel := models.NewUserContacts()
	userContactsModel.UpdateUnread(uid, jsonData.OpenId, 0)

	userMsgModel := models.NewUserMsg()
	msgList, _ := userMsgModel.GetMsgListByUid(uid, jsonData.OpenId, 0, 20)

	retOpenData := &RetOpenData{}
	retOpenData.MsgList = msgList

	retjson, _ := retJson("open", "打开成功", 1, retOpenData)
	returnMsg(req.Header.ClientId, retjson)
}

type DeleteMsgJson struct {
	Id int `json:"id"`
}

//删除消息
func DeleteMsgHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &DeleteMsgJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		retjson, _ := retJson("delete_msg", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "delete_msg").Error(jsonErr)
		return
	}

	uid := helpers.GlobalClientIdBindUid[req.Header.ClientId]
	if uid == 0 {
		//记录错误
		retjson, _ := retJson("delete_msg", "用户没有登录", 90002, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "delete_msg").Error("")
		return
	}

	userMsgModel := models.NewUserMsg()
	_, err := userMsgModel.DeleteMsgById(jsonData.Id)
	if err != nil {
		logger.WithField("controller", "delete_msg").Error(err)
		retjson, _ := retJson("delete_msg", "删除失败", 50001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	retjson, _ := retJson("delete_msg", "删除消息成功", 1, nil)
	returnMsg(req.Header.ClientId, retjson)
}
