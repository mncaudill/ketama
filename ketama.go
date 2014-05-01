package ketama

import (
	"crypto/sha1"
	"sort"
	"strconv"
)

// HashRing is the interface that describes HashRing implementation behavior.
type HashRing interface {
	// Add takes a single element and multiplier. The Bake method must be called when adding elements is complete.
	Add(element string, multiplier int)
	// AddAll adds all of the elements of a given string to int map. It then calls the Bake method.
	AddAll(elements map[string]int)
	// Bake finalizes the internal structures used to prepare the HashRing for reads.
	Bake()
	// Hash looks up the element to be used for a given key.
	Hash(key string) string
}

type ketamaHashRing struct {
	defaultSpots int
	ticks        tickArray
	length       int
}

type node struct {
	node string
	hash uint
}

type tickArray []node

func NewRing(n int) HashRing {
	hashRing := new(ketamaHashRing)
	hashRing.defaultSpots = n
	return hashRing
}

func (p tickArray) Len() int {
	return len(p)
}

func (p tickArray) Less(i, j int) bool {
	return p[i].hash < p[j].hash
}

func (p tickArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p tickArray) Sort() {
	sort.Sort(p)
}

func (h *ketamaHashRing) AddAll(elements map[string]int) {
	for element, multiplier := range elements {
		h.Add(element, multiplier)
	}
	h.Bake()
}

func (h *ketamaHashRing) Add(element string, multiplier int) {
	tSpots := h.defaultSpots * multiplier
	hash := sha1.New()
	for i := 1; i <= tSpots; i++ {
		hash.Write([]byte(element + ":" + strconv.Itoa(i)))
		hashBytes := hash.Sum(nil)

		n := &node{
			node: element,
			hash: uint(hashBytes[19]) | uint(hashBytes[18])<<8 | uint(hashBytes[17])<<16 | uint(hashBytes[16])<<24,
		}

		h.ticks = append(h.ticks, *n)
		hash.Reset()
	}
}

func (h *ketamaHashRing) Bake() {
	h.ticks.Sort()
	h.length = len(h.ticks)
}

func (h *ketamaHashRing) Hash(key string) string {
	hash := sha1.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	v := uint(hashBytes[19]) | uint(hashBytes[18])<<8 | uint(hashBytes[17])<<16 | uint(hashBytes[16])<<24
	i := sort.Search(h.length, func(i int) bool { return h.ticks[i].hash >= v })

	if i == h.length {
		i = 0
	}

	return h.ticks[i].node
}
