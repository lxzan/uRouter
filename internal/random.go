package internal

import (
	"math/rand"
	"sync"
	"time"
)

type RandomString struct {
	mu     sync.Mutex
	r      *rand.Rand
	layout string
}

var (
	AlphabetNumeric = &RandomString{
		layout: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
		mu:     sync.Mutex{},
	}
)

func (c *RandomString) Generate(n int) []byte {
	c.mu.Lock()
	var b = make([]byte, n, n)
	var length = len(c.layout)
	for i := 0; i < n; i++ {
		var idx = c.r.Intn(length)
		b[i] = c.layout[idx]
	}
	c.mu.Unlock()
	return b
}