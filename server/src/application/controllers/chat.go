package controllers

import (
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
		logger.WithField("header", req.Header).Error("json error")
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
		contacts = append(contacts, contact)
	}

	retLoginData := &RetLoginData{}
	retLoginData.Contacts = contacts

	retjson, _ := retJson("login", "登陆成功", 1, retLoginData)
	returnMsg(req.Header.ClientId, retjson)
}

type SendJson struct {
	SendId int    `json:"send_id"`
	RecvId int    `json:"recv_id"`
	Msg    string `json:"msg"`
}

func SendHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	var jsonData SendJson
	jsonErr := json.Unmarshal(req.Content, jsonData)
	if jsonErr != nil {
		//记录错误
		logger.WithField("controller", "send").Error("json error")
		return
	}
}

func OpenHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func AddHandler(client *pillx.Response, protocol pillx.IProtocol) {

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
