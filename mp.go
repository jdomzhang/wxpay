package wxpay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const (
	TRADE_TYPE_JSAPI = "JSAPI"
)

// MP is mini-program
type MP struct {
	client *Client
}

// SetCertPath 设置操作证书，请查看API证书说明 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
func (c *MP) SetCertPath(certPath string) error {
	return c.client.account.SetCertData(certPath)
}

// SetCertBase64 设置操作证书 base64编码
func (c *MP) SetCertBase64(certDataBase64 string) error {
	certData, err := base64.StdEncoding.DecodeString(certDataBase64)
	if err != nil {
		return err
	}
	c.client.account.certData = certData

	return nil
}

// NewMP 小程序支付
func NewMP(appID, mchID, mchKey string) *MP {
	return &MP{
		client: NewClient(NewAccount(appID, mchID, mchKey, false)),
	}
}

// NewSubMP 小程序支付-使用服务商提供的支付商户
func NewSubMP(appID, mchID, mchKey, subAppID, subMchID string) *MP {
	return &MP{
		client: NewClient(NewSubAccount(appID, mchID, mchKey, subAppID, subMchID)),
	}
}

// Prepay 微信统一下单，接口文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (c *MP) Prepay(input *JsapiIn) (*JsapiOut, Params, error) {
	result, err := realError(c.client.UnifiedOrder(input.toMap()))

	if err != nil {
		return nil, result, err
	}

	// create output MP
	appID := c.client.account.appID
	if c.client.account.isSubMode {
		appID = c.client.account.subAppID
	}
	output := JsapiOut{
		AppID:     appID,
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
	return realError(c.client.Refund(input.toMap()))
}

// QueryOrder 查询订单，接口文档 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (c *MP) QueryOrder(outTradeNo string) (string, Params, error) {
	result, err := realError(c.client.OrderQuery(make(Params).SetString("out_trade_no", outTradeNo)))

	if err != nil {
		return "", result, err
	}

	status := result.GetString("trade_state")
	return status, result, nil
}

// ProfitSharingAddReceiver 添加分账接收方 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_3&index=4
func (c *MP) ProfitSharingAddReceiver(receiver *ShareReceiver) (Params, error) {
	old := c.client.account.isProfitSharing
	c.client.account.isProfitSharing = true
	defer func() { c.client.account.isProfitSharing = old }()

	// construct param
	value := make(Params)
	// set receivers with json
	bs, _ := json.Marshal(receiver)
	value.SetString("receiver", string(bs))

	return realError(c.client.ProfitSharingAddReceiver(value))
}

// ProfitSharing 请求单次分账 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_1&index=1
func (c *MP) ProfitSharing(input *SharingInfo) (Params, error) {
	old := c.client.account.isProfitSharing
	c.client.account.isProfitSharing = true
	defer func() { c.client.account.isProfitSharing = old }()

	return realError(c.client.ProfitSharing(input.toMap()))
}

// ProfitSharingMulti 请求多次分账 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_6&index=2
func (c *MP) ProfitSharingMulti(input *SharingInfo) (Params, error) {
	old := c.client.account.isProfitSharing
	c.client.account.isProfitSharing = true
	defer func() { c.client.account.isProfitSharing = old }()

	return realError(c.client.ProfitSharingMulti(input.toMap()))
}

// ProfitSharingFinish 完结分账 https://pay.weixin.qq.com/wiki/doc/api/allocation_sl.php?chapter=25_5&index=6
func (c *MP) ProfitSharingFinish(input *FinishSharing) (Params, error) {
	old := c.client.account.isProfitSharing
	c.client.account.isProfitSharing = true
	defer func() { c.client.account.isProfitSharing = old }()

	return realError(c.client.ProfitSharingFinish(input.toMap()))
}

func realError(result Params, err error) (Params, error) {
	if err != nil {
		return result, err
	}

	if err := result.toError(); err != nil {
		return result, err
	}

	return result, nil
}
