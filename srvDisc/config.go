package srvDisc

import (
	"encoding/xml"
	"math/big"

	jsoniter "github.com/json-iterator/go"
	"github.com/xukgo/gsvcCtrlCenter/constDefine"
	"github.com/xukgo/gsvcCtrlCenter/sm2"

	"github.com/xukgo/gsvcCtrlCenter/models"
)

/*
go mod edit -replace github.com/coreos/bbolt@v1.3.4=go.etcd.io/bbolt@v1.3.4
go mod edit -replace google.golang.org/grpc@v1.29.1=google.golang.org/grpc@v1.26.0
*/
type ConfRoot struct {
	XMLName   xml.Name
	Endpoints []string `xml:"Endpoints>Addr"` //etcd服务器地址, 172.16.0.212:2379
}

type LicRegisterInfo struct {
	LicEncData    string                    `json:"licEnc"`
	LicenseConfig *models.LicenseConfig     `json:"lic"`
	MInfo         *models.MachineUniqueInfo `json:"info"`
	Timestamp     int64                     `json:"timestamp"`
}

func (this *LicRegisterInfo) EncryptJson() ([]byte, error) {
	gson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(this)
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.GetSm2P256V1()
	pub.X, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_X, 16)
	pub.Y, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_Y, 16)

	cipherText, err := sm2.Encrypt(pub, gson, sm2.C1C3C2)
	return cipherText, err
}
func (this *LicRegisterInfo) DecryptJson(data []byte) error {
	priv := new(sm2.PrivateKey)
	priv.Curve = sm2.GetSm2P256V1()
	priv.D, _ = new(big.Int).SetString(constDefine.LIC_SM2_PRIV_D, 16)

	plainText, err := sm2.Decrypt(priv, data, sm2.C1C3C2)
	if err != nil {
		return err
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(plainText, this)
}

type SubscribeSrvConf struct {
	Namespace string `xml:"Namespace"`
	Name      string `xml:"Name"`
	Version   string `xml:"Version"`
}

type SubscribeConf struct {
	Services []SubscribeSrvConf `xml:"Service"`
}

func (this *ConfRoot) FillWithXml(data []byte) error {
	err := xml.Unmarshal(data, this)
	if err != nil {
		return err
	}
	return err
}
