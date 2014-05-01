package ketama

import (
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	ring := NewRing(200)

	nodes := map[string]int{
		"test1.server.com": 1,
		"test2.server.com": 1,
		"test3.server.com": 2,
		"test4.server.com": 5,
	}

	for k, v := range nodes {
		ring.Add(k, v)
	}

	ring.Bake()

	results := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		results[ring.Hash("test value"+strconv.FormatUint(uint64(i), 10))]++
	}

	if compareElements("test4.server.com", "test3.server.com", results, t) {
		return
	}
	if compareElements("test3.server.com", "test2.server.com", results, t) {
		return
	}
	if compareElements("test3.server.com", "test1.server.com", results, t) {
		return
	}
}

func TestAddAll(t *testing.T) {
	ring := NewRing(200)

	ring.AddAll(map[string]int{
		"test1.server.com": 1,
		"test2.server.com": 1,
		"test3.server.com": 2,
		"test4.server.com": 5,
	})

	results := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		results[ring.Hash("test value"+strconv.FormatUint(uint64(i), 10))]++
	}

	if compareElements("test4.server.com", "test3.server.com", results, t) {
		return
	}
	if compareElements("test3.server.com", "test2.server.com", results, t) {
		return
	}
	if compareElements("test3.server.com", "test1.server.com", results, t) {
		return
	}
}

func compareElements(element1, element2 string, results map[string]int, t *testing.T) bool {
	element1Count := countForElement(element1, results)
	element2Count := countForElement(element2, results)
	if element1Count < element2Count {
		t.Errorf("element '%s' (%d) should be accessed more than '%s' (%d).", element1, element1Count, element2, element2Count)
		return true
	}
	return false
}

func countForElement(element string, results map[string]int) int {
	count, hasCount := results[element]
	if hasCount {
		return count
	}
	return 0
}
