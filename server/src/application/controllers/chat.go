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

	"github.com/beck917/pillX/pillx"
	"github.com/bitly/go-simplejson"
)

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
	UserName string `json:"username"`
	PassWord string `json:"password"`
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
		retjson, _ := retJson("login", "数据格式错误", 10003, nil)
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

	var contacts []*Contact
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

type SendJson struct {
	SendId int    `json:"send_id"`
	RecvId int    `json:"recv_id"`
	Msg    string `json:"msg"`
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

	if helpers.GlobalUidBindClientId[jsonData.RecvId] != nil {
		jsonPushData := &SendJson{}
		jsonPushData.SendId = jsonData.SendId
		jsonPushData.RecvId = jsonData.RecvId
		jsonPushData.Msg = jsonData.Msg

		retjson, _ := retJson("pushmsg", "推送信息", 1, jsonPushData)
		returnMsg(helpers.GlobalUidBindClientId[jsonData.RecvId].ClientId, retjson)
	}

	//判断A是否是B的联系人
	userContactsModel := models.NewUserContacts()
	has, _ := userContactsModel.GetContactByUidCid(jsonData.RecvId, uid)
	if has == false {
		userModel := models.NewUser()
		userModel.GetUserById(jsonData.RecvId)
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
	Contact *entities.User `json:"contact"`
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
		retjson, _ := retJson("add", "用户名不存在", 10001, nil)
		returnMsg(req.Header.ClientId, retjson)
		return
	}

	userContactsModel := models.NewUserContacts()
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
	retAddData.Contact = userModel.UserEntity

	retjson, _ := retJson("add", "添加联系人成功", 1, retAddData)
	returnMsg(req.Header.ClientId, retjson)
}

type DeleteJson struct {
	DeleteId int `json:"delete_id"`
}

//删除联系人
func DeleteHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

//删除消息
func DeleteMsgHandler(client *pillx.Response, protocol pillx.IProtocol) {

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
