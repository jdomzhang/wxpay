package wxpay

import (
	"fmt"
	"time"
)

type MP struct {
	client *Client
	secret string
}

func NewMP(appID, secret, mchID, mchKey string) *MP {
	return &MP{
		client: NewClient(NewAccount(appID, mchID, mchKey, false)),
		secret: secret,
	}
}

const (
	TRADE_TYPE_JSAPI = "JSAPI"
)

type InputMP struct {
	OpenID     string
	Body       string
	OutTradeNo string
	TotalFee   int64
	IP         string
	NotifyURL  string
}

// appId,nonceStr,package,signType,timeStamp

type OutputMP struct {
	AppID     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	TimeStamp string `json:"timeStamp"`
	PaySign   string `json:"paySign"`
}

func (c *MP) Prepay(input *InputMP) (*OutputMP, error) {
	var params = make(Params)
	params.SetString("body", input.Body).
		SetString("out_trade_no", input.OutTradeNo).
		SetInt64("total_fee", input.TotalFee).
		SetString("spbill_create_ip", input.IP).
		SetString("notify_url", input.NotifyURL).
		SetString("openid", input.OpenID).
		SetString("trade_type", TRADE_TYPE_JSAPI)

	result, err := c.client.UnifiedOrder(params)

	if err != nil {
		return nil, err
	}

	if ok := result.GetString("return_code") == "SUCCESS" &&
		result.GetString("return_code") == "SUCCESS"; !ok {
		return nil, fmt.Errorf("%v", result)
	}

	// create output MP
	output := OutputMP{
		AppID:     c.client.account.appID,
		NonceStr:  nonceStr(),
		Package:   fmt.Sprintf("prepay_id=%s", result.GetString("prepay_id")),
		SignType:  c.client.signType,
		TimeStamp: fmt.Sprint(time.Now().Unix()),
	}

	output.PaySign = c.client.Sign(output.toMap())

	return &output, nil
}

func (obj *OutputMP) toMap() Params {
	params := make(Params)
	params.SetString("appId", obj.AppID)
	params.SetString("nonceStr", obj.NonceStr)
	params.SetString("package", obj.Package)
	params.SetString("signType", obj.SignType)
	params.SetString("timeStamp", obj.TimeStamp)

	return params
}
