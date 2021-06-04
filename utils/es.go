package utils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"runtime"
	"strings"
	"time"
)

var (
	// ch cache: esCfg.chanSize. set to 10000
	client *elastic.Client
	ch = make(chan *LogData, 10000)
)

func InitEs(address string) (err error) {
	if !strings.HasPrefix(address, "http") {
		address = "http://"+address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return
	}
	fmt.Println("connect to es success.")
	for i:=0;i<runtime.NumCPU();i++{
		go sendToEs()
	}
	return 
}

func SendToEsChan(msg *LogData) {
	ch <- msg
}

func sendToEs() {
	for{
		select {
		case msg := <- ch:
			put, err := client.Index().
				Index(msg.Topic).
				BodyJson(msg.Data).
				//BodyString(msg.Data).
				Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Printf("Indexed data %s to index %s, type %s\n", put.Id, put.Index, put.Index)
		default:
			time.Sleep(time.Second)
		}
	}
}