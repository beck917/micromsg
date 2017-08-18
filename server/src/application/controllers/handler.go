package controllers

import (
	"application/libraries/helpers"
	"application/libraries/logger"
	"encoding/json"
	"time"

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
