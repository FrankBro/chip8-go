package main

import (
	"math/rand"
	"time"
)

type Random interface {
	Int7() uint8
}

type GoRandom struct {
	rand *rand.Rand
}

func NewGoRandom() GoRandom {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)
	random := GoRandom{
		rand: rand,
	}
	return random
}

func (rand GoRandom) Int7() uint8 {
	value := rand.rand.Intn(256)
	return uint8(value)
}
