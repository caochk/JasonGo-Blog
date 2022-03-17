package respUtils

import "encoding/json"

const (
	SUCCESS_CODE   = 200
	ERROR_CODE     = -1
	TOKEN_ERR_CODE = -2
	// Redis
	REDIS_ERR = "redisErr"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r Resp) NewResp(code int, msg string) *Resp {
	return r.NewRespWithData(code, msg, "")
}

// NewRespWithData 创建带data的返回给前端的响应的body部分
func (r Resp) NewRespWithData(code int, msg string, data interface{}) *Resp {
	var rsp = &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	return rsp
}

// ToBytes 需要把对象转换成[]byte类型，以匹配http返回的类型
func (r *Resp) ToBytes() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return []byte("")
	}
	return data
}
