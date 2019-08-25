package wxpay

import (
	"fmt"
	"time"
)

const (
	TRADE_TYPE_JSAPI = "JSAPI"
)

type MP struct {
	client *Client
}

// SetCertPath 设置操作证书，请查看API证书说明 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
func (c *MP) SetCertPath(certPath string) error {
	return c.client.account.SetCertData(certPath)
}

func NewMP(appID, mchID, mchKey string) *MP {
	return &MP{
		client: NewClient(NewAccount(appID, mchID, mchKey, false)),
	}
}

// Prepay 微信统一下单，接口文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (c *MP) Prepay(input *JsapiIn) (*JsapiOut, Params, error) {
	result, err := c.client.UnifiedOrder(input.toMap())

	if err != nil {
		return nil, result, err
	}

	if err := result.toError(); err != nil {
		return nil, result, err
	}

	// create output MP
	output := JsapiOut{
		AppID:     c.client.account.appID,
		NonceStr:  nonceStr(),
		Package:   fmt.Sprintf("prepay_id=%s", result.GetString("prepay_id")),
		SignType:  c.client.signType,
		TimeStamp: fmt.Sprint(time.Now().Unix()),
	}

	output.PaySign = c.client.Sign(output.toMap())
	return &output, result, nil
}

// Refund 微信退款，必须设置证书，接口文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
func (c *MP) Refund(input *RefundIn) (Params, error) {
	result, err := c.client.Refund(input.toMap())

	if err != nil {
		return result, err
	}

	if err := result.toError(); err != nil {
		return result, err
	}

	return result, nil
}

// QueryOrder 查询订单，接口文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (c *MP) QueryOrder(outTradeNo string) (string, Params, error) {
	result, err := c.client.OrderQuery(make(Params).SetString("out_trade_no", outTradeNo))

	if err != nil {
		return "", result, err
	}

	if err := result.toError(); err != nil {
		return "", result, err
	}

	status := result.GetString("trade_state")
	return status, result, nil
}
