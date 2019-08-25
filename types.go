package wxpay

import (
	"fmt"
)

const (
	SUCCESS = "SUCCESS"
)

// JsapiIn 输入参数，支付结果通知文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7&index=8
type JsapiIn struct {
	OpenID     string
	Body       string
	OutTradeNo string
	TotalFee   int64
	IP         string
	NotifyURL  string
}

func (obj *JsapiIn) toMap() Params {
	return make(Params).
		SetString("body", obj.Body).
		SetString("out_trade_no", obj.OutTradeNo).
		SetInt64("total_fee", obj.TotalFee).
		SetString("spbill_create_ip", obj.IP).
		SetString("notify_url", obj.NotifyURL).
		SetString("openid", obj.OpenID).
		SetString("trade_type", TRADE_TYPE_JSAPI)
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

type RefundIn struct {
	OutTradeNo  string `json:"outTradeNo"`
	OutRefundNo string `json:"outRefundNo"`
	TotalFee    int64  `json:"totalFee"`
	RefundFee   int64  `json:"refundFee"`
}

func (obj *RefundIn) toMap() Params {
	return make(Params).
		SetString("out_trade_no", obj.OutTradeNo).
		SetString("out_refund_no", obj.OutRefundNo).
		SetInt64("total_fee", obj.TotalFee).
		SetInt64("refund_fee", obj.RefundFee)
}

func (obj Params) toError() error {

	if obj.GetString("return_code") != SUCCESS {
		return fmt.Errorf(obj.GetString("return_msg"))
	}

	if obj.GetString("result_code") != SUCCESS {
		return fmt.Errorf(obj.GetString("err_code_des"))
	}

	return nil
}
