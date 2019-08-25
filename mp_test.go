package wxpay

import (
	"encoding/json"
	"log"
	"testing"
)

func TestMP_Prepay(t *testing.T) {
	mp := NewMP("aaaa", "cccc", "dddd")

	output, res, err := mp.Prepay(&JsapiIn{
		OpenID:     "---the open id---",
		Body:       "this is a test",
		OutTradeNo: "1234567",
		TotalFee:   1,
		IP:         "127.0.0.1",
		NotifyURL:  "https://example.com/notify",
	})

	log.Println("Result: ", res)

	if err != nil {
		log.Println("Error: ", err)
	} else {
		str, _ := json.Marshal(output)
		log.Println(string(str))
	}
}

func TestMP_Refund(t *testing.T) {
	mp := NewMP("aaaa", "cccc", "dddd")
	if err := mp.SetCertPath("the file path to apiclient_cert.p12"); err != nil {
		log.Fatalln(err)
	}

	res, err := mp.Refund(&RefundIn{
		OutTradeNo:  "1234567",
		OutRefundNo: "1234567-2",
		TotalFee:    1,
		RefundFee:   1,
	})

	log.Println("Result: ", res)

	if err != nil {
		log.Println(err)
	}
}

func TestMP_QueryOrder(t *testing.T) {
	mp := NewMP("aaaa", "cccc", "dddd")

	status, res, err := mp.QueryOrder("1234567")

	log.Println("Result: ", res)

	if err != nil {
		log.Println(err)
	} else {
		log.Println(status)
	}
}
