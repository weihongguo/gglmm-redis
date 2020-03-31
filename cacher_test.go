package redis

import (
	"strconv"
	"testing"
)

type Test struct {
	Value string `json:"value"`
}

func TestCacher(t *testing.T) {
	cacher := NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 10)
	defer cacher.Close()

	err := cacher.Set("key", "value")
	if err != nil {
		t.Fatal(err)
	}

	valueString, err := cacher.GetString("key")
	if err != nil {
		t.Fatal(err)
	}

	if valueString != "value" {
		t.Fatalf("value not match：%s != 'value'" + valueString)
	}

	err = cacher.Set("key", 1)
	if err != nil {
		t.Fatal(err)
	}

	valueInt, err := cacher.GetInt("key")
	if err != nil {
		t.Fatal(err)
	}

	if valueInt != 1 {
		t.Fatalf("value not match %d != 1", valueInt)
	}

	err = cacher.Set("key", 1.0)
	if err != nil {
		t.Fatal(err)
	}

	valueFloat, err := cacher.GetFloat64("key")
	if err != nil {
		t.Fatal(err)
	}

	if valueFloat != 1.0 {
		t.Fatalf("value not match %f != 1.0", valueFloat)
	}

	test := &Test{Value: "value"}

	err = cacher.Set("key", test)
	if err != nil {
		t.Fatal(err)
	}

	objValue := &Test{}
	err = cacher.GetObj("key", objValue)
	if err != nil {
		t.Fatal(err)
	}

	if test.Value != objValue.Value {
		t.Fatal(objValue.Value)
	}

	delCount, err := cacher.Del("key")
	if err != nil {
		t.Fatal(err)
	}
	if delCount != 1 {
		t.Fatalf("del count error: %d != 1", delCount)
	}

	for i := 0; i < 10; i++ {
		cacher.Set("key:"+strconv.Itoa(i), i)
	}
	delCount, err = cacher.DelPattern("key:*")
	if err != nil {
		t.Fatal(err)
	}
	if delCount != 10 {
		t.Fatalf("del count error: %d != 10", delCount)
	}
}
