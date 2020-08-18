package models

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"time"

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
	ExpireTimestamp int64                      `xml:"-" json:"expire"` //单位秒
}

func (this LicenseConfig) ToJson() []byte {
	gson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(this)
	return gson
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
		fmt.Println("节点id数量不允许为空")
		return false
	}
	if len(this.LimitServices) == 0 {
		fmt.Println("限制服务配置不允许为空")
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
		fmt.Printf("授权节点，唯一码=%s，硬盘序列号=%s，cpuID=%s，生成时间=%s\n",
			item.MachineID, item.DiskSerialNumber, item.CpuId, dt.Format("2006-01-02 15:04:05"))
	}
	for _, item := range this.LimitServices {
		fmt.Printf("限制服务数量，%s=%d\n", item.Name, item.Count)
	}
	fmt.Printf("过期时间戳:%d\n", this.ExpireTimestamp)
	dt := time.Unix(this.ExpireTimestamp, 0)
	fmt.Printf("过期时间计算:%s\n", dt.Format("2006-01-02 15:04:05"))
}
