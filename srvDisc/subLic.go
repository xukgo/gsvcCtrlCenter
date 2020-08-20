package srvDisc

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/xukgo/gsvcCtrlCenter/models"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"

	"github.com/xukgo/gsvcCtrlCenter/logUtil"
	"go.uber.org/zap"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var localLicRegDict = make(map[string]*SubLicRegisterInfo)
var localLicRegDictLocker = new(sync.RWMutex)

type SubLicRegisterInfo struct {
	*LicRegisterInfo
	Reversion int64
}

func newSubLicRegisterInfo(reversion int64, regInfo *LicRegisterInfo) *SubLicRegisterInfo {
	model := new(SubLicRegisterInfo)
	model.LicRegisterInfo = regInfo
	model.Reversion = reversion
	return model
}

func (this *Repo) watchLicNode() {
	prefix := formatNodeKeepAlivePrefix()
	for {
		watchChan := this.client.Watch(clientv3.WithRequireLeader(context.TODO()), prefix, clientv3.WithPrefix())
		if watchChan == nil {
			time.Sleep(time.Second)
			continue
		}

		for watchResponse := range watchChan {
			this.updateLicNodeByEvents(watchResponse.Events)
		}
	}
}

func (this *Repo) updateLicNodeByEvents(events []*clientv3.Event) {
	localLicRegDictLocker.Lock()
	defer localLicRegDictLocker.Unlock()

	changeCount := 0
	for _, event := range events {
		logUtil.LoggerCommon.Info("sub lic event", zap.ByteString("key", event.Kv.Key))
		switch event.Type {
		case mvccpb.PUT:
			upsertLicNode(event)
			changeCount++
			break
		case mvccpb.DELETE:
			removeLicNode(event)
			changeCount++
			break
		}
	}

	if changeCount > 0 {
		code, desc := parseLicResult()
		this.clientUpdateLicResult(code, desc)
	}
}

func parseLicResult() (int, string) {
	dictLen := len(localLicRegDict)
	if dictLen == 0 {
		return constDefine.RETCODE_NODE_COUNT_LESS, constDefine.RETDESC_NODE_COUNT_LESS
	}

	encryKey := ""
	var licConfig *models.LicenseConfig = nil
	var validCount = 0
	dtNow := time.Now()
	for _, v := range localLicRegDict {
		if len(encryKey) == 0 {
			encryKey = v.LicEncData
		}
		if encryKey != v.LicEncData {
			return constDefine.RETCODE_NODE_NOTMATCH_NODE, constDefine.RETDESC_NODE_NOTMATCH_NODE
		}
		if dtNow.Sub(time.Unix(v.Timestamp, 0)).Seconds() > constDefine.LIC_EXPIRED_SEC {
			continue
		}
		if licConfig == nil {
			licConfig = v.LicenseConfig
		}
		validCount++
	}
	if licConfig == nil || validCount == 0 {
		return constDefine.RETCODE_NODE_COUNT_LESS, constDefine.RETDESC_NODE_COUNT_LESS
	}
	maxCount := len(licConfig.MachineInfos)
	if validCount <= maxCount/2 {
		return constDefine.RETCODE_NODE_COUNT_LESS, constDefine.RETDESC_NODE_COUNT_LESS
	}
	return constDefine.RETCODE_SUCCESS, constDefine.RETDESC_SUCCESS
}

func removeLicNode(event *clientv3.Event) {
	eventKey := string(event.Kv.Key)
	eventRevision := event.Kv.ModRevision
	v, find := localLicRegDict[eventKey]
	if !find {
		return
	}
	if v.Reversion <= eventRevision {
		delete(localLicRegDict, eventKey)
		return
	}
}

func upsertLicNode(event *clientv3.Event) {
	eventKey := string(event.Kv.Key)
	eventRevision := event.Kv.ModRevision
	v, find := localLicRegDict[eventKey]
	if !find {
		v, err := parseLicValue(event.Kv.Value)
		if err != nil {
			return
		}
		localLicRegDict[eventKey] = newSubLicRegisterInfo(eventRevision, v)
		return
	}

	if v.Reversion <= eventRevision {
		v, err := parseLicValue(event.Kv.Value)
		if err != nil {
			return
		}
		localLicRegDict[eventKey] = newSubLicRegisterInfo(eventRevision, v)
		return
	}
}

func parseLicValue(data []byte) (*LicRegisterInfo, error) {
	data, err := hex.DecodeString(string(data))
	if err != nil {
		logUtil.LoggerCommon.Error("sub lic value decode error", zap.Error(err))
		return nil, err
	}

	model := new(LicRegisterInfo)
	err = model.DecryptJson(data)
	if err != nil {
		logUtil.LoggerCommon.Error("sub lic value LicRegisterInfo DecryptJson  error", zap.Error(err))
		return nil, err
	}
	return model, nil
}
