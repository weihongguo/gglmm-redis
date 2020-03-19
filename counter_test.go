package redis

import "testing"

func TestCounter(t *testing.T) {
	counter := NewCounter("tcp", "127.0.0.1:6379", 5, 10, 3, "counter-test")

	oldValue, err := counter.Zero()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(oldValue)

	value, err := counter.Increase()
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("value error: %d %d\n", value, 1)
	}

	value, err = counter.IncreaseBy(2)
	if err != nil {
		t.Fatal(err)
	}
	if value != 3 {
		t.Fatalf("value error: %d %d\n", value, 3)
	}

	value, err = counter.DecreaseBy(2)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("value error: %d %d\n", value, 1)
	}

	value, err = counter.Decrease()
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
