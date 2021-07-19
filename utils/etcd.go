package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"github.com/muesli/cache2go"
)

var (
	cli *clientv3.Client
	cache = cache2go.Cache("logs")
	ESXHOST = "esxhost"
	REMOTECONSOLE = "remote_console"
	VMCLUSTER = "vmcluster"
	VMSERVICE = "vmservice"
	DATASTORESTORAGE = "datastore_storage"
)

type LogEntry struct {
	Esxhost string `json:"esx_host"`
	Vmcluster string `json:"vmcluster"`
	Remoteconsole string `json:"remote_console"`
	VMServices []string `json:"service"`
	DatastoreStorage []string `json:"datastore_storage"`
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

func etcdGet(hostname string) *LogEntry {
	logentry := &LogEntry{
		Esxhost:       "None",
		Vmcluster:     "None",
		Remoteconsole: "None",
		VMServices: 		[]string{},
		DatastoreStorage: []string{},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	getResp, err := cli.Get(ctx, hostname)
	cancel()
	if err != nil {
		fmt.Printf("Etcd get failed, err : %v\n", err)
	}
	for _, ev := range getResp.Kvs {
		fmt.Println("etcd value: ", string(ev.Value))
		err = json.Unmarshal(ev.Value, &logentry)
	}
	return logentry
}

func ProcessLog(hostname string, dat map[string]interface{}){
	var logentry *LogEntry
	res, err := cache.Value(hostname)

	if err == nil {
		fmt.Println("Found value in cache" )
		logentry = res.Data().(*LogEntry)
		fmt.Printf("--- in cache %v", logentry)

	} else {
		fmt.Println("Error retrieving value from cache:", err)
		logentry = etcdGet(hostname)
		fmt.Printf("--- in etcd %v", logentry)
		cache.Add(hostname, 60*time.Minute, logentry)
	}

	dat[ESXHOST] = logentry.Esxhost
	dat[REMOTECONSOLE] = logentry.Remoteconsole
	dat[VMCLUSTER] = logentry.Vmcluster
	dat[VMSERVICE] = logentry.VMServices
	dat[DATASTORESTORAGE] = logentry.DatastoreStorage
}