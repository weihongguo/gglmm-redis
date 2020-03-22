package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	factory := NewFactory("tcp", "127.0.0.1:6379", 5, 10, 3)
	defer factory.Close()

	counter := factory.NewCounter("test-counter")
	defer counter.Del()
	counter.Zero()
	counter.Incr()
	counterResult, err := counter.Get()
	if err != nil {
		t.Fatal(err)
	} else if counterResult != 1 {
		t.Fatal("counter value error")
	}

	var i int64

	toper := factory.NewToper("test-toper", 5)
	defer toper.Del()
	for i = 0; i < 10; i++ {
		toper.Push(i)
	}
	ids, err := toper.Range()
	if err != nil {
		t.Fatal(err)
	} else {
		if len(ids) != 5 {
			t.Fatal("toper len error")
		} else {
			t.Log(ids)
		}
	}

	hoter := factory.NewHoter("test-hoter")
	defer hoter.Del()
	for i = 0; i < 10; i++ {
		hoter.Add(i, i)
	}
	ids, scores, err := hoter.RevRangeByScore(5)
	if err != nil {
		t.Fatal(err)
	} else {
		if len(ids) != 5 || len(scores) != 5 {
			t.Fatal("hoter len error")
		} else {
			t.Log(ids, scores)
		}
	}

	/*
		cacher := factory.NewCacher(10)
		cacher.Set("test-aa", 1)
		cacherResult, err := cacher.GetInt("test-aa")
		if err != nil {
			t.Fatal(err)
		} else if cacherResult != 1 {
			t.Fatal("cacher value error")
		}

		mq := factory.NewMessageQueue("test-message-queue")
		mq.Push([]byte("test"))
		go mq.BPop(func(message []byte, err error) {
			if err != nil {
				log.Fatal(err)
			} else if string(message) != "test" {
				log.Fatal("mq value error")
			}
		}, 1)
		time.Sleep(time.Second)
	*/
}
