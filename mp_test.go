package wxpay

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMP_Prepay(t *testing.T) {
	mp := NewMP("aaaa", "cccc", "dddd")

	output, err := mp.Prepay(&InputMP{
		OpenID:     "eeee",
		Body:       "this is a test",
		OutTradeNo: "12345",
		TotalFee:   1,
		IP:         "127.0.0.1",
		NotifyURL:  "https://example.com/notify",
	})

	if err != nil {
		fmt.Print(err)
	} else {
		str, _ := json.Marshal(output)
		fmt.Println(string(str))
	}
}
