package srvDisc

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Repo struct {
	config *ConfRoot
	client *clientv3.Client //etcd客户端
}

func (this *Repo) InitFromPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	srvConf := new(ConfRoot)
	err = srvConf.FillWithXml(content)
	if err != nil {
		return err
	}

	this.config = srvConf

	this.client, err = clientv3.New(clientv3.Config{
		Endpoints:   srvConf.Endpoints,
		DialTimeout: time.Duration(2) * time.Second,
	})
	return nil
}

func (this *Repo) StartRegister() error {
	if this.config == nil {
		return fmt.Errorf("register conf is nil")
	}
	go this.Register()
	return nil
}

func (this *Repo) StartSubscribe() error {
	if this.config == nil {
		return fmt.Errorf("register conf is nil")
	}
	err := this.SubScribe()
	return err
}

//
//func (this *Repo) GetSubsNames() []string {
//	subsconf := this.config.SubScribeConf
//	if subsconf == nil {
//		return nil
//	}
//	if len(subsconf.Services) == 0 {
//		return nil
//	}
//
//	arr := make([]string, 0, len(subsconf.Services))
//	for idx := range subsconf.Services {
//		arr = append(arr, subsconf.Services[idx].Name)
//	}
//	return arr
//}

//func (this *Repo) initSubsNodeCache(subSrvInfos []SubBasicInfo) {
//	serviceCount := len(subSrvInfos)
//	if serviceCount <= 0 {
//		return
//	}
//
//	this.subsNodeCache = make(map[string]*SubSrvNodeList)
//	for m := 0; m < serviceCount; m++ {
//		srvNodeList := new(SubSrvNodeList)
//		srvNodeList.SubBasicInfo = *NewSubSrvBasicInfo(subSrvInfos[m].Name, subSrvInfos[m].Version, subSrvInfos[m].Namespace)
//		srvNodeList.NodeInfos = make([]*SrvNodeInfo, 0, 1)
//		this.subsNodeCache[subSrvInfos[m].Name] = srvNodeList
//	}
//}
//
////随机打乱数组
//func randomSortSlice(arr []*LicRegisterInfo) {
//	if len(arr) <= 0 || len(arr) == 1 {
//		return
//	}
//
//	for i := len(arr) - 1; i > 0; i-- {
//		num := randomUtil.NewInt32(0, int32(i+1))
//		arr[i], arr[num] = arr[num], arr[i]
//	}
//}

//func GetSrvDiscover() *Repo {
//	return srvDiscoverInstance
//}
//
//func GetSrvDiscoverConf() *ConfRoot {
//	return &srvDiscoverConf
//}

//func (this *ServiceDiscovery) TriggerRegister() {
//	this.registerHupChan <- true
//}
//

//func (this *Repo) GetConfig() *ConfRoot {
//	return this.config
//}

//func NewSrvDiscover(endpoints []string, options ...SdOption) (*Repo, error) {
//	if len(endpoints) == 0 {
//		return nil, fmt.Errorf("endpoints addrs is empty")
//	}
//
//	serviceDiscovery := &Repo{
//		Endpoints: endpoints,
//	}
//	serviceDiscovery.Timeout = DEFAULT_CONN_TIMEOUT * time.Millisecond
//
//	for _, op := range options {
//		op(serviceDiscovery)
//	}
//
//	var err error
//	serviceDiscovery.client, err = clientv3.New(clientv3.Config{
//		Endpoints:   endpoints,
//		DialTimeout: serviceDiscovery.Timeout,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	serviceDiscovery.subsNodeCache = make(map[string]*SubSrvNodeList)
//	serviceDiscovery.locker = &sync.RWMutex{}
//
//	return serviceDiscovery, nil
//}
