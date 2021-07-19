package utils

import (
	"encoding/json"
	"fmt"
	//"encoding/json"
	"github.com/Shopify/sarama"
	"log"
)

var (
	LOGSOURCE = "logsource"
)

type Consumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		//myMap := Json2map(message.Value)
		//newMsg, _ := json.Marshal(myMap)
		var dat map[string]interface{}
		if err := json.Unmarshal(message.Value, &dat); err != nil {
			fmt.Printf("log format error, can not be unamrsha to map : %v\n", err)
		}

		hostname, ok := dat[LOGSOURCE]
		//hostname, ok := dat["host"]
		if ok {
			ProcessLog(hostname.(string), dat)
		} else {
			fmt.Println("logsource not exists.")
		}
		fmt.Printf("-----%v", dat)
		ld := LogData{
			Data: dat,
			//Data:  string(message.Value),
			Topic: message.Topic,
		}
		//fmt.Println(ld)
		SendToEsChan(&ld)
		session.MarkMessage(message, "")
	}

	return nil
}
