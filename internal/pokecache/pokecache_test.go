package pokecache

import (
	"testing"
	"time"
)

func TestPokeCache(t *testing.T) {
	cases := []struct {
		inputKey string
		inputVal string
		expected string
	}{
		{
			inputKey: "pikachu",
			inputVal: "{}",
			expected: "{}",
		},
		{
			inputKey: "starmie",
			inputVal: "{\"experience\":\"1234\"}",
			expected: "{\"experience\":\"1234\"}",
		},
	}

	for _, c := range cases {
		cache := NewCache(1 * time.Minute)
		cache.Add(c.inputKey, []byte(c.inputVal))
		val, exist := cache.Get(c.inputKey)
		if !exist {
			t.Errorf("Expected key %s doesn't exist", c.inputKey)
		}

		if string(val) != c.expected {
			t.Errorf("Value for key %s differs from expected.\na:%s\ne:%s\n", c.inputKey, string(val), c.expected)
		}
	}
}
