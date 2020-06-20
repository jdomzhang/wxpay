package wxpay

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_UnifiedOrder(t *testing.T) {
	client := NewClient(NewAccount("xxxxx", "xxx", "xxxxx", false))
	params := make(Params)
	params.SetString("body", "test").
		SetString("out_trade_no", "58867657575757").
		SetInt64("total_fee", 1).
		SetString("spbill_create_ip", "127.0.0.1").
		SetString("notify_url", "http://notify.jdomzhang.com/notify").
		SetString("trade_type", "APP")
	t.Log(client.UnifiedOrder(params))
}

func TestSubMch_UnifiedOrder(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	input := &JsapiIn{
		SubOpenID:  "------",
		Body:       "子订单",
		OutTradeNo: fmt.Sprint(time.Now().Unix()),
		TotalFee:   1,
		IP:         "127.0.0.1",
		// TODO: setup notify url
		NotifyURL: "http://unknown.com",
	}

	out, _, err := subMP.Prepay(input)

	// output json
	bs, _ := json.Marshal(out)
	t.Log(string(bs))

	// assert
	assert.NoError(t, err)
}

func TestSubMch_UnifiedOrder_ProfitSharing(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	input := &JsapiIn{
		SubOpenID:  "------",
		Body:       "子订单",
		OutTradeNo: fmt.Sprint(time.Now().Unix()),
		TotalFee:   10,
		IP:         "127.0.0.1",
		// TODO: setup notify url
		NotifyURL:     "http://unknown.com",
		ProfitSharing: true,
	}

	out, _, err := subMP.Prepay(input)

	// output json
	bs, _ := json.Marshal(out)
	t.Log(string(bs))

	// assert
	assert.NoError(t, err)
}

func TestSubMch_QueryOrder(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	status, param, err := subMP.QueryOrder("1592635119")

	// assert
	assert.NoError(t, err)

	// output json
	t.Log("Status", status)
	t.Log("Param", param)
}

func TestSubMch_Refund(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	subMP.SetCertBase64("-----------------------------")
	input := &RefundIn{
		AppID:       appID,
		MchID:       mchID,
		SubAppID:    subAppID,
		SubMchID:    subMchID,
		OutTradeNo:  "1592647167",
		OutRefundNo: "1592647167-R",
		TotalFee:    10,
		RefundFee:   10,
	}
	param, err := subMP.Refund(input)

	// assert
	assert.NoError(t, err)

	// output json
	t.Log("Param", param)
}

func TestSubMch_ProfitSharing(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	subMP.SetCertBase64("-----------------------------")
	param, err := subMP.ProfitSharingAddReceiver(NewDefaultShareReceiver(mchID, "------------------"))

	// assert
	assert.NoError(t, err)
	t.Log("Param", param)
	if err != nil {
		return
	}

	// 单次分账
	transcationID := "4200000590202006201100870586"
	input2 := &SharingInfo{
		TransactionID: transcationID,
		OutOrderNo:    transcationID + "-M",
		Receivers:     []AmountReceiver{*NewDefaultAmountReceiver(mchID, 3, "服务商分成")},
	}
	param, err = subMP.ProfitSharing(input2)

	// assert
	assert.NoError(t, err)
	t.Log("Param", param)
	if err != nil {
		return
	}
}

func TestSubMch_ProfitSharing_Multi(t *testing.T) {
	// 公众号
	appID := "-----"
	mchID := "-----"
	mchKey := "-----"
	subAppID := "-----"
	subMchID := "-----"

	subMP := NewSubMP(appID, mchID, mchKey, subAppID, subMchID)
	subMP.SetCertBase64("-----------------------------")
	param, err := subMP.ProfitSharingAddReceiver(NewDefaultShareReceiver(mchID, "------------------"))

	// assert
	assert.NoError(t, err)
	t.Log("Param", param)
	if err != nil {
		return
	}

	// 开始分账
	input2 := &SharingInfo{
		TransactionID: "4200000591202006209266376722",
		OutOrderNo:    fmt.Sprint(time.Now().Unix()),
		Receivers:     []AmountReceiver{*NewDefaultAmountReceiver(mchID, 3, "服务商分成")},
	}
	param, err = subMP.ProfitSharingMulti(input2)

	// assert
	assert.NoError(t, err)
	t.Log("Param", param)
	if err != nil {
		return
	}

	// 结束分账
	input3 := &FinishSharing{
		TransactionID: "4200000591202006209266376722",
		OutOrderNo:    fmt.Sprint(time.Now().Unix()),
		Description:   "分账结束",
	}
	param, err = subMP.ProfitSharingFinish(input3)

	// assert
	assert.NoError(t, err)

	// output json
	t.Log("Param", param)
}
