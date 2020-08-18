package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"
	"github.com/xukgo/gsvcCtrlCenter/models"
	"github.com/xukgo/gsvcCtrlCenter/sm2"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "c", "", "set config path")
	if len(filePath) == 0 {
		//fmt.Println("请用\"-c\"指定配置文件路径")
		//os.Exit(-1)
		filePath, _ = filepath.Abs("reg.xml")
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件内容失败", err.Error())
		os.Exit(-1)
	}
	lic := new(models.LicenseConfig)
	err = lic.FillWithXml(contents)
	if err != nil {
		fmt.Println("解析文件内容失败", err.Error())
		os.Exit(-1)
	}

	if !lic.CheckValid() {
		os.Exit(-1)
	}

	fmt.Println("解析配置如下：")
	lic.Print()

	gson := lic.ToJson()
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.GetSm2P256V1()
	pub.X, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_X, 16)
	pub.Y, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_Y, 16)

	fmt.Println("原始json：\n", string(gson))
	cipherText, err := sm2.Encrypt(pub, gson, sm2.C1C3C2)
	if err != nil {
		fmt.Println("加密lic出错", err.Error())
		os.Exit(-1)
	}
	fmt.Printf("lic:\n%s\n", hex.EncodeToString(cipherText))
}
