package controllers

import (
	"application/entities"
	"application/libraries/constant"
	"application/libraries/helpers"
	"application/libraries/logger"
	"application/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/beck917/pillX/pillx"
)

type LoginJson struct {
	UserName string `json:"username" valid:"length(5|30)"`
	PassWord string `json:"password" valid:"length(6|30)"`
}

type RetLoginData struct {
	Contacts []*Contact `json:"contacts"`
	Uid      int        `json:"uid"`
}

type Contact struct {
	Cid    int    `json:"cid"`
	Cname  string `json:"cname"`
	Unread int    `json:"unread"`
}

func LoginHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &LoginJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		logger.WithField("controller", "login").Error(jsonErr)
		retjson, _ := retJson("login", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	_, err := govalidator.ValidateStruct(jsonData)
	if err != nil {
		logger.WithField("controller", "login").Error(err)
		retjson, _ := retJson("login", err.Error(), 90004, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	userModel := models.NewUser()
	has, _ := userModel.GetUserByName(jsonData.UserName)

	if has == false {
		//返回错误
		retjson, _ := retJson("login", "用户名不存在", 10001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	if userModel.UserEntity.Password != helpers.Sha1Str(fmt.Sprintf("%s%s", jsonData.PassWord, constant.PASSWORD_SALT)) {
		//返回错误
		retjson, _ := retJson("login", "密码错误", 10002, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	helpers.GlobalUidBindClientId[userModel.UserEntity.Id] = helpers.NewUserData()
	helpers.GlobalUidBindClientId[userModel.UserEntity.Id].ClientId = req.Header.ClientId
	helpers.GlobalUidBindClientId[userModel.UserEntity.Id].Name = userModel.UserEntity.Name
	helpers.GlobalClientIdBindUid[req.Header.ClientId] = userModel.UserEntity.Id

	//获取联系人列表
	userContactsModel := models.NewUserContacts()
	contactsList, _ := userContactsModel.GetContactsListByUid(userModel.UserEntity.Id)

	var contacts []*Contact = make([]*Contact, 0)
	for _, data := range contactsList {
		contact := &Contact{}
		contact.Cid = data.Cid
		contact.Cname = data.Cname
		contact.Unread = data.Unread
		contacts = append(contacts, contact)
	}

	retLoginData := &RetLoginData{}
	retLoginData.Contacts = contacts
	retLoginData.Uid = userModel.UserEntity.Id

	retjson, _ := retJson("login", "登陆成功", 1, retLoginData)
	returnMsg(req.Header.ClientId, retjson)
}

func RegisterHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	jsonData := &LoginJson{}
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		logger.WithField("controller", "register").Error(jsonErr)
		retjson, _ := retJson("register", "数据格式错误", 90001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	_, err := govalidator.ValidateStruct(jsonData)
	if err != nil {
		logger.WithField("controller", "register").Error(err)
		retjson, _ := retJson("register", err.Error(), 90004, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	userModel := models.NewUser()
	has, _ := userModel.GetUserByName(jsonData.UserName)

	if has == true {
		//返回错误
		retjson, _ := retJson("register", "用户名已存在", 10001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	userModel.UserEntity.Id = 0
	userModel.UserEntity.Name = jsonData.UserName
	userModel.UserEntity.Password = helpers.Sha1Str(jsonData.PassWord + constant.PASSWORD_SALT)
	userModel.UserEntity.CreatedTime = entities.JsonTime(time.Now())
	userModel.UserEntity.UpdatedTime = entities.JsonTime(time.Now())
	userModel.Insert(userModel.UserEntity)

	retjson, _ := retJson("register", "注册成功", 1, nil)
	returnMsg(req.Header.ClientId, retjson)
}
