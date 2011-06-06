package ketama

import (
	"fmt"
	"testing"
	"strconv"
)

func TestGetInfo(t *testing.T) {
	ring := NewRing(200)

	nodes := map[string]int{
		"test1.server.com": 1,
		"test2.server.com": 1,
		"test3.server.com": 2,
		"test4.server.com": 5,
	}

	for k, v := range nodes {
		ring.AddNode(k, v)
	}

	ring.Bake()

	m := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		m[ring.Hash("test value"+strconv.Uitoa(uint(i)))]++
	}

	for k, _ := range nodes {
		fmt.Println(k, m[k])
	}
}
