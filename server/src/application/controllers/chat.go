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
	"github.com/bitly/go-simplejson"
)

func OnClientCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	delete(helpers.GlobalUidBindClientId, helpers.GlobalClientIdBindUid[req.Header.ClientId])
	delete(helpers.GlobalClientIdBindUid, req.Header.ClientId)
	logger.Info("客户端断开")
}

func OnClientCloseWorkerHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//清除数据
	delete(helpers.GlobalUidBindClientId, helpers.GlobalClientIdBindUid[req.Header.ClientId])
	delete(helpers.GlobalClientIdBindUid, req.Header.ClientId)
	logger.Info("客户端worker断开")
}

func OnMessageHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)
	jsonObj, jsonErr := simplejson.NewJson(req.Content)
	if jsonErr != nil {
		//记录错误
		logger.WithField("content", req.Content).Error("json error")
		return
	}

	method, err := jsonObj.Get("method").String()

	if err != nil {
		logger.WithField("header", req.Header).Error("method error")
		return
	}

	logger.WithField("header", req.Header).Info("on message")

	cmd := helpers.Crc16([]byte(method))
	req.SetCmd(cmd)
}

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

type RetJson struct {
	ReplyMethod string      `json:"replymethod"`
	Result      int         `json:"result"`
	Msg         string      `json:"msg"`
	Timestamp   int64       `json:"timestamp"`
	Data        interface{} `json:"data"`
}

func retJson(replymethod string, msg string, result int, data interface{}) ([]byte, error) {
	ret := RetJson{}
	ret.Result = result
	ret.ReplyMethod = replymethod
	ret.Msg = msg
	ret.Timestamp = time.Now().Unix()
	ret.Data = data

	jsonret, err := json.Marshal(ret)
	return jsonret, err
}

func returnMsg(clientId uint64, content []byte) {
	gatewayProtocol := &pillx.GateWayProtocol{}
	header := &pillx.GatewayHeader{}
	gatewayProtocol.Header = header

	header.ClientId = clientId
	header.Cmd = 0 //utils.Crc16([]byte(method))
	header.Error = 0
	header.Mark = pillx.PROTO_HEADER_FIRSTCHAR
	header.Version = pillx.GATEWAY_VERSION
	header.Sid = 0
	header.Size = uint16(len(content))
	gatewayProtocol.Content = content

	pillx.SendAllGateWay(pillx.Gateways, gatewayProtocol)
}
