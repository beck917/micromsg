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

type AddJson struct {
	AddName string `json:"add_name"`
}

type RetAddData struct {
	Contact *Contact `json:"contact"`
}

//添加联系人
func AddHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &AddJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		retjson, _ := retJson("add", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "add").Error(jsonErr)
		return
	}

	uid := helpers.GlobalClientIdBindUid[req.Header.ClientId]
	if uid == 0 {
		//记录错误
		retjson, _ := retJson("add", "用户没有登录", 90002, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "add").Error("")
		return
	}

	//获取联系人数据
	userModel := models.NewUser()
	has, _ := userModel.GetUserByName(jsonData.AddName)

	if has == false {
		//返回错误
		retjson, _ := retJson("add", "用户名不存在", 40001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	//判断A是否是B的联系人
	userContactsModel := models.NewUserContacts()
	has, _ = userContactsModel.GetContactByUidCid(uid, userModel.UserEntity.Id)
	if has == true {
		//返回错误
		retjson, _ := retJson("add", "联系人已存在", 40001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	userContactsEntity := userContactsModel.UserContactsEntity
	userContactsEntity.Id = 0
	userContactsEntity.Uid = uid
	userContactsEntity.Cid = userModel.UserEntity.Id
	userContactsEntity.Cname = userModel.UserEntity.Name
	userContactsEntity.Unread = 0
	userContactsEntity.CreatedTime = entities.JsonTime(time.Now())
	userContactsEntity.UpdatedTime = entities.JsonTime(time.Now())
	userContactsModel.Insert(userContactsEntity)

	retAddData := &RetAddData{}
	contact := &Contact{}
	contact.Cid = userModel.UserEntity.Id
	contact.Cname = userModel.UserEntity.Name
	contact.Unread = 0
	retAddData.Contact = contact

	retjson, _ := retJson("add", "添加联系人成功", 1, retAddData)
	returnMsg(req.Header.ClientId, retjson)
}

type DeleteJson struct {
	DeleteId int `json:"delete_id"`
}

//删除联系人
func DeleteHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &DeleteJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		retjson, _ := retJson("add", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "add").Error(jsonErr)
		return
	}

	uid := helpers.GlobalClientIdBindUid[req.Header.ClientId]
	if uid == 0 {
		//记录错误
		retjson, _ := retJson("add", "用户没有登录", 90002, nil)
		returnMsg(req.Header.ClientId, retjson)
		logger.WithField("controller", "add").Error("")
		return
	}
	userContactsModel := models.NewUserContacts()
	userContactsModel.DeleteByUidCid(uid, jsonData.DeleteId)

	retjson, _ := retJson("delete", "删除联系人成功", 1, nil)
	returnMsg(req.Header.ClientId, retjson)
}
