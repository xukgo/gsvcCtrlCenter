package srvDisc

import (
	"math/big"

	jsoniter "github.com/json-iterator/go"
	"github.com/xukgo/gsvcCtrlCenter/constDefine"
	"github.com/xukgo/gsvcCtrlCenter/models"
	"github.com/xukgo/gsvcCtrlCenter/sm2"
)

type LicResultInfo struct {
	Code          int                   `json:"code"`
	Description   string                `json:"description"`
	LicenseConfig *models.LicenseConfig `json:"lic"`
	Timestamp     int64                 `json:"timestamp"`
}

func (this *LicResultInfo) EncryptJson() ([]byte, error) {
	gson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(this)
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.GetSm2P256V1()
	pub.X, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_X, 16)
	pub.Y, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_Y, 16)

	cipherText, err := sm2.Encrypt(pub, gson, sm2.C1C3C2)
	return cipherText, err
}

func (this *LicResultInfo) DecryptJson(data []byte) error {
	priv := new(sm2.PrivateKey)
	priv.Curve = sm2.GetSm2P256V1()
	priv.D, _ = new(big.Int).SetString(constDefine.LIC_SM2_PRIV_D, 16)

	plainText, err := sm2.Decrypt(priv, data, sm2.C1C3C2)
	if err != nil {
		return err
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(plainText, this)
}
