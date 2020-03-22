package redis

import "testing"

func TestToper(t *testing.T) {
	toper := NewToper("tcp", "127.0.0.1:6379", 5, 10, 3, "toper-test", 5)
	defer toper.Close()
	defer toper.Del()

	toper.Push(1)
	toperResult, err := toper.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if toperResult != 1 {
		t.Fatal("toper value error")
	}

	var i int64 = 0
	for i < 10 {
		i++
		toper.Push(i)
	}

	toperResults, err := toper.Range()
	if err != nil {
		t.Fatal(err)
	}
	if len(toperResults) != 5 {
		t.Fatalf("toper value len error %d\n", len(toperResults))
	}
	t.Log(toperResults)
}
