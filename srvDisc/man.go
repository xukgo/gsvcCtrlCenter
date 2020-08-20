package srvDisc

import (
	"time"

	"github.com/xukgo/gsaber/utils/fileUtil"
	"github.com/xukgo/gsvcCtrlCenter/logUtil"
	"go.uber.org/zap"
)

var AppExipreTime time.Time
var AppRegTemplate *LicRegisterInfo

func Start(regTemp *LicRegisterInfo) error {
	fileUrl := fileUtil.GetAbsUrl("conf/service.xml")
	repo := new(Repo)
	err := repo.InitFromPath(fileUrl)
	if err != nil {
		logUtil.LoggerCommon.Error("etcd service init file error", zap.Error(err))
		return err
	}

	AppRegTemplate = regTemp
	AppExipreTime = time.Unix(regTemp.LicenseConfig.ExpireTimestamp, 0)

	repo.StartRegister()
	repo.StartSubscribe()

	return nil
}
