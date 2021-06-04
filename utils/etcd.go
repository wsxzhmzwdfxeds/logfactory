package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Esxhost string `json:"esx_host"`
	Vmcluster string `json:"vmcluster"`
	Remoteconsole string `json:"remote_console"`
}

func InitEtcd(addr string, timeout int) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{addr},
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect etcd success.")
	return
}

func etcdGet(hostname string) LogEntry {
	var logentry LogEntry
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	getResp, err := cli.Get(ctx, hostname)
	cancel()
	if err != nil {
		fmt.Printf("Etcd get failed, err : %v\n", err)
	}
	for _, ev := range getResp.Kvs {
		err = json.Unmarshal(ev.Value, &logentry)
	}
	return logentry
}

func ProcessLog(hostname string, dat map[string]interface{}){
	logentry :=  etcdGet(hostname)
	dat["esxhost"] = logentry.Esxhost
	dat["remote_console"] = logentry.Remoteconsole
	dat["vmcluster"] = logentry.Vmcluster
}