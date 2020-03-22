package redis

import "testing"

func TestCounter(t *testing.T) {
	counter := NewCounter("tcp", "127.0.0.1:6379", 5, 10, 3, "counter-test")
	defer counter.Close()
	defer counter.Del()

	oldValue, err := counter.Zero()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(oldValue)

	value, err := counter.Incr()
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("value error: %d %d\n", value, 1)
	}

	value, err = counter.IncrBy(2)
	if err != nil {
		t.Fatal(err)
	}
	if value != 3 {
		t.Fatalf("value error: %d %d\n", value, 3)
	}

	value, err = counter.DecrBy(2)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("value error: %d %d\n", value, 1)
	}

	value, err = counter.Decr()
	if err != nil {
		t.Fatal(err)
	}
	if value != 0 {
		t.Fatalf("value error: %d %d\n", value, 0)
	}

	value, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if value != 0 {
		t.Fatalf("value error: %d %d\n", value, 0)
	}
}
