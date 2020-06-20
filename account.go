package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID           string
	mchID           string
	apiKey          string
	subAppID        string
	subMchID        string
	certData        []byte
	isSandbox       bool
	isSubMode       bool
	isProfitSharing bool
}

// 创建微信支付账号
func NewAccount(appID string, mchID string, apiKey string, isSanbox bool) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		isSandbox: isSanbox,
	}
}

// 创建微信支付子账号-服务商模式
func NewSubAccount(appID, mchID, apiKey, subAppID, subMchID string) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		subAppID:  subAppID,
		subMchID:  subMchID,
		isSubMode: true,
	}
}

// 设置证书
func (a *Account) SetCertData(certPath string) error {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("读取证书失败")
		return err
	}
	a.certData = certData
	return nil
}
