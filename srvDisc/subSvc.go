package srvDisc

import (
	"context"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var localSvcNodeDict = make(map[string]map[string]int64)
var localSvcNodeDictLocker = new(sync.RWMutex)

func (this *Repo) SubScribe() error {
	go this.watchLicNode()
	for _, item := range AppRegTemplate.LicenseConfig.LimitServices {
		go this.watchSubs(item.Name)
	}
	return nil
}

func (this *Repo) watchSubs(srvName string) {
	servicePrefix := formatServicePrefix(srvName)

	for {
		watchChan := this.client.Watch(clientv3.WithRequireLeader(context.TODO()), servicePrefix, clientv3.WithPrefix())
		if watchChan == nil {
			time.Sleep(time.Second)
			continue
		}

		for watchResponse := range watchChan {
			this.updateServiceNodeByEvents(srvName, watchResponse.Events)
		}
	}
}
func (this *Repo) updateServiceNodeByEvents(srvName string, events []*clientv3.Event) {
	localSvcNodeDictLocker.Lock()
	defer localSvcNodeDictLocker.Unlock()

	for _, event := range events {
		//logUtil.LoggerCommon.Info("sub svc event", zap.ByteString("key", event.Kv.Key))
		switch event.Type {
		case mvccpb.PUT:
			upsertServiceNode(srvName, event)
			break
		case mvccpb.DELETE:
			removeServiceNode(srvName, event)
			break
		}
	}
}

func removeServiceNode(srvName string, event *clientv3.Event) {
	eventKey := string(event.Kv.Key)
	eventRevision := event.Kv.ModRevision
	dict, find := localSvcNodeDict[srvName]
	if !find {
		dict = make(map[string]int64)
		localSvcNodeDict[srvName] = dict
	}
	v, find := dict[eventKey]
	if !find {
		return
	}
	if v <= eventRevision {
		delete(dict, eventKey)
		return
	}
}

func upsertServiceNode(srvName string, event *clientv3.Event) {
	eventKey := string(event.Kv.Key)
	eventRevision := event.Kv.ModRevision
	dict, find := localSvcNodeDict[srvName]
	if !find {
		dict = make(map[string]int64)
		localSvcNodeDict[srvName] = dict
	}
	v, find := dict[eventKey]
	if !find {
		dict[eventKey] = eventRevision
		return
	}
	if v <= eventRevision {
		dict[eventKey] = eventRevision
		return
	}
}
