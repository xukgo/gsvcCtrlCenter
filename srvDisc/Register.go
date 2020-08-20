package srvDisc

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/xukgo/gsvcCtrlCenter/constDefine"

	"go.uber.org/zap"

	"github.com/xukgo/gsvcCtrlCenter/logUtil"

	"github.com/coreos/etcd/clientv3"
)

func (this *Repo) Register() {
	var lease *clientv3.LeaseGrantResponse = nil
	var err error

	for {
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
		lease, err = this.client.Grant(ctx, 6)
		if err != nil || lease == nil {
			time.Sleep(time.Second)
			continue
		}

		this.KeepaliveLease(lease)
	}
}

func (this *Repo) KeepaliveLease(lease *clientv3.LeaseGrantResponse) {
	keepaliveChan, err := this.client.KeepAlive(context.TODO(), lease.ID) //这里需要一直不断，context不允许设置超时
	if err != nil || keepaliveChan == nil {
		time.Sleep(time.Second)
		return
	}
	updateDt := time.Unix(0, 0)
	for {
		select {
		case keepaliveResponse, ok := <-keepaliveChan:
			if !ok || keepaliveResponse == nil {
				logUtil.LoggerCommon.Error(">>>error keepaliveResponse")
				return
			}
			//fmt.Println("keepaliveResponse", keepaliveResponse)
			break
		default:
			dtNow := time.Now()
			if dtNow.Sub(AppExipreTime).Seconds() > 0 {
				logUtil.LoggerCommon.Error("license已过期")
				this.clientUpdateLicResult(constDefine.RETCODE_EXPIRED, constDefine.RETDESC_EXPIRED)
				os.Exit(-1)
			}

			//每隔更新一次
			if dtNow.Sub(updateDt).Seconds() > constDefine.LIC_UPDATE_INTERVAL_SEC {
				err = this.clientUpdateLeaseContent(lease)
				if err == nil {
					updateDt = dtNow
					time.Sleep(100 * time.Millisecond)
				} else {
					time.Sleep(1000 * time.Millisecond)
				}
			}
		}
	}
}

func (this *Repo) clientUpdateLeaseContent(lease *clientv3.LeaseGrantResponse) error {
	key := formatNodeKeepAliveKey()
	AppRegTemplate.Timestamp = time.Now().Unix()
	encryptData, err := AppRegTemplate.EncryptJson()
	if err != nil {
		logUtil.LoggerCommon.Error("AppRegTemplate EncryptJson error", zap.Error(err))
		os.Exit(-1)
	}

	valueStr := hex.EncodeToString(encryptData)
	_, err = this.client.Put(context.TODO(), key, valueStr, clientv3.WithLease(lease.ID))
	return err
}

func (this *Repo) clientUpdateLicResult(code int, desc string) error {
	result := new(LicResultInfo)
	result.Code = code
	result.Description = desc
	result.LicenseConfig = AppRegTemplate.LicenseConfig
	result.Timestamp = time.Now().Unix()
	encryptData, err := result.EncryptJson()
	if err != nil {
		logUtil.LoggerCommon.Error("LicResultInfo EncryptJson error", zap.Error(err))
		return err
	}

	valueStr := hex.EncodeToString(encryptData)
	logUtil.LoggerCommon.Info("update lic result", zap.Int("code", code), zap.String("desc", desc))
	_, err = this.client.Put(context.TODO(), formatLicResultKey(), valueStr)
	return err
}

//
//func (this *Repo) fillRegMoudleInfo(info *LicRegisterInfo, beforeRegisterFunc BeforeRegisterFunc) {
//	if beforeRegisterFunc != nil {
//		beforeRegisterFunc(info)
//	}
//	info.Global.State = "online"
//	info.Global.RefreshTimestamp(time.Now())
//}
