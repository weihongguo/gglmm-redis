package redis

import "testing"

func TestHoter(t *testing.T) {
	hoter := NewHoter("tcp", "127.0.0.1:6379", 5, 10, 3, "hoter-test")
	defer hoter.Close()
	defer hoter.Del()

	var i int64 = 0
	for i < 10 {
		i++
		if _, err := hoter.Add(i, i); err != nil {
			t.Fatal(err)
		}
	}

	if card, err := hoter.Card(); err != nil {
		t.Fatal(err)
	} else if card != 10 {
		t.Fatal("card error")
	}

	if count, err := hoter.Count(3, 7); err != nil {
		t.Fatal(err)
	} else if count != 5 {
		t.Fatal("count error")
	}

	if rem, err := hoter.Rem(5); err != nil {
		t.Fatal(err)
	} else if rem != 1 {
		t.Fatal("rem error")
	}

	if rem, err := hoter.RemRangeByScore(5, 6); err != nil {
		t.Fatal(err)
	} else if rem != 1 {
		t.Fatal("rem error")
	}

	if ids, scores, err := hoter.RangeByScore(5); err != nil {
		t.Fatal(err)
	} else {
		t.Log(ids, scores)
	}

	if ids, scores, err := hoter.RevRangeByScore(5); err != nil {
		t.Fatal(err)
	} else {
		t.Log(ids, scores)
	}
}
