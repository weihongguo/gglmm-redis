package redis

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var stopChan chan<- interface{}

func TestMessageQueue(t *testing.T) {
	mq := NewMessageQueue("tcp", "127.0.0.1:6379", 5, 10, 3, "message-queue-test")
	defer mq.Close()

	sendCount := 5

	for i := 0; i < sendCount; i++ {
		mq.Publish([]byte(fmt.Sprintf("%dtest", i)))
	}

	go handleMessage(mq, 5, "handler1")
	go handleMessage(mq, 5, "handler2")

	for i := 0; i < sendCount; i++ {
		t := time.NewTimer(time.Second * 1)
		select {
		case <-t.C:
			mq.Publish([]byte(fmt.Sprintf("test%d", i)))
		}
	}

	stopChan <- 1
}

func handleMessage(mq *MessageQueue, timeout int, info string) {
	var err error
	stopChan, err = mq.Subscribe(func(message []byte, err error) {
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(info + ": " + string(message))
	}, timeout)
	if err != nil {
		log.Fatal(err)
	}
}
