package apns

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/blackbeans/log4go"
)

//--------------------payload

type Alert struct {
	Body         string        `json:"body,omitempty"`
	ActionLocKey string        `json:"action-loc-key,omitempty"`
	LocKey       string        `json:"loc-key,omitempty"`
	LocArgs      []interface{} `json:"loc-args,omitempty"`
}

type Aps struct {
	Alert string `json:"alert,omitempty"`
	Badge int    `json:"badge,omitempty"` //显示气泡数
	Sound string `json:"sound,omitempty"` //控制push弹出的声音
}

type PayLoad struct {
	aps       Aps
	extParams map[string]interface{} //扩充字段
}

func NewSimplePayLoad(sound string, badge int, body string) *PayLoad {
	aps := Aps{Alert: body, Sound: sound, Badge: badge}
	return &PayLoad{aps: aps, extParams: make(map[string]interface{})}
}

func NewSimplePayLoadWithAps(aps Aps) *PayLoad {
	return &PayLoad{aps: aps, extParams: make(map[string]interface{})}
}

func NewPayLoad(sound string, badge int, alert Alert) *PayLoad {
	data, err := json.Marshal(alert)
	if nil != err {
		log.Error("NEWPAYLOAD|FAIL|ERROR|%s\n", err)
		return nil
	}
	aps := Aps{Alert: string(data), Sound: sound, Badge: badge}
	return &PayLoad{aps: aps, extParams: make(map[string]interface{})}
}

func (self *PayLoad) AddExtParam(key string, val interface{}) *PayLoad {
	self.extParams[key] = val
	return self
}

func (self *PayLoad) Marshal() []byte {

	encoddata := make(map[string]interface{}, 2)
	encoddata["aps"] = self.aps
	for k, v := range self.extParams {
		encoddata[k] = v
	}

	data, err := json.Marshal(encoddata)
	if nil != err {
		log.Error("PAYLOAD|ENCODE|FAIL|%s", err)
		return nil
	}

	return data
}

func WrapPayLoad(payload *PayLoad) (*Item, error) {
	payloadJson := payload.Marshal()
	if nil == payloadJson || len(payloadJson) > 256 {
		return nil, errors.New(fmt.Sprintf("WRAPPAYLOAD|FAIL|%s|len:%d", payloadJson, len(payloadJson)))
	}
	return &Item{id: PAY_LOAD, length: uint16(len(payloadJson)), data: payloadJson}, nil
}

func WrapDeviceToken(token string) (*Item, error) {
	decodeToken, err := hex.DecodeString(token)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("WRAPTOKE|FAIL|INVALID TOKEN|%s|%s", token, err.Error()))

	}
	return &Item{id: DEVICE_TOKEN, length: uint16(len(decodeToken)), data: decodeToken}, nil
}

func WrapNotifyIdentifier(id uint32) *Item {
	return &Item{id: NOTIFY_IDENTIYFIER, length: 4, data: id}
}

func WrapExpirationDate(expirateDate uint32) *Item {
	return &Item{id: EXPIRATED_DATE, length: 4, data: expirateDate}
}

func WrapPriority(priority byte) *Item {
	return &Item{id: PRIORITY, length: 1, data: priority}
}
