package redis

import "testing"

type Test struct {
	Value string `json:"value"`
}

func TestCacher(t *testing.T) {

	cacher := NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 10)
	defer cacher.Close()

	err := cacher.Set("key", "value")
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueString, err := cacher.GetString("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueString != "value" {
		t.Fatalf("value not match：" + valueString)
	}

	err = cacher.Set("key", 1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueInt, err := cacher.GetInt("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueInt != 1 {
		t.Fatalf("value not match")
	}

	err = cacher.Set("key", 1.0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueFloat, err := cacher.GetFloat64("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueFloat != 1.0 {
		t.Fatalf("value not match")
	}

	test := &Test{Value: "value"}

	err = cacher.Set("key", test)
	if err != nil {
		t.Fatalf(err.Error())
	}

	objValue := &Test{}
	err = cacher.GetObj("key", objValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if test.Value != objValue.Value {
		t.Fatalf(objValue.Value)
	}
}
