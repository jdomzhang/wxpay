package wxpay

import "encoding/json"

const (
	SUCCESS = "SUCCESS"
)

// JsapiIn 输入参数，支付结果通知文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7&index=8
type JsapiIn struct {
	OpenID        string
	SubOpenID     string
	Body          string
	OutTradeNo    string
	TotalFee      int64
	IP            string
	NotifyURL     string
	ProfitSharing bool
}

func (obj *JsapiIn) toMap() Params {
	value := make(Params).
		SetString("body", obj.Body).
		SetString("out_trade_no", obj.OutTradeNo).
		SetInt64("total_fee", obj.TotalFee).
		SetString("spbill_create_ip", obj.IP).
		SetString("notify_url", obj.NotifyURL).
		SetString("trade_type", TRADE_TYPE_JSAPI)

	// set sub_openid or openid
	if obj.OpenID != "" {
		value.SetString("openid", obj.OpenID)
	} else {
		value.SetString("sub_openid", obj.SubOpenID)
	}

	// set profit sharing
	// 分账 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=24_3&index=3
	if obj.ProfitSharing {
		value.SetString("profit_sharing", "Y")
	}

	return value
}

// appId,nonceStr,package,signType,timeStamp

type JsapiOut struct {
	AppID     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	TimeStamp string `json:"timeStamp"`
	PaySign   string `json:"paySign"`
}

func (obj *JsapiOut) toMap() Params {
	return make(Params).
		SetString("appId", obj.AppID).
		SetString("nonceStr", obj.NonceStr).
		SetString("package", obj.Package).
		SetString("signType", obj.SignType).
		SetString("timeStamp", obj.TimeStamp)
}

// RefundIn /
// 普通商户：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
// 服务商模式：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_sl_api.php?chapter=9_4
type RefundIn struct {
	AppID       string `json:"appID"`
	MchID       string `json:"mchID"`
	SubAppID    string `json:"subAppID"`
	SubMchID    string `json:"subMchID"`
	OutTradeNo  string `json:"outTradeNo"`
	OutRefundNo string `json:"outRefundNo"`
	TotalFee    int64  `json:"totalFee"`
	RefundFee   int64  `json:"refundFee"`
}

func (obj *RefundIn) toMap() Params {
	value := make(Params).
		SetString("appid", obj.AppID).
		SetString("mch_id", obj.MchID).
		SetString("out_trade_no", obj.OutTradeNo).
		SetString("out_refund_no", obj.OutRefundNo).
		SetInt64("total_fee", obj.TotalFee).
		SetInt64("refund_fee", obj.RefundFee)

	if obj.SubAppID != "" {
		value.SetString("sub_appid", obj.SubAppID)
	}
	if obj.SubMchID != "" {
		value.SetString("sub_mch_id", obj.SubMchID)
	}

	return value
}

/*
	分账 https://pay.weixin.qq.com/wiki/doc/api/sl.html
**/

const (
	RECEIVER_TYPE_MERCHANT_ID         = "MERCHANT_ID"
	RECEIVER_TYPE_PERSONAL_OPENID     = "PERSONAL_OPENID"
	RECEIVER_TYPE_PERSONAL_SUB_OPENID = "PERSONAL_SUB_OPENID"
)

const (
	RELATION_TYPE_SERVICE_PROVIDER = "SERVICE_PROVIDER" // 服务商
	RELATION_TYPE_STORE            = "STORE"            // 门店
	RELATION_TYPE_STAFF            = "STAFF"            // 员工
	RELATION_TYPE_STORE_OWNER      = "STORE_OWNER"      // 店主
	RELATION_TYPE_PARTNER          = "PARTNER"          // 合作伙伴
	RELATION_TYPE_HEADQUARTER      = "HEADQUARTER"      // 总部
	RELATION_TYPE_BRAND            = "BRAND"            // 品牌方
	RELATION_TYPE_DISTRIBUTOR      = "DISTRIBUTOR"      // 分销商
	RELATION_TYPE_USER             = "USER"             // 用户
	RELATION_TYPE_SUPPLIER         = "SUPPLIER"         // 供应商
	RELATION_TYPE_CUSTOM           = "CUSTOM"           // 自定义
)

// ShareReceiver 分账接收方
// https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_3&index=4
type ShareReceiver struct {
	Type         string `json:"type"`          // 分账接收方类型
	Account      string `json:"account"`       // 分账接收方帐号
	Name         string `json:"name"`          // 分账接收方全称，类型是MERCHANT_ID时，是商户全称（必传）
	Description  string `json:"description"`   // 描述
	RelationType string `json:"relation_type"` // 与分账方的关系类型
}

// AmountReceiver 分账金额接收方
// https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_6&index=2
type AmountReceiver struct {
	Type        string `json:"type"`        // 分账接收方类型
	Account     string `json:"account"`     // 分账接收方帐号
	Amount      int    `json:"amount"`      // 分账金额
	Description string `json:"description"` // 描述
}

// NewDefaultShareReceiver /
func NewDefaultShareReceiver(mchID string, name string) *ShareReceiver {
	return &ShareReceiver{
		Type:         RECEIVER_TYPE_MERCHANT_ID,
		Account:      mchID,
		Name:         name,
		RelationType: RELATION_TYPE_SERVICE_PROVIDER,
	}
}

// NewDefaultAmountReceiver /
func NewDefaultAmountReceiver(mchID string, amount int, description string) *AmountReceiver {
	return &AmountReceiver{
		Type:        RECEIVER_TYPE_MERCHANT_ID,
		Account:     mchID,
		Amount:      amount,
		Description: description,
	}
}

// SharingInfo /
// https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_6&index=2
type SharingInfo struct {
	TransactionID string           `json:"transactionID"` // 微信支付订单号
	OutOrderNo    string           `json:"outOrderNo"`    // 商户分账单号
	Receivers     []AmountReceiver `json:"receivers"`     // 分账接收方列表
}

func (obj *SharingInfo) toMap() Params {
	value := make(Params).
		SetString("transaction_id", obj.TransactionID).
		SetString("out_order_no", obj.OutOrderNo)

	bs, _ := json.Marshal(obj.Receivers)

	value.SetString("receivers", string(bs))
	return value
}

// FinishSharing /
// 完结分账 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_5&index=6
type FinishSharing struct {
	TransactionID string `json:"transactionID"` // 微信支付订单号
	OutOrderNo    string `json:"outOrderNo"`    // 商户分账单号
	Description   string `json:"description"`   // 分账完结描述
}

func (obj *FinishSharing) toMap() Params {
	value := make(Params).
		SetString("transaction_id", obj.TransactionID).
		SetString("out_order_no", obj.OutOrderNo).
		SetString("description", obj.Description)

	return value
}
