package main

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/xukgo/gsvcCtrlCenter/srvDisc"

	"go.uber.org/zap"

	"github.com/xukgo/gsvcCtrlCenter/logUtil"

	"github.com/xukgo/gsaber/compon/procUnique"

	"github.com/xukgo/gsvcCtrlCenter/machineCode"

	"github.com/xukgo/gsaber/utils/fileUtil"

	"github.com/xukgo/gsvcCtrlCenter/models"
)

func main() {
	logUtil.LoggerInit()
	var procLocker = procUnique.NewLocker("d9VoiceSvcCtrl")
	err := procLocker.Lock()
	if err != nil {
		logUtil.LoggerCommon.Error("应用不允许多实例运行")
		os.Exit(-1)
	}
	defer procLocker.Unlock()

	licPath := fileUtil.GetAbsUrl("conf/voice.lic")
	licContent, err := ioutil.ReadFile(licPath)
	if err != nil {
		logUtil.LoggerCommon.Error("读取lic文件失败", zap.Error(err))
		os.Exit(-1)
	}

	licStr := string(licContent)
	licStr = strings.ReplaceAll(licStr, " ", "")
	licStr = strings.ReplaceAll(licStr, "\n", "")
	licData, err := hex.DecodeString(licStr)
	if err != nil {
		logUtil.LoggerCommon.Error("lic内容不是十六进制字符串格式", zap.Error(err))
	}

	lic := new(models.LicenseConfig)
	err = lic.DecryptJson(licData)
	if err != nil {
		logUtil.LoggerCommon.Error("解析lic失败", zap.Error(err))
		os.Exit(-1)
	}

	tsNow := time.Now().Unix()
	if tsNow > lic.ExpireTimestamp {
		logUtil.LoggerCommon.Error("license已过期")
		os.Exit(-1)
	}

	uinfo := machineCode.Info()
	if !checkMachineCodeMatch(lic, uinfo) {
		logUtil.LoggerCommon.Error("设备不匹配license，请向厂家咨询")
		os.Exit(-1)
	}
	logUtil.LoggerCommon.Info("lic文件有效")

	regInfo := new(srvDisc.LicRegisterInfo)
	regInfo.LicEncData = licStr
	regInfo.LicenseConfig = lic
	regInfo.MInfo = uinfo
	err = srvDisc.Start(regInfo)
	if err != nil {
		os.Exit(-1)
	}

	exitChan := make(chan bool)
	<-exitChan
}

func checkMachineCodeMatch(lic *models.LicenseConfig, uinfo *models.MachineUniqueInfo) bool {
	for _, item := range lic.MachineInfos {
		if uinfo.MachineID == item.MachineID && uinfo.DiskSerialNumber == item.DiskSerialNumber && uinfo.CpuId == item.CpuId {
			return true
		}
	}
	return false
}
