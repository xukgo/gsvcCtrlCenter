package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"

	"github.com/xukgo/gsvcCtrlCenter/machineCode"
	"github.com/xukgo/gsvcCtrlCenter/sm2"
)

func main() {
	uinfo := machineCode.Info()
	gson := uinfo.ToJson()

	//priv, pub, err := sm2.GenerateRandKey(rand.Reader)
	//if err != nil {
	//	os.Exit(-1)
	//}

	priv := new(sm2.PrivateKey)
	priv.Curve = sm2.GetSm2P256V1()
	priv.D, _ = new(big.Int).SetString(constDefine.MACHINECODE_SM2_PRIV_D, 16)

	pub := new(sm2.PublicKey)
	pub.Curve = sm2.GetSm2P256V1()
	pub.X, _ = new(big.Int).SetString(constDefine.MACHINECODE_SM2_PUB_X, 16)
	pub.Y, _ = new(big.Int).SetString(constDefine.MACHINECODE_SM2_PUB_Y, 16)

	//fmt.Printf("d:%s\n", hex.EncodeToString(priv.D.Bytes()))
	//fmt.Printf("x:%s\n", hex.EncodeToString(pub.X.Bytes()))
	//fmt.Printf("y:%s\n", hex.EncodeToString(pub.Y.Bytes()))

	cipherText, err := sm2.Encrypt(pub, gson, sm2.C1C3C2)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Printf("generate unique id:\n%s\n", hex.EncodeToString(cipherText))

	plainText, err := sm2.Decrypt(priv, cipherText, sm2.C1C3C2)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Println(string(plainText))
}
