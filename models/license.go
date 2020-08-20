package models

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"
	"github.com/xukgo/gsvcCtrlCenter/sm2"

	"github.com/xukgo/gsvcCtrlCenter/util"

	jsoniter "github.com/json-iterator/go"
)

type LimitServiceCountConfig struct {
	Name  string `xml:"name,attr" json:"name"`
	Count int    `xml:"count,attr" json:"count"`
}

type LicenseConfig struct {
	XMLName         xml.Name                   `xml:"Config" json:"-"`
	NodeUidArray    []string                   `xml:"NodeUid" json:"-"`
	MachineInfos    []*MachineUniqueInfo       `xml:"-" json:"machine"`
	LimitServices   []*LimitServiceCountConfig `xml:"LimitService" json:"limitSvc"`
	ExpireString    string                     `xml:"Expire" json:"-"`
	CallParallel    int                        `xml:"CallParallel" json:"CallParallel"`
	ExpireTimestamp int64                      `xml:"-" json:"expire"` //单位秒
}

func (this LicenseConfig) ToPrettyJson() []byte {
	gson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(this, "", "   ")
	return gson
}
func (this LicenseConfig) ToJson() []byte {
	gson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(this)
	return gson
}

func (this *LicenseConfig) EncryptJson() ([]byte, error) {
	gson := this.ToJson()
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.GetSm2P256V1()
	pub.X, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_X, 16)
	pub.Y, _ = new(big.Int).SetString(constDefine.LIC_SM2_PUB_Y, 16)

	cipherText, err := sm2.Encrypt(pub, gson, sm2.C1C3C2)
	return cipherText, err
}

func (this *LicenseConfig) DecryptJson(data []byte) error {
	priv := new(sm2.PrivateKey)
	priv.Curve = sm2.GetSm2P256V1()
	priv.D, _ = new(big.Int).SetString(constDefine.LIC_SM2_PRIV_D, 16)

	plainText, err := sm2.Decrypt(priv, data, sm2.C1C3C2)
	if err != nil {
		return err
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(plainText, this)
}

func (this *LicenseConfig) FillWithXml(xstr []byte) error {
	err := xml.Unmarshal(xstr, this)
	if err != nil {
		return err
	}

	if len(this.ExpireString) > 0 {
		dt, err := util.Str2Time(this.ExpireString)
		if err != nil {
			return fmt.Errorf("解析时间格式出错，%w", err)
		}
		this.ExpireTimestamp = dt.Unix()
	}

	this.MachineInfos = make([]*MachineUniqueInfo, 0, len(this.NodeUidArray))
	for _, item := range this.NodeUidArray {
		item = strings.ReplaceAll(item, " ", "")
		item = strings.ReplaceAll(item, "\n", "")
		buff, err := hex.DecodeString(item)
		if err != nil {
			return fmt.Errorf("NodeUid不是十六进制字符串格式，%w", err)
		}
		info := new(MachineUniqueInfo)
		err = info.DecryptJson(buff)
		if err != nil {
			return fmt.Errorf("NodeUid解密出错，%w", err)
		}
		this.MachineInfos = append(this.MachineInfos, info)
	}
	return nil
}

func (this *LicenseConfig) CheckValid() bool {
	if len(this.NodeUidArray) == 0 {
		fmt.Println("许可节点id数量不允许为空")
		return false
	}
	if len(this.LimitServices) == 0 {
		fmt.Println("限制服务配置不允许为空")
		return false
	}
	if this.CallParallel <= 0 {
		fmt.Println("限制呼叫并发数量必须大于0")
		return false
	}
	if len(this.ExpireString) == 0 {
		fmt.Println("过期时间不允许为空")
		return false
	}

	return true
}

func (this *LicenseConfig) Print() {
	fmt.Printf("节点id数量=%d\n", len(this.NodeUidArray))
	for _, item := range this.MachineInfos {
		dt := time.Unix(0, item.Timestamp)
		fmt.Printf("许可节点，唯一码=%s，硬盘序列号=%s，cpuID=%s，生成时间=%s\n",
			item.MachineID, item.DiskSerialNumber, item.CpuId, dt.Format("2006-01-02 15:04:05"))
	}
	for _, item := range this.LimitServices {
		fmt.Printf("限制服务数量，%s=%d\n", item.Name, item.Count)
	}
	fmt.Printf("呼叫并发:%d\n", this.CallParallel)
	fmt.Printf("过期时间戳:%d\n", this.ExpireTimestamp)
	dt := time.Unix(this.ExpireTimestamp, 0)
	fmt.Printf("过期时间预估:%s\n", dt.Format("2006-01-02 15:04:05"))
}
