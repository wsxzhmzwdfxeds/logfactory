package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

type LogData struct {
	//Data string `json:"data"`
	Data map[string]interface{} `json:"data"`
	Topic string `json:"topic"`
}

func format(msg *sarama.ConsumerMessage) *LogData {
	var ld = new(LogData)
	err := json.Unmarshal(msg.Value, ld)
	if err != nil {
		panic(err)
	}
	return ld
}

func Json2map(jsonStr []byte) map[string]interface{} {
	var mapResult map[string]interface{}
	err := json.Unmarshal(jsonStr, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}
	return mapResult
}
