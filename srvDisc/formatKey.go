package srvDisc

import (
	"fmt"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"

	"github.com/xukgo/gsaber/utils/randomUtil"
)

var KeepAliveRdnId = randomUtil.NewUpperHexString(12)

func formatNodeKeepAliveKey() string {
	key := fmt.Sprintf("lic.voice.%s", KeepAliveRdnId)
	return key
}
func formatNodeKeepAlivePrefix() string {
	return "lic.voice."
}

func formatServicePrefix(name string) string {
	key := fmt.Sprintf("registry.voice.%s", name)
	return key
}

func formatLicResultKey() string {
	return constDefine.LIC_RESULT_KEY
}
