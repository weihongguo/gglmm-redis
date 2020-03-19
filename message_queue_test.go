package redis

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var stop chan<- interface{}

func TestMessageQueue(t *testing.T) {
	mq := NewMessageQueue("tcp", "127.0.0.1:6379", 5, 10, 3, "message-queue-test")
	defer mq.Close()

	sendCount := 5

	for i := 0; i < sendCount; i++ {
		mq.Push([]byte(fmt.Sprintf("%dtest", i)))
	}

	go handlerMessage(mq, 5, "handler1")
	go handlerMessage(mq, 5, "handler2")

	for i := 0; i < sendCount; i++ {
		t := time.NewTimer(time.Second * 1)
		select {
		case <-t.C:
			mq.Push([]byte(fmt.Sprintf("test%d", i)))
		}
	}

	stop <- 1
}

func handlerMessage(mq *MessageQueue, timeout int, info string) {
	var err error
	stop, err = mq.BPop(func(message []byte, err error) {
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
