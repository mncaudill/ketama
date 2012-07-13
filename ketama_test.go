package ketama

import (
	"fmt"
	"strconv"
	"testing"
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
		m[ring.Hash("test value"+strconv.FormatUint(uint64(i), 10))]++
	}

	for k := range nodes {
		fmt.Println(k, m[k])
	}
}
