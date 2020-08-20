package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xukgo/gsvcCtrlCenter/models"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "c", "", "set config path")
	flag.Parse()

	if len(filePath) == 0 {
		fmt.Println("请用\"-c\"指定配置文件路径")
		os.Exit(-1)
		//filePath, _ = filepath.Abs("reg.xml")
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
	fmt.Println("原始json：\n", string(lic.ToPrettyJson()))

	encryptData, err := lic.EncryptJson()
	if err != nil {
		fmt.Println("加密lic数据出错", err.Error())
		os.Exit(-1)
	}
	fmt.Printf("\nlicense:\n%s\n", hex.EncodeToString(encryptData))
}
