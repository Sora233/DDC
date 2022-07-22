package model

type WithGenericErrInfo interface {
	GetGenericErrInfo() *GenericErrInfo
}

type GenericErrInfo struct {
	// 错误码
	ErrorCode int `json:"error_code"`
	// 错误信息
	ErrorMsg string `json:"error_msg,omitempty"`
	// 本次请求requestid
	RequestID string `json:"request_id"`
}

func (gei *GenericErrInfo) GetGenericErrInfo() *GenericErrInfo {
	return gei
}

func (gei *GenericErrInfo) Error(code int) *GenericErrInfo {
	gei.ErrorCode = code
	gei.ErrorMsg = ErrMsg(code)
	return gei
}

func (gei *GenericErrInfo) ErrPrint(code int, args ...interface{}) *GenericErrInfo {
	gei.ErrorCode = code
	gei.ErrorMsg = ErrPrint(code, args...)
	return gei
}

func (gei *GenericErrInfo) ErrPrintf(code int, format string, args ...interface{}) *GenericErrInfo {
	gei.ErrorCode = code
	gei.ErrorMsg = ErrPrintf(code, format, args...)
	return gei
}

type PingReq struct {
}

type PingResp struct {
	GenericErrInfo
	// Pong
	Pong string `json:"pong"`
}

type GetSearchDanmuReq struct {
	Uid int64 `form:"uid" json:"uid" binding:"required" example:"97505"`
}

type Danmu struct {
	// 主播uid
	VTBUid int64 `json:"vtb_uid"`
	// 主播名称
	VTBName string `json:"vtb_name"`
	// 发送时间
	Time string `json:"time"`
	// 弹幕内容
	Text string `json:"text"`
}

type GetSearchDanmuResp struct {
	GenericErrInfo
	// 用户名
	Username string `json:"username"`
	// 用户发送的弹幕
	Damus []*Danmu `json:"damus"`
}

type GetSearchMoneyReq struct {
	Uid int64 `form:"uid" json:"uid" binding:"required" example:"97505"`
}

type SuperChat struct {
	// 主播uid
	VTBUid int64 `json:"vtb_uid"`
	// 主播名称
	VTBName string `json:"vtb_name"`
	// 发送时间
	Time string `json:"time"`
	// superchat留言
	Text string `json:"text"`
	// superchat价格
	Price int `json:"price"`
}

type Gift struct {
	// 主播uid
	VTBUid int64 `json:"vtb_uid"`
	// 主播名称
	VTBName string `json:"vtb_name"`
	// 发送时间
	Time string `json:"time"`
	// 礼物名字
	GiftName string `json:"gift_name"`
	// 礼物数量
	GiftNum int `json:"gift_num"`
	// 礼物价格
	GiftPrice float64 `json:"gift_price"`
}

type GetSearchMoneyResp struct {
	GenericErrInfo
	// 用户名
	Username string `json:"username"`
	// 用户消费的superchat"
	SuperChats []*SuperChat `json:"super_chats"`
	// 用户消费的礼物
	Gifts []*Gift `json:"gifts"`
}
